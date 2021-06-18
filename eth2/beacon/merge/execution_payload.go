package merge

import (
	"fmt"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	. "github.com/protolambda/ztyp/view"
)

type Hash32 = common.Root

const Hash32Type = RootType

type Bytes32 = common.Root

const Bytes32Type = RootType

var ExecutionPayloadHeaderType = ContainerType("ExecutionPayloadHeader", []FieldDef{
	{"block_hash", Hash32Type},
	{"parent_hash", Hash32Type},
	{"coinbase", common.Eth1AddressType},
	{"state_root", Bytes32Type},
	{"number", Uint64Type},
	{"gas_limit", Uint64Type},
	{"gas_used", Uint64Type},
	{"timestamp", common.TimestampType},
	{"receipt_root", Bytes32Type},
	{"logs_bloom", LogsBloomType},
	{"transactions_root", RootType},
})

type ExecutionPayloadHeaderView struct {
	*ContainerView
}

func (v *ExecutionPayloadHeaderView) Raw() (*ExecutionPayloadHeader, error) {
	values, err := v.FieldValues()
	if err != nil {
		return nil, err
	}
	if len(values) != 11 {
		return nil, fmt.Errorf("unexpected number of execution payload header fields: %d", len(values))
	}
	blockHash, err := AsRoot(values[0], err)
	parentHash, err := AsRoot(values[1], err)
	coinbase, err := common.AsEth1Address(values[2], err)
	stateRoot, err := AsRoot(values[3], err)
	number, err := AsUint64(values[4], err)
	gasLimit, err := AsUint64(values[5], err)
	gasUsed, err := AsUint64(values[6], err)
	timestamp, err := common.AsTimestamp(values[7], err)
	receiptRoot, err := AsRoot(values[8], err)
	logsBloomView, err := AsLogsBloom(values[9], err)
	transactionsRoot, err := AsRoot(values[10], err)
	if err != nil {
		return nil, err
	}
	logsBloom, err := logsBloomView.Raw()
	if err != nil {
		return nil, err
	}
	return &ExecutionPayloadHeader{
		BlockHash:        blockHash,
		ParentHash:       parentHash,
		CoinBase:         coinbase,
		StateRoot:        stateRoot,
		Number:           uint64(number),
		GasLimit:         uint64(gasLimit),
		GasUsed:          uint64(gasUsed),
		Timestamp:        timestamp,
		ReceiptRoot:      receiptRoot,
		LogsBloom:        *logsBloom,
		TransactionsRoot: transactionsRoot,
	}, nil
}

func AsExecutionPayloadHeader(v View, err error) (*ExecutionPayloadHeaderView, error) {
	c, err := AsContainer(v, err)
	return &ExecutionPayloadHeaderView{c}, err
}

type ExecutionPayloadHeader struct {
	BlockHash        Hash32             `json:"block_hash" yaml:"block_hash"`
	ParentHash       Hash32             `json:"parent_hash" yaml:"parent_hash"`
	CoinBase         common.Eth1Address `json:"coinbase" yaml:"coinbase"`
	StateRoot        Bytes32            `json:"state_root" yaml:"state_root"`
	Number           uint64             `json:"number" yaml:"number"`
	GasLimit         uint64             `json:"gas_limit" yaml:"gas_limit"`
	GasUsed          uint64             `json:"gas_used" yaml:"gas_used"`
	Timestamp        common.Timestamp   `json:"timestamp" yaml:"timestamp"`
	ReceiptRoot      Bytes32            `json:"receipt_root" yaml:"receipt_root"`
	LogsBloom        LogsBloom          `json:"logs_bloom" yaml:"logs_bloom"`
	TransactionsRoot common.Root        `json:"transactions_root" yaml:"transactions_root"`
}

func (s *ExecutionPayloadHeader) View() *ExecutionPayloadHeaderView {
	br, pr, sr := RootView(s.BlockHash), RootView(s.ParentHash), RootView(s.StateRoot)
	nr, gl, gu := Uint64View(s.Number), Uint64View(s.GasLimit), Uint64View(s.GasUsed)
	t, rcr, txsr := Uint64View(s.Timestamp), RootView(s.ReceiptRoot), RootView(s.TransactionsRoot)
	v, err := AsExecutionPayloadHeader(ExecutionPayloadHeaderType.FromFields(
		&br, &pr, s.CoinBase.View(),
		&sr, nr, gl, gu, t, &rcr, s.LogsBloom.View(), &txsr))
	if err != nil {
		panic(err)
	}
	return v
}

func (s *ExecutionPayloadHeader) Deserialize(dr *codec.DecodingReader) error {
	return dr.FixedLenContainer(&s.BlockHash, &s.ParentHash, &s.CoinBase, &s.StateRoot,
		(*Uint64View)(&s.Number), (*Uint64View)(&s.GasLimit), (*Uint64View)(&s.GasUsed),
		&s.Timestamp, &s.ReceiptRoot, &s.LogsBloom, &s.TransactionsRoot)
}

func (s *ExecutionPayloadHeader) Serialize(w *codec.EncodingWriter) error {
	return w.FixedLenContainer(&s.BlockHash, &s.ParentHash, &s.CoinBase, &s.StateRoot,
		(*Uint64View)(&s.Number), (*Uint64View)(&s.GasLimit), (*Uint64View)(&s.GasUsed),
		&s.Timestamp, &s.ReceiptRoot, &s.LogsBloom, &s.TransactionsRoot)
}

func (s *ExecutionPayloadHeader) ByteLength() uint64 {
	return ExecutionPayloadHeaderType.TypeByteLength()
}

func (b *ExecutionPayloadHeader) FixedLength() uint64 {
	return ExecutionPayloadHeaderType.TypeByteLength()
}

func (s *ExecutionPayloadHeader) HashTreeRoot(hFn tree.HashFn) common.Root {
	return hFn.HashTreeRoot(&s.BlockHash, &s.ParentHash, &s.CoinBase, &s.StateRoot,
		(*Uint64View)(&s.Number), (*Uint64View)(&s.GasLimit), (*Uint64View)(&s.GasUsed),
		&s.Timestamp, &s.ReceiptRoot, &s.LogsBloom, &s.TransactionsRoot)
}

var ExecutionPayloadType = ContainerType("ExecutionPayload", []FieldDef{
	{"block_hash", Hash32Type},
	{"parent_hash", Hash32Type},
	{"coinbase", common.Eth1AddressType},
	{"state_root", Bytes32Type},
	{"number", Uint64Type},
	{"gas_limit", Uint64Type},
	{"gas_used", Uint64Type},
	{"timestamp", common.TimestampType},
	{"receipt_root", Bytes32Type},
	{"logs_bloom", LogsBloomType},
	{"transactions", PayloadTransactionsType},
})

type ExecutionPayloadView struct {
	*ContainerView
}

func AsExecutionPayload(v View, err error) (*ExecutionPayloadView, error) {
	c, err := AsContainer(v, err)
	return &ExecutionPayloadView{c}, err
}

type ExecutionPayload struct {
	BlockHash    Hash32              `json:"block_hash" yaml:"block_hash"`
	ParentHash   Hash32              `json:"parent_hash" yaml:"parent_hash"`
	CoinBase     common.Eth1Address  `json:"coinbase" yaml:"coinbase"`
	StateRoot    Bytes32             `json:"state_root" yaml:"state_root"`
	Number       uint64              `json:"number" yaml:"number"`
	GasLimit     uint64              `json:"gas_limit" yaml:"gas_limit"`
	GasUsed      uint64              `json:"gas_used" yaml:"gas_used"`
	Timestamp    common.Timestamp    `json:"timestamp" yaml:"timestamp"`
	ReceiptRoot  Bytes32             `json:"receipt_root" yaml:"receipt_root"`
	LogsBloom    LogsBloom           `json:"logs_bloom" yaml:"logs_bloom"`
	Transactions PayloadTransactions `json:"transactions" yaml:"transactions"`
}

func (b *ExecutionPayload) Deserialize(spec *common.Spec, dr *codec.DecodingReader) error {
	return dr.Container(&b.BlockHash, &b.ParentHash, &b.CoinBase, &b.StateRoot,
		(*Uint64View)(&b.Number), (*Uint64View)(&b.GasLimit), (*Uint64View)(&b.GasUsed),
		&b.Timestamp, &b.ReceiptRoot, &b.LogsBloom, spec.Wrap(&b.Transactions))
}

func (b *ExecutionPayload) Serialize(spec *common.Spec, w *codec.EncodingWriter) error {
	return w.Container(&b.BlockHash, &b.ParentHash, &b.CoinBase, &b.StateRoot,
		(*Uint64View)(&b.Number), (*Uint64View)(&b.GasLimit), (*Uint64View)(&b.GasUsed),
		&b.Timestamp, &b.ReceiptRoot, &b.LogsBloom, spec.Wrap(&b.Transactions))
}

func (b *ExecutionPayload) ByteLength(spec *common.Spec) uint64 {
	return codec.ContainerLength(&b.BlockHash, &b.ParentHash, &b.CoinBase, &b.StateRoot,
		(*Uint64View)(&b.Number), (*Uint64View)(&b.GasLimit), (*Uint64View)(&b.GasUsed),
		&b.Timestamp, &b.ReceiptRoot, &b.LogsBloom, spec.Wrap(&b.Transactions))
}

func (a *ExecutionPayload) FixedLength(*common.Spec) uint64 {
	// transactions list is not fixed length, so the whole thing is not fixed length.
	return 0
}

func (b *ExecutionPayload) HashTreeRoot(spec *common.Spec, hFn tree.HashFn) common.Root {
	return hFn.HashTreeRoot(&b.BlockHash, &b.ParentHash, &b.CoinBase, &b.StateRoot,
		(*Uint64View)(&b.Number), (*Uint64View)(&b.GasLimit), (*Uint64View)(&b.GasUsed),
		&b.Timestamp, &b.ReceiptRoot, &b.LogsBloom, spec.Wrap(&b.Transactions))
}

func (ep *ExecutionPayload) Header(spec *common.Spec) *ExecutionPayloadHeader {
	return &ExecutionPayloadHeader{
		BlockHash:        ep.BlockHash,
		ParentHash:       ep.ParentHash,
		CoinBase:         ep.CoinBase,
		StateRoot:        ep.StateRoot,
		Number:           ep.Number,
		GasLimit:         ep.GasLimit,
		GasUsed:          ep.GasUsed,
		Timestamp:        ep.Timestamp,
		ReceiptRoot:      ep.ReceiptRoot,
		LogsBloom:        ep.LogsBloom,
		TransactionsRoot: ep.Transactions.HashTreeRoot(spec, tree.GetHashFn()),
	}
}
