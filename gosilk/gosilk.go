package gosilk

// #include <stdio.h>
// #include <silk/silk.h>
// #include <silk/sksite.h>
// #include <silk/skstream.h>
// #include <silk/utils.h>
// #include <silk/silk_types.h>
// #include <silk/rwrec.h>
// #cgo LDFLAGS: -lsilk
import "C"
import "fmt"
import "log"
import "time"
import "net"
import "encoding/binary"

type Header struct {
	compress int
	Invocation string
	annotation string
	Probename string
}

type Rec struct {
	sTime 	time.Time
	eTime   time.Time
	SPort 	uint16
	dPort 	uint16
	proto	uint8
	pkts	uint32
	bytes   uint32
	input   uint16
	output  uint16
	flags   uint8

	SIP     net.IP
	DIP	net.IP
	nhIP	net.IP
}

var written bool
	

func init() {
	C.skAppRegister(C.CString("gosilk"))
}

func New(filename string) *C.skstream_t {
var stream *C.skstream_t
rv:=C.skStreamCreate(&stream,C.SK_IO_WRITE,C.SK_CONTENT_SILK_FLOW)
if rv != 0 {
	log.Fatal("skStreamCreate error")
}
rv=C.skStreamBind(stream, C.CString(filename))
if rv != 0 {
	log.Fatal("skStreamBind error")
}

written=false
return stream
}


func (stream *C.skstream_t) CreateHeader(options *Header) {
	//var stream *C.skstream_t
	var hdr *C.sk_file_header_t
	var hentry *C.sk_header_entry_t
	var rv C.int
	
	hdr=C.skStreamGetSilkHeader(stream)
	//hentry=C.skHentryPackedfileCreate(1164215667,1,5)
	//rv=C.skHeaderAddEntry(hdr,hentry)
	//fmt.Println("AddEntry",rv)

	if options.Probename != "" {
	hentry=C.skHentryProbenameCreate(C.CString(options.Probename))
	rv=C.skHeaderAddEntry(hdr,hentry)
        fmt.Println("AddEntry",rv)
	}

	if options.Invocation != "" {
	hentry=C.skHentryProbenameCreate(C.CString(options.Invocation))
	rv=C.skHeaderAddEntry(hdr,hentry)
	}

	rv= C.skStreamOpen(stream);
	fmt.Println("StreamOpen",rv)
	
	rv=C.skStreamWriteSilkHeader(stream)
	fmt.Println("skStreamWriteSilkHeader",rv)
	

}

func (stream *C.skstream_t) WriteRecord(rec *Rec) {


	var rwrec C.rwRec
	rwrec.proto=C.uint8_t(rec.proto)
	rwrec.sPort=C.uint16_t(rec.SPort)
	rwrec.dPort=C.uint16_t(rec.dPort)
	rwrec.sTime=C.int64_t(rec.sTime.UnixNano()/1000000)
	rwrec.elapsed=C.uint32_t(rec.eTime.Sub(rec.sTime).Nanoseconds()/1000000)
	fmt.Println(rec.SIP)
	C.rwrec_SetDIPv4(&rwrec,C.uint32_t(binary.BigEndian.Uint32([]byte(rec.DIP.To4()))))
	fmt.Println(rec.DIP)
	//C.rwrec_SetSIPv4(&rwrec,C.uint32_t(binary.BigEndian.Uint32([]byte(rec.SIP.To4()))))

	rv:=C.skStreamWriteRecord(stream,&rwrec)
	fmt.Println("rwrec",rv,rwrec.sIP)
	written=true

}

func (stream *C.skstream_t)  CloseStream() {	
	rv:=C.skStreamClose(stream)
	fmt.Println("StreamClose",rv)
	if rv !=0 {
		panic("Disaster")
	}
	
	C.skStreamDestroy(&stream);
	
}

func (stream *C.skstream_t) HaveWritten() bool {

	return written
}

	
	
