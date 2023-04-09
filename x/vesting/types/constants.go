package types

import gridiron "github.com/gridiron-zone/gridiron/types"

const (
	StakingRewardVestingName = "staking_reward_vesting"
	CommunityPoolVestingName = "community_pool_vesting"
	TeamVestingName          = "team_vesting"

	// Strate reserve pool controlled by governance.
	// Not used now, maybe future.
	StrategicReservePoolName = "strategic_reserve_pool"

	StakingRewardVestingTime = gridiron.SecondsPer4Years
	CommunityPoolVestingTime = gridiron.SecondsPer4Years
	TeamVestingTime          = gridiron.SecondsPer4Years

	ClaimVestedPeriod = 10
)
