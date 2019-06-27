// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package patch

import (
	"reflect"
	"strings"
	"testing"

	"github.com/damienr74/kustomize/v3/k8sdeps/kunstruct"
	"github.com/damienr74/kustomize/v3/pkg/resmaptest"
	"github.com/damienr74/kustomize/v3/pkg/resource"
)

var rf = resource.NewFactory(
	kunstruct.NewKunstructuredFactoryImpl())

func TestOverlayRun(t *testing.T) {
	base := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"old-label": "old-value",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx",
							},
						},
					},
				},
			},
		}).ResMap()
	patch := []*resource.Resource{
		rf.FromMap(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"another-label": "foo",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:latest",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "SOMEENV",
										"value": "BAR",
									},
								},
							},
						},
					},
				},
			},
		}),
	}
	expected := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"old-label":     "old-value",
							"another-label": "foo",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:latest",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "SOMEENV",
										"value": "BAR",
									},
								},
							},
						},
					},
				},
			},
		}).ResMap()
	lt, err := NewTransformer(patch, rf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(base, expected) {
		err = expected.ErrorIfNotEqualLists(base)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

func TestMultiplePatches(t *testing.T) {
	base := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx",
							},
						},
					},
				},
			},
		}).ResMap()
	patch := []*resource.Resource{
		rf.FromMap(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:latest",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "SOMEENV",
										"value": "BAR",
									},
								},
							},
						},
					},
				},
			},
		}),
		rf.FromMap(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name": "nginx",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "ANOTHERENV",
										"value": "HELLO",
									},
								},
							},
							map[string]interface{}{
								"name":  "busybox",
								"image": "busybox",
							},
						},
					},
				},
			},
		}),
	}
	expected := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:latest",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "ANOTHERENV",
										"value": "HELLO",
									},
									map[string]interface{}{
										"name":  "SOMEENV",
										"value": "BAR",
									},
								},
							},
							map[string]interface{}{
								"name":  "busybox",
								"image": "busybox",
							},
						},
					},
				},
			},
		}).ResMap()
	lt, err := NewTransformer(patch, rf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(base, expected) {
		err = expected.ErrorIfNotEqualLists(base)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

func TestMultiplePatchesWithConflict(t *testing.T) {
	base := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx",
							},
						},
					},
				},
			},
		}).ResMap()

	patch := []*resource.Resource{
		rf.FromMap(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:latest",
								"env": []interface{}{
									map[string]interface{}{
										"name":  "SOMEENV",
										"value": "BAR",
									},
								},
							},
						},
					},
				},
			},
		}),
		rf.FromMap(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:1.7.9",
							},
						},
					},
				},
			},
		}),
	}

	lt, err := NewTransformer(patch, rf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err == nil {
		t.Fatalf("did not get expected error")
	}
	if !strings.Contains(err.Error(), "conflict") {
		t.Fatalf("expected error to contain %q but get %v", "conflict", err)
	}
}

func TestPatchesWithWrongNamespace(t *testing.T) {
	base := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "deploy1",
				"namespace": "namespace1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx",
							},
						},
					},
				},
			},
		}).ResMap()

	patch := []*resource.Resource{
		rf.FromMap(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "deploy1",
				"namespace": "namespace2",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:1.7.9",
							},
						},
					},
				},
			},
		}),
	}

	lt, err := NewTransformer(patch, rf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err == nil {
		t.Fatalf("did not get expected error")
	}
	if !strings.Contains(err.Error(), "failed to find target for patch") {
		t.Fatalf("expected error to contain %q but get %v", "failed to find target for patch", err)
	}
}

func TestNoSchemaOverlayRun(t *testing.T) {
	base := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"A": "X",
					"B": "Y",
				},
			},
		}).ResMap()
	patch := []*resource.Resource{
		rf.FromMap(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"B": nil,
					"C": "Z",
				},
			},
		}),
	}
	expected := resmaptest_test.NewRmBuilder(t, rf).
		Add(
			map[string]interface{}{
				"apiVersion": "example.com/v1",
				"kind":       "Foo",
				"metadata": map[string]interface{}{
					"name": "my-foo",
				},
				"spec": map[string]interface{}{
					"bar": map[string]interface{}{
						"A": "X",
						"C": "Z",
					},
				},
			}).ResMap()

	lt, err := NewTransformer(patch, rf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err = expected.ErrorIfNotEqualLists(base); err != nil {
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

func TestNoSchemaMultiplePatches(t *testing.T) {
	base := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"A": "X",
					"B": "Y",
				},
			},
		}).ResMap()
	patch := []*resource.Resource{
		rf.FromMap(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"B": nil,
					"C": "Z",
				},
			},
		}),
		rf.FromMap(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"C": "Z",
					"D": "W",
				},
				"baz": map[string]interface{}{
					"hello": "world",
				},
			},
		}),
	}
	expected := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"A": "X",
					"C": "Z",
					"D": "W",
				},
				"baz": map[string]interface{}{
					"hello": "world",
				},
			},
		}).ResMap()

	lt, err := NewTransformer(patch, rf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err = expected.ErrorIfNotEqualLists(base); err != nil {
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

func TestNoSchemaMultiplePatchesWithConflict(t *testing.T) {
	base := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"A": "X",
					"B": "Y",
				},
			},
		}).ResMap()
	patch := []*resource.Resource{
		rf.FromMap(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"B": nil,
					"C": "Z",
				},
			},
		}),
		rf.FromMap(map[string]interface{}{
			"apiVersion": "example.com/v1",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "my-foo",
			},
			"spec": map[string]interface{}{
				"bar": map[string]interface{}{
					"C": "NOT_Z",
				},
			},
		}),
	}

	lt, err := NewTransformer(patch, rf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(base)
	if err == nil {
		t.Fatalf("did not get expected error")
	}
	if !strings.Contains(err.Error(), "conflict") {
		t.Fatalf("expected error to contain %q but get %v", "conflict", err)
	}
}
