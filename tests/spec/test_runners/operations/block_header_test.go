package operations

import (
	"context"
	"github.com/zilm13/zrnt/eth2/beacon/altair"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/merge"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"testing"
)

type BlockHeaderTestCase struct {
	test_util.BaseTransitionTest
	Header *common.BeaconBlockHeader
}

func (c *BlockHeaderTestCase) Load(t *testing.T, forkName test_util.ForkName, readPart test_util.TestPartReader) {
	c.BaseTransitionTest.Load(t, forkName, readPart)
	switch forkName {
	case "phase0":
		var block phase0.BeaconBlock
		test_util.LoadSpecObj(t, "block", &block, readPart)
		c.Header = block.Header(c.Spec)
	case "altair":
		var block altair.BeaconBlock
		test_util.LoadSpecObj(t, "block", &block, readPart)
		c.Header = block.Header(c.Spec)
	case "merge":
		var block merge.BeaconBlock
		test_util.LoadSpecObj(t, "block", &block, readPart)
		c.Header = block.Header(c.Spec)
	default:
		t.Fatalf("unrecognized fork: %s", forkName)
	}
}

func (c *BlockHeaderTestCase) Run() error {
	epc, err := common.NewEpochsContext(c.Spec, c.Pre)
	if err != nil {
		return err
	}
	proposer, err := epc.GetBeaconProposer(c.Header.Slot)
	if err != nil {
		return err
	}
	return common.ProcessHeader(context.Background(), c.Spec, c.Pre, c.Header, proposer)
}

func TestBlockHeader(t *testing.T) {
	test_util.RunTransitionTest(t, test_util.AllForks, "operations", "block_header",
		func() test_util.TransitionTest { return new(BlockHeaderTestCase) })
}
