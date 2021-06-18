package ssz_static

import (
	"bytes"
	"encoding/hex"
	"github.com/golang/snappy"
	"github.com/zilm13/zrnt/eth2/beacon/altair"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/merge"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/eth2/configs"
	"github.com/zilm13/zrnt/tests/spec/test_util"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

type SSZStaticTestCase struct {
	TypeName   string
	Spec       *common.Spec
	Value      interface{}
	Serialized []byte

	Root common.Root
}

func (testCase *SSZStaticTestCase) Run(t *testing.T) {
	// deserialization is the pre-requisite
	{
		r := bytes.NewReader(testCase.Serialized)
		if obj, ok := testCase.Value.(common.SpecObj); ok {
			if err := obj.Deserialize(testCase.Spec, codec.NewDecodingReader(r, uint64(len(testCase.Serialized)))); err != nil {
				t.Fatal(err)
			}
		} else if des, ok := testCase.Value.(codec.Deserializable); ok {
			if err := des.Deserialize(codec.NewDecodingReader(r, uint64(len(testCase.Serialized)))); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatalf("type %s cannot be deserialized", testCase.TypeName)
		}
	}

	t.Run("serialization", func(t *testing.T) {
		var data []byte
		{
			var buf bytes.Buffer
			if obj, ok := testCase.Value.(common.SpecObj); ok {
				if err := obj.Serialize(testCase.Spec, codec.NewEncodingWriter(&buf)); err != nil {
					t.Fatal(err)
				}
			} else if ser, ok := testCase.Value.(codec.Serializable); ok {
				if err := ser.Serialize(codec.NewEncodingWriter(&buf)); err != nil {
					t.Fatal(err)
				}
			} else {
				t.Fatalf("type %s cannot be serialized", testCase.TypeName)
			}
			data = buf.Bytes()
		}

		if len(data) != len(testCase.Serialized) {
			t.Errorf("encoded data has different length: %d (spec) <-> %d (zrnt)\nspec: %x\nzrnt: %x", len(testCase.Serialized), len(data), testCase.Serialized, data)
			return
		}
		for i := 0; i < len(data); i++ {
			if data[i] != testCase.Serialized[i] {
				t.Errorf("byte i: %d differs: %d (spec) <-> %d (zrnt)\nspec: %x\nzrnt: %x", i, testCase.Serialized[i], data[i], testCase.Serialized, data)
				return
			}
		}
	})

	t.Run("hash_tree_root", func(t *testing.T) {
		hfn := tree.GetHashFn()

		var root common.Root
		if obj, ok := testCase.Value.(common.SpecObj); ok {
			root = obj.HashTreeRoot(testCase.Spec, hfn)
		} else if v, ok := testCase.Value.(tree.HTR); ok {
			root = v.HashTreeRoot(hfn)
		} else {
			t.Fatalf("type %s cannot be serialized", testCase.TypeName)
		}
		if root != testCase.Root {
			t.Errorf("hash-tree-roots differ: %s (spec) <-> %s (zrnt)", testCase.Root, root)
			return
		}
	})
}

type ObjAllocator func() interface{}

var objs = map[test_util.ForkName]map[string]ObjAllocator{
	"phase0": {},
	"altair": {},
	"merge":  {},
}

func init() {
	base := map[string]ObjAllocator{
		"AggregateAndProof": func() interface{} { return new(phase0.AggregateAndProof) },
		"Attestation":       func() interface{} { return new(phase0.Attestation) },
		"AttestationData":   func() interface{} { return new(phase0.AttestationData) },
		"AttesterSlashing":  func() interface{} { return new(phase0.AttesterSlashing) },
		"BeaconBlockHeader": func() interface{} { return new(common.BeaconBlockHeader) },
		"Checkpoint":        func() interface{} { return new(common.Checkpoint) },
		"Deposit":           func() interface{} { return new(common.Deposit) },
		"DepositData":       func() interface{} { return new(common.DepositData) },
		//"Eth1Block": func() interface{} { return new(common.Eth1Block) }, // phase0 validator spec remnant
		"Eth1Data":                func() interface{} { return new(common.Eth1Data) },
		"Fork":                    func() interface{} { return new(common.Fork) },
		"ForkData":                func() interface{} { return new(common.ForkData) },
		"HistoricalBatch":         func() interface{} { return new(phase0.HistoricalBatch) },
		"IndexedAttestation":      func() interface{} { return new(phase0.IndexedAttestation) },
		"PendingAttestation":      func() interface{} { return new(phase0.PendingAttestation) },
		"ProposerSlashing":        func() interface{} { return new(phase0.ProposerSlashing) },
		"SignedAggregateAndProof": func() interface{} { return new(phase0.SignedAggregateAndProof) },
		"SignedBeaconBlockHeader": func() interface{} { return new(common.SignedBeaconBlockHeader) },
		"SignedVoluntaryExit":     func() interface{} { return new(phase0.SignedVoluntaryExit) },
		//"SigningData": func() interface{} { return new(common.SigningData) },  // not really encoded/decoded, just HTR
		"Validator":     func() interface{} { return new(phase0.Validator) },
		"VoluntaryExit": func() interface{} { return new(phase0.VoluntaryExit) },
	}
	for k, v := range base {
		objs["phase0"][k] = v
		objs["altair"][k] = v
		objs["merge"][k] = v
	}
	objs["phase0"]["BeaconBlockBody"] = func() interface{} { return new(phase0.BeaconBlockBody) }
	objs["phase0"]["BeaconBlock"] = func() interface{} { return new(phase0.BeaconBlock) }
	objs["phase0"]["BeaconState"] = func() interface{} { return new(phase0.BeaconState) }
	objs["phase0"]["SignedBeaconBlock"] = func() interface{} { return new(phase0.SignedBeaconBlock) }

	objs["altair"]["BeaconBlockBody"] = func() interface{} { return new(altair.BeaconBlockBody) }
	objs["altair"]["BeaconBlock"] = func() interface{} { return new(altair.BeaconBlock) }
	objs["altair"]["BeaconState"] = func() interface{} { return new(altair.BeaconState) }
	objs["altair"]["SignedBeaconBlock"] = func() interface{} { return new(altair.SignedBeaconBlock) }
	objs["altair"]["SyncAggregate"] = func() interface{} { return new(altair.SyncAggregate) }
	//objs["altair"]["LightClientSnapshot"] = func() interface{} { return new(altair.LightClientSnapshot) }
	//objs["altair"]["LightClientUpdate"] = func() interface{} { return new(altair.LightClientUpdate) }
	//objs["altair"]["ContributionAndProof"] = func() interface{} { return new(altair.ContributionAndProof) }
	//objs["altair"]["SignedContributionAndProof"] = func() interface{} { return new(altair.SignedContributionAndProof) }
	//objs["altair"]["SyncAggregatorSelectionData"] = func() interface{} { return new(altair.SyncAggregatorSelectionData) }
	//objs["altair"]["SyncCommittee"] = func() interface{} { return new(altair.SyncCommittee) }
	//objs["altair"]["SyncCommitteeSignature"] = func() interface{} { return new(altair.SyncCommitteeSignature) }

	objs["merge"]["BeaconBlockBody"] = func() interface{} { return new(merge.BeaconBlockBody) }
	objs["merge"]["BeaconBlock"] = func() interface{} { return new(merge.BeaconBlock) }
	objs["merge"]["BeaconState"] = func() interface{} { return new(merge.BeaconState) }
	objs["merge"]["SignedBeaconBlock"] = func() interface{} { return new(merge.SignedBeaconBlock) }
	objs["merge"]["ExecutionPayload"] = func() interface{} { return new(common.ExecutionPayload) }
	objs["merge"]["ExecutionPayloadHeader"] = func() interface{} { return new(common.ExecutionPayloadHeader) }
	//objs["merge"]["PowBlock"] = func() interface{} { return new(merge.PowBlock) }

}

type RootsYAML struct {
	Root string `yaml:"root"`
}

func runSSZStaticTest(fork test_util.ForkName, name string, alloc ObjAllocator, spec *common.Spec) func(t *testing.T) {
	return func(t *testing.T) {
		test_util.RunHandler(t, "ssz_static/"+name, func(t *testing.T, forkName test_util.ForkName, readPart test_util.TestPartReader) {
			c := &SSZStaticTestCase{
				Spec:     readPart.Spec(),
				TypeName: name,
			}

			// Allocate an empty value to decode into later for testing.
			c.Value = alloc()

			// Load the SSZ encoded data as a bytes array. The test will serialize it both ways.
			{
				p := readPart.Part("serialized.ssz_snappy")
				data, err := ioutil.ReadAll(p)
				test_util.Check(t, err)
				uncompressed, err := snappy.Decode(nil, data)
				test_util.Check(t, err)
				test_util.Check(t, p.Close())
				test_util.Check(t, err)
				c.Serialized = uncompressed
			}

			{
				p := readPart.Part("roots.yaml")
				dec := yaml.NewDecoder(p)
				roots := &RootsYAML{}
				test_util.Check(t, dec.Decode(roots))
				test_util.Check(t, p.Close())
				{
					root, err := hex.DecodeString(roots.Root[2:])
					test_util.Check(t, err)
					copy(c.Root[:], root)
				}
			}

			// Run the test case
			c.Run(t)

		}, spec, fork)
	}
}

func TestSSZStatic(t *testing.T) {
	t.Parallel()
	t.Run("minimal", func(t *testing.T) {
		for fork, objByName := range objs {
			t.Run(string(fork), func(t *testing.T) {
				for k, v := range objByName {
					t.Run(k, runSSZStaticTest(fork, k, v, configs.Minimal))
				}
			})
		}
	})
	t.Run("mainnet", func(t *testing.T) {
		for fork, objByName := range objs {
			t.Run(string(fork), func(t *testing.T) {
				for k, v := range objByName {
					t.Run(k, runSSZStaticTest(fork, k, v, configs.Mainnet))
				}
			})
		}
	})
}
