package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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
	VmPeak    int
	VmSize    int
	VmLck     int
	VmPin     int
	VmHWM     int
	VmRSS     int

	VmData int
	VmStk  int
	VmExe  int
	VmLib  int
	VmPTE  int
	VmSwap int

	Threads    int
	NoNewPrivs int
	Seccomp    int
}

type Processes struct {
	pss []Process
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

	return p, nil
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

// /Proc 밑의 프로세스들에서 이름으로 프로세스를 찾습니다.
func (pss *Processes) FindProcessByName(name string) (Process, error) {
	for _, ps := range pss.pss {
		matched, err := regexp.MatchString(name, ps.Name)
		if err != nil {
			return ps, err
		}
		if matched {
			return ps, nil
		}
	}
	return Process{}, fmt.Errorf("프로세스 이름을 찾을 수 없습니다.")
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
		// 2개 이상의 인자값을 갖는 line만 분석합니다.(가정: 1번째 인자는 속성이름, 2번째 인자부터는 값)
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

		case strings.HasPrefix(line, "VmPeak:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmPeak = value * 1024

		case strings.HasPrefix(line, "VmSize:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmSize = value * 1024

		case strings.HasPrefix(line, "VmPin:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmPin = value * 1024

		case strings.HasPrefix(line, "VmHWM:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmHWM = value * 1024

		case strings.HasPrefix(line, "VmRSS:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmRSS = value * 1024

		case strings.HasPrefix(line, "VmData:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmData = value * 1024

		case strings.HasPrefix(line, "VmStk:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmStk = value * 1024

		case strings.HasPrefix(line, "VmExe:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmExe = value * 1024

		case strings.HasPrefix(line, "VmLib:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmLib = value * 1024

		case strings.HasPrefix(line, "VmPTE:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmPTE = value * 1024

		case strings.HasPrefix(line, "VmSwap:"):
			if len(w) != 3 {
				continue
			}
			if w[2] != "kB" {
				continue
			}
			value, err := strconv.Atoi(w[1])
			if err != nil {
				return err
			}
			p.VmSwap = value * 1024

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

// func (lp *Processes) FindProcessByName(name string) ([]LinuxProcess, error) {
// 	out := []LinuxProcess{}
// 	for _, p := range lp.ps {

// 		if ok, err := regexp.MatchString(name, p.Comm); err != nil {
// 			return out, err
// 		} else {
// 			if ok {
// 				out = append(out, p)
// 			}
// 		}
// 	}
// 	return out, nil
// }
