package paych

import (
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	big "github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/cbor"
	v0builtin "github.com/filecoin-project/specs-actors/actors/builtin"

	"github.com/filecoin-project/lotus/chain/actors/adt"
	"github.com/filecoin-project/lotus/chain/types"
)

// Load returns an abstract copy of payment channel state, irregardless of actor version
func Load(store adt.Store, act *types.Actor) (State, error) {
	switch act.Code {
	case v0builtin.PaymentChannelActorCodeID:
		out := v0State{store: store}
		err := store.Get(store.Context(), act.Head, &out)
		if err != nil {
			return nil, err
		}
		return &out, nil
	}
	return nil, xerrors.Errorf("unknown actor code %s", act.Code)
}

// State is an abstract version of payment channel state that works across
// versions
type State interface {
	cbor.Marshaler
	// Channel owner, who has funded the actor
	From() address.Address
	// Recipient of payouts from channel
	To() address.Address

	// Height at which the channel can be `Collected`
	SettlingAt() abi.ChainEpoch

	// Amount successfully redeemed through the payment channel, paid out on `Collect()`
	ToSend() abi.TokenAmount

	// Get total number of lanes
	LaneCount() (uint64, error)

	// Iterate lane states
	ForEachLaneState(cb func(idx uint64, dl LaneState) error) error
}

// LaneState is an abstract copy of the state of a single lane
type LaneState interface {
	Redeemed() big.Int
	Nonce() uint64
}