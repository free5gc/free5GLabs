#!/bin/bash
#
# Configure iptables in N6GW
#
ip route delete default
ip route add default via 10.100.6.1 dev dn0

ip route add 10.60.0.0/16 via 10.100.6.100
ip route add 10.61.0.0/16 via 10.100.6.100

iptables -t nat -A POSTROUTING -o dn0  -j MASQUERADE
iptables -I FORWARD 1 -j ACCEPT