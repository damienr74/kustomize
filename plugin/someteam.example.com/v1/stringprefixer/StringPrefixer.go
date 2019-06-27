// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/damienr74/kustomize/v3/pkg/ifc"
	"github.com/damienr74/kustomize/v3/pkg/resmap"
	"github.com/damienr74/kustomize/v3/pkg/transformers"
	"github.com/damienr74/kustomize/v3/pkg/transformers/config"
	"github.com/damienr74/kustomize/v3/plugin/builtin"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// Add a string prefix to the name.
// A plugin that adapts another plugin.
type plugin struct {
	Metadata metaData `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	t        transformers.Transformer
}

type metaData struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

//nolint: golint
//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) makePrefixSuffixPluginConfig(n string) ([]byte, error) {
	var s struct {
		Prefix     string
		Suffix     string
		FieldSpecs []config.FieldSpec
	}
	s.Prefix = n + "-"
	s.FieldSpecs = []config.FieldSpec{
		{Path: "metadata/name"},
	}
	return yaml.Marshal(s)
}

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	err := yaml.Unmarshal(c, p)
	if err != nil {
		return err
	}
	c, err = p.makePrefixSuffixPluginConfig(p.Metadata.Name)
	if err != nil {
		return err
	}
	prefixer := builtin.NewPrefixSuffixTransformerPlugin()
	err = prefixer.Config(ldr, rf, c)
	if err != nil {
		return errors.Wrapf(
			err, "stringprefixer configure")
	}
	p.t = prefixer
	return nil
}

func (p *plugin) Transform(m resmap.ResMap) error {
	return p.t.Transform(m)
}
