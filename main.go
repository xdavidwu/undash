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
	"k8s.io/apimachinery/pkg/runtime/schema"

	undashctx "github.com/xdavidwu/undash/internal/context"
	"github.com/xdavidwu/undash/internal/dashboard"
	undashhttp "github.com/xdavidwu/undash/internal/http"
	"github.com/xdavidwu/undash/internal/kubernetes"
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

type verberMappedMeta struct {
	singular   string
	namespaced bool
}

var (
	coreNamespacedResources = map[string]kindMeta{
		"services": {
			kind:     "Service",
			singular: "service",
			listKind: "ServiceList",
		},
	}

	verberClientTypeDefaultGV       = schema.GroupVersion{Group: "", Version: "v1"}
	verberClientTypeAppsGV          = schema.GroupVersion{Group: "apps", Version: "v1"}
	verberClientTypeAutoscalingGV   = schema.GroupVersion{Group: "autoscaling", Version: "v1"}
	verberClientTypeNetworkingGV    = schema.GroupVersion{Group: "networking.k8s.io", Version: "v1"}
	verberClientTypeBatchGV         = schema.GroupVersion{Group: "batch", Version: "v1"}
	verberClientTypeBetaBatchGV     = schema.GroupVersion{Group: "batch", Version: "v1beta1"}
	verberClientTypeAPIExtensionsGV = schema.GroupVersion{Group: "apiextensions.k8s.io", Version: "v1"}
	verberClientTypeStorageGV       = schema.GroupVersion{Group: "storage.k8s.io", Version: "v1"}
	verberClientTypeRbacGV          = schema.GroupVersion{Group: "rbac.authorization.k8s.io", Version: "v1"}
	// dashboard extensions custom resource
	verberClientTypePluginsGV = schema.GroupVersion{Group: "dashboard.k8s.io", Version: "v1alpha1"}

	// v2 has its name => client mapping and parse CRD for custom resources,
	// while v3 seems to rely on discovery + dynamic client
	verberMappedResources = map[schema.GroupVersionResource]verberMappedMeta{
		verberClientTypeDefaultGV.WithResource("configmaps"): {
			singular:   "configmap",
			namespaced: true,
		},
		verberClientTypeAppsGV.WithResource("daemonsets"): {
			singular:   "daemonset",
			namespaced: true,
		},
		verberClientTypeAppsGV.WithResource("deployments"): {
			singular:   "deployment",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("events"): {
			singular:   "event",
			namespaced: true,
		},
		verberClientTypeAutoscalingGV.WithResource("horizontalpodautoscalers"): {
			singular:   "horizontalpodautoscaler",
			namespaced: true,
		},
		verberClientTypeNetworkingGV.WithResource("ingresses"): {
			singular:   "ingress",
			namespaced: true,
		},
		verberClientTypeNetworkingGV.WithResource("ingressclasses"): {
			singular:   "ingressclass",
			namespaced: false,
		},
		verberClientTypeBatchGV.WithResource("jobs"): {
			singular:   "job",
			namespaced: true,
		},
		verberClientTypeBetaBatchGV.WithResource("cronjobs"): {
			singular:   "cronjob",
			namespaced: true,
		},
		// XXX singular on dashboard map?
		verberClientTypeDefaultGV.WithResource("limitranges"): {
			singular:   "limitrange",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("namespaces"): {
			singular:   "namespace",
			namespaced: false,
		},
		verberClientTypeDefaultGV.WithResource("nodes"): {
			singular:   "node",
			namespaced: false,
		},
		verberClientTypeDefaultGV.WithResource("persistentvolumeclaims"): {
			singular:   "persistentvolumeclaim",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("persistentvolumes"): {
			singular:   "persistentvolume",
			namespaced: false,
		},
		verberClientTypeAPIExtensionsGV.WithResource("customresourcedefinitions"): {
			singular:   "customresourcedefinition",
			namespaced: false,
		},
		verberClientTypeDefaultGV.WithResource("pods"): {
			singular:   "pod",
			namespaced: true,
		},
		verberClientTypeAppsGV.WithResource("replicasets"): {
			singular:   "replicaset",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("replicationcontrollers"): {
			singular:   "replicationcontroller",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("resourcequotas"): {
			singular:   "resourcequota",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("secrets"): {
			singular:   "secret",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("services"): {
			singular:   "service",
			namespaced: true,
		},
		verberClientTypeDefaultGV.WithResource("serviceaccounts"): {
			singular:   "serviceaccount",
			namespaced: true,
		},
		verberClientTypeAppsGV.WithResource("statefulsets"): {
			singular:   "statefulset",
			namespaced: true,
		},
		verberClientTypeStorageGV.WithResource("storageclasses"): {
			singular:   "storageclass",
			namespaced: false,
		},
		verberClientTypeDefaultGV.WithResource("endpoints"): {
			singular:   "endpoint",
			namespaced: true,
		},
		verberClientTypeNetworkingGV.WithResource("networkpolicies"): {
			singular:   "networkpolicy",
			namespaced: true,
		},
		verberClientTypeRbacGV.WithResource("clusterroles"): {
			singular:   "clusterrole",
			namespaced: false,
		},
		verberClientTypeRbacGV.WithResource("clusterrolebindings"): {
			singular:   "clusterrolebinding",
			namespaced: false,
		},
		verberClientTypeRbacGV.WithResource("roles"): {
			singular:   "role",
			namespaced: true,
		},
		verberClientTypeRbacGV.WithResource("rolebindings"): {
			singular:   "rolebinding",
			namespaced: true,
		},
		verberClientTypePluginsGV.WithResource("plugins"): {
			singular:   "plugin",
			namespaced: true,
		},
		// TODO custom resource via parsing definitions?
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

		if undashhttp.AcceptsExactMediaType(r, undashhttp.MetaV1TableJSON) {
			gvk := res.GroupVersionKind()
			if toTable, ok := kubernetes.UnstructuredListToTableFuncs[gvk]; ok {
				table, err := toTable(res)
				if err != nil {
					return nil, fmt.Errorf("cannot convert list to table: %w", err)
				}

				w.Header().Set("Content-Type", undashhttp.MetaV1TableJSON.String())
				return table, nil
			} else {
				undashctx.GetLogger(ctx).WarnContext(ctx, "list table func not registered", "gvk", gvk.String())
			}
		}

		return res, nil
	})
}

func apiPrefix(gv schema.GroupVersion) string {
	if gv.Group == "" {
		return "/api/" + gv.Version
	}
	return "/apis/" + gv.Group + "/" + gv.Version
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

	for resource := range coreNamespacedResources {
		mux.Handle(
			"/api/v1/namespaces/{namespace}/"+resource,
			coreResourceNamespacedListHandlerFor(resource),
		)
	}

	// TODO advertise these at discovery
	for gvr, meta := range verberMappedResources {
		if meta.namespaced {
			mux.Handle(
				apiPrefix(gvr.GroupVersion())+"/namespaces/{namespace}/"+gvr.Resource+"/{name}",
				&httputil.ReverseProxy{
					Transport: undashhttp.NoExplicitCompression(undashhttp.RequestLog(http.DefaultTransport)),
					Rewrite: func(r *httputil.ProxyRequest) {
						url, err := url.Parse(fmt.Sprintf(
							"%s/api/v1/_raw/%s/namespace/%s/name/%s",
							upstream,
							meta.singular,
							r.In.PathValue("namespace"),
							r.In.PathValue("name"),
						))
						if err != nil {
							panic(err)
						}

						r.Out.URL = url
					},
					ModifyResponse: undashhttp.ChainModifyResponse(
						undashhttp.ErrorResponseAsMetaV1Status,
						undashhttp.RewriteObjectAsTableIfRequested,
					),
				},
			)
		} else {
			mux.Handle(
				apiPrefix(gvr.GroupVersion())+"/"+gvr.Resource+"/{name}",
				&httputil.ReverseProxy{
					Transport: undashhttp.NoExplicitCompression(undashhttp.RequestLog(http.DefaultTransport)),
					Rewrite: func(r *httputil.ProxyRequest) {
						url, err := url.Parse(fmt.Sprintf(
							"%s/api/v1/_raw/%s/name/%s",
							upstream,
							meta.singular,
							r.In.PathValue("name"),
						))
						if err != nil {
							panic(err)
						}

						r.Out.URL = url
					},
					ModifyResponse: undashhttp.ChainModifyResponse(
						undashhttp.ErrorResponseAsMetaV1Status,
						undashhttp.RewriteObjectAsTableIfRequested,
					),
				},
			)
		}
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
