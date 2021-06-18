package operations

import (
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"testing"
)

type VoluntaryExitTestCase struct {
	test_util.BaseTransitionTest
	VoluntaryExit phase0.SignedVoluntaryExit
}

func (c *VoluntaryExitTestCase) Load(t *testing.T, forkName test_util.ForkName, readPart test_util.TestPartReader) {
	c.BaseTransitionTest.Load(t, forkName, readPart)
	test_util.LoadSSZ(t, "voluntary_exit", &c.VoluntaryExit, readPart)
}

func (c *VoluntaryExitTestCase) Run() error {
	epc, err := common.NewEpochsContext(c.Spec, c.Pre)
	if err != nil {
		return err
	}
	return phase0.ProcessVoluntaryExit(c.Spec, epc, c.Pre, &c.VoluntaryExit)
}

func TestVoluntaryExit(t *testing.T) {
	test_util.RunTransitionTest(t, test_util.AllForks, "operations", "voluntary_exit",
		func() test_util.TransitionTest { return new(VoluntaryExitTestCase) })
}
