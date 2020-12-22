package ping2

import (
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
)

// Capturer contains a set of IPs to capture from (or a subnet), otherwise it captures from everything on the wire
type Capturer struct {
	Iface  pcap.Interface
	handle *pcap.Handle
	IPs    []CustomIP
}

// GetDevs prints the devices, so you can pick which one to set the index to.
func (c Capturer) GetDevs() (str string) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	for i, dev := range devs {
		str = fmt.Sprintln(str, i, " : ", dev)
	}
	return
}

// Init sets the capturers interface to read off of, and reads it.
// For simplicity, you can specify the index of the interfaces given from FindAllDevs.
// it is different for every system, to see all devices and associated index, run GetDevs()
func (c *Capturer) Init(index int) error {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return err
	}
	c.Iface = devs[index]
	return nil
}

// CaptureIPs captures from the specified IP addresses in IPs
// it does this asynchronously, and returns captured echo replies to a returned channel
// you can close the channel with the optional parameter channel. Just pass it a number.
func (c Capturer) CaptureIPs(stopChan *chan int) chan net.IP {
	ch := make(chan net.IP)
	go func() {
		handle, err := pcap.OpenLive(c.Iface.Name, 1600, true, pcap.BlockForever)
		if err != nil {
			panic(err)
		}
		defer handle.Close()
		c.handle = handle
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			if stopChan != nil {
				// in the loop, if stopChan is ever given any values, break the loop.
				select {
				case <-*stopChan:
					break
				default:
				}
			}
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
			if icmpLayer != nil && ipLayer != nil {
				icmpLayer, ok := icmpLayer.(*layers.ICMPv4)
				if !ok {
					panic(err)
				}
				for ip := range GenerateIPs(c.IPs) {
					ipLayer, ok := ipLayer.(*layers.IPv4)
					if !ok {
						panic(err)
					}
					if icmpLayer.TypeCode.Type() != 0 || icmpLayer.TypeCode.Code() != 0 {
						break
					}
					if bytes.Equal(ipLayer.SrcIP, ip) {
						var ip2 net.IP
						ip2 = append(ip2, ip...)
						ch <- ip2
					}
				}
			}
		}
		close(ch)
	}()
	return ch
}
