// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package transformer provides transformer factory
package transformer

import (
	"github.com/damienr74/kustomize/v3/k8sdeps/transformer/patch"
	"github.com/damienr74/kustomize/v3/pkg/resource"
	"github.com/damienr74/kustomize/v3/pkg/transformers"
)

// FactoryImpl makes patch transformer and name hash transformer
type FactoryImpl struct{}

// NewFactoryImpl makes a new factoryImpl instance
func NewFactoryImpl() *FactoryImpl {
	return &FactoryImpl{}
}

// MakePatchTransformer makes a new patch transformer
func (p *FactoryImpl) MakePatchTransformer(
	slice []*resource.Resource,
	rf *resource.Factory) (transformers.Transformer, error) {
	return patch.NewTransformer(slice, rf)
}
