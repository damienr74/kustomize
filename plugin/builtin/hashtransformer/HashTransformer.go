// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

//go:generate go run github.com/damienr74/kustomize/v3/cmd/pluginator
package main

import (
	"fmt"

	"github.com/damienr74/kustomize/v3/pkg/ifc"
	"github.com/damienr74/kustomize/v3/pkg/resmap"
)

type plugin struct {
	hasher ifc.KunstructuredHasher
}

//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, config []byte) (err error) {
	p.hasher = rf.RF().Hasher()
	return nil
}

// Transform appends hash to generated resources.
func (p *plugin) Transform(m resmap.ResMap) error {
	for _, res := range m.Resources() {
		if res.NeedHashSuffix() {
			h, err := p.hasher.Hash(res)
			if err != nil {
				return err
			}
			res.SetName(fmt.Sprintf("%s-%s", res.GetName(), h))
		}
	}
	return nil
}
