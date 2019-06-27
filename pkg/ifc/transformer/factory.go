/// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package patch holds miscellaneous interfaces used by kustomize.
package transformer

import (
	"github.com/damienr74/kustomize/v3/pkg/resource"
	"github.com/damienr74/kustomize/v3/pkg/transformers"
)

// Factory makes transformers that require k8sdeps.
type Factory interface {
	MakePatchTransformer(
		slice []*resource.Resource,
		rf *resource.Factory) (transformers.Transformer, error)
}
