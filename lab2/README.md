# Lab 2: Networking Basics in Linux kernel

## Goals of this lab

This document primarily focuses on explaining how the kernel handles various networking behaviors. For practical system implementations, refer to the provided links. The main objective is to understand how the kernel processes network functions and observe how user-space applications interact with the kernel's interfaces to achieve network functionality.

For information on the datapath kernel module [gtp5g](https://ithelp.ithome.com.tw/articles/10302887) in free5GC, you can refer to the following article:

1. [gtp5g Architecture](<https://free5gc.org/guide/Gtp5g/design/#introduction>)
2. [gtp5g PCFP Architecture](<https://free5gc.org/guide/Upf_PFCP/design/>)
3. [Introduction to gtp5g and some kernel concepts ](<https://free5gc.org/blog/20230920/Introduction_of_gtp5g_and_some_kernel_concepts/>)
4. [gtp5g source code explanation ](<https://ithelp.ithome.com.tw/articles/10302887>)

## Important Data Struture in Networking Programming

### sk_buff

> Each packet has its own sk_buff structure.

1. The most frequently allocated and freed structure in the network subsystem.
2. Allocated by skb_init in net/core/sk_buff.c.
3. The data structure is defined in <include/linux/skbuff.h>.
4. This structure is used from L2 to L4, with each layer modifying the fields (adding headers).

#### structure

```
                                ---------------
                               | sk_buff       |
                                ---------------
   ,---------------------------  + head
  /          ,-----------------  + data
 /          /      ,-----------  + tail
|          |      |            , + end
|          |      |           |
v          v      v           v
 -----------------------------------------------
| headroom | data |  tailroom | skb_shared_info |
 -----------------------------------------------
                               + [page frag]
                               + [page frag]
                               + [page frag]
                               + [page frag]       ---------
                               + frag_list    --> | sk_buff |
                                                   ---------

```

[Ref]( https://docs.kernel.org/networking/skbuff.html)

A doubly linked list composed of sk_buff_head and sk_buff represents the network stack's packet queues, such as transmit (TX) and receive (RX) queues. Some fields exist merely to simplify searching, and each sk_buff must be able to quickly find the head of this linked list.
```
struct sk_buff_head {
       /* These two members must be first. */
       struct sk_buff  *next;
       struct sk_buff  *prev;
   
       __u32       qlen; //list 中元素數量
      spinlock_t  lock; // 防止同時存取
  };
```
The detailed definition of sk_buff is influenced by the kernel's configuration settings at compile time. The fields may vary depending on the features enabled (e.g., enabling QoS functionality). Below are descriptions of common fields.

```
struct sk_buff {
    /* These two members must be first. */
    struct sk_buff          *next;
    struct sk_buff          *prev;

    struct sock             *sk; // This field is needed for socket-related data at L4. It is null if the host is neither the destination nor the source.
    struct net_device       *dev;

    char                    cb[48] __aligned(8); // Control block

    unsigned long           _skb_refdst;
    void                    (*destructor)(struct sk_buff *skb);
#ifdef CONFIG_NET_DMA
    struct dma_chan         *dma_chan;
    dma_cookie_t            dma_cookie;
#endif
    /* Data */
    unsigned char           *head; // Points to buffer head
    unsigned char           *data; // Points to data head
    unsigned char           *tail; // Points to data end
    unsigned char           *end;  // Points to buffer end
    unsigned int            len;   // Changes as this structure exists at different layers. When moving up the layers, headers are discarded.
    unsigned int            data_len; // Only contains the size of the data
    unsigned int            mac_len;
    unsigned short          gso_size;
    unsigned short          csum_offset;

    __u32                   priority; // Used for QoS
    __u32                   mark;
    __u16                   queue_mapping;
    __u16                   protocol; // Used by the driver to determine which handler at a higher level should process the packet (each protocol has its own handler)
    __u8                    pkt_type;
    __u8                    ip_summed;
    __u8                    ooo_okay;
    __u8                    nohdr;
    __u8                    nf_trace;
    __u8                    ipvs_property;
    __u8                    peeked;
    __u8                    nfctinfo;
    __u8                    napi_id;

    __be32                  secmark;

    union {
        __wsum              csum;
        struct {
            __u16           csum_start;
            __u16           csum_offset;
        };
    };

    __u16                   vlan_proto;
    __u16                   vlan_tci;

#ifdef CONFIG_NET_CLS_ACT
    __u32                   tc_index;       /* Traffic control index */
#endif

#ifdef CONFIG_NET_SCHED
    __u16                   tc_verd;        /* Traffic control verdict */
#endif

    __u16                   tstamp;
    __u16                   hash;
    __u32                   timestamp;

    struct sk_buff          *frag_list;
    struct sk_buff          *next_frag;

    skb_frag_t              frags[MAX_SKB_FRAGS];
    skb_frag_t              frag_list;

    struct ubuf_info        *ubuf_info;

    struct skb_shared_info  *shinfo;
};

```


1. net_device: The role depends on whether the packet is about to be sent or has just been received.
    1. Receive: The device driver will update this field with a pointer to the data structure representing the receiving interface.
    2. Transmission is more complex 

    Some network functions can aggregate "some devices" into a single virtual interface, serviced by a virtual driver.

    The virtual driver will select a specific device and then set the dev field to point to this net_device structure. Therefore, this value changes during the packet processing.

2. cb: control block
    1. Each layer has its private information storage here, storing temporary data.
    2. For example:
        - IP

            ```c
            struct ip_skb_cb {
                __u32    addr;
                __u32    options;
                // Other temporary information for the IP layer
            };

            void ip_process(struct sk_buff *skb) {
                struct ip_skb_cb *cb = IP_CB(skb);
                cb->addr = ...;
                cb->options = ...;
                // Process the IP packet
            }
            ```

        - TCP

            TCP places tcp_skb_cb data in the cb area

            ```c
            struct tcp_skb_cb {
                union {
                    struct {
                        __u32  seq;       /* TCP sequence number */
                        __u32  end_seq;   /* TCP end sequence number */
                        union {
                            struct {
                                __u16  flag;       /* TCP flags */
                                __u16  sacked;     /* SACKed status */
                            };
                            __u32  ack_seq;        /* Acknowledgment sequence number */
                        };
                    };
                    __u8  header[48];    /* Align to 48 bytes */
                };
                __u32  when;             /* Transmission time */
                __u32  acked;            /* ACKed status */
            };
            ```

            Access using macro

            ```c
            #define TCP_SKB_CB(__skb)   ((struct tcp_skb_cb *)&((__skb)->cb[0]))
            ```

            Example

            ```c
            #include <linux/skbuff.h>
            #include <net/tcp.h>

            void process_tcp_skb(struct sk_buff *skb) {
                struct tcp_skb_cb *tcb = TCP_SKB_CB(skb);

                // Set TCP sequence number
                tcb->seq = 1000;
                // Set TCP end sequence number
                tcb->end_seq = 2000;
                // Set TCP flags
                tcb->flag = 0x10;  // For example, ACK flag

                // Print TCP sequence number and end sequence number
                printk(KERN_INFO "TCP seq: %u, end_seq: %u\n", tcb->seq, tcb->end_seq);
            }
            ```

3. priority
    1. Used for QoS
    2. If the packet is locally generated, the socket layer defines this value.
    3. If the packet is to be forwarded, this field is defined based on the ToS of the IP header.
4. cloned
    1. When an ingress packet needs to be given to multiple receivers, such as protocol handlers, network taps, etc.


## Kernel Space Communitation with User Space
Network programs require communication between user space and kernel space for several reasons:



**Access hardware resources**: Network programs often need to access hardware resources like network interface cards and memory. Since these resources are managed by the operating system kernel only, programs must use system calls to request access from kernel space.

**Transmit and receive network packets**: As packet transmission and reception involve hardware operations, programs must use system calls to send requests to kernel space.

**Networking services like routing and firewalls**: Some network programs provide services like routing and firewalls. These services manage network traffic and packet routing, so programs need to collaborate with kernel space to obtain necessary information and permissions.

#### sysctl
Sysctl is a mechanism in Unix-like operating systems that allows user space programs to read and modify kernel parameters. These parameters, also known as "kernel variables," control various aspects of the kernel's behavior, such as networking, memory management, and process scheduling.

For instance, during the installation of free5GC, the command ```sudo sysctl -w net.ipv4.ip_forward=1``` is used to enable routing functionality on the host machine. These variables are stored in a pseudo-filesystem ,```proc/sys``` . The ip_forward variable is stored in ```/proc/sys/net/ipv4/ip_forward```



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
/* interpret the response */
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

[A comparison between ioctl & Netlink](<https://medium.com/thg-tech-blog/on-linux-netlink-d7af1987f89d>)

[Netlink in free5GC dataplane](<https://free5gc.org/blog/20230920/Introduction_of_gtp5g_and_some_kernel_concepts/#free5gc-upf>)

## ip Command Introduction
The ip command is a versatile tool for managing network interfaces, routing tables, and other networking configurations in Linux systems. 

Here are some of the common functions of ip command.

* ```ip link```: Manages network interfaces, including creating, configuring, and displaying their status.

* ```ip addr```: Manages IP addresses, assigning, deleting, and displaying addresses for network interfaces.

* ```ip route```: Manages routing tables, adding, deleting, and modifying routes to direct network traffic.

* ```ip neigh```: Manages neighbor cache entries, displaying and manipulating ARP and NDISC entries for connected devices.

* ```ip rule```: defines routing policies (which routing table to use) while ip route specifies routing entries (how to forward packets).The route command is restricted to operating on a single routing table, while policy-based routing (PBR) employs multiple routing tables concurrently
    * The differences between ```ip route``` & ```ip rule```

        ip **route** show: 
        ```
        $ ip route
        default via 10.0.2.2 dev enp0s3 proto dhcp src 10.0.2.15 metric 100 
        10.0.2.0/24 dev enp0s3 proto kernel scope link src 10.0.2.15 
        10.0.2.2 dev enp0s3 proto dhcp scope link src 10.0.2.15 metric 100 
        192.168.55.0/24 dev enp0s8 scope link 
        192.168.56.0/24 dev enp0s8 proto kernel scope link src 192.168.56.202
        ```

        ip **rule** show:
        ```
        $ ip rule 
        0:      from all lookup local
        32766:  from all lookup main
        32767:  from all lookup default

        $ ip route show table main
        default via 10.0.2.2 dev enp0s3 proto dhcp src 10.0.2.15 metric 100 
        10.0.2.0/24 dev enp0s3 proto kernel scope link src 10.0.2.15 
        10.0.2.2 dev enp0s3 proto dhcp scope link src 10.0.2.15 metric 100 
        192.168.55.0/24 dev enp0s8 scope link 
        192.168.56.0/24 dev enp0s8 proto kernel scope link src 192.168.56.202        
        ```




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
recvmsg(4, {msg_name={sa_family=AF_NETLINK, nl_pid=0,}})
.
.
.
```

NETLINK_ROUTE:
    Receives routing and link updates and may be used to
    modify the routing tables (both IPv4 and IPv6), IP
    addresses, link parameters, neighbor setups, queueing
    disciplines, traffic classes, and packet classifiers


## How Kernel Receive and Transmit Packet
The following shows the different PDUs associated with three layer of the protocol stack in Kernel.
1. L2: Frame
2. L3: Packet
3. L4: Segment


### Interrupt & Driver
When an interface receives a packet, the driver, device, and kernel work together using interrupts to handle the packet. The detailed process is as follows:

1. Device Receives Packet: The network interface card (NIC) or network device receives the packet and converts it to an electrical signal.
2. Device Generates Interrupt: The device sends an interrupt signal to the CPU to indicate that a packet has arrived.
3. CPU Acknowledges Interrupt: The CPU halts its current tasks and identifies the interrupt source using the Interrupt Vector Table (IVT).
4. CPU Executes ISR: The CPU runs the Interrupt Service Routine (ISR) from the driver layer to handle the interrupt, which processes the incoming packet.
5. Driver Processes Packet: The ISR passes the packet information to the driver, which processes it, such as parsing headers and handling data.
6. Kernel Processes Packet: The driver may pass the packet to the kernel for further processing, such as routing and TCP/IP handling.
7. CPU Resumes Tasks: After packet processing, the CPU resumes its previous tasks.

* IRQ (Interrupt Request): A signal sent by a hardware device to the CPU indicating that it needs attention.

* ISR (Interrupt Service Routine): The code executed by the CPU in response to an IRQ, responsible for handling the interrupt request.

In Linux, the kernel typically does not have direct access to packets stored in a device's queue due to security layers that prevent unauthorized access. When a device receives a packet, it is stored in a queue until the driver processes it. The kernel accesses packet information only through the driver layer, which provides interfaces for packet retrieval.

#### Network Association Interrupt
1. Buffers such as sk_buff need to be allocated.
2. The received data is copied into these buffers.
3. Parameters within the buffer are initialized to inform higher-layer protocols about the type of data.

#### Softirq
**Softirqs** (short for "soft interrupts") are a mechanism in the Linux kernel for deferring the processing of certain tasks from the context of a hardware interrupt handler to a later time when the system is less loaded. This allows the hardware interrupt handler to return quickly, minimizing the time that the CPU is unavailable for other tasks.

- If the interrupt handler needs to perform tasks that are not time-critical or can be safely delayed, it can defer them to softirqs.
- Softirqs are software interrupts that are executed in a separate task context, typically after the interrupt handler has completed its essential tasks.
- Deferring tasks to softirqs allows the interrupt handler to return quickly, minimizing the time the CPU is unavailable for other tasks.
- ksoftirqd is a kernel thread operates on a per-CPU basis, meaning that each CPU has its own dedicated ksoftirqd thread. This allows for parallel processing of soft interrupts, further enhancing the system's responsiveness and ability to handle concurrent workloads.
    - We can observe the ksoftirqd by using px aux | grep soft
        ```
        ps aux | grep soft
        root           9  0.0  0.0      0     0 ?        S    Jul17   0:03 [ksoftirqd/0]
        root          18  0.0  0.0      0     0 ?        S    Jul17   0:04 [ksoftirqd/1]
        root          24  0.0  0.0      0     0 ?        S    Jul17   0:01 [ksoftirqd/2]
        root          30  0.0  0.0      0     0 ?        S    Jul17   0:06 [ksoftirqd/3]
        ```

### L2
#### queue
> Maintained by kernel
* Each queue will contain a pointer to the associated device and a pointer to the ingress/egress sk_buff.
* The loopback device is special as it does not require a queue at all.

#### What Handler (Device Driver) Do?
1. Copy the frame into the sk_buff structure.
2. Initialize the sk_buff settings for higher-level handlers.
    For example, set skb→protocol to specify the upper-layer protocol.
3. Schedule the NET_RX_SOFTIRQ softirq to notify the kernel of the arrival of a new frame.
    Also, place the frame in the CPU's private queue (stored in each CPU's separate ```softnet_data```).

#### How to Notify Kernel
* Old Method: netif_rx
    * Some device drivers use this method.
    * Typically called in an interrupt context, it may temporarily disable CPU interrupts.
* [NAPI (New API)](https://docs.kernel.org/networking/napi.html) 
    * The main idea is to combine interrupt and polling mechanisms.
    * New interrupts are not generated if the kernel has not finished processing other packets in the queue.
    * Reduces CPU loading during high workloads by decreasing the number of interrupts.
#### Processing Frame (Receive)
[```netif_receive_skb```](<https://github.com/torvalds/linux/blob/master/net/core/dev.c>)

Each driver is associated with a specific hardware type (e.g., Ethernet), allowing it to interpret the L2 header and identify the L3 protocol.

Main Functions:

1. Pass the frame copy to the protocol demultiplexer.
2. Forward the frame copy to the L3 handler specified by skb→protocol.
3. Handle layer-specific functions, such as bridging.
#### Transmission Enable & Disable:

When a device driver detects that there is insufficient memory to store a Maximum Transmission Unit (MTU), it calls netif_stop_queue to stop the egress queue. This prevents further transmissions, as the kernel knows it will fail and waste resources.This responsibility lies with the driver.


#### Scheduling for Transmission:

```dev_queue_xmit``` dequeues a frame using one of two methods:

1. Interface with Traffic Control (QoS layer): Through ```qdisc_run```.
2. Directly pass the frame to the device’s hard_start_xmit.

[more details of ```dev_queue_xmit```](<https://abcdxyzk.github.io/blog/2015/08/25/kernel-net-dev_queue_xmit/>)

The sole parameter is a structure containing sk_buff:

1. ```skb→dev``` specifies the outgoing device.
2. ```skb→data``` points to the payload.

When the egress queue is closed, netif_schedule is called to schedule the device:

1. Adds the frame to the head of each CPU’s output queue.
2. Schedules a softirq.

#### Handling NET_TX_SOFTIRQ: ```net_tx_action```
1. Handles tasks that can be deferred.
2. When transmission is complete, dev_kfree_skb_irq is called to notify that the related buffer can be released.

### L3

#### Main Tasks:
1. Sanity Check: Verify checksum and ensure headers are within their specified fields.
2. Firewalling: Netfilter subsystem (used at multiple points in the packet processing flow).
Option:
3. IP Options [details](https://net.academy.lv/lection/net_LS-08ENa_ip-options.pdf)
4. Fragmentation:
Fragmentation and reassembly consume CPU resources, increasing latency. The kernel may use Path MTU Discovery to determine the maximum MTU to avoid fragmentation, updating the routing table with the discovered PMTU.


#### Initial Setup

Handled by ```ip_init```
1. Registers handlers
2. Initializes routing subsystem
3. Sets up the infrastructure for managing IP endpoints

#### Interaction with Netfilter

Netfilter has hooks at various points in the network stack. When a packet or kernel condition matches, the packet passes through these hooks.

Common hooks include:

1. Packet reception
2. Packet forwarding
    *  Prerouting
    *  Postrouting
5. Packet transmission

#### Interaction with Routing Subsystem

* ```ip_route_input```:
Determines the packet's fate: whether to send it to local or forward it.
* ```ip_route_output_flow```:
Returns a gateway and egress net_device.
* ```dst_pmtu```:
Retrieves the PMTU from a specific field in the routing table.
* ```ip_route_… functions```:
    Query the routing table to make decisions based on:
    1. Source IP
    2. Destination IP
    3. Type of Service (ToS)
    4. Receiving device
    5. Sendable devices

#### Handling IP Packets
Handling IP Packets

Registered kernel handler: ```ip_rcv```

After performing sanity checks on the packet, it calls the Netfilter hook.

Other processing is completed by ```ip_rcv_finish```.

```
/*
 * IP receive entry point
 */
int ip_rcv(struct sk_buff *skb, struct net_device *dev, struct packet_type *pt,
           struct net_device *orig_dev)
{
    struct net *net = dev_net(dev);

    skb = ip_rcv_core(skb, net);
    if (skb == NULL)
        return NET_RX_DROP;

    return NF_HOOK(NFPROTO_IPV4, NF_INET_PRE_ROUTING,
                   net, NULL, skb, dev, NULL,
                   ip_rcv_finish);
}

```

```ip_rcv_finish```

1. Determines whether the packet should be forwarded to localhost or to the next hop.
2. Handles some IP options.

[more details](<https://github.com/torvalds/linux/blob/master/net/ipv4/ip_input.c>)

#### Forwarding

Relevant functions are defined in [ip_forward.c](<https://github.com/torvalds/linux/blob/master/net/ipv4/ip_forward.c>)

There are only two functions:

```ip_forward```: Handles all packets with addresses different from local ones.
1. Ensures the packet address can be forwarded.
2. Decreases TTL.
3. Fragmentation (based on MTU).
4. Sends to the outgoing device.

```ip_forward_finish```: At this point, all checks are completed, and the packet is ready to be sent to another system.

## Qeustion

1. What is sk_buff? Please describe in detail the information it stores, its purpose, and the layers that use it.
2. What is Netlink in the context of Linux networking? Describe its purpose, the typical flow of a Netlink communication, the Netlink message format, and how it compares to ioctl.
3. Please explain the main functions of the ip command in Linux network management and through which system calls it accomplishes its tasks.
4. Explain the process of handling network-related interrupts in Linux, including the roles of the device, driver, CPU, and kernel. Additionally, describe what a softirq is and its significance in this context.
5. Describe the process of handling Layer 2 (L2) frames in Linux, including queue management, kernel notification, and frame processing and transmission.
6. Explain the processing of Layer 3 (L3) IP packets in Linux, including initialization, interaction with Netfilter, and packet forwarding.



