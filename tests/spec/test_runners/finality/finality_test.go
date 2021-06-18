package finality

import (
	"context"
	"fmt"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"gopkg.in/yaml.v3"
	"testing"
)

type FinalityTestCase struct {
	test_util.BaseTransitionTest
	Blocks []*phase0.SignedBeaconBlock
}

type BlocksCountMeta struct {
	BlocksCount uint64 `yaml:"blocks_count"`
}

func (c *FinalityTestCase) Load(t *testing.T, readPart test_util.TestPartReader) {
	c.BaseTransitionTest.Load(t, readPart)
	p := readPart.Part("meta.yaml")
	dec := yaml.NewDecoder(p)
	m := &BlocksCountMeta{}
	test_util.Check(t, dec.Decode(&m))
	test_util.Check(t, p.Close())
	loadBlock := func(i uint64) *phase0.SignedBeaconBlock {
		dst := new(phase0.SignedBeaconBlock)
		test_util.LoadSpecObj(t, fmt.Sprintf("blocks_%d", i), dst, readPart)
		return dst
	}
	for i := uint64(0); i < m.BlocksCount; i++ {
		c.Blocks = append(c.Blocks, loadBlock(i))
	}
}

func (c *FinalityTestCase) Run() error {
	epc, err := common.NewEpochsContext(c.Spec, c.Pre)
	if err != nil {
		return err
	}
	state := c.Pre
	valRoot, err := state.GenesisValidatorsRoot()
	if err != nil {
		return err
	}
	digest := common.ComputeForkDigest(c.Spec.GENESIS_FORK_VERSION, valRoot)
	for _, b := range c.Blocks {
		benv := b.Envelope(c.Spec, digest)
		if err := common.StateTransition(context.Background(), c.Spec, epc, state, benv, true); err != nil {
			return err
		}
	}
	return nil
}

func TestBlocks(t *testing.T) {
	test_util.RunTransitionTest(t, "finality", "finality",
		func() test_util.TransitionTest { return new(FinalityTestCase) })
}
