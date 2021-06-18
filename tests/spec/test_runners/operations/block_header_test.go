package operations

import (
	"context"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"testing"
)

type BlockHeaderTestCase struct {
	test_util.BaseTransitionTest
	Block phase0.BeaconBlock
}

func (c *BlockHeaderTestCase) Load(t *testing.T, readPart test_util.TestPartReader) {
	c.BaseTransitionTest.Load(t, readPart)
	test_util.LoadSpecObj(t, "block", &c.Block, readPart)
}

func (c *BlockHeaderTestCase) Run() error {
	epc, err := common.NewEpochsContext(c.Spec, c.Pre)
	if err != nil {
		return err
	}
	proposer, err := epc.GetBeaconProposer(c.Block.Slot)
	if err != nil {
		return err
	}
	return common.ProcessHeader(context.Background(), c.Spec, c.Pre, c.Block.Header(c.Spec), proposer)
}

func TestBlockHeader(t *testing.T) {
	test_util.RunTransitionTest(t, "operations", "block_header",
		func() test_util.TransitionTest { return new(BlockHeaderTestCase) })
}
