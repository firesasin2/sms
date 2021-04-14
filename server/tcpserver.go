package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
)

// Process 구조체
type Process struct {
	Pid       int
	PPid      int
	Name      string
	Umask     string
	State     string
	Tgid      int
	Ngid      int
	TracerPid int
	Uid       []int32
	Gid       []int32
	FDSize    int
	Groups    int
	VmPeak    int
	VmSize    int
	VmLck     int
	VmPin     int
	VmHWM     int
	VmRSS     int

	VmData int
	VmStk  int
	VmExe  int
	VmLib  int
	VmPTE  int
	VmSwap int

	Threads    int
	NoNewPrivs int
	Seccomp    int

	UpTime     int
	Utime      int
	Stime      int
	StartTime  int
	CreateTime int

	TotalCPU   int
	CPUPercent string

	Memory Memory
}

type Memory struct {
	Size int
	Rss  int
	Pss  int
}

func RecieveDataFromServer() {

	port := "1234"
	protocol := "tcp"
	address := "0.0.0.0"

	listen, err := net.Listen(protocol, address+":"+port)
	if err != nil {
		log.Printf("Socket listen port %s failed : %s\n", port, err)
		os.Exit(1)
	}
	defer listen.Close()

	log.Println("Begin listen port: ", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("listen.Accept() failed : ", err)
			continue
		}

		log.Println("conn Accepted")
		go ConnHandler(conn)
	}
}

var MSGHEADERLEN int = 4

func ConnHandler(c net.Conn) {

	jmbytes := [][]byte{}
	other := []byte{}

	readsize := 0
	var err error

	for {
		data := make([]byte, 20000)

		readsize, err = c.Read(data)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		if readsize == 0 {
			continue
		}

		jmbytes, other = SplitStream(append(other, data[:readsize]...))

		for _, jmbyte := range jmbytes {
			//msg := string(jmbyte)
			//log.Println("DATA : ", msg)

			var p Process
			err = json.Unmarshal(jmbyte, &p)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(p)
			}
		}
	}
}

func SplitStream(b []byte) ([][]byte, []byte) {
	bs := bytes.Split(b, []byte{0})
	return bs[:len(bs)-1], bs[len(bs)-1]
}
