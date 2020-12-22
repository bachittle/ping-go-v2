package ping2

import (
	"fmt"
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
func (c Controller) SendAndRecv(timeout time.Duration) {
	c.cap.Init(3)
	chanIP := c.cap.CaptureIPs(nil)
	chanInt := make(chan int)
	n, err := c.sen.SendOne()
	fmt.Println("n:", n)
	if err != nil {
		panic(err)
	}

	dict := make(map[string]bool)

	go func() {
		var ip net.IP
	Loop:
		for {
			select {
			case <-chanInt:
				break Loop
			case ip = <-chanIP:
			}
			fmt.Println("got ip:", ip)
			dict[ip.String()] = true
			n--
		}
	}()
	if n == 0 {
		chanInt <- 1
		return
	}
	time.Sleep(timeout)
	fmt.Println("IPs found:", len(dict))
	fmt.Println(dict)
	chanInt <- 1
	return
}
