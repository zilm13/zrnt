package sanity

import (
	"bytes"
	"fmt"
	"github.com/golang/snappy"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"github.com/protolambda/ztyp/codec"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

type InitializationTestCase struct {
	Spec          *common.Spec
	GenesisState  *phase0.BeaconStateView
	ExpectedState *phase0.BeaconStateView
	Eth1Timestamp common.Timestamp
	Eth1BlockHash common.Root
	Deposits      []common.Deposit
}

type DepositsCountMeta struct {
	DepositsCount uint64 `yaml:"deposits_count"`
}

type Eth1InitData struct {
	Eth1BlockHash common.Root      `yaml:"eth1_block_hash"`
	Eth1Timestamp common.Timestamp `yaml:"eth1_timestamp"`
}

func (c *InitializationTestCase) Load(t *testing.T, readPart test_util.TestPartReader) {
	c.Spec = readPart.Spec()
	{
		p := readPart.Part("state.ssz_snappy")
		if p.Exists() {
			data, err := ioutil.ReadAll(p)
			test_util.Check(t, err)
			test_util.Check(t, p.Close())
			uncompressed, err := snappy.Decode(nil, data)
			test_util.Check(t, err)
			state, err := phase0.AsBeaconStateView(
				phase0.BeaconStateType(c.Spec).Deserialize(
					codec.NewDecodingReader(bytes.NewReader(uncompressed), uint64(len(uncompressed)))))
			test_util.Check(t, err)
			c.ExpectedState = state
		} else {
			// expecting a failed genesis
			c.ExpectedState = nil
		}
	}
	{
		p := readPart.Part("eth1.yaml")
		dec := yaml.NewDecoder(p)
		var eth1Init Eth1InitData
		test_util.Check(t, dec.Decode(&eth1Init))
		test_util.Check(t, p.Close())
		c.Eth1BlockHash = eth1Init.Eth1BlockHash
		c.Eth1Timestamp = eth1Init.Eth1Timestamp
	}
	m := &DepositsCountMeta{}
	{
		p := readPart.Part("meta.yaml")
		dec := yaml.NewDecoder(p)
		test_util.Check(t, dec.Decode(&m))
		test_util.Check(t, p.Close())
	}
	{
		for i := uint64(0); i < m.DepositsCount; i++ {
			var dep common.Deposit
			test_util.LoadSSZ(t, fmt.Sprintf("deposits_%d", i), &dep, readPart)
			c.Deposits = append(c.Deposits, dep)
		}
	}
}

func (c *InitializationTestCase) Run() error {
	res, _, err := phase0.GenesisFromEth1(c.Spec, c.Eth1BlockHash, c.Eth1Timestamp, c.Deposits, false)
	if err != nil {
		return err
	}
	c.GenesisState = res
	return nil
}

func (c *InitializationTestCase) ExpectingFailure() bool {
	return c.ExpectedState == nil
}

func (c *InitializationTestCase) Check(t *testing.T) {
	if c.ExpectingFailure() {
		t.Errorf("was expecting failure, but no error was raised")
	} else {
		diff, err := test_util.CompareStates(c.Spec, c.GenesisState, c.ExpectedState)
		if err != nil {
			t.Fatal(err)
		}
		if diff != "" {
			t.Errorf("genesis result does not match expectation!\n%s", diff)
		}
	}
}

func TestInitialization(t *testing.T) {
	test_util.RunTransitionTest(t, "genesis", "initialization",
		func() test_util.TransitionTest { return new(InitializationTestCase) })
}
