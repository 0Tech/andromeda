package internal

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) CreateAgent(ctx context.Context, creator sdk.AccAddress) (sdk.AccAddress, error) {
	for {
		agentNum, err := k.nextAgent.Next(ctx)
		if err != nil {
			return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		derivationKey := make([]byte, 8)
		binary.BigEndian.PutUint64(derivationKey, agentNum)

		ac, err := authtypes.NewModuleCredential(escrowv1alpha1.ModuleName, []byte("agent"), derivationKey)
		if err != nil {
			return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		address := sdk.AccAddress(ac.Address())
		if k.authKeeper.HasAccount(ctx, address) {
			// collision found. retry.
			continue
		}

		addressStr, err := k.addressBytesToString(address)
		if err != nil {
			return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}
		bac := &authtypes.BaseAccount{
			Address: addressStr,
		}

		if err := bac.SetPubKey(ac); err != nil {
			return nil, errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "failed to set a pubkey")
		}

		acc := k.authKeeper.NewAccount(ctx, bac)
		k.authKeeper.SetAccount(ctx, acc)

		if err := k.setAgent(ctx, address, &escrowv1alpha1.Agent{
			Creator: creator,
		}); err != nil {
			return nil, err
		}

		return address, nil
	}
}

func (k Keeper) GetAgent(ctx context.Context, address sdk.AccAddress) (*escrowv1alpha1.Agent, error) {
	agent, err := k.agents.Get(ctx, address)
	if err != nil {
		if !errors.IsOf(err, collections.ErrNotFound) {
			return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		return nil, escrowv1alpha1.ErrAgentNotFound.Wrap(err.Error())
	}

	return &agent, nil
}

func (k Keeper) setAgent(ctx context.Context, address sdk.AccAddress, agent *escrowv1alpha1.Agent) error {
	if err := k.agents.Set(ctx, address, *agent); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) removeAgent(ctx context.Context, address sdk.AccAddress) error {
	if err := k.agents.Remove(ctx, address); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) iterateAgents(ctx context.Context, fn func(address sdk.AccAddress, agent escrowv1alpha1.Agent) error) error {
	iter, err := k.agents.Iterate(ctx, nil)
	if err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		value, err := iter.Value()
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		address := key
		agent := value

		if err := fn(address, agent); err != nil {
			return err
		}
	}

	return nil
}
