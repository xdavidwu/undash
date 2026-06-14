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

type verberMappedMeta struct {
	singular   string
	namespaced bool
	kind       string
	listKind   string
	itemsKey   string // resource if unset
}

var (
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
			kind:       "ConfigMap",
			listKind:   "ConfigMapList",
			itemsKey:   "items",
		},
		verberClientTypeAppsGV.WithResource("daemonsets"): {
			singular:   "daemonset",
			namespaced: true,
			kind:       "DaemonSet",
		},
		verberClientTypeAppsGV.WithResource("deployments"): {
			singular:   "deployment",
			namespaced: true,
			kind:       "Deployment",
		},
		verberClientTypeDefaultGV.WithResource("events"): {
			singular:   "event",
			namespaced: true,
			kind:       "Event",
			listKind:   "EventList",
		},
		verberClientTypeAutoscalingGV.WithResource("horizontalpodautoscalers"): {
			singular:   "horizontalpodautoscaler",
			namespaced: true,
			kind:       "HorizontalPodAutoscaler",
		},
		verberClientTypeNetworkingGV.WithResource("ingresses"): {
			singular:   "ingress",
			namespaced: true,
			kind:       "Ingress",
		},
		verberClientTypeNetworkingGV.WithResource("ingressclasses"): {
			singular:   "ingressclass",
			namespaced: false,
			kind:       "IngressClass",
		},
		verberClientTypeBatchGV.WithResource("jobs"): {
			singular:   "job",
			namespaced: true,
			kind:       "Job",
		},
		verberClientTypeBetaBatchGV.WithResource("cronjobs"): {
			singular:   "cronjob",
			namespaced: true,
			kind:       "CronJob",
		},
		// XXX singular on dashboard map?
		verberClientTypeDefaultGV.WithResource("limitranges"): {
			singular:   "limitrange",
			namespaced: true,
			kind:       "LimitRange",
		},
		verberClientTypeDefaultGV.WithResource("namespaces"): {
			singular:   "namespace",
			namespaced: false,
			kind:       "Namespace",
			listKind:   "NamespaceList",
		},
		verberClientTypeDefaultGV.WithResource("nodes"): {
			singular:   "node",
			namespaced: false,
			kind:       "Node",
			listKind:   "NodeList",
		},
		verberClientTypeDefaultGV.WithResource("persistentvolumeclaims"): {
			singular:   "persistentvolumeclaim",
			namespaced: true,
			kind:       "PersistentVolumeClaim",
			listKind:   "PersistentVolumeClaimList",
			itemsKey:   "items",
		},
		verberClientTypeDefaultGV.WithResource("persistentvolumes"): {
			singular:   "persistentvolume",
			namespaced: false,
			kind:       "PersistentVolume",
			listKind:   "PersistentVolumeList",
			itemsKey:   "items",
		},
		verberClientTypeAPIExtensionsGV.WithResource("customresourcedefinitions"): {
			singular:   "customresourcedefinition",
			namespaced: false,
			kind:       "CustomResourceDefinition",
		},
		verberClientTypeDefaultGV.WithResource("pods"): {
			singular:   "pod",
			namespaced: true,
			kind:       "Pod",
			listKind:   "PodList",
		},
		verberClientTypeAppsGV.WithResource("replicasets"): {
			singular:   "replicaset",
			namespaced: true,
			kind:       "ReplicaSet",
		},
		verberClientTypeDefaultGV.WithResource("replicationcontrollers"): {
			singular:   "replicationcontroller",
			namespaced: true,
			kind:       "ReplicationController",
			listKind:   "ReplicationControllerList",
			itemsKey:   "replicationControllers",
		},
		verberClientTypeDefaultGV.WithResource("resourcequotas"): {
			singular:   "resourcequota",
			namespaced: true,
			kind:       "ResourceQuota",
		},
		verberClientTypeDefaultGV.WithResource("secrets"): {
			singular:   "secret",
			namespaced: true,
			kind:       "Secret",
			listKind:   "SecretList",
		},
		verberClientTypeDefaultGV.WithResource("services"): {
			singular:   "service",
			namespaced: true,
			kind:       "Service",
			listKind:   "ServiceList",
		},
		verberClientTypeDefaultGV.WithResource("serviceaccounts"): {
			singular:   "serviceaccount",
			namespaced: true,
			kind:       "ServiceAccount",
			listKind:   "ServiceAccountList",
			itemsKey:   "items",
		},
		verberClientTypeAppsGV.WithResource("statefulsets"): {
			singular:   "statefulset",
			namespaced: true,
			kind:       "StatefulSet",
		},
		verberClientTypeStorageGV.WithResource("storageclasses"): {
			singular:   "storageclass",
			namespaced: false,
			kind:       "StorageClass",
		},
		verberClientTypeDefaultGV.WithResource("endpoints"): {
			singular:   "endpoint", // XXX not really?
			namespaced: true,
			kind:       "Endpoints",
		},
		verberClientTypeNetworkingGV.WithResource("networkpolicies"): {
			singular:   "networkpolicy",
			namespaced: true,
			kind:       "NetworkPolicy",
		},
		verberClientTypeRbacGV.WithResource("clusterroles"): {
			singular:   "clusterrole",
			namespaced: false,
			kind:       "ClusterRole",
		},
		verberClientTypeRbacGV.WithResource("clusterrolebindings"): {
			singular:   "clusterrolebinding",
			namespaced: false,
			kind:       "ClusterRoleBinding",
		},
		verberClientTypeRbacGV.WithResource("roles"): {
			singular:   "role",
			namespaced: true,
			kind:       "Role",
		},
		verberClientTypeRbacGV.WithResource("rolebindings"): {
			singular:   "rolebinding",
			namespaced: true,
			kind:       "RoleBinding",
		},
		verberClientTypePluginsGV.WithResource("plugins"): {
			singular:   "plugin",
			namespaced: true,
			kind:       "Plugin",
		},
		// TODO custom resource via parsing definitions?
	}
)

func listHandlerFor(gvr schema.GroupVersionResource) http.Handler {
	return undashhttp.JSONHandler(func(w http.ResponseWriter, r *http.Request) (runtime.Object, error) {
		ns := r.PathValue("namespace")
		meta := verberMappedResources[gvr]
		client := undashhttp.NewDefaultClient()
		ctx := r.Context()

		path := ""
		if ns != "" {
			path = fmt.Sprintf("%s/api/v1/%s/%s", upstream, meta.singular, ns)
		} else {
			path = fmt.Sprintf("%s/api/v1/%s", upstream, meta.singular)
		}

		listRes, err := client.Call(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot list objects: %w", err)
		}
		defer listRes.Body.Close()

		body, err := io.ReadAll(listRes.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read object list: %w", err)
		}

		itemsKey := gvr.Resource
		if meta.itemsKey != "" {
			itemsKey = meta.itemsKey
		}

		listObj := dashboard.ListUnmarshaler{ItemsKey: itemsKey}
		if err := json.Unmarshal(body, &listObj); err != nil {
			return nil, fmt.Errorf("cannot unmarshal object list: %w", err)
		}

		res := &unstructured.UnstructuredList{}
		res.SetAPIVersion("v1")
		res.SetKind(meta.listKind)
		for _, obj := range listObj.ObjectMetas {
			path := ""
			if obj.Namespace != "" {
				path = fmt.Sprintf("%s/api/v1/_raw/%s/namespace/%s/name/%s", upstream, meta.singular, obj.Namespace, obj.Name)
			} else {
				path = fmt.Sprintf("%s/api/v1/_raw/%s/name/%s", upstream, meta.singular, obj.Name)
			}

			realObjRes, err := client.Call(ctx, http.MethodGet, path, nil)
			// TODO fight with MSG_DASHBOARD_EXCLUSIVE_RESOURCE_ERROR?
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

// XXX supporting apidiscoveryv2 (1.30+) since it's easier
// but dashboard targets 1.25
func buildAPIGroupDiscoveries() map[string]*apidiscoveryv2.APIGroupDiscovery {
	versionDiscoveries := map[schema.GroupVersion]*apidiscoveryv2.APIVersionDiscovery{}
	for gvr, meta := range verberMappedResources {
		versionDiscovery, ok := versionDiscoveries[gvr.GroupVersion()]
		if !ok {
			versionDiscoveries[gvr.GroupVersion()] = &apidiscoveryv2.APIVersionDiscovery{
				Version:   gvr.Version,
				Freshness: apidiscoveryv2.DiscoveryFreshnessCurrent,
			}
			versionDiscovery = versionDiscoveries[gvr.GroupVersion()]
		}

		verbs := []string{"get", "delete", "update"}
		if meta.listKind != "" {
			verbs = append(verbs, "list")
		}

		scope := apidiscoveryv2.ScopeNamespace
		if !meta.namespaced {
			scope = apidiscoveryv2.ScopeCluster
		}

		// TODO shortnames
		// codegen api types => internal types => storage func (r *REST) ShortNames() []string

		versionDiscovery.Resources = append(versionDiscovery.Resources, apidiscoveryv2.APIResourceDiscovery{
			Resource: gvr.Resource,
			ResponseKind: &metav1.GroupVersionKind{
				Kind: meta.kind,
			},
			Scope:            scope,
			SingularResource: meta.singular,
			Verbs:            verbs,
		})
	}

	res := map[string]*apidiscoveryv2.APIGroupDiscovery{}
	for gv, versionDiscovery := range versionDiscoveries {
		groupDiscovery, ok := res[gv.Group]
		if !ok {
			res[gv.Group] = &apidiscoveryv2.APIGroupDiscovery{
				ObjectMeta: metav1.ObjectMeta{
					Name: gv.Group,
				},
			}
			groupDiscovery = res[gv.Group]
		}

		// we don't really need to sort here since every kind in dashboard is of single api version
		groupDiscovery.Versions = append(groupDiscovery.Versions, *versionDiscovery)
	}

	return res
}

func buildAPIGroupDiscoveryLists() (core, apis apidiscoveryv2.APIGroupDiscoveryList) {
	gvk := schema.GroupVersionKind{
		Group:   "apidiscovery.k8s.io",
		Version: "v2",
		Kind:    "APIGroupDiscoveryList",
	}
	core.SetGroupVersionKind(gvk)
	apis.SetGroupVersionKind(gvk)

	for group, discovery := range buildAPIGroupDiscoveries() {
		if group == "" {
			core.Items = append(core.Items, *discovery)
		} else {
			apis.Items = append(apis.Items, *discovery)
		}
	}
	return
}

func main() {
	logger := slog.Default()
	ctx := context.Background()

	core, apis := buildAPIGroupDiscoveryLists()

	mux := http.NewServeMux()
	mux.Handle(
		"/api",
		undashhttp.JSONHandler(func(w http.ResponseWriter, r *http.Request) (*apidiscoveryv2.APIGroupDiscoveryList, error) {
			w.Header().Set("Content-Type", apidiscoveryv2MimeType)
			return &core, nil
		}),
	)
	mux.Handle(
		"/apis",
		undashhttp.JSONHandler(func(w http.ResponseWriter, r *http.Request) (*apidiscoveryv2.APIGroupDiscoveryList, error) {
			w.Header().Set("Content-Type", apidiscoveryv2MimeType)
			return &apis, nil
		}),
	)

	for gvr, meta := range verberMappedResources {
		prefix := apiPrefix(gvr.GroupVersion())
		if meta.listKind != "" {
			mux.Handle(
				"GET "+prefix+"/"+gvr.Resource,
				listHandlerFor(gvr),
			)
		}

		if meta.namespaced {
			mux.Handle(
				prefix+"/namespaces/{namespace}/"+gvr.Resource+"/{name}",
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

			if meta.listKind != "" {
				mux.Handle(
					"GET "+prefix+"/namespaces/{namespace}/"+gvr.Resource,
					listHandlerFor(gvr),
				)
			}
		} else {
			mux.Handle(
				prefix+"/"+gvr.Resource+"/{name}",
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
