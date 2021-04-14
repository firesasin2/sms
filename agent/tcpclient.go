package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// tcp 소켓서버에 데이터를 전달하기 위한 함수
func WriteDataToServer(conn net.Conn) {

	for {
		// channel에 들어오면 tcp 소켓서버에 보냅니다.
		p := <-q2

		var buffer bytes.Buffer
		buffer.WriteString(MakeJSONStringFromProcess(p))

		_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
		_, err := conn.Write(append(buffer.Bytes(), 0))
		if err != nil {
			log.Printf("failed to send data : %v", err)
		}
	}
}

func MakeJSONStringFromProcess(p Process) string {

	pbyte, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	value := string(pbyte)

	return value
}
