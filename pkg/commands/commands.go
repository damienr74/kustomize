// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package commands holds the CLI glue mapping textual commands/args to method calls.
package commands

import (
	"flag"
	"os"

	"github.com/damienr74/kustomize/v3/k8sdeps/kunstruct"
	"github.com/damienr74/kustomize/v3/k8sdeps/transformer"
	"github.com/damienr74/kustomize/v3/k8sdeps/validator"
	"github.com/damienr74/kustomize/v3/pkg/commands/build"
	"github.com/damienr74/kustomize/v3/pkg/commands/edit"
	"github.com/damienr74/kustomize/v3/pkg/commands/misc"
	"github.com/damienr74/kustomize/v3/pkg/fs"
	"github.com/damienr74/kustomize/v3/pkg/pgmconfig"
	"github.com/damienr74/kustomize/v3/pkg/resmap"
	"github.com/damienr74/kustomize/v3/pkg/resource"
	"github.com/spf13/cobra"
)

// NewDefaultCommand returns the default (aka root) command for kustomize command.
func NewDefaultCommand() *cobra.Command {
	fSys := fs.MakeRealFS()
	stdOut := os.Stdout

	c := &cobra.Command{
		Use:   pgmconfig.ProgramName,
		Short: "Manages declarative configuration of Kubernetes",
		Long: `
Manages declarative configuration of Kubernetes.
See https://github.com/damienr74/kustomize
`,
	}

	uf := kunstruct.NewKunstructuredFactoryImpl()
	rf := resmap.NewFactory(resource.NewFactory(uf))
	v := validator.NewKustValidator()
	c.AddCommand(
		build.NewCmdBuild(
			stdOut, fSys, v,
			rf, transformer.NewFactoryImpl()),
		edit.NewCmdEdit(fSys, v, uf),
		misc.NewCmdConfig(fSys),
		misc.NewCmdVersion(stdOut),
	)
	c.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
	return c
}
