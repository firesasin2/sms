package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Process 구조체
type Process struct {
	Pid       int
	PPid      int
	Name      string
	Cmdline   string
	Umask     string
	State     string
	Tgid      int
	Ngid      int
	TracerPid int
	Uid       []int
	UserName  string
	Gid       []int
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

	UpTime     int
	Utime      int
	Stime      int
	StartTime  int
	CreateTime int

	TotalCPU   int
	CPUPercent string

	Memory Memory
}

type Stat struct {
	Cpu       string
	User      int
	Nice      int
	System    int
	Idle      int
	Iowait    int
	Irq       int
	Softirq   int
	Steal     int
	Guest     int
	Guestnice int

	Total int
}

var (
	Hertz    = 100
	ProcStat = []Stat{}
)

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

// /Proc 밑의 프로세스들에서 이름으로 프로세스를 찾습니다.
func (pss *Processes) FindProcessByName(name string) ([]Process, error) {
	fpss := []Process{}
	for _, ps := range pss.pss {
		matched, err := regexp.MatchString(name, ps.Name)
		if err != nil {
			return fpss, err
		}
		if matched {
			fpss = append(fpss, ps)
		}
	}
	return fpss, nil
}

// 프로세스 상태를 가져옵니다.(/proc/stat)
func GetProcStat() (Stat, error) {

	s := Stat{}

	f, err := os.Open(fmt.Sprintf("/proc/stat"))
	if err != nil {
		return s, err
	}
	defer f.Close()

	w := bufio.NewScanner(f)
	for w.Scan() {

		line := w.Text()

		w := strings.Fields(line)
		if len(w) < 2 {
			continue
		}
		// cpu만 분석합니다.
		if !strings.HasPrefix(w[0], "cpu") {
			continue
		}

		s.Cpu = w[0]

		value, err := strconv.Atoi(w[1])
		if err != nil {
			return s, err
		}
		s.User = value

		value, err = strconv.Atoi(w[2])
		if err != nil {
			return s, err
		}
		s.Nice = value

		value, err = strconv.Atoi(w[3])
		if err != nil {
			return s, err
		}
		s.System = value

		value, err = strconv.Atoi(w[4])
		if err != nil {
			return s, err
		}
		s.Idle = value

		value, err = strconv.Atoi(w[5])
		if err != nil {
			return s, err
		}
		s.Iowait = value

		value, err = strconv.Atoi(w[6])
		if err != nil {
			return s, err
		}
		s.Irq = value

		value, err = strconv.Atoi(w[7])
		if err != nil {
			return s, err
		}
		s.Softirq = value

		value, err = strconv.Atoi(w[8])
		if err != nil {
			return s, err
		}
		s.Steal = value

		s.Total = s.User + s.Nice + s.System + s.Idle + s.Iowait + s.Irq + s.Softirq + s.Steal
	}

	return s, nil
}

// 프로세스 상태를 가져옵니다.(/proc/{pid}/stat)
func (p *Process) GetProcessStat() error {

	// /proc/{pid}/stat를 파싱합니다.
	d, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", p.Pid))
	if err != nil {
		return err
	}
	w := strings.Fields(string(d))

	if len(w) < 22 {
		return fmt.Errorf("stat 파일 파싱 오류 : 예상한 길이보다 작습니다.")
	}

	value, err := strconv.Atoi(w[0])
	if err != nil {
		return err
	}
	p.UpTime = value

	value, err = strconv.Atoi(w[13])
	if err != nil {
		return err
	}
	p.Utime = value

	value, err = strconv.Atoi(w[14])
	if err != nil {
		return err
	}
	p.Stime = value

	value, err = strconv.Atoi(w[21])
	if err != nil {
		return err
	}
	p.StartTime = value

	// uptime을 구합니다.
	uptimeFileBytes, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return err
	}
	uptimeFileString := string(uptimeFileBytes)
	uptimeString := strings.Split(uptimeFileString, " ")[0]
	fvalue, err := strconv.ParseFloat(uptimeString, 64)
	if err != nil {
		return err
	}
	uptime := fvalue

	// createTime을 구합니다.
	now := int(time.Now().Unix())
	p.CreateTime = now - int(uptime) + (p.StartTime / Hertz)

	//log.Println(p.CreateTime, uptime, Hertz, p.StartTime/100, p.UpTime, p.Utime, p.Stime)

	return nil
}

// 프로세스 cmdline를 가져옵니다.(/proc/{pid}/cmdline)
func (p *Process) GetProcessCmdline() error {
	bvalue, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", p.Pid))
	if err != nil {
		return err
	}
	// 널 문자를 띄어 쓰기로 바꿈
	f := func(r rune) bool {
		if r == '\u0000' {
			return true
		}
		return false
	}
	arr := strings.FieldsFunc(string(bvalue), f)

	p.Cmdline = strings.Join(arr, " ")

	return nil
}

// 프로세스 상태를 가져옵니다.(/proc/{pid}/status)
func (p *Process) GetProcessStatus(acs []Account) error {
	// /proc/{pid}/status를 파싱합니다.
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
				p.Uid = make([]int, 0, 4)
				for i := 1; i <= 4; i++ {
					value, err := strconv.Atoi(w[i])
					if err != nil {
						return err
					}
					p.Uid = append(p.Uid, value)
				}
				// effective UID
				value, err := findAccountNameFromUid(p.Uid[1], acs)
				if err != nil {
					return err
				}
				p.UserName = value
			}

		case strings.HasPrefix(line, "Gid:"):
			if len(w) > 4 {
				p.Gid = make([]int, 0, 4)
				for i := 1; i <= 4; i++ {
					value, err := strconv.Atoi(w[i])
					if err != nil {
						return err
					}
					p.Gid = append(p.Gid, value)
				}
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

// Hertz를 구합니다.
func GetHertz() (int, error) {
	// Hertz를 구합니다.
	HertzOutput, err := exec.Command("getconf", "CLK_TCK").Output()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	HertzArray := strings.Split(string(HertzOutput), "\n")[0]
	value, err := strconv.Atoi(strings.Fields(HertzArray)[0])
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return value, nil
}

func (p *Process) GetTotalCPU() error {

	// /proc/stat를 가져옵니다.
	ProcStat, err := GetProcStat()
	if err != nil {
		return err
	}

	p.TotalCPU = ProcStat.Total

	return nil
}

// CPU 사용량(%)을 계산합니다.
func calculateCPUPercent(p Process, oldProcess Process) (float64, error) {
	diff := float64(p.TotalCPU - oldProcess.TotalCPU)
	percent := 100.0 * ((float64(p.Utime+p.Stime) - float64(oldProcess.Utime+oldProcess.Stime)) / diff)

	return percent, nil
}

// 프로세스를 모니터한후, channel에 결과를 전송합니다.
func MonitorProcess(p Process) {

	var err error

	// 전체 계정을 가져옵니다.
	acs, err := GetAccounts()
	if err != nil {
		log.Println(err)
	}

	oldProcess := p

	// 프로세스 상태 정보를 가져옵니다.
	err = p.GetProcessStatus(acs)
	if err != nil {
		log.Println(err)
	}
	err = p.GetProcessStat()
	if err != nil {
		log.Println(err)
	}
	err = p.GetTotalCPU()
	if err != nil {
		log.Println(err)
	}
	err = p.GetProcessCmdline()
	if err != nil {
		log.Println(err)
	}

	// 메모리 상태정보를 가져옵니다.
	err = p.Memory.GetProcessMomory(p.Pid)
	if err != nil {
		log.Println(err)
	}

	p.CPUPercent = "0.00"
	// CPU 사용 퍼센트를 계산합니다.
	percent, err := calculateCPUPercent(p, oldProcess)
	if err != nil {
		log.Println(err)
	}
	p.CPUPercent = fmt.Sprintf("%.2f", percent)

	// CSV에 프로세스 현재 상태를 씁니다.
	q <- p
	q2 <- p

	// test 로그
	log.Println(p.Name, p.Pid, p.CPUPercent, p.Memory.Pss)
}
