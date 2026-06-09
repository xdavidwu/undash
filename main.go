package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	apidiscoveryv2 "k8s.io/api/apidiscovery/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/xdavidwu/undash/internal/dashboard"
	undashhttp "github.com/xdavidwu/undash/internal/http"
	"github.com/xdavidwu/undash/internal/kubernetes"
)

const (
	upstream = "http://localhost:9090"

	apidiscoveryv2MimeType = runtime.ContentTypeJSON + ";g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList"
)

type kindMeta struct {
	kind        string
	singular    string
	listKind    string
	listToTable func(*unstructured.UnstructuredList) (*metav1.Table, error)
}

var (
	coreNamespacedResources = map[string]kindMeta{
		"services": {
			kind:        "Service",
			singular:    "service",
			listKind:    "ServiceList",
			listToTable: kubernetes.UnstructuredListToTableFunc(kubernetes.V1ServiceListToTable),
		},
	}
)

func coreResourceNamespacedListHandlerFor(resource string) http.Handler {
	return undashhttp.JSONHandler(func(w http.ResponseWriter, r *http.Request) (runtime.Object, error) {
		ns := r.PathValue("namespace")
		kind := coreNamespacedResources[resource]
		client := undashhttp.NewDefaultClient()
		ctx := r.Context()

		listRes, err := client.Call(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/v1/%s/%s", upstream, kind.singular, ns),
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot list objects: %w", err)
		}
		defer listRes.Body.Close()

		body, err := io.ReadAll(listRes.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read object list: %w", err)
		}

		listObj := dashboard.ListUnmarshaler{Resource: resource}
		if err := json.Unmarshal(body, &listObj); err != nil {
			return nil, fmt.Errorf("cannot unmarshal object list: %w", err)
		}

		res := &unstructured.UnstructuredList{}
		res.SetAPIVersion("v1")
		res.SetKind(kind.listKind)
		for _, obj := range listObj.ObjectMetas {
			realObjRes, err := client.Call(
				ctx,
				http.MethodGet,
				fmt.Sprintf("%s/api/v1/_raw/%s/namespace/%s/name/%s", upstream, kind.singular, ns, obj.Name),
				nil,
			)
			if err != nil {
				return nil, fmt.Errorf("cannot get real object: %w", err)
			}
			defer realObjRes.Body.Close()

			realObjBytes, err := io.ReadAll(realObjRes.Body)
			if err != nil {
				return nil, fmt.Errorf("cannot read real object: %w", err)
			}

			obj, _, err := unstructured.UnstructuredJSONScheme.Decode(realObjBytes, nil, &unstructured.Unstructured{})
			if err != nil {
				return nil, fmt.Errorf("cannot decode real object: %w", err)
			}

			res.Items = append(res.Items, *obj.(*unstructured.Unstructured))
		}

		if undashhttp.AcceptsExactMediaType(r, undashhttp.MetaV1TableJSON) && kind.listToTable != nil {
			table, err := kind.listToTable(res)
			if err != nil {
				return nil, fmt.Errorf("cannot convert list to table: %w", err)
			}

			w.Header().Set("Content-Type", undashhttp.MetaV1TableJSON.String())
			return table, nil
		}

		return res, nil
	})
}

func main() {
	logger := slog.Default()
	ctx := context.Background()

	mux := http.NewServeMux()
	mux.Handle(
		"/api",
		undashhttp.JSONHandler(func(w http.ResponseWriter, r *http.Request) (*apidiscoveryv2.APIGroupDiscoveryList, error) {
			// XXX supporting apidiscoveryv2 (1.30+) since it's easier
			// but dashboard targets 1.25
			discoveries := []apidiscoveryv2.APIResourceDiscovery{}
			for resource, kind := range coreNamespacedResources {
				discoveries = append(discoveries, apidiscoveryv2.APIResourceDiscovery{
					Resource: resource,
					ResponseKind: &metav1.GroupVersionKind{
						Kind: kind.kind,
					},
					Scope:            apidiscoveryv2.ScopeNamespace,
					SingularResource: kind.singular,
					Verbs:            []string{"list", "get", "delete", "update"},
				})
			}
			list := &apidiscoveryv2.APIGroupDiscoveryList{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apidiscovery.k8s.io/v2",
					Kind:       "APIGroupDiscoveryList",
				},
				Items: []apidiscoveryv2.APIGroupDiscovery{
					{
						Versions: []apidiscoveryv2.APIVersionDiscovery{
							{
								Version:   "v1",
								Resources: discoveries,
								Freshness: apidiscoveryv2.DiscoveryFreshnessCurrent,
							},
						},
					},
				},
			}

			w.Header().Set("Content-Type", apidiscoveryv2MimeType)
			return list, nil
		}),
	)
	mux.Handle(
		"/apis",
		undashhttp.JSONHandler(func(w http.ResponseWriter, r *http.Request) (*apidiscoveryv2.APIGroupDiscoveryList, error) {
			// XXX supporting apidiscoveryv2 (1.30+) since it's easier
			// but dashboard targets 1.25
			list := &apidiscoveryv2.APIGroupDiscoveryList{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apidiscovery.k8s.io/v2",
					Kind:       "APIGroupDiscoveryList",
				},
				Items: []apidiscoveryv2.APIGroupDiscovery{},
			}

			w.Header().Set("Content-Type", apidiscoveryv2MimeType)
			return list, nil
		}),
	)

	for resource, kind := range coreNamespacedResources {
		mux.Handle(
			"/api/v1/namespaces/{namespace}/"+resource,
			coreResourceNamespacedListHandlerFor(resource),
		)

		mux.Handle(
			"/api/v1/namespaces/{namespace}/"+resource+"/{name}",
			&httputil.ReverseProxy{
				Transport: undashhttp.NoExplicitCompression(undashhttp.RequestLog(http.DefaultTransport)),
				Rewrite: func(r *httputil.ProxyRequest) {
					url, err := url.Parse(fmt.Sprintf(
						"%s/api/v1/_raw/%s/namespace/%s/name/%s",
						upstream,
						kind.singular,
						r.In.PathValue("namespace"),
						r.In.PathValue("name"),
					))
					if err != nil {
						panic(err)
					}

					r.Out.URL = url
				},
				ModifyResponse: undashhttp.ErrorResponseAsMetaV1Status,
			},
		)
	}

	l, err := net.Listen("tcp", "localhost:9091")
	if err != nil {
		logger.ErrorContext(ctx, "cannot listen", "error", err)
		os.Exit(1)
	}
	logger.InfoContext(ctx, "listening", "addr", l.Addr())

	if err := http.Serve(l, undashhttp.InjectLogger(undashhttp.AccessLog(mux), logger)); err != nil {
		logger.ErrorContext(ctx, "cannot serve http", "error", err)
		os.Exit(1)
	}
}
