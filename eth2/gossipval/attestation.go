package gossipval

import (
	"context"
	"errors"
	"fmt"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/zilm13/zrnt/eth2/util/bls"
	"github.com/protolambda/ztyp/tree"
	"time"
)

const MAXIMUM_GOSSIP_CLOCK_DISPARITY = 500 * time.Millisecond
const ATTESTATION_PROPAGATION_SLOT_RANGE = 32

type AttestationValBackend interface {
	BadBlockValidator
	Spec
	SlotAfter
	Chain
	DomainGetter
	// Checks if the (target epoch, voter) pair was seen, does not do any tracking.
	SeenAttestation(targetEpoch common.Epoch, voter common.ValidatorIndex) bool
}

const catchupTimeout = time.Second * 2

func ValidateAttestation(ctx context.Context, subnet uint64, att *phase0.Attestation,
	attVal AttestationValBackend) GossipValidatorResult {
	spec := attVal.Spec()

	targetSlot, err := spec.EpochStartSlot(att.Data.Target.Epoch)
	if err != nil {
		return GossipValidatorResult{REJECT, fmt.Errorf("cannot get start slot of attestation target epoch %d: %w", att.Data.Target.Epoch, err)}
	}

	// [IGNORE] attestation.data.slot is within the last ATTESTATION_PROPAGATION_SLOT_RANGE slots
	// (within a MAXIMUM_GOSSIP_CLOCK_DISPARITY allowance) --
	// i.e. attestation.data.slot + ATTESTATION_PROPAGATION_SLOT_RANGE >= current_slot >= attestation.data.slot

	// overflow check
	if att.Data.Slot+ATTESTATION_PROPAGATION_SLOT_RANGE < att.Data.Slot {
		return GossipValidatorResult{REJECT, fmt.Errorf("attestation slot overflow: %d", att.Data.Slot)}
	}
	// check minimum, with account for clock disparity
	if minSlot := attVal.SlotAfter(-MAXIMUM_GOSSIP_CLOCK_DISPARITY); att.Data.Slot+ATTESTATION_PROPAGATION_SLOT_RANGE < minSlot {
		return GossipValidatorResult{IGNORE, fmt.Errorf("attestation slot %d is too old, minimum slot is %d", att.Data.Slot, minSlot)}
	}
	// check maximum, with account for clock disparity
	if maxSlot := attVal.SlotAfter(MAXIMUM_GOSSIP_CLOCK_DISPARITY); att.Data.Slot > maxSlot {
		return GossipValidatorResult{IGNORE, fmt.Errorf("attestation slot %d is too new, maximum slot is %d", att.Data.Slot, maxSlot)}
	}

	// [REJECT] The attestation's epoch matches its target --
	// i.e. attestation.data.target.epoch == compute_epoch_at_slot(attestation.data.slot)
	attEpoch := spec.SlotToEpoch(att.Data.Slot)
	if att.Data.Target.Epoch != attEpoch {
		return GossipValidatorResult{REJECT, fmt.Errorf("attestation slot %d is epoch %d and does not match target %d", att.Data.Slot, attEpoch, att.Data.Target.Epoch)}
	}

	// [REJECT] The attestation is unaggregated -- that is, it has exactly one participating validator
	if participants := att.AggregationBits.OnesCount(); participants != 1 {
		return GossipValidatorResult{REJECT, fmt.Errorf("attestation has too many participants set, expected 1, got %d", participants)}
	}

	// [REJECT] The block being voted for (attestation.data.beacon_block_root) passes validation.
	if attVal.IsBadBlock(att.Data.BeaconBlockRoot) {
		return GossipValidatorResult{REJECT, errors.New("attestation voted for invalid block")}
	}

	ch := attVal.Chain()
	// [IGNORE] The block being voted for (attestation.data.beacon_block_root) has been seen
	// (via both gossip and non-gossip sources) (a client MAY queue aggregates for processing once block is retrieved).
	blockRef, ok := ch.ByBlock(att.Data.BeaconBlockRoot)
	if !ok {
		return GossipValidatorResult{IGNORE, errors.New("attestation voted for unknown block")}
	}
	// TODO: this is a nice sanity check, but not strictly necessary if forkchoice handles it anyway.
	if refSlot := blockRef.Step().Slot(); refSlot > att.Data.Slot {
		return GossipValidatorResult{REJECT, errors.New("attestation voted for block in the future")}
	}

	// [REJECT] The attestation's target block is an ancestor of the block named in the LMD vote --
	// i.e. get_ancestor(store, attestation.data.beacon_block_root, compute_start_slot_at_epoch(attestation.data.target.epoch))
	//        == attestation.data.target.root
	if unknown, inSubtree := ch.InSubtree(att.Data.Target.Root, att.Data.BeaconBlockRoot); unknown {
		return GossipValidatorResult{IGNORE, errors.New("unknown block and/or target, cannot check if in subtree")}
	} else if !inSubtree {
		return GossipValidatorResult{REJECT, errors.New("block not in subtree of target")}
	}

	// [REJECT] The current finalized_checkpoint is an ancestor of the block defined
	// by attestation.data.beacon_block_root --
	// i.e. get_ancestor(store, attestation.data.beacon_block_root, compute_start_slot_at_epoch(store.finalized_checkpoint.epoch))
	//        == store.finalized_checkpoint.root
	fin := ch.FinalizedCheckpoint()
	if att.Data.BeaconBlockRoot != fin.Root {
		if unknown, inSubtree := ch.InSubtree(fin.Root, att.Data.BeaconBlockRoot); unknown {
			return GossipValidatorResult{IGNORE, errors.New("unknown block, cannot check if in subtree")}
		} else if !inSubtree {
			return GossipValidatorResult{REJECT, errors.New("block not in subtree of finalized root")}
		}
	} else if fin.Epoch >= att.Data.Target.Epoch {
		return GossipValidatorResult{REJECT, errors.New("cannot vote for finalized root as target")}
	}

	// TODO: additional validation of data.source?

	towardsCtx, cancel := context.WithTimeout(ctx, catchupTimeout)
	defer cancel()
	targetRef, err := ch.Towards(towardsCtx, att.Data.Target.Root, targetSlot)
	if err != nil {
		return GossipValidatorResult{IGNORE, fmt.Errorf("unknown target root %s: %w", att.Data.Target.Root, err)}
	}

	targetEpc, err := targetRef.EpochsContext(ctx)
	if err != nil {
		return GossipValidatorResult{IGNORE, fmt.Errorf("unavailable target epc %s: %w", att.Data.Target.Root, err)}
	}

	// [REJECT] The committee index is within the expected range --
	// i.e. data.index < get_committee_count_per_slot(state, data.target.epoch).
	committeeCountPerSlot, err := targetEpc.GetCommitteeCountAtSlot(att.Data.Slot)
	if err != nil {
		return GossipValidatorResult{REJECT, fmt.Errorf("cannot get commitee count for slot %d: %w", att.Data.Slot, err)}
	}
	if uint64(att.Data.Index) >= committeeCountPerSlot {
		return GossipValidatorResult{REJECT, fmt.Errorf("committee index %d out of range %d", att.Data.Index, committeeCountPerSlot)}
	}

	// [REJECT] The attestation is for the correct subnet --
	// i.e. compute_subnet_for_attestation(committees_per_slot, attestation.data.slot, attestation.data.index)
	//   == subnet_id, where committees_per_slot = get_committee_count_per_slot(state, attestation.data.target.epoch)
	assignedSubnet, err := phase0.ComputeSubnetForAttestation(spec, committeeCountPerSlot, att.Data.Slot, att.Data.Index)
	if err != nil {
		return GossipValidatorResult{REJECT, fmt.Errorf("cannot get subnet for attestation (slot %d, committee index %d): %w", att.Data.Slot, att.Data.Index, err)}
	}
	if subnet != assignedSubnet {
		return GossipValidatorResult{REJECT, fmt.Errorf("attestation (slot %d, committee index %d) received on subnet %d, but should be on subnet %d", att.Data.Slot, att.Data.Index, subnet, assignedSubnet)}
	}

	// [REJECT] The number of aggregation bits matches the committee size -- i.e. len(attestation.aggregation_bits) == len(get_beacon_committee(state, data.slot, data.index))
	committee, err := targetEpc.GetBeaconCommittee(att.Data.Slot, att.Data.Index)
	if err != nil {
		return GossipValidatorResult{REJECT, fmt.Errorf("attestation was validated, but committee is not available: %w", err)}
	}

	if bl := att.AggregationBits.BitLen(); bl != uint64(len(committee)) {
		return GossipValidatorResult{REJECT, fmt.Errorf("attestation has bitlength %d, but expected %d bits", bl, len(committee))}
	}

	// [IGNORE] There has been no other valid attestation seen on an attestation subnet that has an identical attestation.data.target.epoch and participating validator index.
	voter, err := att.AggregationBits.SingleParticipant(committee)
	if err != nil {
		return GossipValidatorResult{REJECT, fmt.Errorf("attestation was expected to have a single voter, but failed: %w", err)}
	}
	if attVal.SeenAttestation(att.Data.Target.Epoch, voter) {
		return GossipValidatorResult{IGNORE, errors.New("attestation vote was already seen (this attestation may be slashable if signature is valid!)")}
	}

	// [REJECT] The signature of attestation is valid.

	// We already know that the voter is part of the committee in the target epoch,
	// we can just hit the cache without further checking the validator index.
	pubkey, ok := targetEpc.PubkeyCache.Pubkey(voter)
	if !ok {
		return GossipValidatorResult{IGNORE, errors.New("failed to find pubkey for voter, cache is wrong")}
	}
	dom, err := attVal.GetDomain(spec.DOMAIN_BEACON_ATTESTER, att.Data.Target.Epoch)
	if err != nil {
		return GossipValidatorResult{IGNORE, errors.New("failed to get domain info for signature check")}
	}
	sigRoot := common.ComputeSigningRoot(att.Data.HashTreeRoot(tree.GetHashFn()), dom)
	if !bls.Verify(pubkey, sigRoot, att.Signature) {
		return GossipValidatorResult{REJECT, errors.New("invalid attestation signature")}
	}

	return GossipValidatorResult{ACCEPT, nil}
}
