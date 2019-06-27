// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package plugins_test

import (
	"testing"

	"github.com/damienr74/kustomize/v3/internal/loadertest"
	"github.com/damienr74/kustomize/v3/k8sdeps/kunstruct"
	. "github.com/damienr74/kustomize/v3/pkg/plugins"
	"github.com/damienr74/kustomize/v3/pkg/resmap"
	"github.com/damienr74/kustomize/v3/pkg/resource"
)

const (
	secretGenerator = `
apiVersion: builtin
kind: SecretGenerator
metadata:
  name: secretGenerator
name: mySecret
behavior: merge
envFiles:
- a.env
- b.env
valueFiles:
- longsecret.txt
literals:
- FRUIT=apple
- VEGETABLE=carrot
`
	someServiceGenerator = `
apiVersion: someteam.example.com/v1
kind: SomeServiceGenerator
metadata:
  name: myServiceGenerator
service: my-service
port: "12345"
`
)

func TestLoader(t *testing.T) {
	tc := NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"builtin", "", "SecretGenerator")
	tc.BuildGoPlugin(
		"someteam.example.com", "v1", "SomeServiceGenerator")

	rmF := resmap.NewFactory(resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl()))

	l := NewLoader(ActivePluginConfig(), rmF)
	if l == nil {
		t.Fatal("expect non-nil loader")
	}

	ldr := loadertest.NewFakeLoader("/foo")

	m, err := rmF.NewResMapFromBytes([]byte(
		someServiceGenerator + "---\n" + secretGenerator))
	if err != nil {
		t.Fatal(err)
	}

	_, err = l.LoadGenerators(ldr, m)
	if err != nil {
		t.Fatal(err)
	}
}
