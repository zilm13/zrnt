package merge

import (
	"context"
	"fmt"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
)

func (state *BeaconStateView) ProcessEpoch(ctx context.Context, spec *common.Spec, epc *common.EpochsContext) error {
	vals, err := state.Validators()
	if err != nil {
		return err
	}
	flats, err := common.FlattenValidators(vals)
	if err != nil {
		return err
	}
	//// TODO: generalize phase0 attestation processing to support merge.
	//attesterData, err := phase0.ComputeEpochAttesterData(ctx, spec, epc, flats, state)
	//if err != nil {
	//	return err
	//}
	//just := phase0.JustificationStakeData{
	//	CurrentEpoch:                  epc.CurrentEpoch.Epoch,
	//	TotalActiveStake:              epc.TotalActiveStake,
	//	PrevEpochUnslashedTargetStake: attesterData.PrevEpochUnslashedStake.TargetStake,
	//	CurrEpochUnslashedTargetStake: attesterData.CurrEpochUnslashedTargetStake,
	//}
	//if err := phase0.ProcessEpochJustification(ctx, spec, &just, state); err != nil {
	//	return err
	//}
	// TODO: generalize phase0 reward processing to support merge.
	//if err := phase0.ProcessEpochRewardsAndPenalties(ctx, spec, epc, attesterData, state); err != nil {
	//	return err
	//}
	if err := phase0.ProcessEpochRegistryUpdates(ctx, spec, epc, flats, state); err != nil {
		return err
	}
	if err := phase0.ProcessEpochSlashings(ctx, spec, epc, flats, state); err != nil {
		return err
	}
	if err := phase0.ProcessEffectiveBalanceUpdates(ctx, spec, epc, flats, state); err != nil {
		return err
	}
	if err := phase0.ProcessEth1DataReset(ctx, spec, epc, state); err != nil {
		return err
	}
	if err := phase0.ProcessSlashingsReset(ctx, spec, epc, state); err != nil {
		return err
	}
	if err := phase0.ProcessRandaoMixesReset(ctx, spec, epc, state); err != nil {
		return err
	}
	if err := phase0.ProcessHistoricalRootsUpdate(ctx, spec, epc, state); err != nil {
		return err
	}
	// TODO: generalize phase0 attestation processing to support merge.
	//if err := phase0.ProcessParticipationRecordUpdates(ctx, spec, epc, state); err != nil {
	//	return err
	//}
	return nil
}

func (state *BeaconStateView) ProcessBlock(ctx context.Context, spec *common.Spec, epc *common.EpochsContext, benv *common.BeaconBlockEnvelope) error {
	signedBlock, ok := benv.SignedBlock.(*SignedBeaconBlock)
	if !ok {
		return fmt.Errorf("unexpected block type %T in Merge ProcessBlock", benv.SignedBlock)
	}
	block := &signedBlock.Message
	slot, err := state.Slot()
	if err != nil {
		return err
	}
	proposerIndex, err := epc.GetBeaconProposer(slot)
	if err != nil {
		return err
	}
	if err := common.ProcessHeader(ctx, spec, state, block.Header(spec), proposerIndex); err != nil {
		return err
	}
	body := &block.Body
	if err := phase0.ProcessRandaoReveal(ctx, spec, epc, state, body.RandaoReveal); err != nil {
		return err
	}
	if err := phase0.ProcessEth1Vote(ctx, spec, epc, state, body.Eth1Data); err != nil {
		return err
	}
	// Safety checks, in case the user of the function provided too many operations
	if err := body.CheckLimits(spec); err != nil {
		return err
	}

	if err := phase0.ProcessProposerSlashings(ctx, spec, epc, state, body.ProposerSlashings); err != nil {
		return err
	}
	if err := phase0.ProcessAttesterSlashings(ctx, spec, epc, state, body.AttesterSlashings); err != nil {
		return err
	}
	// TODO: generalize phase0 attestation processing to support merge.
	//if err := phase0.ProcessAttestations(ctx, spec, epc, state, body.Attestations); err != nil {
	//	return err
	//}
	if err := phase0.ProcessDeposits(ctx, spec, epc, state, body.Deposits); err != nil {
		return err
	}
	if err := phase0.ProcessVoluntaryExits(ctx, spec, epc, state, body.VoluntaryExits); err != nil {
		return err
	}
	// TODO: implement execution-payload processing.
	return nil
}
