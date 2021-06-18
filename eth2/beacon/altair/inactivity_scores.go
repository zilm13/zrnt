package altair

import (
	"context"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/tree"
	. "github.com/protolambda/ztyp/view"
)

type InactivityScores []Uint64View

func (a *InactivityScores) Deserialize(spec *common.Spec, dr *codec.DecodingReader) error {
	return dr.List(func() codec.Deserializable {
		i := len(*a)
		*a = append(*a, Uint64View(0))
		return &(*a)[i]
	}, Uint64Type.TypeByteLength(), spec.VALIDATOR_REGISTRY_LIMIT)
}

func (a InactivityScores) Serialize(spec *common.Spec, w *codec.EncodingWriter) error {
	return w.List(func(i uint64) codec.Serializable {
		return &a[i]
	}, Uint64Type.TypeByteLength(), uint64(len(a)))
}

func (a InactivityScores) ByteLength(spec *common.Spec) (out uint64) {
	return uint64(len(a)) * Uint64Type.TypeByteLength()
}

func (a *InactivityScores) FixedLength(spec *common.Spec) uint64 {
	return 0 // it's a list, no fixed length
}

func (li InactivityScores) HashTreeRoot(spec *common.Spec, hFn tree.HashFn) common.Root {
	length := uint64(len(li))
	return hFn.Uint64ListHTR(func(i uint64) uint64 {
		return uint64(li[i])
	}, length, spec.VALIDATOR_REGISTRY_LIMIT)
}

func InactivityScoresType(spec *common.Spec) *BasicListTypeDef {
	return BasicListType(Uint64Type, spec.VALIDATOR_REGISTRY_LIMIT)
}

type InactivityScoresView struct {
	*BasicListView
}

func AsInactivityScores(v View, err error) (*InactivityScoresView, error) {
	c, err := AsBasicList(v, err)
	return &InactivityScoresView{c}, err
}

func (v *InactivityScoresView) GetScore(index common.ValidatorIndex) (uint64, error) {
	s, err := AsUint64(v.Get(uint64(index)))
	return uint64(s), err
}

func (v *InactivityScoresView) SetScore(index common.ValidatorIndex, score uint64) error {
	return v.Set(uint64(index), Uint64View(score))
}

func ProcessInactivityUpdates(ctx context.Context, spec *common.Spec, attesterData *EpochAttesterData, state *BeaconStateView) error {
	inactivityScores, err := state.InactivityScores()
	if err != nil {
		return err
	}
	finalized, err := state.FinalizedCheckpoint()
	if err != nil {
		return err
	}
	finalityDelay := attesterData.PrevEpoch - finalized.Epoch
	isInactivityLeak := finalityDelay > spec.MIN_EPOCHS_TO_INACTIVITY_PENALTY

	for _, vi := range attesterData.EligibleIndices {
		if !attesterData.Flats[vi].Slashed && (attesterData.PrevParticipation[vi]&TIMELY_TARGET_FLAG != 0) {
			score, err := inactivityScores.GetScore(vi)
			if err != nil {
				return err
			}
			if score > 0 {
				score -= 1
				if err := inactivityScores.SetScore(vi, score); err != nil {
					return err
				}
			}
		} else if isInactivityLeak {
			score, err := inactivityScores.GetScore(vi)
			if err != nil {
				return err
			}
			if err := inactivityScores.SetScore(vi, score+spec.INACTIVITY_SCORE_BIAS); err != nil {
				return err
			}
		}
	}
	return nil
}
