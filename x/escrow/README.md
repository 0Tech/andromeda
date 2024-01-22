# `x/escrow`

## Abstract

This modules provides means to execute messages whose signers come from
different parties. Typical example would be a trade, by which two individual
account can exchange assets. Without this module, they have to use a tx which
include the messages that requires signatures from both of them. Or, they may
create a contract into x/wasm. The former requires a tool to safely exchange
the half-signed transaction (which x/group tries to avoid), while the latter
involves cumbersome preparation for certain tailer made contract code. The
trade itself would be quite simple and obvious, hence this module tries to
reduce such efforts on those tedious use cases.


## Contents

* [Concepts](#concepts)
* [State](#state)
* [Msg Service](#msg-service)
    * [Msg/UpdateParams](#msgupdateparams)
    * [Msg/CreateAgent](#msgcreateagent)
    * [Msg/SubmitProposal](#msgsubmitproposal)
    * [Msg/Exec](#msgexec)
* [Events](#events)
    * [EventUpdateParams](#eventupdateparams)
    * [EventCreateAgent](#eventcreateagent)
    * [EventSubmitProposal](#eventsubmitproposal)
    * [EventExec](#eventexec)
* [Client](#client)
    * [CLI](#cli)
    * [gRPC](#grpc)


## Concepts

### Agent

In abstract manner, the purpose of agents is to provide a place holder for
proposals. In real life, an agent is vital component for safe transactions,
which arranges the procedure as a neutral actor.

Anyone can create agents as many as they want.

### Proposal

There are two parties involved in a transaction. The first one, proposer, is
who suggests a transaction by broadcasting Msg/SubmitProposal. The second one,
executor, is who accepts the transaction, by broadcasting Msg/Exec.

A transaction consists of pre-actions, actions and post-actions. All three
actions will be executed in sequence, and all the messages must be successful
to complete the transaction.

Pre-actions are steps for transferring the control on certain assets from the
proposer to the agent, so determined by the proposer on Msg/SubmitProposal. To
ensure the successful transfer at Msg/Exec, pre-actions are executed at the
proposal submission.

Post-actions are steps for transferring the control on certain assets from the
agent to the proposer. Post-actions ensure the proposer gets the proper reward,
so it is determined by the proposer on Msg/SubmitProposal.

Actions are steps for transferring the control on certain assets from the
executor to agent and vice versa, so determined by the executor on Msg/Exec.
The messages included in actions would take assets on the agent (reserved by
the proposer), and reserve the reward for the proposer into the agent. The
executor MUST place the reward as stated in Msg/SubmitProposal or the
post-actions would fail, resulting the failure of the transaction.

#### Submitting Proposals

A transaction begins with a certain proposer submitting a proposal. It is
triggered by broadcasting Msg/SubmitProposal. The message has information of
the proposer, the agent, the pre-actions and post-actions.
An agent MUST be prepared by the proposer, in prior to the submission. The
signers included in the pre-actions and post-actions MUST be either the
proposer or the agent. As mentioned above, pre-actions will be executed in this
step.

After successful submission of the proposal, the agent would be pruned from the
state.

#### Executing Proposals

A proposal would be executed by a certain executor who is interested in. It is
triggered by broadcasting Msg/Exec. The message has information of the
executor, the agent and the actions. The agent MUST be indentical with that of
the proposal. The signers included in the actions MUST be either the executor
or the agent.

After successful execution of the proposal, the proposal would be pruned from
the state.


## State

### Params

* Params: `0x00 -> ProtocolBuffer(Params)`

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/types.proto#L6-L7
```

### NextAgent

* NextAgent: `0x10 -> uint64`

### Agents

* Agents: `0x11 | creator_address | agent_address -> ProtocolBuffer(Agent)`

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/types.proto#L9-L11
```

### NextProposal

* NextProposal: `0x20 -> uint64`

### Proposals

* Proposals: `0x21 | proposer_address | proposal_id -> ProtocolBuffer(Proposal)`

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/types.proto#L13-L23
```


## Msg Service

### Msg/UpdateParams

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L25-L38
```

### Msg/CreateAgent

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L43-L49
```

### Msg/SubmitProposal

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L57-L75
```

### Msg/Exec

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L83-L99
```


## Events

### EventUpdateParams

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L7-L9
```

### EventCreateAgent

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L11-L18
```

### EventSubmitProposal

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L20-L36
```

### EventExec

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L38-L48
```


## Client

### CLI

#### Query

```bash
and query escrow --help
```

```bash
Querying commands for the escrow module

Usage:
  and query escrow [flags]
  and query escrow [command]

Available Commands:
  agent       queries an agent.
  agents      queries all the agents.
  params      queries the module params.
  proposal    queries a proposal.
  proposals   queries all the proposals.
```

##### params

```bash
and query escrow params --help
```

```bash
```

##### agent

```bash
and query escrow agent --help
```

```bash
queries an agent.

Usage:
  and query escrow agent --address [address] [flags]

Examples:
$ and query escrow agent --address cosmos1aaa...
agent:
  address: cosmos1aaa...
  creator: cosmos1...
```

##### agents

```bash
and and query escrow agents --help
```

```bash
queries all the agents.

Usage:
  and query escrow agents [flags]

Examples:
$ and query escrow agents
agents:
- address: cosmos1...
  creator: cosmos1...
- address: cosmos1...
  creator: cosmos1...
- address: cosmos1...
  creator: cosmos1...
pagination:
  total: "3"
```

##### proposal

```bash
and query escrow proposal --help
```

```bash
queries a proposal.

Usage:
  and query escrow proposal --id [id] [flags]

Examples:
$ and query escrow proposal --id 3
proposal:
  agent: cosmos1...
  id: "3"
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
```

##### proposals

```bash
and query escrow proposals --help
```

```bash
queries all the proposals.

Usage:
  and query escrow proposals [flags]

Examples:
$ and query escrow proposals
pagination:
  total: "2"
proposals:
- agent: cosmos1...
  id: "3"
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
  proposer: cosmos1...
```

#### Transactions

```bash
and tx escrow --help
```

```bash
Transactions commands for the escrow module

Usage:
  and tx escrow [flags]
  and tx escrow [command]

Available Commands:
  create-agent    creates an escrow agent for a proposal.
  exec            executes a proposal.
  submit-proposal submits a proposal.
  update-params   updates the module parameters.
```

##### create-agent

```bash
and tx escrow create-agent --help
```

```bash
creates an escrow agent for a proposal.

Usage:
  and tx escrow create-agent --from [creator] [flags]

Examples:
$ and tx escrow create-agent --from cosmos1ccc...
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
confirm transaction before signing and broadcasting [y/N]:
```

##### submit-proposal

```bash
and tx escrow submit-proposal --help
```

```bash
submits a proposal.

Note:
  agent:
    it must be created by the proposer.
  pre-actions:
    the signer of each message must be either the proposer or the agent.
  post-actions:
    the signer of each message must be either the proposer or the agent.

Usage:
  and tx escrow submit-proposal --from [proposer] --agent [agent] --pre-actions [pre-actions] --post-actions [post-actions] [flags]

Examples:
$ and tx escrow submit-proposal --from cosmos1ppp... --agent cosmos1aaa... \
    --pre-actions '{"@type": "/cosmos.nft.v1beta1.MsgSend",
                    "class_id": "cat",
                    "id": "octocat",
                    "sender": "cosmos1ppp...",
                    "receiver": "cosmos1aaa..."}' \
    --post-actions '{"@type": "/cosmos.bank.v1beta1.MsgSend",
                     "from_address": "cosmos1aaa",
                     "to_address": "cosmos1ppp...",
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
  - '@type': /andromeda.escrow.v1alpha1.MsgSubmitProposal
    agent: cosmos1aaa...
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
confirm transaction before signing and broadcasting [y/N]:
```

##### exec

```bash
and tx escrow exec --help
```

```bash
executes a proposal.

Note:
  actions:
    the signer of each message must be either the executor or the agent.

Usage:
  and tx escrow exec --proposal [proposal] --from [executor] --agent [agent] --actions [actions] [flags]

Examples:
$ and tx escrow exec --proposal 2 --from cosmos1eee... --agent cosmos1... \
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
confirm transaction before signing and broadcasting [y/N]:
```

### gRPC

```bash
grpcurl -plaintext \
  localhost:9090 list andromeda.escrow.v1alpha1.Query
```

```bash
andromeda.escrow.v1alpha1.Query.Agent
andromeda.escrow.v1alpha1.Query.Agents
andromeda.escrow.v1alpha1.Query.Params
andromeda.escrow.v1alpha1.Query.Proposal
andromeda.escrow.v1alpha1.Query.Proposals
```

#### andromeda.escrow.v1alpha1.Query.Params

```bash
grpcurl -plaintext \
  localhost:9090 describe andromeda.escrow.v1alpha1.Query.Params
```

Example:

```bash
grpcurl -plaintext \
  localhost:9090 andromeda.escrow.v1alpha1.Query.Params
```

Example Output:

```bash
{
  "params": {}
}
```

#### andromeda.escrow.v1alpha1.Query.Agent

```bash
grpcurl -plaintext \
  localhost:9090 describe andromeda.escrow.v1alpha1.Query.Agent
```

Example:

```bash
grpcurl -plaintext \
  -d '{"address": "cosmos1aaa..."}' \
  localhost:9090 andromeda.escrow.v1alpha1.Query.Agent
```

Example Output:

```bash
{
  "agent": {
    "address": "cosmos1aaa...",
    "creator": "cosmos1..."
  }
}
```

### andromeda.escrow.v1alpha1.Query.Agents

```bash
grpcurl -plaintext \
  localhost:9090 describe andromeda.escrow.v1alpha1.Query.Agents
```

Example:

```bash
grpcurl -plaintext \
  localhost:9090 andromeda.escrow.v1alpha1.Query.Agents
```

Example Output:

```bash
{
  "agents": [
    {
      "address": "cosmos1...",
      "creator": "cosmos1..."
    },
    {
      "address": "cosmos1...",
      "creator": "cosmos1..."
    },
    {
      "address": "cosmos1...",
      "creator": "cosmos1..."
    }
  ],
  "pagination": {
    "total": "3"
  }
}
```

#### andromeda.escrow.v1alpha1.Query.Proposal

```bash
grpcurl -plaintext \
  localhost:9090 describe andromeda.escrow.v1alpha1.Query.Proposal
```

Example:

```bash
grpcurl -plaintext \
  -d '{"id": "3"}' \
  localhost:9090 andromeda.escrow.v1alpha1.Query.Proposal
```

Example Output:

```bash
{
  "proposal": {
    "id": "3",
    "proposer": "cosmos1...",
    "agent": "cosmos1...",
    "preActions": [
      {
        "@type": "/cosmos.nft.v1beta1.MsgSend",
        "classId": "...",
        "id": "...",
        "sender": "cosmos1...",
        "receiver": "cosmos1..."
      }
    ],
    "postActions": [
      {
        "@type": "/cosmos.bank.v1beta1.MsgSend",
        "amount": [
          {
            "denom": "stake",
            "amount": "42"
          }
        ],
        "fromAddress": "cosmos1...",
        "toAddress": "cosmos1..."
      }
    ]
  }
}
```

### andromeda.escrow.v1alpha1.Query.Proposals

```bash
grpcurl -plaintext \
  localhost:9090 describe andromeda.escrow.v1alpha1.Query.Proposals
```

Example:

```bash
grpcurl -plaintext \
  localhost:9090 andromeda.escrow.v1alpha1.Query.Proposals
```

Example Output:

```bash
{
  "proposals": [
    {
      "id": "3",
      "proposer": "cosmos1...",
      "agent": "cosmos1...",
      "preActions": [
        {
          "@type": "/cosmos.nft.v1beta1.MsgSend",
          "classId": "...",
          "id": "...",
          "sender": "cosmos1...",
          "receiver": "cosmos1..."
        }
      ],
      "postActions": [
        {
          "@type": "/cosmos.bank.v1beta1.MsgSend",
          "amount": [
            {
              "denom": "stake",
              "amount": "42"
            }
          ],
          "fromAddress": "cosmos1...",
          "toAddress": "cosmos1..."
        }
      ]
    },
    {
      "id": "4",
      "proposer": "cosmos1...",
      "agent": "cosmos1...",
      "preActions": [
        {
          "@type": "/cosmos.nft.v1beta1.MsgSend",
          "classId": "...",
          "id": "...",
          "sender": "cosmos1...",
          "receiver": "cosmos1..."
        }
      ],
      "postActions": [
        {
          "@type": "/cosmos.bank.v1beta1.MsgSend",
          "amount": [
            {
              "denom": "stake",
              "amount": "42"
            }
          ],
          "fromAddress": "cosmos1...",
          "toAddress": "cosmos1..."
        },
        {
          "@type": "/cosmos.nft.v1beta1.MsgSend",
          "classId": "...",
          "id": "...",
          "sender": "cosmos1...",
          "receiver": "cosmos1..."
        }
      ]
    },
  ],
  "pagination": {
    "total": "2"
  }
}
```
