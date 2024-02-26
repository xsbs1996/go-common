package gotool

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

// IsIPv4 判断是否ipv4地址
func IsIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ".")
}

// IsIPv6 判断是否ipv6地址
func IsIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ":")
}

// IPv4toUint32 ip格式转uint32
func IPv4toUint32(ip string) (uint32, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return 0, errors.New("ParseIP error")
	}
	i = i.To4()
	return binary.BigEndian.Uint32(i), nil
}

// Uint32toIPv4 Uint32转ip格式
func Uint32toIPv4(ipInt uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	return ip.String()
}

// Uint32toIP Uint32转成net.IP
func Uint32toIP(ipInt uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipInt)
	return ip
}

// GetLocalIP 获取本机网卡IP
func GetLocalIP() ([]string, error) {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, v := range addrList {
		if ipNet, ok := v.(*net.IPNet); ok {
			ip := ipNet.IP
			if ip.To4() != nil && !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips, nil
}

// CIDRToUint32 将CIDR转成数字,如  1.0.0.0/24 转成 16777216 16777471
func CIDRToUint32(cidr string) (start uint32, end uint32, err error) {
	i, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return
	}
	mask, bit := n.Mask.Size()
	start, err = IPv4toUint32(i.String())
	if err != nil {
		return
	}
	network, err := IPv4toUint32(n.IP.String())
	if err != nil {
		return
	}
	end = 1<<(uint32(bit-mask)) + network - 1
	return
}

// CIDRToIPRange 将CIDR转成起始IP-结束IP,如  192.168.0.0/24 转成 192.168.0.0 192.168.0.255
func CIDRToIPRange(cidr string) (startIp string, endIp string, err error) {
	start, end, err := CIDRToUint32(cidr)
	if err != nil {
		return
	}
	startIp = Uint32toIPv4(start)
	endIp = Uint32toIPv4(end)
	return
}

// IPRangeToCIDR 将起始IP-结束IP转成CIDR,如 192.168.0.0 192.168.0.255  转成 192.168.0.0/24
func IPRangeToCIDR(startIp string, endIp string) (cidr string, err error) {
	start, err := IPv4toUint32(startIp)
	if err != nil {
		return
	}
	end, err := IPv4toUint32(endIp)
	if err != nil {
		return
	}
	//取主机位长度
	bit := len(fmt.Sprintf("%b", start^end))
	//起始地址（网络号）
	ipInt := (start >> uint32(bit)) << uint32(bit)
	cidr = fmt.Sprintf("%v/%v", Uint32toIPv4(ipInt), 32-bit)
	return
}

// ipNetMaskBit 用startIp,endIp计算子网掩码长度，如192.168.0.0 192.168.0.255 返回 24
func ipNetMaskBit(startIp string, endIp string) (bit int, err error) {
	start, err := IPv4toUint32(startIp)
	if err != nil {
		return
	}
	end, err := IPv4toUint32(endIp)
	if err != nil {
		return
	}
	bit = 32 - len(fmt.Sprintf("%b", start^end))
	return
}

// CIDRToIPMask 将 1.1.1.0/24  转成1.1.1.0 255.255.255.0
func CIDRToIPMask(cidr string) (string, string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", err
	}
	bytesToIP := func(b []byte) string {
		i := binary.BigEndian.Uint32(b)
		return Uint32toIPv4(i)
	}
	ip := ipNet.IP.String()
	mask := bytesToIP(ipNet.Mask)
	return ip, mask, nil
}

// IPMaskToCIDR 将 1.1.1.0 255.255.255.0   转成 1.1.1.0/24
func IPMaskToCIDR(ip string, mask string) string {
	IP := net.ParseIP(ip)
	if IP == nil {
		return ""
	}
	return fmt.Sprintf("%s/%v", IP, MaskLength(mask))
}

// InverseMask 计算反掩码  将255.255.255.0 转成0.0.0.255
func InverseMask(mask string) string {
	i, err := IPv4toUint32(mask)
	if err != nil {
		return ""
	}
	return Uint32toIPv4(^i)
}

// MaskLength 计算掩码长度  255.255.255.0 得出 24
func MaskLength(mask string) int {
	ip := net.ParseIP(mask).To4()
	if ip == nil {
		return 0
	}
	i := []byte(ip)
	ipMask := net.IPv4Mask(i[0], i[1], i[2], i[3])
	ones, _ := ipMask.Size()
	return ones
}

// PrivateIpNet Well-known IPv4 Private addresses
var PrivateIpNet = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

// IsPrivateIP 是否为私网ip
func IsPrivateIP(ip string) bool {
	for _, ipNet := range PrivateIpNet {
		_, n, _ := net.ParseCIDR(ipNet)
		if n.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}
