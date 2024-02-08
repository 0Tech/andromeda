package internal_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
)

func TestValidateGenesisParams(t *testing.T) {
	_, _, k := setupEscrowKeeper(t)

	tester := func(subject escrowv1alpha1.GenesisState_Params) error {
		gs := k.DefaultGenesis()
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
	cdc, _, k := setupEscrowKeeper(t)
	addressCodec := cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addressBytesToString := func(address []byte) string {
		addressStr, err := addressCodec.BytesToString(address)
		assert.NoError(t, err)
		return addressStr
	}

	addresses := simtestutil.CreateIncrementalAccounts(2)
	creatorStr := addressBytesToString(createRandomAccounts(1)[0])

	tester := func(subject []*escrowv1alpha1.GenesisState_Agent) error {
		gs := k.DefaultGenesis()
		gs.Agents = subject
		return k.ValidateGenesis(gs)
	}
	cases := []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
		{
			"": {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
					*subject = []*escrowv1alpha1.GenesisState_Agent{}
				},
			},
		},
	}
	for _, address := range addresses {
		addressStr := addressBytesToString(address)

		added := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
			{
				"[no agent": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
						added = false
					},
				},
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

		// add duplicate proposal.
		cases = append(cases, map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
			"no duplicate agent": {},
			"duplicate agent": {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
					if !added {
						return
					}
					*subject = append(*subject, &escrowv1alpha1.GenesisState_Agent{
						Address: addressStr,
						Creator: creatorStr,
					})
				},
				Error: func() error {
					if !added {
						return nil
					}
					return escrowv1alpha1.ErrDuplicateEntry
				},
			},
		})
	}

	testutil.DoTest(t, tester, cases)
}

func TestValidateGenesisProposals(t *testing.T) {
	cdc, _, k := setupEscrowKeeper(t)
	addressCodec := cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addressBytesToString := func(address []byte) string {
		addressStr, err := addressCodec.BytesToString(address)
		assert.NoError(t, err)
		return addressStr
	}

	agents := createRandomAccounts(2)
	proposerStr := addressBytesToString(createRandomAccounts(1)[0])

	tester := func(subject []*escrowv1alpha1.GenesisState_Proposal) error {
		gs := k.DefaultGenesis()
		gs.Proposals = subject
		return k.ValidateGenesis(gs)
	}
	cases := []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
		{
			"": {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
					*subject = []*escrowv1alpha1.GenesisState_Proposal{}
				},
			},
		},
	}
	for _, agent := range agents {
		agentStr := addressBytesToString(agent)

		added := false
		cases = append(cases, []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
			{
				"[no proposal": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						added = false
					},
				},
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
				"nil post_actions": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid post_actions": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].PostActions = []*codectypes.Any{}
					},
				},
			},
			{
				"nil metadata": {
					Error: func() error {
						if !added {
							return nil
						}
						return escrowv1alpha1.ErrUnimplemented
					},
				},
				"valid metadata": {
					Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
						if !added {
							return
						}
						(*subject)[len(*subject)-1].Metadata = randomString(int(k.DefaultGenesis().Params.MaxMetadataLength) - 1)
					},
				},
			},
		}...)

		// add duplicate proposal.
		cases = append(cases, map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
			"no duplicate proposal": {},
			"duplicate proposal": {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
					if !added {
						return
					}
					*subject = append(*subject, &escrowv1alpha1.GenesisState_Proposal{
						Agent:       agentStr,
						Proposer:    proposerStr,
						PreActions:  []*codectypes.Any{},
						PostActions: []*codectypes.Any{},
						Metadata:    randomString(int(k.DefaultGenesis().Params.MaxMetadataLength) - 1),
					})
				},
				Error: func() error {
					if !added {
						return nil
					}
					return escrowv1alpha1.ErrDuplicateEntry
				},
			},
		})
	}

	testutil.DoTest(t, tester, cases)
}

func TestValidateGenesis(t *testing.T) {
	_, _, k := setupEscrowKeeper(t)

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

func TestInitExportGenesisParams(t *testing.T) {
	_, ctxBefore, k := setupEscrowKeeper(t)

	tester := func(subject escrowv1alpha1.GenesisState_Params) error {
		gsInput := k.DefaultGenesis()
		gsInput.Params = &subject
		assert.NoError(t, k.ValidateGenesis(gsInput))

		ctxAfter, _ := sdk.UnwrapSDKContext(ctxBefore).CacheContext()
		err := k.InitGenesis(ctxAfter, gsInput)
		if err != nil {
			return err
		}

		paramsBefore, err := k.GetParams(ctxBefore)
		assert.Error(t, err)
		assert.Nil(t, paramsBefore)

		paramsAfter, err := k.GetParams(ctxAfter)
		assert.NoError(t, err)
		assert.NotNil(t, paramsAfter)
		assert.Equal(t, subject.MaxMetadataLength, paramsAfter.MaxMetadataLength)

		gsOutput, err := k.ExportGenesis(ctxAfter)
		assert.NoError(t, err)
		assert.NotNil(t, gsOutput)
		assert.Equal(t, *gsInput, *gsOutput)

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.GenesisState_Params]{
		{
			"valid max_metadata_length": {
				Malleate: func(subject *escrowv1alpha1.GenesisState_Params) {
					subject.MaxMetadataLength = k.DefaultGenesis().Params.MaxMetadataLength
				},
			},
		},
	}

	testutil.DoTest(t, tester, cases)
}

func TestInitExportGenesisAgents(t *testing.T) {
	cdc, ctxBefore, k := setupEscrowKeeper(t)
	addressCodec := cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addressBytesToString := func(address []byte) string {
		addressStr, err := addressCodec.BytesToString(address)
		assert.NoError(t, err)
		return addressStr
	}
	addressStringToBytes := func(address string) sdk.AccAddress {
		addressBytes, err := addressCodec.StringToBytes(address)
		assert.NoError(t, err)
		return addressBytes
	}

	addresses := simtestutil.CreateIncrementalAccounts(4)
	creatorStr := addressBytesToString(createRandomAccounts(1)[0])

	tester := func(subject []*escrowv1alpha1.GenesisState_Agent) error {
		gsInput := k.DefaultGenesis()
		gsInput.Agents = subject
		assert.NoError(t, k.ValidateGenesis(gsInput))

		ctxAfter, _ := sdk.UnwrapSDKContext(ctxBefore).CacheContext()
		err := k.InitGenesis(ctxAfter, gsInput)
		if err != nil {
			return err
		}

		for i, agent := range subject {
			address := addressStringToBytes(agent.Address)
			creator := addressStringToBytes(agent.Creator)

			agentBefore, err := k.GetAgent(ctxBefore, address)
			assert.Error(t, err, i)
			assert.Nil(t, agentBefore, i)

			agentAfter, err := k.GetAgent(ctxAfter, address)
			assert.NoError(t, err, i)
			assert.NotNil(t, agentAfter, i)
			assert.Equal(t, creator, sdk.AccAddress(agentAfter.Creator), i)
		}

		gsOutput, err := k.ExportGenesis(ctxAfter)
		assert.NoError(t, err)
		assert.NotNil(t, gsOutput)
		assert.Equal(t, *gsInput, *gsOutput)

		return nil
	}
	cases := []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
		{
			"": {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
					*subject = []*escrowv1alpha1.GenesisState_Agent{}
				},
			},
		},
	}
	for i, address := range addresses {
		addressStr := addressBytesToString(address)

		cases = append(cases, map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Agent]{
			"": {},
			fmt.Sprintf("agent %d", i): {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Agent) {
					*subject = append(*subject, &escrowv1alpha1.GenesisState_Agent{
						Address: addressStr,
						Creator: creatorStr,
					})
				},
			},
		})
	}

	testutil.DoTest(t, tester, cases)
}

func TestInitExportGenesisProposals(t *testing.T) {
	cdc, ctxBefore, k := setupEscrowKeeper(t)
	addressCodec := cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addressBytesToString := func(address []byte) string {
		addressStr, err := addressCodec.BytesToString(address)
		assert.NoError(t, err)
		return addressStr
	}
	addressStringToBytes := func(address string) sdk.AccAddress {
		addressBytes, err := addressCodec.StringToBytes(address)
		assert.NoError(t, err)
		return addressBytes
	}

	agents := simtestutil.CreateIncrementalAccounts(4)
	proposerStr := addressBytesToString(createRandomAccounts(1)[0])

	tester := func(subject []*escrowv1alpha1.GenesisState_Proposal) error {
		gsInput := k.DefaultGenesis()
		gsInput.Proposals = subject
		assert.NoError(t, k.ValidateGenesis(gsInput))

		ctxAfter, _ := sdk.UnwrapSDKContext(ctxBefore).CacheContext()
		err := k.InitGenesis(ctxAfter, gsInput)
		if err != nil {
			return err
		}

		for i, proposal := range subject {
			agent := addressStringToBytes(proposal.Agent)

			proposalBefore, err := k.GetProposal(ctxBefore, agent)
			assert.Error(t, err, i)
			assert.Nil(t, proposalBefore, i)

			proposalAfter, err := k.GetProposal(ctxAfter, agent)
			assert.NoError(t, err, i)
			assert.NotNil(t, proposalAfter, i)
			assert.Equal(t, proposal.Proposer, addressBytesToString(proposalAfter.Proposer), i)
			assert.Equal(t, proposal.PreActions, proposalAfter.PreActions, i)
			assert.Equal(t, proposal.PostActions, proposalAfter.PostActions, i)
			assert.Equal(t, proposal.Metadata, proposalAfter.Metadata, i)
		}

		gsOutput, err := k.ExportGenesis(ctxAfter)
		assert.NoError(t, err)
		assert.NotNil(t, gsOutput)
		assert.Equal(t, *gsInput, *gsOutput)

		return nil
	}
	cases := []map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
		{
			"": {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
					*subject = []*escrowv1alpha1.GenesisState_Proposal{}
				},
			},
		},
	}
	for i, agent := range agents {
		agentStr := addressBytesToString(agent)

		cases = append(cases, map[string]testutil.Case[[]*escrowv1alpha1.GenesisState_Proposal]{
			"": {},
			fmt.Sprintf("proposal %d", i): {
				Malleate: func(subject *[]*escrowv1alpha1.GenesisState_Proposal) {
					*subject = append(*subject, &escrowv1alpha1.GenesisState_Proposal{
						Agent:       agentStr,
						Proposer:    proposerStr,
						PreActions:  []*codectypes.Any{},
						PostActions: []*codectypes.Any{},
						Metadata:    "metadata",
					})
				},
			},
		})
	}

	testutil.DoTest(t, tester, cases)
}

func TestInitExportGenesis(t *testing.T) {
	_, ctxBefore, k := setupEscrowKeeper(t)

	tester := func(subject escrowv1alpha1.GenesisState) error {
		assert.NoError(t, k.ValidateGenesis(&subject))

		ctxAfter, _ := sdk.UnwrapSDKContext(ctxBefore).CacheContext()
		err := k.InitGenesis(ctxAfter, &subject)
		if err != nil {
			return err
		}

		gsOutput, err := k.ExportGenesis(ctxAfter)
		assert.NoError(t, err)
		assert.NotNil(t, gsOutput)
		assert.Equal(t, subject, *gsOutput)

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.GenesisState]{
		{
			"valid params": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.Params = k.DefaultGenesis().Params
				},
			},
		},
		{
			"valid next_agent": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.NextAgent = k.DefaultGenesis().NextAgent
				},
			},
		},
		{
			"valid agents": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.Agents = k.DefaultGenesis().Agents
				},
			},
		},
		{
			"valid proposals": {
				Malleate: func(subject *escrowv1alpha1.GenesisState) {
					subject.Proposals = k.DefaultGenesis().Proposals
				},
			},
		},
	}

	testutil.DoTest(t, tester, cases)
}
