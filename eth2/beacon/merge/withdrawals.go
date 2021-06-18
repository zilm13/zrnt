package merge

import (
	"fmt"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	. "github.com/protolambda/ztyp/view"
	"github.com/zilm13/zrnt/eth2/beacon/common"
)

type Root = tree.Root

type Bytes32 = Root

const Bytes32Type = RootType

const (
	// TODO: pass REGISTRY_LIMIT from yaml
	WITHDRAWAL_REGISTRY_LIMIT = 1099511627776
)

func WithdrawalRegistryType(spec *common.Spec) ListTypeDef {
	return ComplexListType(WithdrawalType, WITHDRAWAL_REGISTRY_LIMIT)
}

var WithdrawalType = ContainerType("Withdrawal", []FieldDef{
	{"validator_index", common.ValidatorIndexType},
	{"withdrawal_credentials", Bytes32Type},
	{"withdrawn_epoch", Uint64Type},
	{"amount", Uint64Type},
})

type WithdrawalView struct {
	*ContainerView
}

func (v *WithdrawalView) Raw() (*Withdrawal, error) {
	values, err := v.FieldValues()
	if err != nil {
		return nil, err
	}
	if len(values) != 11 {
		return nil, fmt.Errorf("unexpected number of withdrawal fields: %d", len(values))
	}
	validatorIndex, err := common.AsValidatorIndex(values[0], err)
	withdrawalCredentials, err := AsRoot(values[1], err)
	withdrawnEpoch, err := AsUint64(values[2], err)
	amount, err := AsUint64(values[3], err)
	return &Withdrawal{
		ValidatorIndex:        validatorIndex,
		WithdrawalCredentials: withdrawalCredentials,
		WithdrawnEpoch:        uint64(withdrawnEpoch),
		Amount:                uint64(amount),
	}, nil
}

func AsWithdrawal(v View, err error) (*WithdrawalView, error) {
	c, err := AsContainer(v, err)
	return &WithdrawalView{c}, err
}

type Withdrawal struct {
	ValidatorIndex        common.ValidatorIndex `json:"validator_index" yaml:"validator_index"`
	WithdrawalCredentials Bytes32               `json:"withdrawal_credentials" yaml:"withdrawal_credentials"`
	WithdrawnEpoch        uint64                `json:"withdrawn_epoch" yaml:"withdrawn_epoch"`
	Amount                uint64                `json:"amount" yaml:"amount"`
}

func (s *Withdrawal) Deserialize(dr *codec.DecodingReader) error {
	return dr.FixedLenContainer(&s.ValidatorIndex, &s.WithdrawalCredentials,
		(*Uint64View)(&s.WithdrawnEpoch), (*Uint64View)(&s.Amount))
}

func (s *Withdrawal) Serialize(w *codec.EncodingWriter) error {
	return w.FixedLenContainer(&s.ValidatorIndex, &s.WithdrawalCredentials,
		(*Uint64View)(&s.WithdrawnEpoch), (*Uint64View)(&s.Amount))
}

func (s *Withdrawal) ByteLength() uint64 {
	return WithdrawalType.TypeByteLength()
}

func (b *Withdrawal) FixedLength() uint64 {
	return WithdrawalType.TypeByteLength()
}

func (s *Withdrawal) HashTreeRoot(hFn tree.HashFn) common.Root {
	return hFn.HashTreeRoot(&s.ValidatorIndex, &s.WithdrawalCredentials,
		(*Uint64View)(&s.WithdrawnEpoch), (*Uint64View)(&s.Amount))
}

type WithdrawalRegistry []*Withdrawal

func (a *WithdrawalRegistry) Deserialize(dr *codec.DecodingReader) error {
	return dr.List(func() codec.Deserializable {
		i := len(*a)
		*a = append(*a, &Withdrawal{})
		return (*a)[i]
	}, WithdrawalType.TypeByteLength(), WITHDRAWAL_REGISTRY_LIMIT)
}

func (a WithdrawalRegistry) Serialize(w *codec.EncodingWriter) error {
	return w.List(func(i uint64) codec.Serializable {
		return a[i]
	}, WithdrawalType.TypeByteLength(), uint64(len(a)))
}

func (a WithdrawalRegistry) ByteLength() (out uint64) {
	return uint64(len(a)) * WithdrawalType.TypeByteLength()
}

func (a WithdrawalRegistry) FixedLength() uint64 {
	return 0 // it's a list, no fixed length
}

func (li WithdrawalRegistry) HashTreeRoot(hFn tree.HashFn) common.Root {
	length := uint64(len(li))
	return hFn.ComplexListHTR(func(i uint64) tree.HTR {
		if i < length {
			return li[i]
		}
		return nil
	}, length, WITHDRAWAL_REGISTRY_LIMIT)
}

type WithdrawalRegistryView struct{ *ComplexListView }

func AsWithdrawalRegistry(v View, err error) (*WithdrawalRegistryView, error) {
	c, err := AsComplexList(v, err)
	return &WithdrawalRegistryView{c}, nil
}

func (registry *WithdrawalRegistryView) WithdrawalCount() (uint64, error) {
	return registry.Length()
}

//
//func (registry *WithdrawalRegistryView) Withdrawal(index uint64) (Withdrawal, error) {
//	return AsWithdrawal(registry.Get(index))
//}
//
//func (registry *WithdrawalRegistryView) Iter() (next func() (val Withdrawal, ok bool, err error)) {
//	iter := registry.ReadonlyIter()
//	return func() (val Withdrawal, ok bool, err error) {
//		elem, ok, err := iter.Next()
//		if err != nil || !ok {
//			return nil, ok, err
//		}
//		w, err := AsWithdrawal(elem, nil)
//		return w, true, err
//	}
//}

func (registry *WithdrawalRegistryView) IsValidIndex(index uint64) (valid bool, err error) {
	count, err := registry.Length()
	if err != nil {
		return false, err
	}
	return uint64(index) < count, nil
}
