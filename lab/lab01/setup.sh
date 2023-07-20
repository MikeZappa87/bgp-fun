#!/bin/bash

clab deploy --topo ./lab01.clab.yaml

ip netns exec clab-lab01-leaf01 ip link add cni0 type bridge
ip netns exec clab-lab01-leaf01 ip link set cni0 up
ip netns exec clab-lab01-leaf01 ip link set swp10 master cni0

ip netns exec clab-lab01-leaf02 ip link add cni0 type bridge
ip netns exec clab-lab01-leaf02 ip link set cni0 up
ip netns exec clab-lab01-leaf02 ip link set swp11 master cni0

ip netns exec clab-lab01-pc01 ip addr add 192.168.1.3/24 dev eth1
ip netns exec clab-lab01-pc02 ip addr add 192.168.1.4/24 dev eth1

ip netns exec clab-lab01-leaf01 ip addr add 192.168.1.1/24 dev cni0
ip netns exec clab-lab01-leaf02 ip addr add 192.168.1.1/24 dev cni0

ip netns exec clab-lab01-leaf01 ip addr add 172.16.1.10/32 dev lo
ip netns exec clab-lab01-leaf02 ip addr add 172.20.1.11/32 dev lo
ip netns exec clab-lab01-leaf02 ip addr add 172.20.1.12/32 dev lo

ip netns exec clab-lab01-pc01 ip route del default
ip netns exec clab-lab01-pc02 ip route del default
ip netns exec clab-lab01-leaf01 ip route del default
ip netns exec clab-lab01-leaf02 ip route del default
ip netns exec clab-lab01-spine01 ip route del default

ip netns exec clab-lab01-pc01 ip route add default via 192.168.1.1 dev eth1
ip netns exec clab-lab01-pc02 ip route add default via 192.168.1.1 dev eth1