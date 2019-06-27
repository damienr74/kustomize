// Code generated by pluginator on SecretGenerator; DO NOT EDIT.
package builtin

import (
	"github.com/damienr74/kustomize/v3/pkg/ifc"
	"github.com/damienr74/kustomize/v3/pkg/resmap"
	"github.com/damienr74/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

type SecretGeneratorPlugin struct {
	ldr ifc.Loader
	rf  *resmap.Factory
	types.GeneratorOptions
	types.SecretArgs
}

//noinspection GoUnusedGlobalVariable
func NewSecretGeneratorPlugin() *SecretGeneratorPlugin {
	return &SecretGeneratorPlugin{}
}

func (p *SecretGeneratorPlugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, config []byte) (err error) {
	p.GeneratorOptions = types.GeneratorOptions{}
	p.SecretArgs = types.SecretArgs{}
	err = yaml.Unmarshal(config, p)
	p.ldr = ldr
	p.rf = rf
	return
}

func (p *SecretGeneratorPlugin) Generate() (resmap.ResMap, error) {
	return p.rf.FromSecretArgs(p.ldr, &p.GeneratorOptions, p.SecretArgs)
}
