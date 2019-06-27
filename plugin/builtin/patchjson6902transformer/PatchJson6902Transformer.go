// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

//go:generate go run github.com/damienr74/kustomize/v3/cmd/pluginator
package main

import (
	"github.com/damienr74/kustomize/v3/pkg/ifc"
	"github.com/damienr74/kustomize/v3/pkg/patch/transformer"
	"github.com/damienr74/kustomize/v3/pkg/resmap"
	"github.com/damienr74/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	ldr     ifc.Loader
	Patches []types.PatchJson6902 `json:"patches,omitempty" yaml:"patches,omitempty"`
}

//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) (err error) {
	p.ldr = ldr
	p.Patches = nil
	return yaml.Unmarshal(c, p)
}

func (p *plugin) Transform(m resmap.ResMap) error {
	t, err := transformer.NewPatchJson6902Factory(p.ldr).
		MakePatchJson6902Transformer(p.Patches)
	if err != nil {
		return err
	}
	return t.Transform(m)
}
