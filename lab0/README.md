# Lab 0: Network Programming with Go

## Introduction

In Lab 0, you will learn how to build network program (e.g. TCP client/server) with Go.

## Goals of this lab

- Understand network programming.
- Understand how to use `go test`.

## Environment setup

The free5GC project adopts Go 1.21 for development.
Thus, we recommend you use Go 1.21+ to complete all of the labs.

- [go.dev > Download and install > Go installation](https://go.dev/doc/install)

## What is `0.0.0.0`?

`0.0.0.0` is a special IP address often used to indicate that the server should listen on all available network interfaces. This is useful when you want to listen on all available network interfaces, such as when you want to listen on both the loopback interface and the network interface.

## What is `localhost`?

`127.0.0.1` is the loopback IP address, also known as `localhost`. It is used to establish an IP connection to the same machine or computer being used by the end-user. This is useful when you want to test network applications on the same machine without having to connect to a remote server.

## IPv4

IPv4 is the fourth version of the Internet Protocol (IP). It is one of the core protocols of standards-based internetworking methods in the Internet and other packet-switched networks. IPv4 was the first version deployed for production in the ARPANET in 1983. It still routes most Internet traffic today, despite the ongoing deployment of a successor protocol, IPv6.

## IPv6

IPv6 is the most recent version of the Internet Protocol (IP), the communications protocol that provides an identification and location system for computers on networks and routes traffic across the Internet. IPv6 was developed by the Internet Engineering Task Force (IETF) to address the long-anticipated problem of IPv4 address exhaustion. IPv6 is intended to replace IPv4.

## TCP and UDP

TCP (Transmission Control Protocol) and UDP (User Datagram Protocol) are two of the core protocols of the Internet Protocol (IP) suite. They are used to transport data between applications over the Internet. TCP is a connection-oriented protocol that provides reliable, ordered, and error-checked delivery of data between applications. UDP is a connectionless protocol that provides unreliable, unordered, and error-checked delivery of data between applications.

### Congestion Control in TCP

Congestion control is a mechanism that prevents network congestion by regulating the rate at which data is transmitted over the network. It ensures that the network is not overloaded with data.

Congestion control in TCP is implemented using a variety of algorithms, such as slow start, congestion avoidance, fast retransmit, and fast recovery. These algorithms work together to ensure that data is transmitted at an optimal rate and that network congestion is avoided.

Moreover, BBR (Bottleneck Bandwidth and Round-trip propagation time) is a congestion control algorithm developed by Google that is designed to improve network performance by optimizing the rate at which data is transmitted over the network. BBR is designed to be more efficient than traditional TCP congestion control algorithms, such as Reno and Cubic, by taking into account the bandwidth and round-trip propagation time of the network.

## Exercise: Implement a TCP server

- Please implement the pre-defined functions `TcpListener()` and `TcpHandler()`.
- Expected behaviors:
    - Support handling multiple connections simultaneously.
    - After a connection is established, the TCP server should always respond with what it received.
    - For example, the TCP server will respond with "OK" if the client sends "OK" to the server.
- `TcpListener()`
    - Please follow the `listenerInterface` interface.
    - Used for listening on a specific IP + port, depending on the parameters passed into this function.
    - Once the listener accepts a new connection, it should delegate the connection to `TcpHandler()`.
- `TcpHandler()`
    - Should follow the `handlerInterface` interface.
    - Used to handle a single connection.

After the implementation is completed, you can use the command below for validation:

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
2024/06/14 10:18:04 Client [127.0.0.1:43502] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43494] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43480] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43468] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43460] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43510] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43428] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43440] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43446] Error: EOF
2024/06/14 10:18:04 Client [127.0.0.1:43456] Error: EOF
ok      github.com/ianchen0119/free5GLab/lab0   6.050s
```
