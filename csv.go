package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func WriteCSVLine(fileName string) (*csv.Writer, error) {

	csvfile := fmt.Sprintf(fileName)
	f, err := os.OpenFile(csvfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	return csv.NewWriter(bufio.NewWriter(f)), nil
}

// CSV파일 헤더를 만듭니다.
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
