package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Process 구조체
type Process struct {
	Name      string
	Umask     string
	State     string
	Tgid      int
	Ngid      int
	Pid       int
	PPid      int
	TracerPid int
	Uid       []int32
	Gid       []int32
	FDSize    int
	Groups    int

	Threads    int
	NoNewPrivs int
	Seccomp    int
}

// 프로세스 정보를 가져옵니다.
func NewProcess(pid int) (Process, error) {
	p := Process{
		Pid: pid,
	}

	// pid가 존재하는지 검사합니다.
	exist, err := PidExist(pid)
	if err != nil {
		return p, err
	}
	if !exist {
		return p, fmt.Errorf("pid의 프로세스가 없습니다. : %d", pid)
	}

	p.GetProcessStatus()

	return p, nil
}

// 프로세스 상태를 가져옵니다.(/proc/{pid}/status)
func (p *Process) GetProcessStatus() error {
	// Process 정보를 얻기 위해 /proc/{pid}/status를 파싱합니다.
	f, err := os.Open(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewScanner(f)
	for w.Scan() {
		line := w.Text()
		w := strings.Fields(line)
		if len(w) < 2 {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Name:"):
			p.Name = w[1]

		case strings.HasPrefix(line, "Umask:"):
			p.Umask = w[1]

		case strings.HasPrefix(line, "State:"):
			p.State = w[1]

		case strings.HasPrefix(line, "Tgid:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.Tgid = value

		case strings.HasPrefix(line, "Ngid:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.Ngid = value

		case strings.HasPrefix(line, "Pid:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.Pid = value

		case strings.HasPrefix(line, "PPid:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.PPid = value

		case strings.HasPrefix(line, "TracerPid:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.TracerPid = value

		case strings.HasPrefix(line, "Uid:"):
			if len(w) > 4 {
				// 작성 중
			}

		case strings.HasPrefix(line, "Gid:"):
			if len(w) > 4 {
				// 작성 중
			}

		case strings.HasPrefix(line, "FDSize:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.FDSize = value

		case strings.HasPrefix(line, "Groups:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.Groups = value

		case strings.HasPrefix(line, "Threads:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.Threads = value

		case strings.HasPrefix(line, "NoNewPrivs:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.NoNewPrivs = value

		case strings.HasPrefix(line, "Seccomp:"):
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.Seccomp = value
		}
	}

	return nil
}

// 프로세스 ID가 존재하는지 검사합니다.
func PidExist(pid int) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("pid가 0보다 작습니다 : %d", pid)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}

	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil
	}

	return false, err
}
