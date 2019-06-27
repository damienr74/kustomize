// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	"github.com/damienr74/kustomize/v3/pkg/kusttest"
	"github.com/damienr74/kustomize/v3/pkg/plugins"
)

func TestLabelTransformer(t *testing.T) {
	tc := plugins.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"builtin", "", "LabelTransformer")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	rm := th.LoadAndRunTransformer(`
apiVersion: builtin
kind: LabelTransformer
metadata:
  name: notImportantHere
labels:
  app: myApp
  env: production
fieldSpecs:
  - path: metadata/labels
    create: true
`, `
apiVersion: v1
kind: Service
metadata:
  name: myService
spec:
  ports:
  - port: 7002
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: Service
metadata:
  labels:
    app: myApp
    env: production
  name: myService
spec:
  ports:
  - port: 7002
`)
}
