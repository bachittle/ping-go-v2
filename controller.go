package ping2

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// Controller is the intermediary between capturer and sender
// 	  it sends data with the sender and captures packets and gets specific IP packets with the capturer
type Controller struct {
	SrcIP  net.IP
	DstIPs []CustomIP
	cap    Capturer
	sen    Sender
}

// Init initiallizes the controller.
func (c *Controller) Init() {
	c.cap.IPs = c.DstIPs
	c.sen.SrcIP = c.SrcIP
	c.sen.DstIPs = c.DstIPs

}

// SendAndRecv sends data with the sender and receives data with the capturer
func (c Controller) SendAndRecv(timeout time.Duration) map[string]bool {
	c.cap.Init(3)
	id := rand.Intn(65535)
	chanInt := make(chan int)
	chanIP := c.cap.CaptureIPs(nil, uint16(id))
	n, err := c.sen.SendOne(id)
	if err != nil {
		panic(err)
	}
	//fmt.Println("n:", n)

	dict := make(map[string]bool)

	go func() {
		var ip net.IP
	Loop:
		for {
			//fmt.Println("in loop...")
			select {
			case <-chanInt:
				break Loop
			case ip = <-chanIP:
				//fmt.Println("ip:", ip)
				dict[ip.String()] = true
				n--
				break
			}
		}
	}()
	if n == 0 {
		chanInt <- 1
		return dict
	}
	time.Sleep(timeout)
	fmt.Println("IPs found:", len(dict))
	//fmt.Println(dict)
	chanInt <- 1
	return dict
}
