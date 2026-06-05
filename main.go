package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	apidiscoveryv2 "k8s.io/api/apidiscovery/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	upstream = "http://localhost:9090"

	apidiscoveryv2MimeType = runtime.ContentTypeJSON + ";g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList"
)

type kindMeta struct {
	kind     string
	singular string
	listKind string
}

var (
	coreNamespacedResources = map[string]kindMeta{
		"services": {
			kind:     "Service",
			singular: "service",
			listKind: "ServiceList",
		},
	}
)

type jsonHandler[T any] func(w http.ResponseWriter, r *http.Request) (T, error)

func (j jsonHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	res, err := j(w, r)
	if err != nil {
		var errRes metav1.Status

		if apiStatus, ok := err.(apierrors.APIStatus); ok {
			errRes = apiStatus.Status()
		} else {
			errRes = metav1.Status{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Status",
				},
				Status:  metav1.StatusFailure,
				Message: err.Error(),
			}
		}
		headers.Set("Content-Type", runtime.ContentTypeJSON)

		if errRes.Code != 0 {
			w.WriteHeader(int(errRes.Code))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(errRes)
	}

	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", runtime.ContentTypeJSON)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

func coreResourceNamespacedListHandlerFor(resource string) http.Handler {
	return jsonHandler[*metav1.List](func(w http.ResponseWriter, r *http.Request) (*metav1.List, error) {
		ns := r.PathValue("namespace")
		kind := coreNamespacedResources[resource]

		listRes, err := http.Get(fmt.Sprintf("%s/api/v1/%s/%s", upstream, kind.singular, ns))
		if err != nil {
			return nil, fmt.Errorf("cannot list objects: %w", err)
		}
		defer listRes.Body.Close()

		body, err := io.ReadAll(listRes.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read object list: %w", err)
		}

		listObj := map[string][]struct {
			ObjectMeta struct {
				Name string `json:"name"`
			} `json:"objectMeta"`
		}{}
		json.Unmarshal(body, &listObj)

		res := &metav1.List{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       kind.listKind,
			},
			Items: []runtime.RawExtension{},
		}
		for _, obj := range listObj[resource] {
			name := obj.ObjectMeta.Name

			realObjRes, err := http.Get(fmt.Sprintf("%s/api/v1/_raw/%s/namespace/%s/name/%s", upstream, kind.singular, ns, name))
			if err != nil {
				return nil, fmt.Errorf("cannot get real object: %w", err)
			}
			defer realObjRes.Body.Close()

			realObj, err := io.ReadAll(realObjRes.Body)
			if err != nil {
				return nil, fmt.Errorf("cannot read real object: %w", err)
			}

			res.Items = append(res.Items, runtime.RawExtension{Raw: realObj})
		}

		return res, nil
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle(
		"/api",
		jsonHandler[*apidiscoveryv2.APIGroupDiscoveryList](func(w http.ResponseWriter, r *http.Request) (*apidiscoveryv2.APIGroupDiscoveryList, error) {
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
		jsonHandler[*apidiscoveryv2.APIGroupDiscoveryList](func(w http.ResponseWriter, r *http.Request) (*apidiscoveryv2.APIGroupDiscoveryList, error) {
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
			},
		)
	}

	l, _ := net.Listen("tcp", "localhost:9091")
	http.Serve(l, mux)
}
