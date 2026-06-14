package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

type APIStructType int

const (
	APIStructUnknown APIStructType = iota
	APIStructObject
	APIStructList
)

func GetAPIStructType(st *types.Struct) APIStructType {
	hasObjectMeta := false
	hasListMeta := false
FieldsLoop:
	for field := range st.Fields() {
		if !field.Exported() || !field.Embedded() {
			continue
		}
		switch field.Name() {
		case "ObjectMeta":
			hasObjectMeta = true
			break FieldsLoop
		case "ListMeta":
			hasListMeta = true
			break FieldsLoop
		}
	}
	if hasObjectMeta {
		return APIStructObject
	}
	if hasListMeta {
		return APIStructList
	}
	return APIStructUnknown
}

func hasDeepCopyObject(named *types.Named) bool {
	for method := range named.Methods() {
		if method.Name() == "DeepCopyObject" {
			return true
		}
	}

	return false
}

type conversionFn struct {
	fn        *types.Func
	fromPtrTo *types.Named
}

type info struct {
	pkg               string
	objectConversions []conversionFn
	listConversions   []conversionFn
}

func segment(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func pkgToIdentifier(path string) string {
	comps := strings.Split(path, "/")
	groupver := comps[len(comps)-2:]
	return segment(groupver[0]) + segment(groupver[1])
}

func tableConversionFnName(apiType *types.Named) string {
	tn := apiType.Obj()
	api := pkgToIdentifier(tn.Pkg().Path())
	return api + tn.Name() + "ToTable"
}

func localImportName(path string) string {
	comps := strings.Split(path, "/")
	groupver := comps[len(comps)-2:]
	return groupver[0] + groupver[1]
}

// find pkg/apis/ -type d -mindepth 2 -maxdepth 2 -name 'v[0-9]*' | sort | sed 's|^|"k8s.io/kubernetes/|;s|$|",|'
var pkgs = []string{
	"k8s.io/kubernetes/pkg/apis/abac/v0",
	"k8s.io/kubernetes/pkg/apis/abac/v1beta1",
	"k8s.io/kubernetes/pkg/apis/admission/v1",
	"k8s.io/kubernetes/pkg/apis/admission/v1beta1",
	"k8s.io/kubernetes/pkg/apis/admissionregistration/v1",
	"k8s.io/kubernetes/pkg/apis/admissionregistration/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/admissionregistration/v1beta1",
	"k8s.io/kubernetes/pkg/apis/apidiscovery/v2",
	"k8s.io/kubernetes/pkg/apis/apidiscovery/v2beta1",
	"k8s.io/kubernetes/pkg/apis/apiserverinternal/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/apps/v1",
	"k8s.io/kubernetes/pkg/apis/apps/v1beta1",
	"k8s.io/kubernetes/pkg/apis/apps/v1beta2",
	"k8s.io/kubernetes/pkg/apis/authentication/v1",
	"k8s.io/kubernetes/pkg/apis/authentication/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/authentication/v1beta1",
	"k8s.io/kubernetes/pkg/apis/authorization/v1",
	"k8s.io/kubernetes/pkg/apis/authorization/v1beta1",
	"k8s.io/kubernetes/pkg/apis/autoscaling/v1",
	"k8s.io/kubernetes/pkg/apis/autoscaling/v2",
	"k8s.io/kubernetes/pkg/apis/autoscaling/v2beta1",
	"k8s.io/kubernetes/pkg/apis/autoscaling/v2beta2",
	"k8s.io/kubernetes/pkg/apis/batch/v1",
	"k8s.io/kubernetes/pkg/apis/batch/v1beta1",
	"k8s.io/kubernetes/pkg/apis/certificates/v1",
	"k8s.io/kubernetes/pkg/apis/certificates/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/certificates/v1beta1",
	"k8s.io/kubernetes/pkg/apis/coordination/v1",
	"k8s.io/kubernetes/pkg/apis/coordination/v1beta1",
	"k8s.io/kubernetes/pkg/apis/core/v1",
	"k8s.io/kubernetes/pkg/apis/discovery/v1",
	"k8s.io/kubernetes/pkg/apis/discovery/v1beta1",
	"k8s.io/kubernetes/pkg/apis/events/v1",
	"k8s.io/kubernetes/pkg/apis/events/v1beta1",
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1",
	"k8s.io/kubernetes/pkg/apis/flowcontrol/v1",
	"k8s.io/kubernetes/pkg/apis/flowcontrol/v1beta1",
	"k8s.io/kubernetes/pkg/apis/flowcontrol/v1beta2",
	"k8s.io/kubernetes/pkg/apis/flowcontrol/v1beta3",
	"k8s.io/kubernetes/pkg/apis/imagepolicy/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/networking/v1",
	"k8s.io/kubernetes/pkg/apis/networking/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/networking/v1beta1",
	"k8s.io/kubernetes/pkg/apis/node/v1",
	"k8s.io/kubernetes/pkg/apis/node/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/node/v1beta1",
	"k8s.io/kubernetes/pkg/apis/policy/v1",
	"k8s.io/kubernetes/pkg/apis/policy/v1beta1",
	"k8s.io/kubernetes/pkg/apis/rbac/v1",
	"k8s.io/kubernetes/pkg/apis/rbac/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/rbac/v1beta1",
	"k8s.io/kubernetes/pkg/apis/resource/v1alpha2",
	"k8s.io/kubernetes/pkg/apis/scheduling/v1",
	"k8s.io/kubernetes/pkg/apis/scheduling/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/scheduling/v1beta1",
	"k8s.io/kubernetes/pkg/apis/storage/v1",
	"k8s.io/kubernetes/pkg/apis/storage/v1alpha1",
	"k8s.io/kubernetes/pkg/apis/storage/v1beta1",
	"k8s.io/kubernetes/pkg/apis/storagemigration/v1alpha1",
}

var ignores = []string{
	// prefer Convert_v1_ReplicationController_To_apps_ReplicaSet
	"Convert_v1_ReplicationController_To_core_ReplicationController",
}

func writeAST(path string, file *ast.File) error {
	f := &bytes.Buffer{}
	if err := format.Node(f, token.NewFileSet(), file); err != nil {
		return err
	}

	// XXX hack for map literal in generateUnstructuredTable
	s := f.String()
	s = strings.ReplaceAll(s, "{", "{\n")
	s = strings.ReplaceAll(s, ")}", "),\n}")
	s = strings.ReplaceAll(s, "),", "),\n")

	// 2nd pass for newlines
	b, err := format.Source([]byte(s))
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, b, 0o644); err != nil {
		return err
	}

	return nil
}

func extractInfos(pkgs []*packages.Package) ([]info, error) {
	infos := []info{}
	for _, pkg := range pkgs {
		info := info{
			pkg: pkg.PkgPath,
		}

		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if !obj.Exported() {
				continue
			}
			if fn, ok := obj.(*types.Func); ok {
				if strings.HasPrefix(fn.Name(), "Convert_"+pkg.Name+"_") && !slices.Contains(ignores, fn.Name()) {
					params := fn.Signature().Params()
					if params.Len() > 0 {
						inType := params.At(0).Type()

						if ptr, ok := inType.(*types.Pointer); ok {
							if named, ok := ptr.Elem().(*types.Named); ok {
								// XXX check if implements runtime.Object?
								if !hasDeepCopyObject(named) {
									continue
								}

								if st, ok := named.Underlying().(*types.Struct); ok {
									switch GetAPIStructType(st) {
									case APIStructObject:
										info.objectConversions = append(info.objectConversions, conversionFn{
											fn:        fn,
											fromPtrTo: named,
										})
									case APIStructList:
										info.listConversions = append(info.listConversions, conversionFn{
											fn:        fn,
											fromPtrTo: named,
										})
									case APIStructUnknown:
									}
								}
							}
						}
					}
				}
			}
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func generateTables(infos []info) error {
	importSpecs := []ast.Spec{}
	varDecls := []ast.Decl{}
	for _, info := range infos {
		if len(info.objectConversions) == 0 && len(info.listConversions) == 0 {
			continue
		}

		localImport := localImportName(info.pkg)
		importSpecs = append(importSpecs, &ast.ImportSpec{
			Name: ast.NewIdent(localImport),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(info.pkg),
			},
		})

		for _, convs := range []struct {
			typ         string
			conversions []conversionFn
		}{
			{
				typ:         "objects",
				conversions: info.objectConversions,
			},
			{
				typ:         "lists",
				conversions: info.listConversions,
			},
		} {
			if len(convs.conversions) == 0 {
				continue
			}

			varSpecs := []ast.Spec{}
			for _, conv := range convs.conversions {
				varSpecs = append(varSpecs, &ast.ValueSpec{
					Names: []*ast.Ident{ast.NewIdent(tableConversionFnName(conv.fromPtrTo))},
					Values: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("ToTableFunc"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent(localImport),
									Sel: ast.NewIdent(conv.fn.Name()),
								},
							},
						},
					},
				})
			}

			varDecls = append(varDecls, &ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: fmt.Sprintf("// %s %s to metav1.Table functions", localImport, convs.typ),
						},
					},
				},
				Tok:   token.VAR,
				Specs: varSpecs,
			})
		}

	}

	file := &ast.File{
		Name: ast.NewIdent("kubernetes"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: importSpecs,
			},
		},
	}
	file.Decls = append(file.Decls, varDecls...)

	return writeAST("zz_generated.table.go", file)
}

//var ListToTableFuncs = map[schema.GroupVersionKind]func(*unstructured.UnstructuredList) (*metav1.Table, error){
//	corev1.SchemeGroupVersion.WithKind("ServiceList"): UnstructuredListToTableFunc(CoreV1ServiceListToTable),
//}

func generateUnstructuredTable(infos []info) error {
	importSpecs := []ast.Spec{
		&ast.ImportSpec{
			Name: ast.NewIdent("metav1"),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("k8s.io/apimachinery/pkg/apis/meta/v1"),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("k8s.io/apimachinery/pkg/runtime/schema"),
			},
		},
	}
	listMapKVE := []ast.Expr{}
	objectMapKVE := []ast.Expr{}
	for _, info := range infos {
		if len(info.listConversions) == 0 && len(info.objectConversions) == 0 {
			continue
		}

		localImport := localImportName(info.pkg)
		importSpecs = append(importSpecs, &ast.ImportSpec{
			Name: ast.NewIdent(localImport),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(info.pkg),
			},
		})

		for _, conv := range info.objectConversions {
			objectMapKVE = append(objectMapKVE, &ast.KeyValueExpr{
				Key: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent(localImport),
							Sel: ast.NewIdent("SchemeGroupVersion"),
						},
						Sel: ast.NewIdent("WithKind"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote(conv.fromPtrTo.Obj().Name()),
						},
					},
				},
				Value: &ast.CallExpr{
					Fun: ast.NewIdent("UnstructuredToTableFunc"),
					Args: []ast.Expr{
						ast.NewIdent(tableConversionFnName(conv.fromPtrTo)),
					},
				},
			})
		}

		for _, conv := range info.listConversions {
			listMapKVE = append(listMapKVE, &ast.KeyValueExpr{
				Key: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent(localImport),
							Sel: ast.NewIdent("SchemeGroupVersion"),
						},
						Sel: ast.NewIdent("WithKind"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote(conv.fromPtrTo.Obj().Name()),
						},
					},
				},
				Value: &ast.CallExpr{
					Fun: ast.NewIdent("UnstructuredListToTableFunc"),
					Args: []ast.Expr{
						ast.NewIdent(tableConversionFnName(conv.fromPtrTo)),
					},
				},
			})
		}

	}

	file := &ast.File{
		Name: ast.NewIdent("kubernetes"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: importSpecs,
			},
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("UnstructuredListToTableFuncs"),
						},
						Values: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.MapType{
									Key: &ast.SelectorExpr{
										X:   ast.NewIdent("schema"),
										Sel: ast.NewIdent("GroupVersionKind"),
									},
									Value: &ast.FuncType{
										Params: &ast.FieldList{
											List: []*ast.Field{
												{
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("unstructured"),
															Sel: ast.NewIdent("UnstructuredList"),
														},
													},
												},
											},
										},
										Results: &ast.FieldList{
											List: []*ast.Field{
												{
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("metav1"),
															Sel: ast.NewIdent("Table"),
														},
													},
												},
												{
													Type: ast.NewIdent("error"),
												},
											},
										},
									},
								},
								Elts: listMapKVE,
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("UnstructuredToTableFuncs"),
						},
						Values: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.MapType{
									Key: &ast.SelectorExpr{
										X:   ast.NewIdent("schema"),
										Sel: ast.NewIdent("GroupVersionKind"),
									},
									Value: &ast.FuncType{
										Params: &ast.FieldList{
											List: []*ast.Field{
												{
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("unstructured"),
															Sel: ast.NewIdent("Unstructured"),
														},
													},
												},
											},
										},
										Results: &ast.FieldList{
											List: []*ast.Field{
												{
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("metav1"),
															Sel: ast.NewIdent("Table"),
														},
													},
												},
												{
													Type: ast.NewIdent("error"),
												},
											},
										},
									},
								},
								Elts: objectMapKVE,
							},
						},
					},
				},
			},
		},
	}

	return writeAST("zz_generated.unstructured.go", file)
}

func main() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax,
	}, pkgs...)
	if err != nil {
		panic(err)
	}

	infos, err := extractInfos(pkgs)
	if err != nil {
		panic(err)
	}

	if err := generateTables(infos); err != nil {
		panic(fmt.Errorf("cannot generate tables: %w", err))
	}

	if err := generateUnstructuredTable(infos); err != nil {
		panic(fmt.Errorf("cannot generate tables: %w", err))
	}
}
