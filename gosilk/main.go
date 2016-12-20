package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"io/ioutil"
	"gosilk"
	"path"

	"github.com/calmh/ipfix"
)

type IP []byte

type UDPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

type UDP4Addr struct {
	IP   [4]byte
	Port int
}

func writeData(f os.File, ifs []ipfix.InterpretedField) {

	//var mymap = map[string]interface{}
	mymap := make(map[string]interface{})

	for k, iif := range ifs {
		fmt.Println(iif.Name)
		mymap[iif.Name] = iif.Value
		fmt.Println(k, iif)

	}

}

func createFileandTemp(addr *net.UDPAddr) (finalname string, tempname string) {
        now:= time.Now()
        filename := fmt.Sprintf("%d%d%d%d%d%d.%v:%d.", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), addr.IP, addr.Port)
        prefix := "/tmp"
        f,err:=ioutil.TempFile(prefix,filename)
        checkError(err)
        finalname=f.Name()
        fmt.Println("Opening ",finalname)
        f.Close()
        tempname=path.Dir(finalname)+"/."+path.Base(finalname)
	fmt.Println("Temp name",tempname)
	return finalname,tempname
}

func process(addr *net.UDPAddr, in <-chan []byte, tick <-chan time.Time) {
	fmt.Println("plop")
	s := ipfix.NewSession()
	i := ipfix.NewInterpreter(s)

	opts := gosilk.Header{Probename: "P4444", Invocation: "i ran this file"}

	finalname,tempname:=createFileandTemp(addr)
	stream := gosilk.New(tempname)
	stream.CreateHeader(&opts)

	for {

		select {
		case data := <-in:
			fmt.Println("Length of data ", len(data))
			msg, err := s.ParseBuffer(data)
			if err != nil {
				fmt.Println("Error", err)
				continue
			}

			fmt.Println("Header ", msg.Header)
			//.DomainID," Sequence",msg.Header.SequenceNumber)
			fmt.Println("DataRecords ", len(msg.DataRecords), "Templates ", len(msg.TemplateRecords))
			//fmt.Println("Writing to ",f.Name())
			for j, record := range msg.DataRecords {
				fmt.Printf("Record number ", j)
				//		//	writeData(*f,record)
				ifs := i.Interpret(record)
				//fmt.Println(ifs)
				//	writeData(*f,ifs)
  				var rec gosilk.Rec
			        fmt.Println("ifs length ",len(ifs))	
				for k, iif := range ifs {
					//fmt.Fprintln(f,iif)
					fmt.Println(k, iif)
					if iif.EnterpriseID == 0 {
						switch iif.FieldID {
						case 1:
							fmt.Println("bytes")
						case 2:
							fmt.Println("packets")
						case 4:
							fmt.Println("protocol")
						case 6:
							fmt.Println("flags")
						case 7:
							fmt.Println("sport")
							rec.SPort=2222
						case 8:
							fmt.Println("4sIP")
							rec.SIP=net.IPv4(1,2,3,4)
						case 10:
							fmt.Println("in")
						case 11:
							fmt.Println("dport")
						case 12:
							fmt.Println("4dIP")
							rec.DIP=net.IPv4(1,2,3,4)
						case 14:
							fmt.Println("out")
						case 15:
							fmt.Println("4nhip")
						case 153:
							fmt.Println("flowEndMilliseconds")
						case 154:
							fmt.Println("flowStartMilliseconds")
						case 136:
							fmt.Println("flowendReason")
						default: 
							fmt.Println("Unhandled field",iif.Name,iif.FieldID)
							
						}
					}
				}
				rec.SIP=net.IPv4(1,2,3,4)
				rec.DIP=net.IPv4(1,2,3,4)
				stream.WriteRecord(&rec)
			}
			//f.Sync()
			fmt.Println("synced to disk")

		case <-tick:
			//Timer ticked.  Open another file
			stream.CloseStream()
			if stream.HaveWritten() {
			err:=os.Rename(tempname,finalname)
			checkError(err)
			} else {
			err:=os.Remove(tempname)
			checkError(err)
			err=os.Remove(finalname)
                        checkError(err)
			}

			finalname,tempname=createFileandTemp(addr)
			stream = gosilk.New(tempname)
        		stream.CreateHeader(&opts)
			// Do we need the above?
			//f,err=ioutil.TempFile("/tmp/",prefix)
			//checkError(err)
			//defer f.Close()
			fmt.Println("Tick")
			//writeHeader(*f)
		}

	}
}

func main() {
	log.Println("Wang")
	raddr, err := net.ResolveUDPAddr("udp", ":47000")
	checkError(err)
	rsock, err := net.ListenUDP("udp", raddr)
	defer rsock.Close()
	checkError(err)
	interval := 5

	//buf := make([]byte, 1500)
	var channels = make(map[UDP4Addr]chan []byte)
	tick := time.Tick(time.Duration(interval) * time.Second)

	for {
		buf := make([]byte, 1500)
		n, addr, err := rsock.ReadFromUDP(buf)
		fmt.Println("Received ", n, " bytes from ", addr.IP)
		//addrS:=addr.IP.String()

		// Check to see if we already have a session open for this source ip/port traffic
		addr4 := UDP4Addr{Port: addr.Port}
		copy(addr4.IP[:], addr.IP.To4()[0:4])
		fmt.Println(addr4)
		_, ok := channels[addr4]

		if ok {
			// We've seen this before
			fmt.Println("Session seen before")
			channels[addr4] <- buf[:n]
		} else {
			fmt.Println("New session")
			channels[addr4] = make(chan []byte)
			// Make a new goroutine to receive this data
			go process(addr, channels[addr4], tick)
			channels[addr4] <- buf[:n]

		}

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
