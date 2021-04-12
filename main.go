package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	// Hertz를 구합니다.
	var err error
	Hertz, err = GetHertz()
	if err != nil {
		log.Fatal(err)
	}

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

	ticker := time.NewTicker(time.Duration(flagInterval) * time.Second)
	go func() {
		for _ = range ticker.C {
			for _, ps := range fpss {
				go MonitorProcess(ps, q)
			}
		}
	}()

	// 종료대기 : SIGINT (Ctrl+C) 신호를 받을때까지 기다립니다.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
