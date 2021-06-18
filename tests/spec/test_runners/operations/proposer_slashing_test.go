package operations

import (
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"testing"
)

type ProposerSlashingTestCase struct {
	test_util.BaseTransitionTest
	ProposerSlashing phase0.ProposerSlashing
}

func (c *ProposerSlashingTestCase) Load(t *testing.T, forkName test_util.ForkName, readPart test_util.TestPartReader) {
	c.BaseTransitionTest.Load(t, forkName, readPart)
	test_util.LoadSSZ(t, "proposer_slashing", &c.ProposerSlashing, readPart)
}

func (c *ProposerSlashingTestCase) Run() error {
	epc, err := common.NewEpochsContext(c.Spec, c.Pre)
	if err != nil {
		return err
	}
	return phase0.ProcessProposerSlashing(c.Spec, epc, c.Pre, &c.ProposerSlashing)
}

func TestProposerSlashing(t *testing.T) {
	test_util.RunTransitionTest(t, test_util.AllForks, "operations", "proposer_slashing",
		func() test_util.TransitionTest { return new(ProposerSlashingTestCase) })
}
