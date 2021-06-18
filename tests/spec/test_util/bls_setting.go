package test_util

import (
	"github.com/zilm13/zrnt/eth2/util/bls"
	"gopkg.in/yaml.v3"
	"testing"
)

const (
	BlsOptional = 0
	BlsRequired = 1
	BlsIgnored  = 2
)

type BLSMeta struct {
	BlsSetting int `yaml:"bls_setting"`
}

func HandleBLS(testRunner CaseRunner) CaseRunner {
	return func(t *testing.T, readPart TestPartReader) {
		part := readPart.Part("meta.yaml")
		if part.Exists() {
			meta := BLSMeta{}
			dec := yaml.NewDecoder(part)
			dec.KnownFields(false)
			Check(t, dec.Decode(&meta))
			Check(t, part.Close())
			if meta.BlsSetting == BlsRequired && !bls.BLS_ACTIVE {
				t.Skip("skipping BLS-required test because BLS is disabled")
				return
			}
			if meta.BlsSetting == BlsIgnored && bls.BLS_ACTIVE {
				t.Skip("skipping BLS-ignored test because BLS is enabled")
				return
			}
		}
		testRunner(t, readPart)
	}
}
