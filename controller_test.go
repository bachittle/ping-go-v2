package ping2

import (
	"net"
	"testing"
	"time"
)

func TestController(t *testing.T) {
	con := Controller{
		SrcIP: net.IP{192, 168, 50, 77},
	}
	con.DstIPs = append(con.DstIPs, CustomIP{net.IP{192, 168, 50, 1}, nil})

	con.Init()
	con.SendAndRecv(100 * time.Millisecond)

	con = Controller{
		SrcIP: net.IP{192, 168, 50, 77},
	}
	n := 24
	con.DstIPs = append(con.DstIPs, CustomIP{net.IP{192, 168, 50, 1}, &n})

	con.Init()
	con.SendAndRecv(5 * time.Second)
}
