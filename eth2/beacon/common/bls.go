package common

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zilm13/zrnt/eth2/util/bls"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	. "github.com/protolambda/ztyp/view"
)

type CachedPubkey = bls.CachedPubkey

type BLSPubkey = bls.BLSPubkey

func ViewPubkey(pub *BLSPubkey) *BLSPubkeyView {
	v, _ := BLSPubkeyType.Deserialize(codec.NewDecodingReader(bytes.NewReader(pub[:]), 48))
	return &BLSPubkeyView{v.(*BasicVectorView)}
}

var BLSPubkeyType = BasicVectorType(ByteType, 48)

type BLSSignature = bls.BLSSignature

func ViewSignature(sig *BLSSignature) *BLSSignatureView {
	v, _ := BLSSignatureType.Deserialize(codec.NewDecodingReader(bytes.NewReader(sig[:]), 48))
	return &BLSSignatureView{v.(*BasicVectorView)}
}

var BLSSignatureType = BasicVectorType(ByteType, 96)

// Mixed into a BLS domain to define its type
type BLSDomainType [4]byte

func (p BLSDomainType) MarshalText() ([]byte, error) {
	return []byte("0x" + hex.EncodeToString(p[:])), nil
}

func (p BLSDomainType) String() string {
	return "0x" + hex.EncodeToString(p[:])
}

func (p *BLSDomainType) UnmarshalText(text []byte) error {
	if p == nil {
		return errors.New("cannot decode into nil BLSDomainType")
	}
	if len(text) >= 2 && text[0] == '0' && (text[1] == 'x' || text[1] == 'X') {
		text = text[2:]
	}
	if len(text) != 8 {
		return fmt.Errorf("unexpected length string '%s'", string(text))
	}
	_, err := hex.Decode(p[:], text)
	return err
}

// Sometimes a beacon state is not available, or too much for what it is good for.
// Functions that just need a specific BLS domain can use this function.
type BLSDomainFn func(typ BLSDomainType, epoch Epoch) (BLSDomain, error)

// BLS domain (8 bytes): fork version (32 bits) concatenated with BLS domain type (32 bits)
type BLSDomain [32]byte

func (dom *BLSDomain) Deserialize(dr *codec.DecodingReader) error {
	_, err := dr.Read(dom[:])
	return err
}

func (a *BLSDomain) FixedLength(*Spec) uint64 {
	return 32
}

func (dom BLSDomain) HashTreeRoot(hFn tree.HashFn) Root {
	return Root(dom) // just convert to root type (no hashing involved)
}

func (dom BLSDomain) String() string {
	return "0x" + hex.EncodeToString(dom[:])
}

func ComputeDomain(domainType BLSDomainType, forkVersion Version, genesisValidatorsRoot Root) (out BLSDomain) {
	copy(out[0:4], domainType[:])
	forkDataRoot := ComputeForkDataRoot(forkVersion, genesisValidatorsRoot)
	copy(out[4:32], forkDataRoot[0:28])
	return
}

type SigningData struct {
	ObjectRoot Root
	Domain     BLSDomain
}

func (d *SigningData) HashTreeRoot(hFn tree.HashFn) Root {
	return hFn.HashTreeRoot(d.ObjectRoot, d.Domain)
}

func ComputeSigningRoot(msgRoot Root, dom BLSDomain) Root {
	withDomain := SigningData{
		ObjectRoot: msgRoot,
		Domain:     dom,
	}
	return withDomain.HashTreeRoot(tree.GetHashFn())
}

// For pubkeys/signatures in state, a tree-representation is used. (TODO: cache optimized deserialized/parsed bls points)

type BLSPubkeyView struct {
	*BasicVectorView
}

func AsBLSPubkey(v View, err error) (BLSPubkey, error) {
	if err != nil {
		return BLSPubkey{}, err
	}
	bv, err := AsBasicVector(v, nil)
	if err != nil {
		return BLSPubkey{}, err
	}
	pub := BLSPubkeyView{BasicVectorView: bv}
	var out BLSPubkey
	buf := bytes.NewBuffer(out[:0])
	if err := pub.Serialize(codec.NewEncodingWriter(buf)); err != nil {
		return BLSPubkey{}, nil
	}
	copy(out[:], buf.Bytes())
	return out, nil
}

type BLSSignatureView struct {
	*BasicVectorView
}

func AsBLSSignature(v View, err error) (BLSSignature, error) {
	if err != nil {
		return BLSSignature{}, err
	}
	bv, err := AsBasicVector(v, nil)
	if err != nil {
		return BLSSignature{}, err
	}
	pub := BLSSignatureView{BasicVectorView: bv}
	var out BLSSignature
	buf := bytes.NewBuffer(out[:0])
	if err := pub.Serialize(codec.NewEncodingWriter(buf)); err != nil {
		return BLSSignature{}, nil
	}
	copy(out[:], buf.Bytes())
	return out, nil
}
