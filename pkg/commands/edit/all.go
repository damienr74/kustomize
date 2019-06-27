// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package edit

import (
	"github.com/damienr74/kustomize/v3/pkg/commands/edit/add"
	"github.com/damienr74/kustomize/v3/pkg/commands/edit/fix"
	"github.com/damienr74/kustomize/v3/pkg/commands/edit/remove"
	"github.com/damienr74/kustomize/v3/pkg/commands/edit/set"
	"github.com/damienr74/kustomize/v3/pkg/fs"
	"github.com/damienr74/kustomize/v3/pkg/ifc"
	"github.com/damienr74/kustomize/v3/pkg/loader"
	"github.com/spf13/cobra"
)

// NewCmdEdit returns an instance of 'edit' subcommand.
func NewCmdEdit(
	fSys fs.FileSystem, v ifc.Validator, kf ifc.KunstructuredFactory) *cobra.Command {
	c := &cobra.Command{
		Use:   "edit",
		Short: "Edits a kustomization file",
		Long:  "",
		Example: `
	# Adds a configmap to the kustomization file
	kustomize edit add configmap NAME --from-literal=k=v

	# Sets the nameprefix field
	kustomize edit set nameprefix <prefix-value>

	# Sets the namesuffix field
	kustomize edit set namesuffix <suffix-value>
`,
		Args: cobra.MinimumNArgs(1),
	}

	c.AddCommand(
		add.NewCmdAdd(fSys, loader.NewFileLoaderAtCwd(v, fSys), kf),
		set.NewCmdSet(fSys, v),
		fix.NewCmdFix(fSys),
		remove.NewCmdRemove(fSys, loader.NewFileLoaderAtCwd(v, fSys)),
	)
	return c
}
