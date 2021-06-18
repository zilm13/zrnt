package operations

import (
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"testing"
)

type DepositTestCase struct {
	test_util.BaseTransitionTest
	Deposit common.Deposit
}

func (c *DepositTestCase) Load(t *testing.T, readPart test_util.TestPartReader) {
	c.BaseTransitionTest.Load(t, readPart)
	test_util.LoadSSZ(t, "deposit", &c.Deposit, readPart)
}

func (c *DepositTestCase) Run() error {
	epc, err := common.NewEpochsContext(c.Spec, c.Pre)
	if err != nil {
		return err
	}
	return phase0.ProcessDeposit(c.Spec, epc, c.Pre, &c.Deposit, false)
}

func TestDeposit(t *testing.T) {
	test_util.RunTransitionTest(t, "operations", "deposit",
		func() test_util.TransitionTest { return new(DepositTestCase) })
}
