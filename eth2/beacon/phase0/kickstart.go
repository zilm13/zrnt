package phase0

import (
	"errors"
	hbls "github.com/herumi/bls-eth-go-binary/bls"
	"github.com/zilm13/zrnt/eth2/beacon/common"
)

type KickstartValidatorData struct {
	Pubkey                common.BLSPubkey
	WithdrawalCredentials common.Root
	Balance               common.Gwei
}

// To build a genesis state without Eth 1.0 deposits, i.e. directly from a sequence of minimal validator data.
func KickStartState(spec *common.Spec, eth1BlockHash common.Root, time common.Timestamp, validators []KickstartValidatorData) (*BeaconStateView, *common.EpochsContext, error) {
	deps := make([]common.Deposit, len(validators), len(validators))

	for i := range validators {
		v := &validators[i]
		d := &deps[i]
		d.Data = common.DepositData{
			Pubkey:                v.Pubkey,
			WithdrawalCredentials: v.WithdrawalCredentials,
			Amount:                v.Balance,
			Signature:             common.BLSSignature{},
		}
	}

	state, epc, err := GenesisFromEth1(spec, eth1BlockHash, 0, deps, true)
	if err != nil {
		return nil, nil, err
	}
	if err := state.SetGenesisTime(time); err != nil {
		return nil, nil, err
	}
	return state, epc, nil
}

// To build a genesis state without Eth 1.0 deposits, i.e. directly from a sequence of minimal validator data.
func KickStartStateWithSignatures(spec *common.Spec, eth1BlockHash common.Root, time common.Timestamp, validators []KickstartValidatorData, keys [][32]byte) (*BeaconStateView, *common.EpochsContext, error) {
	deps := make([]common.Deposit, len(validators), len(validators))

	for i := range validators {
		v := &validators[i]
		d := &deps[i]
		d.Data = common.DepositData{
			Pubkey:                v.Pubkey,
			WithdrawalCredentials: v.WithdrawalCredentials,
			Amount:                v.Balance,
			Signature:             common.BLSSignature{},
		}
		var secKey hbls.SecretKey
		if err := secKey.Deserialize(keys[i][:]); err != nil {
			return nil, nil, err
		}
		dom := common.ComputeDomain(common.DOMAIN_DEPOSIT, spec.GENESIS_FORK_VERSION, common.Root{})
		msg := common.ComputeSigningRoot(d.Data.MessageRoot(), dom)
		sig := secKey.SignHash(msg[:])
		var p common.BLSPubkey
		copy(p[:], secKey.GetPublicKey().Serialize())
		if p != d.Data.Pubkey {
			return nil, nil, errors.New("privkey invalid, expected different pubkey")
		}
		copy(d.Data.Signature[:], sig.Serialize())
	}

	state, epc, err := GenesisFromEth1(spec, eth1BlockHash, 0, deps, true)
	if err != nil {
		return nil, nil, err
	}
	if err := state.SetGenesisTime(time); err != nil {
		return nil, nil, err
	}
	return state, epc, nil
}
