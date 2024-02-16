# `x/test`

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=0tech_andromeda_x-test&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=0tech_andromeda_x-test)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=0tech_andromeda_x-test&metric=coverage)](https://sonarcloud.io/summary/new_code?id=0tech_andromeda_x-test)

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

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/test/proto/andromeda/test/v1alpha1/types.proto#L4-L5

## Msg Service

### Msg/Create

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/test/proto/andromeda/test/v1alpha1/tx.proto#L18-L28

### Msg/Send

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/test/proto/andromeda/test/v1alpha1/tx.proto#L33-L46


## Events

### EventCreate

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/test/proto/andromeda/test/v1alpha1/event.proto#L6-L13

### EventSend

https://github.com/0Tech/andromeda/blob/f405ccd9e13c31233f4d34d46b500a05eb8ef8e7/x/test/proto/andromeda/test/v1alpha1/event.proto#L15-L25
