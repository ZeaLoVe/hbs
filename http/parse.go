package http

import (
	"fmt"
	"github.com/open-falcon/hbs/g"
	"net"
	"strings"
)

var privateBlocks []*net.IPNet

func init() {
	// Add each private block
	privateBlocks = make([]*net.IPNet, 3)
	_, block, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		panic(fmt.Sprintf("Bad cidr. Got %v", err))
	}
	privateBlocks[0] = block

	_, block, err = net.ParseCIDR("172.16.0.0/12")
	if err != nil {
		panic(fmt.Sprintf("Bad cidr. Got %v", err))
	}
	privateBlocks[1] = block

	_, block, err = net.ParseCIDR("192.168.0.0/16")
	if err != nil {
		panic(fmt.Sprintf("Bad cidr. Got %v", err))
	}
	privateBlocks[2] = block
}

func isPrivateIP(ip_str string) bool {
	ip := net.ParseIP(ip_str)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

func ParseIP(ip string, block *net.IPNet) string {
	if block == nil {
		return ip
	}
	origin_ip := strings.Split(ip, ".")
	target_ip := strings.Split(block.IP.String(), ".")
	target_ip[3] = origin_ip[3]
	return fmt.Sprintf("%v.%v.%v.%v", target_ip[0], target_ip[1], target_ip[2], target_ip[3])
}

func PrivateIP(ip string, Nats []g.NAT) string {
	for _, nat := range Nats {
		//简单的过滤掉 不需要判断的网段
		if ip[1] != nat.PublicIP[1] && ip[2] != nat.PublicIP[2] {
			continue
		}
		_, public, err := net.ParseCIDR(nat.PublicIP)
		if err != nil {
			panic(fmt.Sprintf("Bad cidr. Got %v", err))
		}
		if public.Contains(net.ParseIP(ip)) {
			_, private, err := net.ParseCIDR(nat.PrivateIP)
			if err != nil {
				panic(fmt.Sprintf("Bad cidr. Got %v", err))
			}
			return ParseIP(ip, private)
		}
	}
	return ip
}
