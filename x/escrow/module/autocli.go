package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/api/andromeda/escrow/v1alpha1"
)

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: escrowv1alpha1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Short:     "queries the module params.",
					Use:       "params",
				},
				{
					RpcMethod: "Agent",
					Short:     "queries an agent.",
					Use:       "agent --address [address]",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"address": {
							Usage: "the address of an agent",
						},
					},
				},
				{
					RpcMethod: "Agents",
					Short:     "queries all the agents.",
				},
				{
					RpcMethod: "Proposal",
					Short:     "queries a proposal.",
					Use:       "proposal --id [id]",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"id": {
							Usage: "the identifier of a proposal",
						},
					},
				},
				{
					RpcMethod: "Proposals",
					Short:     "queries all the proposals.",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: escrowv1alpha1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Short:     "updates the module parameters.",
					Long: `updates the module parameters.

Note:
  params:
    All parameters must be supplied.`,
					Use: "update-params --from [authority] --params [params]",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"params": {
							Usage: "the parameters to update",
						},
					},
				},
				{
					RpcMethod: "CreateAgent",
					Short:     "creates an escrow agent for a proposal.",
					Use:       "create-agent --from [creator]",
				},
				{
					RpcMethod: "SubmitProposal",
					Short:     "submits a proposal.",
					Long: `submits a proposal.

Note:
  agent:
    it must be created by the proposer.
  pre-actions:
    the signer of each message must be either the proposer or the agent.
  post-actions:
    the signer of each message must be either the proposer or the agent.`,
					Use: "submit-proposal --from [proposer] --agent [agent] --pre-actions [pre-actions] --post-actions [post-actions]",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"agent": {
							Usage: "the address of the agent in charge",
						},
						"pre_actions": {
							Usage:        "the messages which will be executed on the submission",
							DefaultValue: "{}",
						},
						"post_actions": {
							Usage:        "the messages which will be executed after the actions included in Msg/Exec",
							DefaultValue: "{}",
						},
					},
				},
				{
					RpcMethod: "Exec",
					Short:     "executes a proposal.",
					Long: `executes a proposal.

Note:
  actions:
    the signer of each message must be either the executor or the agent.`,
					Use: "exec --proposal [proposal] --from [executor] --agent [agent] --actions [actions]",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"proposal": {
							Usage: "the identifier of the proposal",
						},
						"agent": {
							Usage: "the address of the agent in charge",
						},
						"actions": {
							Usage:        "the messages which will be executed on the execution",
							DefaultValue: "{}",
						},
					},
				},
			},
		},
	}
}
