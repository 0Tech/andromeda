# `x/test`

## Abstract

Some modules have messages which may use other messages. Typical examples would
be:

* [`x/escrow`](../escrow/README.md)
* [`x/authz`](https://github.com/cosmos/cosmos-sdk/x/authz/README.md)
* [`x/gov`](https://github.com/cosmos/cosmos-sdk/x/gov/README.md)
* [`x/group`](https://github.com/cosmos/cosmos-sdk/x/group/README.md)

Testing those kind of messages requires messages with ownership semantics,
however, the corresponding module often lacks such kind of messages. While
implementing (or mocking) such messages by each module would be a good way to
avoid compatibility issues on the test, it would be also resonable to prepare a
simple, stable module only for the tests. That's the rationale of this module.


## Contents

* [Concepts](#concepts)
* [State](#state)
* [Msg Service](#msg-service)
    * [Msg/Create](#msgcreate)
    * [Msg/Send](#msgsend)
* [Events](#events)
    * [EventCreate](#eventcreate)
    * [EventSend](#eventsend)


## Concepts

### Asset

Assets would be identified by its name in string. You cannot stack same assets
on the same account. Anyone can create assets by their own.


## State

### Assets

One can change the prefix through the argument of the keeper.

* Assets: `0xff | owner_address | asset_name -> ProtocolBuffer(Asset)`

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/test/proto/andromeda/test/v1alpha1/types.proto#L3-L5
```


## Msg Service

### Msg/Create

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/test/proto/andromeda/test/v1alpha1/tx.proto#L18-L27
```

### Msg/Send

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/test/proto/andromeda/test/v1alpha1/tx.proto#L32-L44
```


## Events

### EventCreate

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/test/proto/andromeda/test/v1alpha1/event.proto#L6-L13
```

### EventSend

```protobuf reference
https://github.com/0tech/andromeda/blob/main/x/test/proto/andromeda/test/v1alpha1/event.proto#L15-L25
```
