// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

//go:generate go run github.com/damienr74/kustomize/v3/cmd/pluginator
package main

import (
	"github.com/damienr74/kustomize/v3/pkg/ifc"
	"github.com/damienr74/kustomize/v3/pkg/resmap"
	"github.com/damienr74/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	ldr ifc.Loader
	rf  *resmap.Factory
	types.GeneratorOptions
	types.ConfigMapArgs
}

//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, config []byte) (err error) {
	p.GeneratorOptions = types.GeneratorOptions{}
	p.ConfigMapArgs = types.ConfigMapArgs{}
	err = yaml.Unmarshal(config, p)
	p.ldr = ldr
	p.rf = rf
	return
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	return p.rf.FromConfigMapArgs(p.ldr, &p.GeneratorOptions, p.ConfigMapArgs)
}
