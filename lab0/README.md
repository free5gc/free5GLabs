# Lab 0: Network Programming with Go

## Introduction

In Lab 0, you will learn how to build network program (e.g. TCP client/server) with Go.

## Goals of this lab

- Understand network programming
- Understand how to use `go test`

## Environment setup

free5GC project adopts Go 1.21 for develoment.
Thus we recommand you use Go 1.21+ to complete all of labs.
- [go.dev>Download and install>Go installation](https://go.dev/doc/install)

## What is `0.0.0.0`?

## What is `localhost`?

## IPv4 and IPv6

## TCP and UDP

## Exercise: Implement a TCP server

- Please implement pre-defined functions `TcpListener()` and `TcpHandler()`.
- Expected behaviours:
    - Support to handle multiple connection simultaneously.
    - After connection established, TCP server always respond what it received.
    - For example, TCP server will respond "OK" if client sent "OK" to Server.
- `TcpListener()`
    - Please follow `listenerInterface` interface.
    - Used to listening on specific IP + port, it depends on what parameters passed into this function.
    - Once the listener accepts new connection, listener should delegate the connection to `TcpHandler()`.
- `TcpHandler()`
    - Should follow `handlerInterface` interface.
    - Used to handle single connection.

After implementation completed, you can use the command below for validation:
```sh
make test
```
And the expected result looks like:
```sh
go test -v -race -timeout 30s ./...
?       github.com/ianchen0119/free5GLab/lab0/ans       [no test files]
=== RUN   TestTcpFunction
2024/06/14 10:17:59 TCP is listening on 127.0.0.1:8080
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43428
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43428]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43440
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43440]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43446
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43446]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43456
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43456]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43460
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43460]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43468
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43468]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43480
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43480]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43494
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43494]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43502
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43502]
2024/06/14 10:18:04 new client accepted: 127.0.0.1:43510
2024/06/14 10:18:04 Handle Request from [127.0.0.1:43510]
--- PASS: TestTcpFunction (5.03s)
PASS
2024/06/14 10:18:04 Client [127.0.0.1:43502[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43494[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43480[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43468[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43460[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43510[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43428[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43440[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43446[] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43456[] Error: EOF
ok      github.com/ianchen0119/free5GLab/lab0   6.050s
```