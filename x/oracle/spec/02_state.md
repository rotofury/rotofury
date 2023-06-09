<!--
order: 2
-->

# State

## AggregateExchangeRatePrevote

`AggregateExchangeRatePrevote` containing validator voter's aggregated prevote for all denoms for the current `VotePeriod`.

- AggregateExchangeRatePrevote: `0x06<valAddress_Bytes> -> ProtocolBuffer(AggregateExchangeRatePrevote)`

```go
// AggregateVoteHash is hash value to hide vote exchange rates
// which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")
type AggregateVoteHash []byte

type AggregateExchangeRatePrevote struct {
	Hash        AggregateVoteHash // Vote hex hash to protect centralize data source problem
	Voter       sdk.ValAddress    // Voter val address
	SubmitBlock int64
}
```

## AggregateExchangeRateVote

`AggregateExchangeRateVote` containing validator voter's aggregate vote for all denoms for the current `VotePeriod`.

- AggregateExchangeRateVote: `0x07<valAddress_Bytes> -> ProtocolBuffer(AggregateExchangeRateVote)`

```go
type ExchangeRateTuple struct {
	Denom        string  `json:"denom"`
	ExchangeRate sdk.Dec `json:"exchange_rate"`
}

type ExchangeRateTuples []ExchangeRateTuple

type AggregateExchangeRateVote struct {
	ExchangeRateTuples ExchangeRateTuples // ExchangeRates of vote targets against USD
	Voter              sdk.ValAddress     // voter val address of validator
}
```

## ExchangeRate

An `sdk.Dec` that stores the current exchange rate of a given denom against USD, which is used by the Maker module for pricing swaps etc.

- ExchangeRate: `0x03<denom_Bytes> -> ProtocolBuffer(sdk.Dec)`

## FeederDelegation

An `sdk.AccAddress` (`mer-` account) address of `operator`'s delegated price feeder.

- FeederDelegation: `0x04<valAddress_Bytes> -> ProtocolBuffer(sdk.AccAddress)`

## MissCounter

An `int64` representing the number of `VotePeriods` that validator `operator` missed during the current `SlashWindow`.

- MissCounter: `0x05<valAddress_Bytes> -> ProtocolBuffer(int64)`
