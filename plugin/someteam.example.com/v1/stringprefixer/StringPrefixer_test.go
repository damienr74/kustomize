// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	"github.com/damienr74/kustomize/v3/pkg/kusttest"
	"github.com/damienr74/kustomize/v3/pkg/plugins"
)

func TestStringPrefixerPlugin(t *testing.T) {
	tc := plugins.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"someteam.example.com", "v1", "StringPrefixer")
	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	m := th.LoadAndRunTransformer(`
apiVersion: someteam.example.com/v1
kind: StringPrefixer
metadata:
  name: wowsa
`,
		`apiVersion: apps/v1
kind: MeatBall
metadata:
  name: meatball
`)

	th.AssertActualEqualsExpected(m, `
apiVersion: apps/v1
kind: MeatBall
metadata:
  name: wowsa-meatball
`)
}
