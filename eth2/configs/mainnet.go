package configs

import (
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

var Mainnet = &common.Spec{
	CONFIG_NAME: "mainnet",
	Phase0Config: common.Phase0Config{
		MAX_COMMITTEES_PER_SLOT:               64,
		TARGET_COMMITTEE_SIZE:                 128,
		MAX_VALIDATORS_PER_COMMITTEE:          2048,
		MIN_PER_EPOCH_CHURN_LIMIT:             4,
		CHURN_LIMIT_QUOTIENT:                  1 << 16,
		SHUFFLE_ROUND_COUNT:                   90,
		MIN_GENESIS_ACTIVE_VALIDATOR_COUNT:    1 << 14,
		MIN_GENESIS_TIME:                      1606824000,
		HYSTERESIS_QUOTIENT:                   4,
		HYSTERESIS_DOWNWARD_MULTIPLIER:        1,
		HYSTERESIS_UPWARD_MULTIPLIER:          5,
		SAFE_SLOTS_TO_UPDATE_JUSTIFIED:        8,
		ETH1_FOLLOW_DISTANCE:                  2048,
		TARGET_AGGREGATORS_PER_COMMITTEE:      16,
		RANDOM_SUBNETS_PER_VALIDATOR:          1,
		EPOCHS_PER_RANDOM_SUBNET_SUBSCRIPTION: 256,
		SECONDS_PER_ETH1_BLOCK:                14,
		DEPOSIT_CHAIN_ID:                      1,
		DEPOSIT_NETWORK_ID:                    1,
		DEPOSIT_CONTRACT_ADDRESS:              [20]byte{0x00, 0x00, 0x00, 0x00, 0x21, 0x9a, 0xb5, 0x40, 0x35, 0x6c, 0xBB, 0x83, 0x9C, 0xbe, 0x05, 0x30, 0x3d, 0x77, 0x05, 0xFa},
		MIN_DEPOSIT_AMOUNT:                    1000_000_000,
		MAX_EFFECTIVE_BALANCE:                 32_000_000_000,
		EJECTION_BALANCE:                      16_000_000_000,
		EFFECTIVE_BALANCE_INCREMENT:           1_000_000_000,
		GENESIS_FORK_VERSION:                  common.Version{0x00, 0x00, 0x00, 0x00},
		BLS_WITHDRAWAL_PREFIX:                 [1]byte{0x00},
		GENESIS_DELAY:                         604800,
		SECONDS_PER_SLOT:                      12,
		MIN_ATTESTATION_INCLUSION_DELAY:       1,
		SLOTS_PER_EPOCH:                       32,
		MIN_SEED_LOOKAHEAD:                    1,
		MAX_SEED_LOOKAHEAD:                    4,
		EPOCHS_PER_ETH1_VOTING_PERIOD:         64,
		SLOTS_PER_HISTORICAL_ROOT:             8192,
		MIN_VALIDATOR_WITHDRAWABILITY_DELAY:   256,
		SHARD_COMMITTEE_PERIOD:                256,
		MIN_EPOCHS_TO_INACTIVITY_PENALTY:      4,
		EPOCHS_PER_HISTORICAL_VECTOR:          1 << 16,
		EPOCHS_PER_SLASHINGS_VECTOR:           1 << 13,
		HISTORICAL_ROOTS_LIMIT:                1 << 24,
		VALIDATOR_REGISTRY_LIMIT:              1 << 40,
		BASE_REWARD_FACTOR:                    64,
		PROPORTIONAL_SLASHING_MULTIPLIER:      1,
		WHISTLEBLOWER_REWARD_QUOTIENT:         512,
		PROPOSER_REWARD_QUOTIENT:              8,
		INACTIVITY_PENALTY_QUOTIENT:           1 << 26,
		MIN_SLASHING_PENALTY_QUOTIENT:         128,
		MAX_PROPOSER_SLASHINGS:                16,
		MAX_ATTESTER_SLASHINGS:                2,
		MAX_ATTESTATIONS:                      128,
		MAX_DEPOSITS:                          16,
		MAX_VOLUNTARY_EXITS:                   16,
		DOMAIN_BEACON_PROPOSER:                common.BLSDomainType{0x00, 0x00, 0x00, 0x00},
		DOMAIN_BEACON_ATTESTER:                common.BLSDomainType{0x01, 0x00, 0x00, 0x00},
		DOMAIN_RANDAO:                         common.BLSDomainType{0x02, 0x00, 0x00, 0x00},
		DOMAIN_DEPOSIT:                        common.BLSDomainType{0x03, 0x00, 0x00, 0x00},
		DOMAIN_VOLUNTARY_EXIT:                 common.BLSDomainType{0x04, 0x00, 0x00, 0x00},
		DOMAIN_SELECTION_PROOF:                common.BLSDomainType{0x05, 0x00, 0x00, 0x00},
		DOMAIN_AGGREGATE_AND_PROOF:            common.BLSDomainType{0x06, 0x00, 0x00, 0x00},
	},
	AltairConfig: common.AltairConfig{
		ALTAIR_FORK_EPOCH:                       ^common.Epoch(0),
		ALTAIR_FORK_VERSION:                     common.Version{0x01, 0x00, 0x00, 0x00},
		INACTIVITY_PENALTY_QUOTIENT_ALTAIR:      3 * (1 << 24),
		MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR:    64,
		PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR: 2,
		SYNC_COMMITTEE_SIZE:                     512,
		EPOCHS_PER_SYNC_COMMITTEE_PERIOD:        256,
		INACTIVITY_SCORE_BIAS:                   4,
		INACTIVITY_SCORE_RECOVERY_RATE:          16,
		MIN_SYNC_COMMITTEE_PARTICIPANTS:         1,
		DOMAIN_SYNC_COMMITTEE:                   common.BLSDomainType{0x07, 0x00, 0x00, 0x00},
		DOMAIN_SYNC_COMMITTEE_SELECTION_PROOF:   common.BLSDomainType{0x08, 0x00, 0x00, 0x00},
		DOMAIN_CONTRIBUTION_AND_PROOF:           common.BLSDomainType{0x09, 0x00, 0x00, 0x00},
	},
	MergeConfig: common.MergeConfig{
		MERGE_FORK_EPOCH:   ^common.Epoch(0),
		MERGE_FORK_VERSION: common.Version{0x02, 0x00, 0x00, 0x00},
	},
	ShardingConfig: common.ShardingConfig{
		SHARDING_FORK_EPOCH:             ^common.Epoch(0),
		SHARDING_FORK_VERSION:           common.Version{0x03, 0x00, 0x00, 0x00},
		MAX_SHARDS:                      1024,
		INITIAL_ACTIVE_SHARDS:           64,
		GASPRICE_ADJUSTMENT_COEFFICIENT: 8,
		MAX_SHARD_PROPOSER_SLASHINGS:    16,
		MAX_SHARD_HEADERS_PER_SHARD:     4,
		SHARD_STATE_MEMORY_SLOTS:        256,
		MAX_SAMPLES_PER_BLOCK:           2048,
		TARGET_SAMPLES_PER_BLOCK:        1024,
		MAX_GASPRICE:                    1 << 33,
		MIN_GASPRICE:                    8,
		DOMAIN_SHARD_PROPOSER:           common.BLSDomainType{0x80, 0x00, 0x00, 0x00},
		DOMAIN_SHARD_COMMITTEE:          common.BLSDomainType{0x81, 0x00, 0x00, 0x00},
	},
}
