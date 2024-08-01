# Lab 5: 5G protocol stack observation (tcpdump)

## Introduction

In Lab5, you will learn how to use the Linux command `tcpdump` to capture packets passing through a network interface.

## Goals of this lab

- Learn how to use tcpdump capture packets
- Learn how to use wireshark observing captured packets on N2, N3 and N4

## Preparation

- Install [Wireshark](https://www.wireshark.org/download.html) to install for Windows and Mac, or using `sudo apt install wireshark` for ubuntu
- Install free5GC: refer [free5GC Install](https://free5gc.org/guide/3-install-free5gc/) 
- Install Packetrusher as our UE/gNB simulator: refer [PacketRusher Install](https://free5gc.org/blog/20240110/20240110/) to install

## What is Promiscuous Mode ?

Network interface cards (NICs) have two configuration modes: **Normal Mode** and **Promiscuous Mode**.  

In the standard state, NICs operate in Normal Mode, where the network card only accepts data from the network port if the destination address is specifically directed to it. However, when analyzing network traffic, identifying network packets, and troubleshooting network issues, Promiscuous Mode is often enabled. In this mode, the NIC will receive and process all packets passing through the network interface, regardless of their destination.

## What is Packet Caputure (PCAP) File ?

PCAP file is a file format used to capture and store network traffic data, the data including details like source and destination IP addresses, ports, protocols, and the packet's payload.
PCAP file can be opened and analyzed using tools like Wireshark and tcpdump respectively. These tools allow users to inspect the network traffic in detail, apply filters, and decode protocols.

> Common Uses for PCAP file :  
> Network Troubleshooting, Security Analysis, Network Protoco Development

## How to Use Tcpdump ?

### Basic usage

> Tips : tcpdump requires root privileges.

#### Capture all packet on {interface}
```
$ tcpdump -i {interface name}
```

#### Write out capture result to PCAP file
```
$ tcpdump -i {interface name} -w {XXX}.pcap
```

#### Capture packets from a specific source IP or port
```
$ tcpdump host {host ip}
$ tcpdump port {port number}
$ tcpdump host {host ip} and port {port number}
```

#### Capture the first n packets
```
$ tcpdump -c {n}
```

#### Capture packets from a specific protocol
```
$ tcpdump tcp
$ tcpdump udp
$ tcpdump icmp
```

### Simple Example
![Schematic1](./images/Schematic.png)
#### Assuming we have a network configuration as shown in the diagram, and we want to capture the first 10 **ICMP packets** from the **Server** on the **Client's eth0** interface, and save the result to icmp.pcap file, we can use the following command.

```
$ tcpdump icmp host 192.168.0.3 -i eth0 -c 10 -w icmp.pcap
```

## Tcpdump Caputure Packet Passing Through N2 & N3

#### In this part, we will capture the packets passing through the N2 and N3 paths, and in the next part, we will use Wireshark to observe which messages are transmitted using NGAP packets when the UE connects and disconnects to the core network.

- N2: handles control signaling between the `gNB` and `AMF`
- N3: manages user data transmission between the `gNB` and `UPF`
- NGAP: Next Generation Application Protocol, Used for signaling interactions between base stations (gNB) and the 5G core network (5G Core) in 5G networks. It operates over the **N2 interface** within the 5G network architecture, primarily responsible for control plane message exchange.

### <u>Step 1 : </u>
#### Check your AMF & UPF IP address (On free5GC)
> For my example, it's 192.168.56.102 

![Schematic2](./images/ipa.png)

### <u>Step 2 : </u>
#### Check your UE & gNB IP address (On PacketRusher)
> For my example, it's 192.168.56.103

![Schematic3](./images/ipa2.png)

### <u>Step 3 : </u>
#### Run free5GC (On free5GC)
```
$ cd free5gc
$ ./run.sh
```

### <u>Step 4 : </u>
#### Start capturing packets and save to N2N3.pcap (On free5GC or PacketRusher)
```
$ tcpdump -i enp0s8 -w N2N3.pcap
```

### <u>Step 5 : </u>
#### Start the UE connection (On PacketRusher)
```
$ cd PacketRusher
$ ./packetrusher ue
```

### <u>Step 6 : </u>
#### Since the UE needs to go through the N3 interface to communicate with the Data Network, we can send ICMP packets externally to facilitate observation later.
On PacketRusher : 

```
$ ip vrf exec vrf0000000003 ping 8.8.8.8
```


### <u>Step 7 : </u>
#### Shutdown UE.

### <u>Step 8 : </u>
#### Stop tcpdump, and you can get the PCAP file.

## Wireshark Packet filter Skills

In this part, we will use Wireshark to observe which messages are transmitted and the protocol stack using NGAP packets when the UE connects and disconnects to the core network.  
In Wireshark, you can enter conditions at the top to filter out the specific packets you need from all the captured packets.

#### Common filtering conditions includes : 
- Specific IP : ```ip.addr == {specific IP}```
- Source IP : ```ip.src == {src IP}```
- Source Port : ```tcp.srcport == {port number}```
- Destination IP : ```ip.dst == {dst IP}```
- Destination Port : ```tcp.dstport == {port number}```
- MAC Address : ```eth.addr == {MAC Address}```
- Protocol types : ```ngap``` or ```icmp``` or ```tcp``` or ```sctp``` or others 

### <u>Step 1 : </u>
#### Open N2N3.pcap on Wireshark application

### <u>Step 2 : (N2 observation)</u> 
#### Enter filtering conditions
Since we want to observe packets between the gNB and the AMF, we can set the source and destination IP addresses as filtering conditions, as well as filter for NGAP packet types.

![wsfiltercondition1](./images/wireshark_filter1.png)

#### Sorted by time
You can observe from the 'Info' field what setup procedures the UE performed with the AMF through the gNB.

![wspackets](./images/wireshark_packets.png)

#### Explanation  
- NGSetupRequest & NGSetupResponse : After the gNB is initialized, it communicates with the AMF to perform NGAP protocol setup.
- 


### <u>Step 3 : (N3 observation)</u> 


## Exercise: PFCP protocol stack observation on N4
- N4: controls user plane configuration and session management between the `SMF` and `UPF`

In this exercise, 


## Reference
