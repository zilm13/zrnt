package sharding

import (
	"github.com/zilm13/zrnt/eth2/beacon/common"
	. "github.com/protolambda/ztyp/view"
)

type DataCommitmentsSlotVec []DataCommitment

func DataCommitmentsSlotVecType(spec *common.Spec) *ComplexVectorTypeDef {
	return ComplexVectorType(DataCommitmentType, spec.MAX_SHARDS)
}

type DataCommitmentsEpochTable []DataCommitmentsSlotVec

func DataCommitmentsEpochTableType(spec *common.Spec) *ComplexVectorTypeDef {
	return ComplexVectorType(DataCommitmentsSlotVecType(spec), uint64(spec.SLOTS_PER_EPOCH))
}
