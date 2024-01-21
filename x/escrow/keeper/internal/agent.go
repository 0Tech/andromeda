package internal

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

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
			return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "failed to set a pubkey")
		}

		acc := k.authKeeper.NewAccount(ctx, bac)
		k.authKeeper.SetAccount(ctx, acc)

		if err := k.setAgent(ctx, address, creator, &escrowv1alpha1.Agent{}); err != nil {
			// TODO: invariant broken
			return nil, err
		}

		return address, nil
	}
}

func (k Keeper) getAgentKey(ctx context.Context, address sdk.AccAddress) (*collections.Pair[sdk.AccAddress, sdk.AccAddress], error) {
	key, err := k.agents.Indexes.address.MatchExact(ctx, address)
	if err != nil {
		if !errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		return nil, escrowv1alpha1.ErrAgentNotFound.Wrap(err.Error())
	}

	return &key, nil
}

func (k Keeper) GetAgent(ctx context.Context, address sdk.AccAddress) (sdk.AccAddress, *escrowv1alpha1.Agent, error) {
	key, err := k.getAgentKey(ctx, address)
	if err != nil {
		return nil, nil, err
	}

	creator := key.K1()

	agent, err := k.agents.Get(ctx, *key)
	if err != nil {
		return nil, nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return creator, &agent, nil
}

func (k Keeper) HasAgent(ctx context.Context, address, creator sdk.AccAddress) error {
	has, err := k.agents.Has(ctx, collections.Join(creator, address))
	if err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	if !has {
		addrStr, err := k.addressBytesToString(address)
		if err != nil {
			return err
		}

		return escrowv1alpha1.ErrAgentNotFound.Wrap(addrStr)
	}

	return nil
}

func (k Keeper) setAgent(ctx context.Context, address, creator sdk.AccAddress, agent *escrowv1alpha1.Agent) error {
	if err := k.agents.Set(ctx, collections.Join(creator, address), *agent); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) removeAgent(ctx context.Context, address sdk.AccAddress) error {
	key, err := k.getAgentKey(ctx, address)
	if err != nil {
		return err
	}

	if err := k.agents.Remove(ctx, *key); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) iterateAgents(ctx context.Context, fn func(address, creator sdk.AccAddress, agent escrowv1alpha1.Agent) error) error {
	iter, err := k.agents.Indexes.address.Iterate(ctx, nil)
	if err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, err := iter.PrimaryKey()
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		value, err := k.agents.Get(ctx, key)
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		creator := key.K1()
		address := key.K2()
		agent := value

		if err := fn(address, creator, agent); err != nil {
			return err
		}
	}

	return nil
}
