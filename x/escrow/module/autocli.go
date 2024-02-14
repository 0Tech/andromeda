package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/client/v2/autocli"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/api/andromeda/escrow/v1alpha1"
)

var _ autocli.HasAutoCLIConfig = (*AppModule)(nil)

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: autoCLIQuery(),
		Tx:    autoCLITx(),
	}
}

func autoCLIQuery() *autocliv1.ServiceCommandDescriptor {
	return &autocliv1.ServiceCommandDescriptor{
		Service: escrowv1alpha1.Query_ServiceDesc.ServiceName,
		RpcCommandOptions: []*autocliv1.RpcCommandOptions{
			autoCLIQueryParams(),
			autoCLIQueryAgent(),
			autoCLIQueryAgentsByCreator(),
			autoCLIQueryAgents(),
			autoCLIQueryProposal(),
			autoCLIQueryProposalsByProposer(),
			autoCLIQueryProposals(),
		},
	}
}

func autoCLIQueryParams() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "Params",
		Short:     "queries the module parameters.",
		Use:       "params",
		Example: `$ and query escrow params
max_metadata_length: "42"`,
	}
}

func autoCLIQueryAgent() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "Agent",
		Short:     "queries an agent.",
		Use:       "agent --agent [agent]",
		FlagOptions: map[string]*autocliv1.FlagOptions{
			"agent": {
				Usage: "the address of an agent",
			},
		},
		Example: `$ and query escrow agent --agent cosmos1aaa...
agent:
  address: cosmos1aaa...
  creator: cosmos1...`,
	}
}

func autoCLIQueryAgentsByCreator() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "AgentsByCreator",
		Short:     "queries all the agents by its creator.",
		Use:       "agents-by-creator --creator [creator]",
		FlagOptions: map[string]*autocliv1.FlagOptions{
			"creator": {
				Usage: "the address of a creator",
			},
		},
		Example: `$ and query escrow agents-by-creator --creator cosmos1ccc...
agents:
- address: cosmos1...
  creator: cosmos1ccc...
- address: cosmos1...
  creator: cosmos1ccc...
pagination:
  total: "2"`,
	}
}

func autoCLIQueryAgents() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
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
	}
}

func autoCLIQueryProposal() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "Proposal",
		Short:     "queries a proposal.",
		Use:       "proposal --agent [agent]",
		FlagOptions: map[string]*autocliv1.FlagOptions{
			"agent": {
				Usage: "the address of an agent in charge",
			},
		},
		Example: `$ and query escrow proposal --agent cosmos1aaa...
proposal:
  agent: cosmos1aaa...
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
	}
}

func autoCLIQueryProposalsByProposer() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "ProposalsByProposer",
		Short:     "queries all the proposals by its proposer.",
		Use:       "proposals-by-proposer --proposer [proposer]",
		FlagOptions: map[string]*autocliv1.FlagOptions{
			"proposer": {
				Usage: "the address of a proposer",
			},
		},
		Example: `$ and query escrow proposals-by-proposer --proposer cosmos1ppp...
pagination:
  total: "1"
proposals:
- agent: cosmos1...
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
  proposer: cosmos1ppp...`,
	}
}

func autoCLIQueryProposals() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "Proposals",
		Short:     "queries all the proposals.",
		Example: `$ and query escrow proposals
pagination:
  total: "2"
proposals:
- agent: cosmos1...
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
	}
}

func autoCLITx() *autocliv1.ServiceCommandDescriptor {
	return &autocliv1.ServiceCommandDescriptor{
		Service: escrowv1alpha1.Msg_ServiceDesc.ServiceName,
		RpcCommandOptions: []*autocliv1.RpcCommandOptions{
			autoCLITxUpdateParams(),
			autoCLITxCreateAgent(),
			autoCLITxSubmitProposal(),
			autoCLITxExec(),
		},
	}
}

func autoCLITxUpdateParams() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "UpdateParams",
		Short:     "updates the module parameters.",
		Use:       "update-params --from [authority] --max-metadata-length [max-metadata-length]",
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
	}
}

func autoCLITxCreateAgent() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
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
	}
}

func autoCLITxSubmitProposal() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
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
				Usage: "the messages which will be executed on the submission",
			},
			"post_actions": {
				Usage: "the messages which will be executed after the actions included in Msg/Exec",
			},
			"metadata": {
				Usage: "any arbitrary metadata attached to the proposal",
			},
		},
		Example: `$ and tx escrow submit-proposal --from cosmos1ppp... --agent cosmos1aaa... \
    --pre-actions '{"@type": "/cosmos.nft.v1beta1.MsgSend",
                    "class_id": "cat",
                    "id": "leopardcat",
                    "sender": "cosmos1ppp...",
                    "receiver": "cosmos1aaa..."}' \
    --post-actions '{"@type": "/cosmos.bank.v1beta1.MsgSend",
                     "from_address": "cosmos1aaa",
                     "to_address": "cosmos1ppp...",
                     "amount": [{"amount": "42", "denom": "stake"}]}' \
    --metadata "sell leopardcat for 42stake"
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
    metadata: sell leopardcat for 42stake
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
        id: leopardcat
        receiver: cosmos1aaa...
        sender: cosmos1ppp...
      proposer: cosmos1ppp...
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
confirm transaction before signing and broadcasting [y/N]:`,
	}
}

func autoCLITxExec() *autocliv1.RpcCommandOptions {
	return &autocliv1.RpcCommandOptions{
		RpcMethod: "Exec",
		Short:     "executes a proposal.",
		Long: `executes a proposal.

Note:
  actions:
    the signer of each message must be either the executor or one of the agents.`,
		Use: "exec --from [executor] --agents [agents] --actions [actions]",
		FlagOptions: map[string]*autocliv1.FlagOptions{
			"agents": {
				Usage: "the addresses of the agents in charge",
			},
			"actions": {
				Usage: "the messages which will be executed on the execution",
			},
		},
		Example: `$ and tx escrow exec --from cosmos1eee... --agents cosmos1aaa... \
    --actions '{"@type": "/cosmos.nft.v1beta1.MsgSend",
                "class_id": "cat",
                "id": "leopardcat",
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
      id: leopardcat
      receiver: cosmos1eee...
      sender: cosmos1aaa...
    - '@type': /cosmos.bank.v1beta1.MsgSend
      amount:
      - amount: "42"
        denom: stake
      from_address: cosmos1eee...
      to_address: cosmos1aaa...
    agents:
    - cosmos1aaa...
    executor: cosmos1eee...
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
confirm transaction before signing and broadcasting [y/N]:`,
	}
}
