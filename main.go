package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {

	// 전체 프로세스를 생성합니다.
	pss, err := NewProcesses()
	if err != nil {
		log.Fatal(err)
	}

	// 입력받은 이름으로부터 프로세스들을 찾습니다.
	fpss, err := pss.FindProcessByName(flagPName)
	if err != nil {
		log.Fatal(err)
	}

	// CSV파일에 Header를 씁니다.
	name := os.Args[0] + ".csv"
	w, err := WriteCSVHeader(name)
	if err != nil {
		log.Fatal(err)
	}

	// 동시에 CSV파일에 쓰기 위해 channel 생성
	q := make(chan Process)

	// q에 요청이 들어오면, CSV파일에 q내용을 씁니다.
	go WriteCSVBody(w, q)

	// 찾은 모든 프로세스들을 주기적으로 모니터링합니다.
	for _, fps := range fpss {
		go MonitorProcess(fps, q)
	}

	// 종료대기 : SIGINT (Ctrl+C) 신호를 받을때까지 기다립니다.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
