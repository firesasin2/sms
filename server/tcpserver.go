package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
)

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
			//log.Println("DATA : ", msg)

			var p Process
			err = json.Unmarshal(jmbyte, &p)
			if err != nil {
				log.Println(err)
			}

			// csv 파일에 정보를 넣기 위해 q에 프로세스를 전달합니다.
			q <- p
		}
	}
}

func SplitStream(b []byte) ([][]byte, []byte) {
	bs := bytes.Split(b, []byte{0})
	return bs[:len(bs)-1], bs[len(bs)-1]
}
