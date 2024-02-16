# `x/escrow`

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=0tech_andromeda_x-escrow&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=0tech_andromeda_x-escrow)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=0tech_andromeda_x-escrow&metric=coverage)](https://sonarcloud.io/summary/new_code?id=0tech_andromeda_x-escrow)

## Abstract

This module provides means to execute messages whose signers come from
different parties. Typical example would be a trade, by which two individual
account can exchange assets. Without this module, they have to use a tx which
include the messages that requires signatures from both of them. Or, they may
create a contract into x/wasm. The former requires a tool to safely exchange
the half-signed transaction (which x/group tries to avoid), while the latter
involves cumbersome preparation for certain tailer made contract code. The
trade itself would be quite simple and obvious, hence this module tries to
reduce such efforts on common use cases.

Taking it a step further, this module allows us to build an ecosystem for
healthy, speedy and convenient transactions. One can execute multiple proposals
together, which lowers the hurdles to becoming a retailer or broker. Hence, if
it's a quite reasonable deal, some user will trigger your proposal (together
with other proposals) for their profit, meaning you don't even need to find
matching proposals by yourself. Just submit a proposal and forget about it.


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
* [Examples](#examples)


## Concepts

### Agent

In abstract manner, the purpose of agents is to provide a place holder for
proposals. In real life, an agent is vital component for safe transactions,
which arranges the procedure as a neutral actor.

Anyone can create agents as many as they want.

### Proposal

There are two parties involved in a transaction. The first one, proposer,
suggests a transaction by broadcasting `Msg/SubmitProposal`. The second one,
executor, accepts the transaction, by broadcasting `Msg/Exec`.

A transaction consists of pre-actions, actions and post-actions. All three
actions will be executed in sequence, and all the messages MUST be successful
to complete the transaction.

Pre-actions are steps for transferring the control on certain assets from the
proposer to the agent, so determined by the proposer on `Msg/SubmitProposal`.
To ensure the successful transfer at `Msg/Exec`, pre-actions are executed at
the proposal submission.

Post-actions are steps for transferring the control on certain assets from the
agent to the proposer. Post-actions ensure the proposer gets the proper reward,
so it is determined by the proposer on `Msg/SubmitProposal`.

Actions are steps for transferring the control on certain assets from the
executor to agent and vice versa, so determined by the executor on `Msg/Exec`.
The messages included in actions would take assets on the agent (reserved by
the proposer), and reserve the reward for the proposer into the agent. The
executor must place the reward as stated in `Msg/SubmitProposal` or the
post-actions would fail, resulting the failure of the transaction.

#### Submitting Proposals

A transaction begins with a certain proposer submitting a proposal. It is
triggered by broadcasting `Msg/SubmitProposal`. The message has information of
the proposer, the agent, the pre-actions and post-actions.
An agent MUST be prepared by the proposer, in prior to the submission. The
signer of each message included in the pre-actions and post-actions MUST be
either the proposer or the agent. As mentioned above, pre-actions will be
executed in this step.

After successful submission of the proposal, the agent would be pruned from the
state.

#### Executing Proposals

A proposal would be executed by a certain executor who is interested in. It is
triggered by broadcasting `Msg/Exec`. The message has information of the
executor, the agents and the actions. Each agent MUST have the corresponding
proposal. The signer of each message included in the actions MUST be either the
executor or one of the agents. The execution order of the post-actions is same
as the inclusion order of the agents.

After successful execution of the proposal, the proposal would be pruned from
the state.


## State

### Params

* Params: `0x00 -> ProtocolBuffer(Params)`

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/types.proto#L6-L10

### NextAgent

* NextAgent: `0x10 -> uint64`

### Agents

* Agents: `0x11 | agent_address -> ProtocolBuffer(Agent)`
* AgentsByCreator: `0x12 | creator_address | agent_address`

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/types.proto#L12-L16

### Proposals

* Proposals: `0x20 | agent_address -> ProtocolBuffer(Proposal)`
* ProposalsByProposer: `0x21 | proposer_address | agent_address`

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/types.proto#L18-L31


## Msg Service

### Msg/UpdateParams

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L25-L34

### Msg/CreateAgent

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L39-L45

### Msg/SubmitProposal

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L53-L74

### Msg/Exec

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/tx.proto#L79-L93


## Events

### EventUpdateParams

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L7-L14

### EventCreateAgent

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L16-L23

### EventSubmitProposal

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L25-L41

### EventExec

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/escrow/proto/andromeda/escrow/v1alpha1/event.proto#L43-L53


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
  agent                 queries an agent.
  agents                queries all the agents.
  agents-by-creator     queries all the agents by its creator.
  params                queries the module parameters.
  proposal              queries a proposal.
  proposals             queries all the proposals.
  proposals-by-proposer queries all the proposals by its proposer.
```

##### params

```bash
and query escrow params --help
```

```bash
queries the module parameters.

Usage:
  and query escrow params [flags]

Examples:
$ and query escrow params
max_metadata_length: "42"
```

##### agent

```bash
and query escrow agent --help
```

```bash
queries an agent.

Usage:
  and query escrow agent --agent [agent] [flags]

Examples:
$ and query escrow agent --agent cosmos1aaa...
agent:
  address: cosmos1aaa...
  creator: cosmos1...
```

##### agents-by-creator

```bash
and query escrow agents-by-creator --help
```

```bash
queries all the agents by its creator.

Usage:
  and query escrow agents-by-creator --creator [creator] [flags]

Examples:
$ and query escrow agents-by-creator --creator cosmos1ccc...
agents:
- address: cosmos1...
  creator: cosmos1ccc...
- address: cosmos1...
  creator: cosmos1ccc...
pagination:
  total: "2"
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
  and query escrow proposal --agent [agent] [flags]

Examples:
$ and query escrow proposal --agent cosmos1aaa...
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
  proposer: cosmos1...
```

##### proposals-by-proposer

```bash
and query escrow proposals-by-proposer --help
```

```bash
queries all the proposals by its proposer.

Usage:
  and query escrow proposals-by-proposer --proposer [proposer] [flags]

Examples:
$ and query escrow proposals-by-proposer --proposer cosmos1ppp...
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
  proposer: cosmos1ppp...
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

##### update-params

```bash
and tx escrow update-params --help
```

```bash
updates the module parameters.

Usage:
  and tx escrow update-params --from [authority] --max-metadata-length [max-metadata-length] [flags]

Examples:
$ and tx escrow update-params --from cosmos1aaa... --max-metadata-length 42
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
confirm transaction before signing and broadcasting [y/N]:
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
  and tx escrow submit-proposal --from [proposer] --agent [agent] --pre-actions [pre-actions] --post-actions [post-actions] --metadata [metadata] [flags]

Examples:
$ and tx escrow submit-proposal --from cosmos1ppp... --agent cosmos1aaa... \
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
    the signer of each message must be either the executor or one of the agents.

Usage:
  and tx escrow exec --from [executor] --agents [agents] --actions [actions] [flags]

Examples:
$ and tx escrow exec --from cosmos1eee... --agents cosmos1aaa... \
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
andromeda.escrow.v1alpha1.Query.AgentsByCreator
andromeda.escrow.v1alpha1.Query.Params
andromeda.escrow.v1alpha1.Query.Proposal
andromeda.escrow.v1alpha1.Query.Proposals
andromeda.escrow.v1alpha1.Query.ProposalsByProposer
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
  "maxMetadataLength": "42"
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
  -d '{"agent": "cosmos1aaa..."}' \
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

#### andromeda.escrow.v1alpha1.Query.AgentsByCreator

```bash
grpcurl -plaintext \
  localhost:9090 describe andromeda.escrow.v1alpha1.Query.AgentsByCreator
```

Example:

```bash
grpcurl -plaintext \
  -d '{"creator": "cosmos1ccc..."}' \
  localhost:9090 andromeda.escrow.v1alpha1.Query.AgentsByCreator
```

Example Output:

```bash
{
  "agents": [
    {
      "address": "cosmos1...",
      "creator": "cosmos1ccc..."
    },
    {
      "address": "cosmos1...",
      "creator": "cosmos1ccc..."
    }
  ],
  "pagination": {
    "total": "2"
  }
}
```

#### andromeda.escrow.v1alpha1.Query.Agents

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
  -d '{"agent": "cosmos1aaa..."}' \
  localhost:9090 andromeda.escrow.v1alpha1.Query.Proposal
```

Example Output:

```bash
{
  "proposal": {
    "agent": "cosmos1aaa...",
    "proposer": "cosmos1...",
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
    ],
    "metadata": "very good deal"
  }
}
```

#### andromeda.escrow.v1alpha1.Query.ProposalsByProposer

```bash
grpcurl -plaintext \
  localhost:9090 describe andromeda.escrow.v1alpha1.Query.ProposalsByProposer
```

Example:

```bash
grpcurl -plaintext \
  -d '{"proposer": "cosmos1ppp..."}' \
  localhost:9090 andromeda.escrow.v1alpha1.Query.ProposalsByProposer
```

Example Output:

```bash
{
  "proposals": [
    {
      "agent": "cosmos1...",
      "proposer": "cosmos1ppp...",
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
      ],
      "metadata": "limited time offer for you"
    },
  ],
  "pagination": {
    "total": "1"
  }
}
```

#### andromeda.escrow.v1alpha1.Query.Proposals

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
      "agent": "cosmos1...",
      "proposer": "cosmos1...",
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
      ],
      "metadata": "very good deal"
    },
    {
      "agent": "cosmos1...",
      "proposer": "cosmos1...",
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
      ],
      "metadata": "limited time offer for you"
    },
  ],
  "pagination": {
    "total": "2"
  }
}
```


## Examples

### Sale of an NFT for coins

It would be the most trivial usage of this module. In this example, a proposer
`cosmos1ppp...` sells an x/nft token `cat:leopardcat` for x/bank coins
`42stake`:

```yaml
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
```

The proposal would be executed by `cosmos1eee...`, who is willing to pay
`42stake` for `cat:leopardcat`:

```yaml
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
```

### Sale of coins for an NFT

It may sound strange, but one can sell coins for a certain NFT. In this
example, a proposer `cosmos1ppp...` sells x/bank coins `42stake` for an x/nft
token `cat:leopardcat`:

```yaml
agent: cosmos1aaa...
metadata: sell 42stake for leopardcat
post_actions:
- '@type': /cosmos.nft.v1beta1.MsgSend
  value:
    class_id: cat
    id: leopardcat
    receiver: cosmos1ppp...
    sender: cosmos1aaa...
pre_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "42"
    denom: stake
  from_address: cosmos1ppp...
  to_address: cosmos1aaa...
proposer: cosmos1ppp...
```

The proposal would be executed by `cosmos1eee...`, who is willing to pay
`cat:leopardcat` for `42stake`:

```yaml
actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "42"
    denom: stake
  from_address: cosmos1aaa...
  to_address: cosmos1eee...
- '@type': /cosmos.nft.v1beta1.MsgSend
  class_id: cat
  id: leopardcat
  receiver: cosmos1aaa...
  sender: cosmos1eee...
agents:
- cosmos1aaa...
executor: cosmos1eee...
```

### Sale of an NFT to a specific account

One may want to offer deals to a specific account. In this example, a proposer
`cosmos1ppp...` sells an x/nft token `cat:leopardcat` for x/bank coins
`42stake` to `cosmos1eee...`:

```yaml
agent: cosmos1aaa...
metadata: sell leopardcat for 42stake
post_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "42"
    denom: stake
  from_address: cosmos1aaa...
  to_address: cosmos1ppp...
- '@type': /cosmos.nft.v1beta1.MsgSend
  value:
    class_id: cat
    id: leopardcat
    receiver: cosmos1eee...
    sender: cosmos1aaa...
pre_actions:
- '@type': /cosmos.nft.v1beta1.MsgSend
  value:
    class_id: cat
    id: leopardcat
    receiver: cosmos1aaa...
    sender: cosmos1ppp...
proposer: cosmos1ppp...
```

Technically, the proposal could be executed by anyone. However, it would not be
a wise choice executing this proposal for the others, because the messages
included in the pre-actions, actions and post-actions are all-or-none. And
this proposal has a specific receipient of the NFT, so the proposer can assure
the designated recipient will receive the NFT eventually.

The recipient, `cosmos1eee...`, who has an interest on this deal, will execute
the proposal:

```yaml
actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "42"
    denom: stake
  from_address: cosmos1eee...
  to_address: cosmos1aaa...
agents:
- cosmos1aaa...
executor: cosmos1eee...
```

### Broker

One can broker multiple proposals. In this example, one proposer
`cosmos1ppp...` sells x/bank coins `1notscam` for x/bank coins `42stake`:

```yaml
agent: cosmos1aaa...
metadata: sell 1notscam for 42stake
post_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "42"
    denom: stake
  from_address: cosmos1aaa...
  to_address: cosmos1ppp...
pre_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "1"
    denom: notscam
  from_address: cosmos1ppp...
  to_address: cosmos1aaa...
proposer: cosmos1ppp...
```

The other proposer `cosmos1qqq...` sells x/bank coins `4242stake` for x/bank
coins `1notscam`:

```yaml
agent: cosmos1bbb...
metadata: sell 1notscam for 4242stake
post_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "1"
    denom: notscam
  from_address: cosmos1bbb...
  to_address: cosmos1qqq...
pre_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "4242"
    denom: stake
  from_address: cosmos1qqq...
  to_address: cosmos1bbb...
proposer: cosmos1qqq...
```

In real life, this is a typical type of scam. If you execute the first proposal
to buy `1notscam`, hoping to sell it to `cosmos1qqq...` to earn the difference
`4200stake`, they will execute the second proposal immediately. But on a cosmos
chain, it won't work, because you can send multiple messages in a tx, ensuring
they would be executed in all-or-none manner.

Or, you can broker the proposals, meaning executing multiple proposals in one
`Msg/Exec`:

```yaml
actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "4200"
    denom: stake
  from_address: cosmos1bbb...
  to_address: cosmos1eee...
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "42"
    denom: stake
  from_address: cosmos1bbb...
  to_address: cosmos1aaa...
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "1"
    denom: notscam
  from_address: cosmos1aaa...
  to_address: cosmos1bbb...
agents:
- cosmos1aaa...
- cosmos1bbb...
executor: cosmos1eee...
```

In this way, you don't even need to prepare `42stake` to trigger the first
proposal.

### Retail

Wholesalers often sell large amount of assets at a reduced price. One can be a
retailer when certain amount of relevant proposals are accumulated on the
chain. In this example, a wholesaler `cosmos1wholesaler...` sells x/bank coins
`10000egg` for x/bank coins `10000000stake`:

```yaml
agent: cosmos1wholesaleragent...
metadata: sell 10000egg for 10000000stake
post_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "10000000"
    denom: stake
  from_address: cosmos1wholesaleragent...
  to_address: cosmos1wholesaler...
pre_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "10000"
    denom: egg
  from_address: cosmos1wholesaler...
  to_address: cosmos1wholesaleragent...
proposer: cosmos1wholesaler...
```

For simplicity, suppose there were 100 proposals selling x/bank coins
`150000stake` for x/bank coins `100egg`:

```yaml
agent: cosmos1consumerzeroagent...
metadata: sell 150000stake for 100egg
post_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "100"
    denom: egg
  from_address: cosmos1consumerzeroagent...
  to_address: cosmos1consumerzero...
pre_actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "150000"
    denom: stake
  from_address: cosmos1consumerzero...
  to_address: cosmos1consumerzeroagent...
proposer: cosmos1consumerzero...
```

In this case, literally anyone can be a retailer, executing the proposals:

```yaml
actions:
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "50000"
    denom: stake
  from_address: cosmos1consumerzeroagent...
  to_address: cosmos1retailer...
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "100000"
    denom: stake
  from_address: cosmos1consumerzeroagent...
  to_address: cosmos1wholesaleragent...
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "100"
    denom: egg
  from_address: cosmos1wholesaleragent...
  to_address: cosmos1consumerzeroagent...
...
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "50000"
    denom: stake
  from_address: cosmos1consumerninetynineagent...
  to_address: cosmos1retailer...
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "100000"
    denom: stake
  from_address: cosmos1consumerninetynineagent...
  to_address: cosmos1wholesaleragent...
- '@type': /cosmos.bank.v1beta1.MsgSend
  amount:
  - amount: "100"
    denom: egg
  from_address: cosmos1wholesaleragent...
  to_address: cosmos1consumerninetynineagent...
agents:
- cosmos1wholesaleragent...
- cosmos1consumerzeroagent...
...
- cosmos1consumerninetynineagent...
executor: cosmos1retailer...
```

Obviously, the actions in this example are not optimal, but it's sufficient for
the demonstration.
