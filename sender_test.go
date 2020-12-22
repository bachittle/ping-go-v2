package ping2

import (
	"net"
	"testing"
)

func TestSend(t *testing.T) {
	s := Sender{
		SrcIP: net.IP{192, 168, 50, 77},
	}
	s.DstIPs = append(s.DstIPs, CustomIP{net.IP{192, 168, 50, 1}, nil})

	s.SendOne()

	s = Sender{
		SrcIP: net.IP{192, 168, 50, 77},
	}
	n := 24
	s.DstIPs = append(s.DstIPs, CustomIP{net.IP{192, 168, 50, 1}, &n})

	s.SendOne()
}
