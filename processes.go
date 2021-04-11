package main

import (
	"io/ioutil"
	"strconv"
)

type Processes struct {
	pss []Process
}

// /Proc 밑의 프로세스들을 가지고 옵니다.
func NewProcesses() (Processes, error) {
	var pss Processes

	fss, err := ioutil.ReadDir("/proc/")
	if err != nil {
		return pss, err
	}

	// /proc Dir를 순회합니다.
	for _, fs := range fss {
		if !fs.IsDir() {
			return pss, err
		}

		Pid, err := strconv.Atoi(fs.Name())
		if err != nil {
			continue
		}
		// 프로세스 객체를 생성합니다.
		p, err := NewProcess(Pid)
		if err != nil {
			continue
		}

		// 프로세스 상태 정보를 가져옵니다.
		err = p.GetProcessStatus()
		if err != nil {
			continue
		}

		pss.pss = append(pss.pss, p)
	}
	return pss, nil
}
