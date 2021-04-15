package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
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
