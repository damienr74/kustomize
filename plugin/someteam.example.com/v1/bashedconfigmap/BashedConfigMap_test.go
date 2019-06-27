// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	"github.com/damienr74/kustomize/v3/pkg/kusttest"
	"github.com/damienr74/kustomize/v3/pkg/plugins"
)

func TestBashedConfigMapPlugin(t *testing.T) {
	tc := plugins.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildExecPlugin(
		"someteam.example.com", "v1", "BashedConfigMap")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	m := th.LoadAndRunGenerator(`
apiVersion: someteam.example.com/v1
kind: BashedConfigMap
metadata:
  name: whatever
argsOneLiner: alice myMomsMaidenName
`)
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  password: myMomsMaidenName
  username: alice
kind: ConfigMap
metadata:
  name: example-configmap-test
`)
}
