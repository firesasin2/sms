package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"reflect"
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

	// 첫번째 필드는 현재 시간을 기록합니다.
	line = append(line, "Time")

	// Process 구조체의 네임필드를 순회합니다.
	e := reflect.ValueOf(&p).Elem()
	fieldNum := e.NumField()
	for i := 0; i < fieldNum; i++ {
		t := e.Type().Field(i)
		line = append(line, t.Name)
	}

	return line
}

// CSV파일 데이터를 만듭니다.(1라인)
func MakeCSVValueFromProcess(p Process) []string {
	line := []string{}

	// 첫번째 필드는 현재 시간을 기록합니다.
	line = append(line, strconv.FormatInt(time.Now().Unix(), 10))

	// Process 구조체의 값필드를 순회합니다.
	e := reflect.ValueOf(&p).Elem()
	fieldNum := e.NumField()
	for i := 0; i < fieldNum; i++ {
		v := e.Field(i).Interface()

		switch t := reflect.TypeOf(v); t.Kind() {
		case reflect.Int:
			line = append(line, strconv.Itoa(v.(int)))
		case reflect.String:
			line = append(line, v.(string))
		default:
			line = append(line, "")
		}
	}

	return line
}

// CSV Body를 만듭니다.
func WriteCSVBody(w *csv.Writer, q chan Process) {

	for {
		// CSV파일에 프로세스 값을 씁니다.
		p := <-q
		w.Write(MakeCSVValueFromProcess(p))
		w.Flush()
	}
}
