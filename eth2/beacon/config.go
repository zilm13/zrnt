package beacon

import "math/big"

type Phase0Config struct {

	// Misc.
	MAX_COMMITTEES_PER_SLOT      uint64
	TARGET_COMMITTEE_SIZE        uint64
	MAX_VALIDATORS_PER_COMMITTEE uint64
	MIN_PER_EPOCH_CHURN_LIMIT    uint64
	CHURN_LIMIT_QUOTIENT         uint64
	SHUFFLE_ROUND_COUNT          uint8

	// Genesis.
	MIN_GENESIS_ACTIVE_VALIDATOR_COUNT uint64
	MIN_GENESIS_TIME                   Timestamp

	// Balance math
	HYSTERESIS_QUOTIENT            uint64
	HYSTERESIS_DOWNWARD_MULTIPLIER uint64
	HYSTERESIS_UPWARD_MULTIPLIER   uint64

	// Phase0 tweaks
	PROPORTIONAL_SLASHING_MULTIPLIER uint64

	// Fork Choice
	SAFE_SLOTS_TO_UPDATE_JUSTIFIED uint64

	// Validator
	ETH1_FOLLOW_DISTANCE                  uint64
	TARGET_AGGREGATORS_PER_COMMITTEE      uint64
	RANDOM_SUBNETS_PER_VALIDATOR          uint64
	EPOCHS_PER_RANDOM_SUBNET_SUBSCRIPTION uint64
	SECONDS_PER_ETH1_BLOCK                uint64

	// Deposit contract
	DEPOSIT_CHAIN_ID         uint64
	DEPOSIT_NETWORK_ID       uint64
	DEPOSIT_CONTRACT_ADDRESS [20]byte // TODO eth1 address type

	// Gwei values
	MIN_DEPOSIT_AMOUNT          Gwei
	MAX_EFFECTIVE_BALANCE       Gwei
	EJECTION_BALANCE            Gwei
	EFFECTIVE_BALANCE_INCREMENT Gwei

	// Initial values
	GENESIS_FORK_VERSION  Version
	BLS_WITHDRAWAL_PREFIX [1]byte

	// Time parameters
	GENESIS_DELAY                       uint64
	SECONDS_PER_SLOT                    uint64
	MIN_ATTESTATION_INCLUSION_DELAY     uint64
	SLOTS_PER_EPOCH                     uint64
	MIN_SEED_LOOKAHEAD                  uint64
	MAX_SEED_LOOKAHEAD                  uint64
	EPOCHS_PER_ETH1_VOTING_PERIOD       uint64
	SLOTS_PER_HISTORICAL_ROOT           uint64
	MIN_VALIDATOR_WITHDRAWABILITY_DELAY uint64
	SHARD_COMMITTEE_PERIOD              uint64
	MIN_EPOCHS_TO_INACTIVITY_PENALTY    uint64

	// State vector lengths
	EPOCHS_PER_HISTORICAL_VECTOR uint64
	EPOCHS_PER_SLASHINGS_VECTOR  uint64
	HISTORICAL_ROOTS_LIMIT       uint64
	VALIDATOR_REGISTRY_LIMIT     uint64

	// Reward and penalty quotients
	BASE_REWARD_FACTOR            uint64
	WHISTLEBLOWER_REWARD_QUOTIENT uint64
	PROPOSER_REWARD_QUOTIENT      uint64
	INACTIVITY_PENALTY_QUOTIENT   uint64
	MIN_SLASHING_PENALTY_QUOTIENT uint64

	// Max operations per block
	MAX_PROPOSER_SLASHINGS uint64
	MAX_ATTESTER_SLASHINGS uint64
	MAX_ATTESTATIONS       uint64
	MAX_DEPOSITS           uint64
	MAX_VOLUNTARY_EXITS    uint64

	// Signature domains
	DOMAIN_BEACON_PROPOSER     BLSDomain
	DOMAIN_BEACON_ATTESTER     BLSDomain
	DOMAIN_RANDAO              BLSDomain
	DOMAIN_DEPOSIT             BLSDomain
	DOMAIN_VOLUNTARY_EXIT      BLSDomain
	DOMAIN_SELECTION_PROOF     BLSDomain
	DOMAIN_AGGREGATE_AND_PROOF BLSDomain
}

type Phase1Config struct {
	// phase1-fork
	PHASE_1_FORK_VERSION  Version
	PHASE_1_FORK_SLOT     uint64
	INITIAL_ACTIVE_SHARDS uint64

	// beacon-chain
	MAX_SHARDS                      uint64
	LIGHT_CLIENT_COMMITTEE_SIZE     uint64
	GASPRICE_ADJUSTMENT_COEFFICIENT uint64

	// Shard block configs
	MAX_SHARD_BLOCK_SIZE             uint64
	TARGET_SHARD_BLOCK_SIZE          uint64
	SHARD_BLOCK_OFFSETS              []uint64
	MAX_SHARD_BLOCKS_PER_ATTESTATION uint64
	BYTES_PER_CUSTODY_CHUNK          uint64
	CUSTODY_RESPONSE_DEPTH           uint64

	// Gwei values
	MAX_GASPRICE uint64
	MIN_GASPRICE uint64

	// Time parameters
	ONLINE_PERIOD                 uint64
	LIGHT_CLIENT_COMMITTEE_PERIOD uint64

	// Max operations per block
	MAX_CUSTODY_CHUNK_CHALLENGE_RECORDS uint64

	// Domain types
	DOMAIN_SHARD_PROPOSAL  BLSDomain
	DOMAIN_SHARD_COMMITTEE BLSDomain
	DOMAIN_LIGHT_CLIENT    BLSDomain

	// custody-game domains
	DOMAIN_CUSTODY_BIT_SLASHING      BLSDomain
	DOMAIN_LIGHT_SELECTION_PROOF     BLSDomain
	DOMAIN_LIGHT_AGGREGATE_AND_PROOF BLSDomain

	// Custody game
	RANDAO_PENALTY_EPOCHS                          uint64
	EARLY_DERIVED_SECRET_PENALTY_MAX_FUTURE_EPOCHS uint64
	EPOCHS_PER_CUSTODY_PERIOD                      uint64
	CUSTODY_PERIOD_TO_RANDAO_PADDING               uint64
	MAX_CHUNK_CHALLENGE_DELAY                      uint64

	CUSTODY_PRIME                *big.Int
	CUSTODY_SECRETS              uint64
	BYTES_PER_CUSTODY_ATOM       uint64
	CUSTODY_PROBABILITY_EXPONENT uint64

	// Max operations
	MAX_CUSTODY_KEY_REVEALS          uint64
	MAX_EARLY_DERIVED_SECRET_REVEALS uint64
	MAX_CUSTODY_CHUNK_CHALLENGES     uint64
	MAX_CUSTODY_CHUNK_CHALLENGE_RESP uint64
	MAX_CUSTODY_SLASHINGS            uint64

	// Reward and penalty quotients
	EARLY_DERIVED_SECRET_REVEAL_SLOT_REWARD_MULTIPLE uint64
	MINOR_REWARD_QUOTIENT                            uint64
}

type Spec struct {
	Phase0Config
	Phase1Config
}
