package ping2

import (
	"io/ioutil"
	"runtime"
	"testing"
)

func TestGetDevs(t *testing.T) {
	var c Capturer
	t.Log(c.GetDevs())
}

func TestInit(t *testing.T) {
	var c Capturer
	outcome := "eth0"
	switch runtime.GOOS {
	case "windows":
		c.Init(3)

		b, err := ioutil.ReadFile("secret.txt")
		if err != nil {
			panic(err)
		}
		outcome = string(b)
		break
	default:
		c.Init(0)
		break
	}
	if c.Iface.Name != outcome {
		t.Error("Values do not match!")
		t.Error(c.Iface.Name, "!=", outcome)
	}
}

// starts and stops right away
func TestCapture(t *testing.T) {
	var c Capturer
	c.Init(3)
	ch := make(chan int, 1)
	c.CaptureIPs(&ch, 1)
	ch <- 1
}
