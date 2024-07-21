# Lab 2: Networking Basics in Linux kernel

## Introduction



## Goals of this lab

## Terminology

## ip Command Introduction


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

## Forwarding

### L2

### L3

### L4

### routing subsystem

### Netfilter