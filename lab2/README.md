# Lab 2: Networking Basics in Linux kernel

## Introduction



## Goals of this lab

## Important Data Struture in Networking Programming

## Kernel Space Communitation with User Space
Network programs require communication between user space and kernel space for several reasons:

**Access hardware resources**: Network programs often need to access hardware resources like network interface cards and memory. Since these resources are managed by the operating system kernel only, programs must use system calls to request access from kernel space.

**Transmit and receive network packets**: As packet transmission and reception involve hardware operations, programs must use system calls to send requests to kernel space.

**Networking services like routing and firewalls**: Some network programs provide services like routing and firewalls. These services manage network traffic and packet routing, so programs need to collaborate with kernel space to obtain necessary information and permissions.

#### sysctl
Sysctl is a mechanism in Unix-like operating systems that allows user space programs to read and modify kernel parameters. These parameters, also known as "kernel variables," control various aspects of the kernel's behavior, such as networking, memory management, and process scheduling.

For instance, during the installation of free5gc, the command ```sudo sysctl -w net.ipv4.ip_forward=1``` is used to enable routing functionality on the host machine. These variables are stored in a pseudo-filesystem ,```proc/sys``` . The ip_forward variable is stored in ```/proc/sys/net/ipv4/ip_forward```



#### ioctl (input/output control)
ioctl is a system call used to perform device-specific input/output operations and other operations that cannot be expressed by regular file semantics. It is a powerful mechanism for user-space programs to interact with device drivers and control hardware devices.

Ioctl works by sending a request code and any necessary data to the kernel. The kernel then processes the request and returns a response. The specific behavior of an ioctl call depends on the device driver and the request code.

[man page of ioctl](<https://man7.org/linux/man-pages/man2/ioctl.2.html>)


#### Netlink
The goal of netlink was to provide a better way of modifying network related settings and transferring network related information between userspace and kernel.
When you are writing a linux application that needs either kernel to userspace communications or userspace to kernel communications, the typical answer is to use `ioctl` and sockets.

This is a simple mechanism for sending information down from userspace into the kernel to make requests for info, or to direct the kernel to perform an operation on behalf of the userspace application.

Traditionally, ioctls were used for communication between user space and the kernel, particularly for network-related tasks. However, ioctls have limitations for complex data transfer and can be cumbersome. Netlink emerged as a more versatile solution specifically designed for network communication. It offers a message-based approach with features like reliable delivery, multicast support, and flexible data structures, making it a superior choice for exchanging network data and managing network settings.

A very simplified flow of a Netlink “call” will therefore look something like:

``` bash
fd = socket(AF_NETLINK, SOCK_RAW, NETLINK_GENERIC);

/* format the request */
send(fd, &request, **sizeof**(request));
n = recv(fd, &response, RSP_BUFFER_SIZE);
/* interpret the response n */
```
#### Netlinl message format

``` 
Netlink Message
   |-- Netlink Message Header
   |   |-- Message Length (nlmsg_len)
   |   |-- Message Type (nlmsg_type)
   |   |-- Message Flags (nlmsg_flags)
   |   |-- Sequence Number (nlmsg_seq)
   |   |-- Sender Process ID (nlmsg_pid)
   |   |-- Message Group ID (nlmsg_group)
   |-- Netlink Message Payload
   |   |-- (Data specific to the message type)
   |-- Padding

```

```
struct nlmsghdr {
 __u32  nlmsg_len; /* Length of message including header */
 __u16  nlmsg_type; /* Message content */
 __u16  nlmsg_flags; /* Additional flags */
 __u32  nlmsg_seq; /* Sequence number */
 __u32  nlmsg_pid; /* Sending process port ID */
};
```

A message can contain a second header defining the type of netlink message; the most common of these are:

- `NETLINK_ROUTE` for modifying routing tables, queuing, traffic classifiers etc.
- `NETLINK_NETFILTER` for netfilter related information
- `NETLINK_KOBJECT_UEVENT` for communications from kernel to userspace (for an application to subscribe to kernel events)
- `NETLINK_GENERIC` for users to develop application specific messages

[a comparison between ioctl & Netlink](<https://medium.com/thg-tech-blog/on-linux-netlink-d7af1987f89d>)

[Netlink in free5gc dataplane](<https://free5gc.org/blog/20230920/Introduction_of_gtp5g_and_some_kernel_concepts/#free5gc-upf>)

## ip Command Introduction
The ip command is a versatile tool for managing network interfaces, routing tables, and other networking configurations in Linux systems. 

Here are some of the common functions of ip command.

* ip link: Manages network interfaces, including creating, configuring, and displaying their status.

* ip addr: Manages IP addresses, assigning, deleting, and displaying addresses for network interfaces.

* ip route: Manages routing tables, adding, deleting, and modifying routes to direct network traffic.

* ip neigh: Manages neighbor cache entries, displaying and manipulating ARP and NDISC entries for connected devices.

Next, We use strace to observe how the ip command uses system calls.

The kernel maintains various data structures to store information about network interfaces, their state, and associated resources:

- **struct net_device:** This is the fundamental data structure that represents a network interface. It encapsulates a wealth of information, including the interface name, type, MAC address, MTU, flags, statistics, and pointers to other relevant data structures.
- **struct net_device_stats:** This structure holds statistical information about a network interface, such as the number of packets transmitted and received, errors encountered, and bytes transferred.

### ip link
some common command
- `ip link`: Displays a list of all network interfaces.
- `ip link show dev <interface_name>`: Shows detailed information about a specific interface.
- `ip link set dev <interface_name> up`: Brings up an interface.
- `ip link set dev <interface_name> down`: Takes down an interface.
- `ip link set dev <interface_name> mtu <new_mtu>`: Changes the MTU of an interface.
- `ip link add name <virtual_interface_name> type veth peer name <peer_interface_name>`: Creates a virtual Ethernet (veth) pair.
- `ip link add name <bridge_interface_name> type bridge`: Creates a bridge interface.

#### ip link show
```
$strace -e getsockopt,setsockopt,bind,sendto,recvmsg ip link show
setsockopt(3, SOL_SOCKET, SO_SNDBUF, [32768], 4) = 0
setsockopt(3, SOL_SOCKET, SO_RCVBUF, [1048576], 4) = 0
setsockopt(3, SOL_NETLINK, NETLINK_EXT_ACK, [1], 4) = 0
bind(3, {sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, 12) = 0
setsockopt(3, SOL_NETLINK, NETLINK_DUMP_STRICT_CHK, [1], 4) = 0
sendto(3, {{len=32, type=RTM_NEWLINK, flags=NLM_F_REQUEST|NLM_F_ACK, seq=0, pid=0}, {ifi_family=AF_UNSPEC, ifi_ty
pe=ARPHRD_NETROM, ifi_index=0, ifi_flags=0, ifi_change=0}}, 32, 0, NULL, 0) = 32                                 recvmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base={{l
en=52, type=NLMSG_ERROR, flags=0, seq=0, pid=17245}, {error=-EPERM, msg={{len=32, type=RTM_NEWLINK, flags=NLM_F_REQUEST|NLM_F_ACK, seq=0, pid=0}, {ifi_family=AF_UNSPEC, ifi_type=ARPHRD_NETROM, ifi_index=0, ifi_flags=0, ifi_change=0}}}}, iov_len=16384}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 52                                sendto(3, {{len=40, type=RTM_GETLINK, flags=NLM_F_REQUEST|NLM_F_DUMP, seq=1721545369, pid=0}, {ifi_family=AF_PACK
ET, ifi_type=ARPHRD_NETROM, ifi_index=0, ifi_flags=0, ifi_change=0}, {{nla_len=8, nla_type=IFLA_EXT_MASK}, 1}}, 40, 0, NULL, 0) = 40                                                                                              recvmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}}, ...) 
. 
.
.
.                                                                          
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: enp0s3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP mode DEFAULT group default qlen 100
0                                                                                                                    link/ether 08:00:27:f8:ca:f5 brd ff:ff:ff:ff:ff:ff
3: enp0s8: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP mode DEFAULT group default qlen 100
0                       
```


* `getsockopt`: This system call retrieves socket options, which are often used for Netlink configuration.
    ```
    int getsockopt(int sock, int level, int optname, void *optval, socklen_t *optlen);
    ```
* `setsockopt`: This system call sets socket options, also commonly used for Netlink configuration.
    ```
        int setsockopt(int sock, int level, int optname, const void *optval, socklen_t optlen);
    ```
* `bind`: This system call binds a socket to a network address, which is sometimes necessary for Netlink communication.
* `sendto`: This system call sends data to a specified socket, including Netlink sockets.
* `recvmsg`: This system call receives data from a socket, including Netlink sockets.

#### IP link set enp0s3 up

```
$ sudo strace -e getsockopt,setsockopt,bind,sendto,recvmsg ip link set enp0s3 up
setsockopt(3, SOL_SOCKET, SO_SNDBUF, [32768], 4) = 0
setsockopt(3, SOL_SOCKET, SO_RCVBUF, [1048576], 4) = 0 //set receive buffer size
setsockopt(3, SOL_NETLINK, NETLINK_EXT_ACK, [1], 4) = 0
bind(3, {sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, 12) = 0
setsockopt(3, SOL_NETLINK, NETLINK_DUMP_STRICT_CHK, [1], 4) = 0
sendto(3, {{len=32, type=RTM_NEWLINK, flags=NLM_F_REQUEST|NLM_F_ACK, seq=0, pid=0}, {ifi_family=AF_UNSPEC, ifi_type=ARPHRD_NETROM, ifi_index=0, ifi_flags=0, ifi_change=0}}, 32, 0, NULL, 0) = 32
recvmsg(3, {msg_name={sa_family=AF_NETLINK,....}})
.
.
.

```
`sendto(3, ..., {type=RTM_NEWLINK, flags=NLM_F_REQUEST|NLM_F_ACK}, ...)`
This call sends a Netlink message of type `RTM_NEWLINK` with request and acknowledgement flags.
`ip link set` focuses on sending a specific message (RTM_NEWLINK) to modify the state of an interface.

### ip addr

Manages IP addresses and link-layer addresses (LLAs) associated with network interfaces.

- `ip addr show`: Displays all assigned IP addresses and LLAs.
- `ip addr show dev <interface_name>`: Shows IP addresses and LLAs for a specific interface.
- `ip addr add <ip_address>/<subnet_mask> dev <interface_name>`: Assigns an IP address to an interface.
- `ip addr del <ip_address>/<subnet_mask> dev <interface_name>`: Removes an IP address from an interface.
- `ip addr add lladdr <lladdr> dev <interface_name>`: Adds an LLA (e.g., MAC address) to an interface.
- `ip addr del lladdr <lladdr> dev <interface_name>`: Removes an LLA from an interface.

```
$ sudo strace -e socket,setsockopt,bind,sendto,recvmsg ip addr show
socket(AF_NETLINK, SOCK_RAW|SOCK_CLOEXEC, NETLINK_ROUTE) = 3
setsockopt(3, SOL_SOCKET, SO_SNDBUF, [32768], 4) = 0
setsockopt(3, SOL_SOCKET, SO_RCVBUF, [1048576], 4) = 0
setsockopt(3, SOL_NETLINK, NETLINK_EXT_ACK, [1], 4) = 0
bind(3, {sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, 12) = 0
setsockopt(3, SOL_NETLINK, NETLINK_DUMP_STRICT_CHK, [1], 4) = 0
sendto(3, {{len=40, type=RTM_GETLINK, flags=NLM_F_REQUEST|NLM_F_DUMP, seq=1721546151, pid=0}, {ifi_family=AF_UNSPEC, ifi_type=ARPHRD_NETROM, ifi_index=0, ifi_flags=0, ifi_change=0}, {{nla_len=8, nla_type=IFLA_EXT_MASK}, 1}}, 40, 0, NULL, 0) = 40
recvmsg(3, {msg_name={sa_family=AF_NETLINK, ......}}...)
.
.
.
.
```
The messages include details like:
- Interface name (`lo`, `enp0s3`)
- Interface type (loopback, ethernet)
- Flags (UP, RUNNING, BROADCAST)
- MTU (Maximum Transmission Unit)
- MAC address (`\x08\x00\x27\xf8\xca\xf5` for `enp0s3`)
- Statistics (packets transmitted/received, bytes transmitted/received)


### ip route
> Manages routing tables, controlling how packets are routed through the network.

- `ip route`: Displays the current routing table.
- `ip route add <destination_network>/<subnet_mask> via <gateway_ip>`: Adds a route to a specific network.
- `ip route del <destination_network>/<subnet_mask>`: Deletes a route to a network.
- `ip route show dev <interface_name>`: Shows routes associated with a specific interface.
- `ip route change <destination_network>/<subnet_mask> metric <new_metric>`: Changes the metric (priority) of a route.

#### ip route add 
This command will add a new route to the routing table for the 192.168.56.0/24 subnet. It will specify that traffic destined for this subnet should be sent through the enp0s8 network interface.

```
$ sudo strace -e 'setsockopt,socket,getsockname,sendmsg,recvmsg' ip route add 192.168.55.0/24 de
v enp0s8
socket(AF_NETLINK, SOCK_RAW|SOCK_CLOEXEC, NETLINK_ROUTE) = 3
setsockopt(3, SOL_SOCKET, SO_SNDBUF, [32768], 4) = 0
setsockopt(3, SOL_SOCKET, SO_RCVBUF, [1048576], 4) = 0
setsockopt(3, SOL_NETLINK, NETLINK_EXT_ACK, [1], 4) = 0
getsockname(3, {sa_family=AF_NETLINK, nl_pid=18468, nl_groups=00000000}, [12]) = 0
setsockopt(3, SOL_NETLINK, NETLINK_DUMP_STRICT_CHK, [1], 4) = 0
socket(AF_NETLINK, SOCK_RAW|SOCK_CLOEXEC, NETLINK_ROUTE) = 4
setsockopt(4, SOL_SOCKET, SO_SNDBUF, [32768], 4) = 0
setsockopt(4, SOL_SOCKET, SO_RCVBUF, [1048576], 4) = 0
setsockopt(4, SOL_NETLINK, NETLINK_EXT_ACK, [1], 4) = 0
getsockname(4, {sa_family=AF_NETLINK, nl_pid=-804631643, nl_groups=00000000}, [12]) = 0
sendmsg(4, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base={{len=52, type=RTM_GETLINK, flags=NLM_F_REQUEST, seq=1721550802, pid=0}, {ifi_family=AF_UNSPEC, ifi_type=ARPHRD_NETROM, ifi_index=0, ifi_flags=0, ifi_change=0}, [{{nla_len=8, nla_type=IFLA_EXT_MASK}, 9}, {{nla_len=11, nla_type=IFLA_IFNAME}, "enp0s8"}]}, iov_len=52}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 52
recvmsg(4, {msg_name={sa_family=AF_NETLINK, nl_pid=0,}}
.
.
.
```

NETLINK_ROUTE:
    Receives routing and link updates and may be used to
    modify the routing tables (both IPv4 and IPv6), IP
    addresses, link parameters, neighbor setups, queueing
    disciplines, traffic classes, and packet classifiers


## Forwarding

### L2

### L3

### L4

### routing subsystem

### Netfilter