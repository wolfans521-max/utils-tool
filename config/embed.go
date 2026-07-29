package config

import (
	_ "embed"
)

//go:embed log.yaml
var LogYaml []byte
