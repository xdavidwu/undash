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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	upstream = "http://localhost:9090"

	apidiscoveryv2MimeType = runtime.ContentTypeJSON + ";g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList"
)

type kindMeta struct {
	kind string
	singular string
	listKind string
}

var (
	coreNamespacedResources = map[string]kindMeta{
		"services": {
			kind: "Service",
			singular: "service",
			listKind: "ServiceList",
		},
	}
)

func coreResourceNamespacedListHandlerFor(resource string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ns := r.PathValue("namespace")
		kind := coreNamespacedResources[resource]

		listRes, err := http.Get(fmt.Sprintf("%s/api/v1/%s/%s", upstream, kind.singular, ns))
		if err != nil {
			panic(err)
		}
		defer listRes.Body.Close()
		body, _ := io.ReadAll(listRes.Body)
		listObj := map[string]any{}
		json.Unmarshal(body, &listObj)
		fmt.Printf("%+v\n", listObj)

		nativeList := metav1.List{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind: kind.listKind,
			},
			Items: []runtime.RawExtension{},
		}
		for _, obj := range listObj[resource].([]any) {
			name := obj.(map[string]any)["objectMeta"].(map[string]any)["name"].(string)
			fmt.Printf("found %s\n", name)
			realObjRes, _ := http.Get(fmt.Sprintf("%s/api/v1/_raw/%s/namespace/%s/name/%s", upstream, kind.singular, ns, name))
			defer realObjRes.Body.Close()
			realObj, _ := io.ReadAll(realObjRes.Body)
			nativeList.Items = append(nativeList.Items, runtime.RawExtension{Raw: realObj})
		}
		w.Header().Set("Content-Type", runtime.ContentTypeJSON)

		encoder := json.NewEncoder(w)
		encoder.Encode(nativeList)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/api",
		func(w http.ResponseWriter, r *http.Request) {
			// XXX supporting apidiscoveryv2 (1.30+) since it's easier
			// but dashboard targets 1.25
			discoveries := []apidiscoveryv2.APIResourceDiscovery{}
			for resource, kind := range coreNamespacedResources {
				discoveries = append(discoveries, apidiscoveryv2.APIResourceDiscovery{
					Resource: resource,
					ResponseKind: &metav1.GroupVersionKind{
						Kind: kind.kind,
					},
					Scope: apidiscoveryv2.ScopeNamespace,
					SingularResource: kind.singular,
					Verbs: []string{"list", "get", "delete", "update"},
				})
			}
			list := apidiscoveryv2.APIGroupDiscoveryList{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apidiscovery.k8s.io/v2",
					Kind: "APIGroupDiscoveryList",
				},
				Items: []apidiscoveryv2.APIGroupDiscovery{
					{
						Versions: []apidiscoveryv2.APIVersionDiscovery{
							{
								Version: "v1",
								Resources: discoveries,
								Freshness: apidiscoveryv2.DiscoveryFreshnessCurrent,
							},
						},
					},
				},
			}
			w.Header().Set("Content-Type", apidiscoveryv2MimeType)
			encoder := json.NewEncoder(w)
			encoder.Encode(list)
		},
	)
	mux.HandleFunc(
		"/apis",
		func(w http.ResponseWriter, r *http.Request) {
			// XXX supporting apidiscoveryv2 (1.30+) since it's easier
			// but dashboard targets 1.25
			list := apidiscoveryv2.APIGroupDiscoveryList{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apidiscovery.k8s.io/v2",
					Kind: "APIGroupDiscoveryList",
				},
				Items: []apidiscoveryv2.APIGroupDiscovery{},
			}
			w.Header().Set("Content-Type", apidiscoveryv2MimeType)
			encoder := json.NewEncoder(w)
			encoder.Encode(list)
		},
	)

	for resource, kind := range coreNamespacedResources {
		mux.Handle(
			"/api/v1/namespaces/{namespace}/" + resource,
			coreResourceNamespacedListHandlerFor(resource),
		)

		mux.Handle(
			"/api/v1/namespaces/{namespace}/" + resource + "/{name}",
			&httputil.ReverseProxy{
				Rewrite: func(r *httputil.ProxyRequest) {
					url, _ := url.Parse(fmt.Sprintf(
						"%s/api/v1/_raw/%s/namespace/%s/name/%s",
						upstream,
						kind.singular,
						r.In.PathValue("namespace"),
						r.In.PathValue("name"),
					))
					fmt.Printf("%+v\n", url)
					r.Out.URL = url
				},
				ModifyResponse: func(resp *http.Response) error {
					fmt.Printf("%+v\n", resp.Request)
					return nil
				},
			},
		)
	}

	l, _ := net.Listen("tcp", "localhost:9091")
	http.Serve(l, mux)
}
