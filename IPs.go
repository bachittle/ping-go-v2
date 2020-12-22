package ping2

import (
	"math"
	"net"
)

// CustomIP contains an IP and (optionally) a subnet. If optional, is nil.
type CustomIP struct {
	IP     net.IP
	Subnet *int
}

// GenerateIPs generates all the IP addresses
// given the subnet prefix (if there is one) and destination IP in CustomIP
func GenerateIPs(IPs []CustomIP) chan net.IP {
	c := make(chan net.IP)
	go func() {
		for _, ip2 := range IPs {
			var ip net.IP
			ip = append(ip, ip2.IP...)
			if ip2.Subnet != nil {
				maxSuffix := uint64(math.Pow(2, float64(32-*ip2.Subnet)))
				for i := uint64(0); i < maxSuffix; i++ {
					num := uint64(math.Pow(2, 24))
					if i >= num {
						ip[0] = ip2.IP[0] + byte(int(i/num)%256)
					}
					num = uint64(math.Pow(2, 16))
					if i >= num {
						ip[1] = ip2.IP[1] + byte(int(i/num)%256)
					}
					num = uint64(math.Pow(2, 8))
					if i >= num {
						ip[2] = byte(int(ip2.IP[2]) + int(i/num)%256)
					}
					ip[3] = byte(i % 256)
					c <- ip
				}
			} else {
				c <- ip2.IP
			}
		}
		close(c)
	}()
	return c
}
