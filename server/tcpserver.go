package main

import (
	"encoding/json"
	"log"
)

func RecieveDataFromClient() {

	port := "1234"
	address := "0.0.0.0"

	s, err := NewServer(address+":"+port, "../etc/server.pem", false)
	if err != nil {
		log.Println(err)
	}
	go s.Handle(ConnHandler)
}

//서버
func ConnHandler(peer string, c *Client) {
	defer c.Close()

	for {
		msgs, err := c.Read()
		if err != nil {
			log.Println(err)
			break
		}
		for _, msg := range msgs {
			//log.Println("[%d번째 메시지] From:%s Msg:%s\n", i, peer, string(msg))
			var p Process
			err = json.Unmarshal(msg, &p)
			if err != nil {
				log.Println(err)
			}

			// csv 파일에 정보를 넣기 위해 q에 프로세스를 전달합니다.
			q <- p
		}
	}
}
