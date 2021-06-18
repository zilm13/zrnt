package phase0

import (
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	. "github.com/protolambda/ztyp/view"
	"github.com/zilm13/zrnt/eth2/beacon/common"
)

type Root = tree.Root

type Bytes32 = Root

type Withdrawal struct {
	ValidatorIndex        common.ValidatorIndex `json:"validator_index" yaml:"validator_index"`
	WithdrawalCredentials Bytes32               `json:"withdrawal_credentials" yaml:"withdrawal_credentials"`
	WithdrawnEpoch        common.Epoch          `json:"withdrawn_epoch" yaml:"withdrawn_epoch"`
	Amount                common.Gwei           `json:"amount" yaml:"amount"`
}

func (w *Withdrawal) Deserialize(dr *codec.DecodingReader) error {
	return dr.FixedLenContainer(&w.ValidatorIndex, &w.WithdrawalCredentials,
		&w.WithdrawnEpoch, &w.Amount)
}

func (w *Withdrawal) Serialize(dw *codec.EncodingWriter) error {
	return dw.FixedLenContainer(&w.ValidatorIndex, &w.WithdrawalCredentials,
		&w.WithdrawnEpoch, &w.Amount)
}

func (a *Withdrawal) ByteLength() uint64 {
	return WithdrawalType.TypeByteLength()
}

func (*Withdrawal) FixedLength() uint64 {
	return WithdrawalType.TypeByteLength()
}

func (w *Withdrawal) HashTreeRoot(hFn tree.HashFn) common.Root {
	return hFn.HashTreeRoot(w.ValidatorIndex, w.WithdrawalCredentials, w.WithdrawnEpoch,
		w.Amount)
}

func (w *Withdrawal) View() *WithdrawalView {
	wCred := RootView(w.WithdrawalCredentials)
	c, _ := WithdrawalType.FromFields(
		Uint64View(w.ValidatorIndex),
		&wCred,
		Uint64View(w.WithdrawnEpoch),
		Uint64View(w.Amount),
	)
	return &WithdrawalView{c}
}

var WithdrawalType = ContainerType("Withdrawal", []FieldDef{
	{"validator_index", common.ValidatorIndexType},
	{"withdrawal_credentials", common.Bytes32Type}, // Commitment to pubkey for withdrawals
	{"withdrawn_epoch", common.EpochType},
	{"amount", common.GweiType},
})

const (
	_withdrawalValidatorIndex = iota
	_withdrawalWithdrawalCredentials
	_withdrawalWithdrawnEpoch
	_withdrawalAmount
)

type WithdrawalView struct {
	*ContainerView
}

var _ common.Withdrawal = (*WithdrawalView)(nil)

func NewWithdrawalView() *WithdrawalView {
	return &WithdrawalView{ContainerView: WithdrawalType.New()}
}

func AsWithdrawal(v View, err error) (*WithdrawalView, error) {
	c, err := AsContainer(v, err)
	return &WithdrawalView{c}, err
}

func (v *WithdrawalView) ValidatorIndex() (common.ValidatorIndex, error) {
	return common.AsValidatorIndex(v.Get(_withdrawalValidatorIndex))
}
func (v *WithdrawalView) WithdrawalCredentials() (out common.Root, err error) {
	return AsRoot(v.Get(_withdrawalWithdrawalCredentials))
}
func (v *WithdrawalView) WithdrawnEpoch() (common.Epoch, error) {
	return common.AsEpoch(v.Get(_withdrawalWithdrawnEpoch))
}
func (v *WithdrawalView) Amount() (common.Gwei, error) {
	return common.AsGwei(v.Get(_withdrawalAmount))
}
