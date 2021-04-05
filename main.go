package main

import (
	"log"
	"os"
	"time"
)

func main() {

	// 전체 프로세스를 생성합니다.
	pss, err := NewProcesses()
	if err != nil {
		log.Fatal(err)
	}

	// 프로세스 이름을 찾습니다.
	p, err := pss.FindProcessByName(flagPName)
	if err != nil {
		log.Fatal(err)
	}

	// 프로세스 객체를 생성합니다.
	// p, err := NewProcess(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// CSV파일 Writer를 생성합니다.
	fileName := os.Args[0] + ".csv"
	CSVWriter, err := WriteCSVLine(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// CSV파일에 헤더를 씁니다.
	CSVWriter.Write(MakeCSVHeaderFromProcess(p))
	CSVWriter.Flush()

	for {
		// 프로세스 상태를 최신으로 유지합니다.
		err = p.GetProcessStatus()
		if err != nil {
			log.Fatal(err)
		}

		// test 로그
		log.Println(p)

		// CSV파일에 프로세스 값을 씁니다.
		CSVWriter.Write(MakeCSVValueFromProcess(p))
		CSVWriter.Flush()

		// 지정 주기만큼 sleep합니다.
		time.Sleep(time.Duration(flagInterval) * time.Second)
	}
}
