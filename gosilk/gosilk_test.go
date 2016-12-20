package gosilk

import "testing"
import "time"
import "net"

func TestInit(t *testing.T) {

stream:=New("wang")

stream.CreateHeader(&Header{Probename: "P4444", Invocation: "i ran this file"})

stream.WriteRecord(&Rec{ sTime: time.Now(),
			 eTime: time.Now().Add(10000000),
			 sPort: 3333,
			 dPort: 4444,
			 sIP:  net.ParseIP("10.10.10.1"),
			 dIP:  net.ParseIP("192.168.1.1"),
			 proto: 6})

stream.CloseStream()


}
