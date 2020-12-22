package ping2

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
)

// Sender sends an echo request to the specified IP address
type Sender struct {
	SrcIP  net.IP
	DstIPs []CustomIP
}

// SendOne sends one packet to the specified SrcIP and DstIPs
// returns the amount of successful IPs it sent to, and the random ID generated for checking
func (s Sender) SendOne(id int) (int, error) {
	conn, err := icmp.ListenPacket("ip:icmp", fmt.Sprint(s.SrcIP))
	n := 0
	if err != nil {
		return n, err
	}

	reqMsg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  1,
			Data: []byte(""),
		},
	}

	reqBinary, err := reqMsg.Marshal(nil)
	if err != nil {
		return n, err
	}

	for ip := range GenerateIPs(s.DstIPs) {
		ipAddr := &net.IPAddr{IP: ip, Zone: ""}
		_, err = conn.WriteTo(reqBinary, ipAddr)
		if err != nil {
			return n, err
		}
		n++
	}
	return n, err
}
