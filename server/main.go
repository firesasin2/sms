package main

import (
	"os"
	"os/signal"
)

func main() {

	// 클라이언트에서 데이터를 수신합니다.
	go RecieveDataFromServer()

	// 종료대기 : SIGINT (Ctrl+C) 신호를 받을때까지 기다립니다.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
