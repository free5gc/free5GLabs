# Lab 3: Deploying free5GC

## Introduction

In Lab3, you will learn how to deploy free5GC with docker and set up interface network.

## Goals of this lab

- Understand listening address 
- Learn how to deploy free5GC with docker 
- Learn how to use free5GC

## Preparation

* Install GTP5G: refer [GTP5G github](https://github.com/free5gc/gtp5g) to install

* Install Docker Engine and Docker Compose: refer [Docker Website](https://docs.docker.com/engine/install/ubuntu/) to install 

## Listening Address

Listening Address the IP address and port used by a server to listen for connections from clients. These addresses and ports are used to accept requests from clients. ex: 192.168.100.101:12345

In free5GC, each NF (Network Function) has its own listening addresses used to receive and process requests sent by other NFs.

## N2 & N3 & N4 interface
![architecture](./images/architecture.png)
* N2: handles control signaling between the `gNB` and `AMF`
* N3: manages user data transmission between the `gNB` and `UPF`
* N4: controls user plane configuration and session management between the `SMF` and `UPF`

These three interfaces are the most important interfaces in the 5G system. Therefore, this lab will teach you how to configure the addresses for these interfaces.

### N6 NAT
Interface to data network.

UPF will performs NAT on packets output through the N6 interface. Rules are set up in `upf-iptables.sh`.

### N9
Interface for two UPFs communication.

### Docker bridge network
Docker Bridge Network is one of the most commonly used network modes in Docker. It allows Docker containers to communicate with each other through a virtual bridge. This mode is mainly used for communication between containers and between containers and the host machine.

After installing Docker Engine, there will be a default bridge network named `docker0`, with a subnet of `172.17.0.0/16`. You can use `brctl show` to see the bridge information.

```sh
~$ brctl show
bridge name     bridge id               STP enabled     interfaces
docker0         8000.02426fa2174a       no
```

If we run a container, a pair of veth interfaces will be created, with one placed on `docker0` and the other inside the container. This allows the container to communicate with other containers on the bridge or with the host.

```sh
~$ docker run -dit --name alpine1 alpine ash
7ef26c2ccd704b35ec168baa498239a773e7fd8ed0cd8db9ebb24638c14e46b2

// docker0 bridge has one veth to container
~$ brctl show
bridge name     bridge id               STP enabled     interfaces
docker0         8000.02426fa2174a       no              veth43f5e85
```

Since the Docker bridge is implemented based on the Linux bridge, we can use `brctl`, a Linux utility or docker cli for managing bridges. Or use docker cli to manage it. [brctl document](https://man7.org/linux/man-pages/man8/brctl.8.html)

## Exercise: Configure N2 & N3 & N4 interface in Docker Compose
In this exercise, we will use docker bridge network to set up these three interfaces.

In bottom of exercise/deploy_exercise.yaml, you can find network setting.
```yaml
networks:
  privnet:
    ipam:
      driver: default
      config:
        - subnet: 10.100.200.0/24
    driver_opts:
      com.docker.network.bridge.name: br-free5gc
```
For example, privnet is the bridge network for NFs internal communicaion. ex: Nnrf, Nudm...

Each NF, except for the `UPF`, is assigned an IP address within a `privnet` for internal communication. And assign it an alias for ease of use.
```yaml
networks:
    privnet:
    aliases:
        - udr.free5gc.org
```
We create three bridge networks named `n1net`, `n2net`, `n3net` and `n6net`. Our goal is respectively assigned network to `AMF`, `SMF`, `UPF` and `UERANSIM`.

```yaml
n3net:
  ipam:
    driver: default
    config:
      - subnet: 10.100.3.0/24
  driver_opts:
    com.docker.network.bridge.name: br-n3

n4net:
  ipam:
    driver: default
    config:
      - subnet: 10.100.4.0/24
  driver_opts:
    com.docker.network.bridge.name: br-n4

n6net:
  ipam:
    driver: default
    config:
      - subnet: 10.100.6.0/24
  driver_opts:
    com.docker.network.bridge.name: br-n6
    com.docker.network.container_iface_prefix: dn
```

Here is the sample of assign `n3net` to upf :
```yaml
networks:
  n3net:
    aliases:
      - upf.n3.org
    ipv4_address: 10.100.3.100
```

After assigning the networks, it is necessary to update the IP addresses in the corresponding NF configurations.

In `gnbcfg.yaml`, `amfcfg.yaml`, `smfcfg.yaml` and `upfcfg.yaml`, you can find the settings to configure IP addresses as follows.
```yaml
  pfcp: # the IP address of N4 interface on this SMF (PFCP)
    # addr config is deprecated in smf config v1.0.3, please use the following config
    nodeID: update here # the Node ID of this SMF
    listenAddr: update here # the IP/FQDN of N4 interface on this SMF (PFCP)
    externalAddr: update here # the IP/FQDN of N4 interface on this SMF (PFCP)
```
Please replace `update here` with the configured N2, N3, and N4 addresses.

Tips: 
In `smfcfg.yaml`, you will configure the `UPF` N3 interface address because it is required for setting up sessions during SM context creation. If you only use an alias when configuring this address, it may cause DNS resolution issues. Therefore, in `deploy_exercise.yaml`, you should set a static IP address for the `UPF` N3 network and use it here.

After configuring, clone [free5gc-compose](https://github.com/free5gc/free5gc-compose). Then move `free5GLab/lab3/exercise/deploy_exercise.yaml` to `free5gc-compose/` and copy the contents of the files from the `free5GLab/lab3/exercise/config` directory to `free5gc-compose/config`.

You can use these commands to start or stop docker compose.
```sh
ce ~/free5gc-compose

// start
docker compose -f deploy_exercise.yaml up

// remove
docker compose -f deploy_exercise.yaml down
```

Refer [Create Subscriber via Webconsole](https://free5gc.org/guide/Webconsole/Create-Subscriber-via-webconsole/#5-add-new-subscriber) to create subscriber 

And, attach to ueransim container and run ue
```sh
// attach to container
docker exec -it ueransim bash

// run ue
./nr-ue -c config/uecfg.yaml
```
And the expected result looks like:
```sh
2024-07-05 12:29:35.111] [nas] [info] UE switches to state [MM-DEREGISTERED/PLMN-SEARCH]
[2024-07-05 12:29:35.111] [rrc] [debug] New signal detected for cell[1], total [1] cells in coverage
[2024-07-05 12:29:36.572] [nas] [info] Selected plmn[208/93]
[2024-07-05 12:29:36.573] [rrc] [info] Selected cell plmn[208/93] tac[1] category[SUITABLE]
[2024-07-05 12:29:36.573] [nas] [info] UE switches to state [MM-DEREGISTERED/PS]
[2024-07-05 12:29:36.573] [nas] [info] UE switches to state [MM-DEREGISTERED/NORMAL-SERVICE]
[2024-07-05 12:29:36.574] [nas] [debug] Initial registration required due to [MM-DEREG-NORMAL-SERVICE]
[2024-07-05 12:29:36.575] [nas] [debug] UAC access attempt is allowed for identity[0], category[MO_sig]
[2024-07-05 12:29:36.576] [nas] [debug] Sending Initial Registration
[2024-07-05 12:29:36.576] [rrc] [debug] Sending RRC Setup Request
[2024-07-05 12:29:36.579] [rrc] [info] RRC connection established
[2024-07-05 12:29:36.592] [rrc] [info] UE switches to state [RRC-CONNECTED]
[2024-07-05 12:29:36.593] [nas] [info] UE switches to state [MM-REGISTER-INITIATED]
[2024-07-05 12:29:36.596] [nas] [info] UE switches to state [CM-CONNECTED]
[2024-07-05 12:29:36.691] [nas] [debug] Authentication Request received
[2024-07-05 12:29:36.692] [nas] [debug] Received SQN [000000000027]
[2024-07-05 12:29:36.692] [nas] [debug] SQN-MS [000000000000]
[2024-07-05 12:29:36.735] [nas] [debug] Security Mode Command received
[2024-07-05 12:29:36.735] [nas] [debug] Selected integrity[2] ciphering[0]
[2024-07-05 12:29:36.961] [nas] [debug] Registration accept received
[2024-07-05 12:29:36.962] [nas] [info] UE switches to state [MM-REGISTERED/NORMAL-SERVICE]
[2024-07-05 12:29:36.962] [nas] [debug] Sending Registration Complete
[2024-07-05 12:29:36.962] [nas] [info] Initial Registration is successful
[2024-07-05 12:29:36.963] [nas] [debug] Sending PDU Session Establishment Request
[2024-07-05 12:29:36.964] [nas] [debug] UAC access attempt is allowed for identity[0], category[MO_sig]
[2024-07-05 12:29:36.964] [nas] [debug] Sending PDU Session Establishment Request
[2024-07-05 12:29:36.964] [nas] [debug] UAC access attempt is allowed for identity[0], category[MO_sig]
[2024-07-05 12:29:37.173] [nas] [debug] Configuration Update Command received
[2024-07-05 12:29:37.645] [nas] [debug] PDU Session Establishment Accept received
[2024-07-05 12:29:37.645] [nas] [info] PDU Session establishment is successful PSI[1]
[2024-07-05 12:29:37.688] [nas] [debug] PDU Session Establishment Accept received
[2024-07-05 12:29:37.691] [nas] [info] PDU Session establishment is successful PSI[2]
```
And use `ping` to test it can reach date network
```sh
ping -I uesimtun0 8.8.8.8
```

If you encounter any issues during the exercise, you can refer to the `free5GLab/lab3/ans` folder.

## Reference
* [3GPP TS 23.501](https://portal.3gpp.org/desktopmodules/Specifications/SpecificationDetails.aspx?specificationId=3144)