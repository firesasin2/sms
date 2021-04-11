package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Memory struct {
	Size int
	Rss  int
	Pss  int
}

// 프로세스 메모리의 자세한 상태를 가져옵니다.(/proc/{pid}/smaps)
func (m *Memory) GetProcessMomory(Pid int) error {
	// 프로세스 메모리 정보를 얻기 위해 /proc/{pid}/smaps를 파싱합니다.
	f, err := os.Open(fmt.Sprintf("/proc/%d/smaps", Pid))
	if err != nil {
		return err
	}
	defer f.Close()

	m.Size = 0
	m.Rss = 0
	m.Pss = 0

	w := bufio.NewScanner(f)
	for w.Scan() {
		line := w.Text()
		w := strings.Fields(line)
		// 2개 이상의 인자값을 갖는 line만 분석합니다.(가정: 1번째 인자는 속성이름, 2번째 인자부터는 값)
		if len(w) < 3 {
			continue
		}
		// 3번째 인자 값은 kB 여야함
		if w[2] != "kB" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Size:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			m.Size += (value * 1024)

		// Resident set size. 해당 프로세스에서 사용 중인 물리적 페이지 크기(공유된 페이지 크기를 알 수 없음)
		case strings.HasPrefix(line, "Rss:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			m.Rss += (value * 1024)

		// Proportional set size. 해당 프로세스에서만 사용하는 고유한 페이지 수
		case strings.HasPrefix(line, "Pss:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			m.Pss += (value * 1024)

		}
	}

	return nil
}
