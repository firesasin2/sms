package util

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

// CSV 헤더를 파일에 씁니다.
func WriteCSVHeader(fileName string) (*csv.Writer, error) {

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(bufio.NewWriter(f))

	// CSV파일에 헤더를 씁니다.
	p, err := NewProcess(os.Getpid())
	if err != nil {
		log.Fatal(err)
	}

	w.Write(MakeCSVHeaderFromProcess(p))
	w.Flush()

	return w, nil
}

// CSV파일 헤더를 만듭니다.(Process로부터)
func MakeCSVHeaderFromProcess(p Process) []string {
	line := []string{}

	for _, field := range flagfieldparsed {
		line = append(line, field)
	}

	return line
}

// CSV파일 데이터를 만듭니다.(1라인)
func MakeCSVValueFromProcess(p Process) []string {
	line := []string{}

	for _, field := range flagfieldparsed {
		switch field {
		case "TIME":
			line = append(line, strconv.FormatInt(time.Now().Unix(), 10))
		case "CPU":
			line = append(line, p.CPUPercent)
		case "MEMORYBYTES":
			line = append(line, strconv.Itoa(p.Memory.Pss))
		case "CMD1":
			line = append(line, p.Name)
		case "CMD2":
			line = append(line, p.Cmdline)
		case "PID":
			line = append(line, strconv.Itoa(p.Pid))
		case "PPID":
			line = append(line, strconv.Itoa(p.PPid))
		case "USER":
			line = append(line, p.UserName)
		case "CREATETIME":
			line = append(line, strconv.Itoa(p.CreateTime))
		}
	}

	return line
}

// CSV Body를 만듭니다.
func WriteCSVBody(w *csv.Writer) {

	for {
		// CSV파일에 프로세스 값을 씁니다.
		p := <-q
		w.Write(MakeCSVValueFromProcess(p))
		w.Flush()
	}
}
