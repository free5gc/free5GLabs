# Lab 3: Deploying free5GC

## Introduction

In Lab3, you will learn how to deploy free5GC with docker and set up interface network.

## Goals of this lab

- Understand listening address 
- Learn how to deploy free5GC with docker 
- Learn how to use free5GC

## Preparation

### Install GTP5G

refer [GTP5G github](https://github.com/free5gc/gtp5g) to install

### Install Docker Engine and Docker Compose 
refer [Docker Website](https://docs.docker.com/engine/install/ubuntu/) to instll 

## Listening Address

Listening Address the IP address and port used by a server to listen for connections from clients. These addresses and ports are used to accept requests from clients. ex: 192.168.100.101:12345

In free5GC, each NF (Network Function) has its own listening address used to receive and process requests sent by other NFs.