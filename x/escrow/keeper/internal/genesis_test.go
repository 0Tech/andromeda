package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
)

func TestValidateGenesisParams(t *testing.T) {
	_, _, k, _, _ := setupKeepers(t) //nolint:dogsled

	gs := k.DefaultGenesis()
	tester := func(subject escrowv1alpha1.GenesisState_Params) error {
		gs.Params = &subject
		return k.ValidateGenesis(gs)
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.GenesisState_Params]{
		{
			"nil max_metadata_length": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid max_metadata_length": {
				Malleate: func(subject *escrowv1alpha1.GenesisState_Params) {
					subject.MaxMetadataLength = k.DefaultGenesis().Params.MaxMetadataLength
				},
			},
		},
	}

	testutil.DoTest(t, tester, cases)
}

func TestValidateGenesisAgents(t *testing.T) {
	cdc, _, k, _, _ := setupKeepers(t) //nolint:dogsled
	addressCodec := cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addressBytesToString := func(address []byte) string {
		addressStr, err := addressCodec.BytesToString(address)
		assert.NoError(t, err)
		return addressStr
	}

	addresses := simtestutil.CreateIncrementalAccounts(4)
	creatorStr := addressBytesToString(createRandomAccounts(1)[0])

	gs := k.DefaultGenesis()
	tester := func(subject []*escrowv1alpha1.GenesisState_Agent) error {
		gs.Agents = subject
		return k.ValidateGenesis(gs)
	}
	cases := []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{}
	for i := 0; i < len(addresses)/2; i++ {
		addressStr := addressBytesToString(addresses[i*2+1])
		descendingAddressStr := addressBytesToString(addresses[i*2])

		added := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
			{
				"[nil agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						added = false
						*subject = append(*subject, nil)
					},
					Error: func() error {
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"[non-nil agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						added = true
						*subject = append(*subject, &escrowv1alpha1.GenesisState_Agent{})
					},
				},
			},
			{
				"nil address": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid address": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Address = addressStr
					},
				},
				"invalid address": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Address = notInBech32
					},
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrInvalidAddress
					},
				},
			},
			{
				"nil creator]": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid creator]": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Creator = creatorStr
					},
				},
				"invalid creator]": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Creator = notInBech32
					},
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrInvalidAddress
					},
				},
			},
		}...)

		addedDuplicate := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
			{
				"[no duplicate agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						addedDuplicate = false
					},
				},
				"[duplicate agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !added {
							return
						}
						addedDuplicate = true
						*subject = append(*subject, &escrowv1alpha1.GenesisState_Agent{})
					},
					Error: func() error {
						if addedDuplicate {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"valid address": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject)-1].Address = addressStr
					},
				},
			},
			{
				"valid creator]": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject)-1].Creator = creatorStr
					},
				},
			},
		}...)

		addedDescending := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
			{
				"[no descending agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						addedDescending = false
					},
				},
				"[descending agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !added {
							return
						}
						addedDescending = true
						*subject = append(*subject, &escrowv1alpha1.GenesisState_Agent{})
					},
					Error: func() error {
						if addedDescending {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"valid address": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !addedDescending {
							return
						}
						(*subject)[len(*subject)-1].Address = descendingAddressStr
					},
				},
			},
			{
				"valid creator]": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						if !addedDescending {
							return
						}
						(*subject)[len(*subject)-1].Creator = creatorStr
					},
				},
			},
		}...)
	}

	testutil.DoTest(t, tester, cases)
}

func TestValidateGenesisProposals(t *testing.T) {
	cdc, _, k, _, _ := setupKeepers(t) //nolint:dogsled
	addressCodec := cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addressBytesToString := func(address []byte) string {
		addressStr, err := addressCodec.BytesToString(address)
		assert.NoError(t, err)
		return addressStr
	}

	ids := make([]uint64, 4)
	for i := range ids {
		ids[i] = uint64(i) + 1
	}

	var proposerStr, agentStr string
	for _, addrStr := range []*string{
		&proposerStr,
		&agentStr,
	} {
		*addrStr = addressBytesToString(createRandomAccounts(1)[0])
	}

	gs := k.DefaultGenesis()
	tester := func(subject []*escrowv1alpha1.GenesisState_Proposal) error {
		gs.Proposals = subject
		return k.ValidateGenesis(gs)
	}
	cases := []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{}
	for i := 0; i < len(ids)/2; i++ {
		id := ids[i*2+1]
		descendingID := ids[i*2]

		added := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
			{
				"[nil proposal": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						added = false
						*subject = append(*subject, nil)
					},
					Error: func() error {
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"[non-nil proposal": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						added = true
						*subject = append(*subject, &escrowv1alpha1.GenesisState_Proposal{})
					},
				},
			},
			{
				"nil id": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid id": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Id = id
					},
				},
			},
			{
				"nil proposer": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid proposer": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Proposer = proposerStr
					},
				},
				"invalid proposer": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Proposer = notInBech32
					},
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrInvalidAddress
					},
				},
			},
			{
				"nil agent": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Agent = agentStr
					},
				},
				"invalid agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Agent = notInBech32
					},
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrInvalidAddress
					},
				},
			},
			{
				"nil pre_actions": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid pre_actions": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].PreActions = []*codectypes.Any{}
					},
				},
			},
			{
				"nil post_actions]": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid post_actions]": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].PostActions = []*codectypes.Any{}
					},
				},
			},
		}...)

		addedDuplicate := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
			{
				"[no duplicate proposal": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						addedDuplicate = false
					},
				},
				"[duplicate proposal": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						addedDuplicate = true
						*subject = append(*subject, &escrowv1alpha1.GenesisState_Proposal{})
					},
					Error: func() error {
						if addedDuplicate {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"valid id": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject)-1].Id = id
					},
				},
			},
			{
				"valid proposer": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject)-1].Proposer = proposerStr
					},
				},
			},
			{
				"valid agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject)-1].Agent = agentStr
					},
				},
			},
			{
				"valid pre_actions": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject)-1].PreActions = []*codectypes.Any{}
					},
				},
			},
			{
				"valid post_actions]": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject)-1].PostActions = []*codectypes.Any{}
					},
				},
			},
		}...)

		addedDescending := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
			{
				"[no descending proposal": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						addedDescending = false
					},
				},
				"[descending proposal": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						addedDescending = true
						*subject = append(*subject, &escrowv1alpha1.GenesisState_Proposal{})
					},
					Error: func() error {
						if addedDescending {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"valid id": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDescending {
							return
						}
						(*subject)[len(*subject)-1].Id = descendingID
					},
				},
			},
			{
				"valid proposer": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDescending {
							return
						}
						(*subject)[len(*subject)-1].Proposer = proposerStr
					},
				},
			},
			{
				"valid agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDescending {
							return
						}
						(*subject)[len(*subject)-1].Agent = agentStr
					},
				},
			},
			{
				"valid pre_actions": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDescending {
							return
						}
						(*subject)[len(*subject)-1].PreActions = []*codectypes.Any{}
					},
				},
			},
			{
				"valid post_actions]": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !addedDescending {
							return
						}
						(*subject)[len(*subject)-1].PostActions = []*codectypes.Any{}
					},
				},
			},
		}...)
	}

	testutil.DoTest(t, tester, cases)
}

func TestValidateGenesis(t *testing.T) {
	_, _, k, _, _ := setupKeepers(t) //nolint:dogsled

	tester := func(subject escrowv1alpha1.GenesisState) error {
		return k.ValidateGenesis(&subject)
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.GenesisState]{
		{
			"nil params": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid params": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.Params = k.DefaultGenesis().Params
				},
			},
		},
		{
			"nil next_agent": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid next_agent": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.NextAgent = k.DefaultGenesis().NextAgent
				},
			},
		},
		{
			"nil agents": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid agents": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.Agents = k.DefaultGenesis().Agents
				},
			},
		},
		{
			"nil next_proposal": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid next_proposal": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.NextProposal = k.DefaultGenesis().NextProposal
				},
			},
		},
		{
			"nil proposals": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid proposals": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.Proposals = k.DefaultGenesis().Proposals
				},
			},
		},
	}

	testutil.DoTest(t, tester, cases)
}
