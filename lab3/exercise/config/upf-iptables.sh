#!/bin/bash
#
# Configure iptables in UPF
#
ip route delete default via 10.100.3.1
ip route add default via 10.100.6.1 dev dn0

iptables -t nat -A POSTROUTING -o dn0  -j MASQUERADE
iptables -I FORWARD 1 -j ACCEPT
