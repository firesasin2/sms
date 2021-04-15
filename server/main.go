package main

import (
	"log"
	"os"
	"os/signal"
)

// 동시에 CSV파일에 쓰기 위해 channel
var (
	q               chan Process
	flagfieldparsed = []string{"TIME", "CPU", "MEMORYBYTES", "CMD1", "CMD2", "PID", "PPID", "USER", "CREATETIME"}
)

func main() {

	// channel 초기화
	q = make(chan Process)

	// CSV파일에 Header를 씁니다.
	name := os.Args[0] + ".csv"
	w, err := WriteCSVHeader(name)
	if err != nil {
		log.Fatal(err)
	}

	// q에 요청이 들어오면, CSV파일에 q내용을 씁니다.
	go WriteCSVBody(w)

	// 클라이언트에서 데이터를 수신합니다.
	go RecieveDataFromServer()

	// 종료대기 : SIGINT (Ctrl+C) 신호를 받을때까지 기다립니다.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
