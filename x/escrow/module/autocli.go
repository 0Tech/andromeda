package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/client/v2/autocli"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/api/andromeda/escrow/v1alpha1"
)

var _ autocli.HasAutoCLIConfig = (*AppModule)(nil)

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: escrowv1alpha1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Short:     "queries the module parameters.",
					Use:       "params",
					Example: `$ and query escrow params
max_metadata_length: "42"`,
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
					Example: `$ and query escrow agent --address cosmos1aaa...
agent:
  address: cosmos1aaa...
  creator: cosmos1...`,
				},
				{
					RpcMethod: "Agents",
					Short:     "queries all the agents.",
					Example: `$ and query escrow agents
agents:
- address: cosmos1...
  creator: cosmos1...
- address: cosmos1...
  creator: cosmos1...
- address: cosmos1...
  creator: cosmos1...
pagination:
  total: "3"`,
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
					Example: `$ and query escrow proposal --id 3
proposal:
  agent: cosmos1...
  id: "3"
  metadata: very good deal
  post_actions:
  - type: cosmos-sdk/MsgSend
    value:
      amount:
      - amount: "42"
        denom: stake
      from_address: cosmos1...
      to_address: cosmos1...
  pre_actions:
  - type: /cosmos.nft.v1beta1.MsgSend
    value:
      class_id: ...
      id: ...
      receiver: cosmos1...
      sender: cosmos1...
  proposer: cosmos1...`,
				},
				{
					RpcMethod: "Proposals",
					Short:     "queries all the proposals.",
					Example: `$ and query escrow proposals
pagination:
  total: "2"
proposals:
- agent: cosmos1...
  id: "3"
  metadata: very good deal
  post_actions:
  - type: cosmos-sdk/MsgSend
    value:
      amount:
      - amount: "42"
        denom: stake
      from_address: cosmos1...
      to_address: cosmos1...
  pre_actions:
  - type: /cosmos.nft.v1beta1.MsgSend
    value:
      class_id: ...
      id: ...
      receiver: cosmos1...
      sender: cosmos1...
  proposer: cosmos1...
- agent: cosmos1...
  id: "4"
  metadata: limited time offer for you
  post_actions:
  - type: cosmos-sdk/MsgSend
    value:
      amount:
      - amount: "42"
        denom: stake
      from_address: cosmos1...
      to_address: cosmos1...
  - type: /cosmos.nft.v1beta1.MsgSend
    value:
      class_id: ...
      id: ...
      receiver: cosmos1...
      sender: cosmos1...
  pre_actions:
  - type: /cosmos.nft.v1beta1.MsgSend
    value:
      class_id: ...
      id: ...
      receiver: cosmos1...
      sender: cosmos1...
  proposer: cosmos1...`,
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
  max-metadata-length:
    it must be greater than or equal to the current's.`,
					Use: "update-params --from [authority] --max-metadata-length [max-metadata-length]",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"max_metadata_length": {
							Usage: "the maximum length allowed for metadata",
						},
					},
					Example: `$ and tx escrow update-params --from cosmos1aaa... --max-metadata-length 42
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /andromeda.escrow.v1alpha1.MsgUpdateParams
    authority: cosmos1aaa...
    max_metadata_length: "42"
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
confirm transaction before signing and broadcasting [y/N]:`,
				},
				{
					RpcMethod: "CreateAgent",
					Short:     "creates an escrow agent for a proposal.",
					Use:       "create-agent --from [creator]",
					Example: `$ and tx escrow create-agent --from cosmos1ccc...
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /andromeda.escrow.v1alpha1.MsgCreateAgent
    creator: cosmos1ccc...
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
confirm transaction before signing and broadcasting [y/N]:`,
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
					Use: "submit-proposal --from [proposer] --agent [agent] --pre-actions [pre-actions] --post-actions [post-actions] --metadata [metadata]",
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
						"metadata": {
							Usage: "any arbitrary metadata attached to the proposal",
						},
					},
					Example: `$ and tx escrow submit-proposal --from cosmos1ppp... --agent cosmos1aaa... \
    --pre-actions '{"@type": "/cosmos.nft.v1beta1.MsgSend",
                    "class_id": "cat",
                    "id": "octocat",
                    "sender": "cosmos1ppp...",
                    "receiver": "cosmos1aaa..."}' \
    --post-actions '{"@type": "/cosmos.bank.v1beta1.MsgSend",
                     "from_address": "cosmos1aaa",
                     "to_address": "cosmos1ppp...",
                     "amount": [{"amount": "42", "denom": "stake"}]}' \
    --metadata "sell octocat for 42stake"
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /andromeda.escrow.v1alpha1.MsgSubmitProposal
    agent: cosmos1aaa...
    metadata: sell octocat for 42stake
    post_actions:
    - '@type': /cosmos.bank.v1beta1.MsgSend
      amount:
      - amount: "42"
        denom: stake
      from_address: cosmos1aaa...
      to_address: cosmos1ppp...
    pre_actions:
    - '@type': /cosmos.nft.v1beta1.MsgSend
      value:
        class_id: cat
        id: octocat
        receiver: cosmos1aaa...
        sender: cosmos1ppp...
      proposer: cosmos1ppp...
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
confirm transaction before signing and broadcasting [y/N]:`,
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
					Example: `$ and tx escrow exec --proposal 2 --from cosmos1eee... --agent cosmos1... \
    --actions '{"@type": "/cosmos.nft.v1beta1.MsgSend",
                "class_id": "cat",
                "id": "octocat",
                "sender": "cosmos1aaa...",
                "receiver": "cosmos1eee..."}' \
    --actions '{"@type": "/cosmos.bank.v1beta1.MsgSend",
                "from_address": "cosmos1eee...",
                "to_address": "cosmos1aaa...",
                "amount": [{"amount": "42", "denom": "stake"}]}'
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /andromeda.escrow.v1alpha1.MsgExec
    actions:
    - '@type': /cosmos.nft.v1beta1.MsgSend
      class_id: cat
      id: octocat
      receiver: cosmos1eee...
      sender: cosmos1aaa...
    - '@type': /cosmos.bank.v1beta1.MsgSend
      amount:
      - amount: "42"
        denom: stake
      from_address: cosmos1eee...
      to_address: cosmos1aaa...
    agent: cosmos1aaa...
    executor: cosmos1eee...
    proposal: "2"
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
confirm transaction before signing and broadcasting [y/N]:`,
				},
			},
		},
	}
}
