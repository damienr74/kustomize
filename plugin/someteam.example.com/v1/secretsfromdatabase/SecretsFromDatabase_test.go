// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	"github.com/damienr74/kustomize/v3/pkg/kusttest"
	"github.com/damienr74/kustomize/v3/pkg/plugins"
)

func TestSecretsFromDatabasePlugin(t *testing.T) {
	tc := plugins.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"someteam.example.com", "v1", "SecretsFromDatabase")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	m := th.LoadAndRunGenerator(`
apiVersion: someteam.example.com/v1
kind: SecretsFromDatabase
metadata:
  name: mySecretGenerator
name: forbiddenValues
namespace: production
keys:
- ROCKET
- VEGETABLE
`)
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  ROCKET: U2F0dXJuVg==
  VEGETABLE: Y2Fycm90
kind: Secret
metadata:
  name: forbiddenValues
  namespace: production
type: Opaque
`)
}
