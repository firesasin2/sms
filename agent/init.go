package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	flagHelp        bool
	flagVersion     bool
	flagPid         int
	flagPName       string
	flagInterval    int
	flagfield       string
	flagfieldparsed []string
)

// 동시에 CSV파일에 쓰기 위해 channel
var (
	q  chan Process
	q2 chan Process
)

func init() {
	flag.BoolVar(&flagHelp, "h", false, "도움말")
	flag.BoolVar(&flagVersion, "v", false, "버전")
	//flag.IntVar(&flagPid, "pid", 1, "프로세스 아이디")
	flag.StringVar(&flagPName, "p", "nginx", "프로세스 이름")
	flag.IntVar(&flagInterval, "i", 20, "수집 주기")
	flag.StringVar(&flagfield, "f", "TIME,CPU,MEMORYBYTES,CMD1,CMD2,PID,PPID,USER,CREATETIME", "필드")
	flag.Parse()

	if flagHelp {
		PrintHelp()
		os.Exit(0)
	}
	if flagVersion {
		PrintVersion()
		os.Exit(0)
	}

	// channel 초기화
	q = make(chan Process)
	q2 = make(chan Process)

	// csv에서 어떤 컬럼을 쓸지 인자값 받음
	parsed, err := ParseFlag(flagfield)
	if err != nil {
		log.Fatal(err)
	}
	flagfieldparsed = parsed
}

// help
func PrintHelp() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Println(`  -v`)
	fmt.Println(`  버전`)
	fmt.Println(`  -h`)
	fmt.Println(`  도움말`)
	// fmt.Println(`  -pid`)
	// fmt.Println(`  프로세스 아이디`)
	fmt.Println(`  -p`)
	fmt.Println(`  프로세스 이름`)
	fmt.Println(`  -i`)
	fmt.Println(`  수집 주기`)
}

// version info
func PrintVersion() {
	fmt.Println(`Component Name: sms V0.1`)
	fmt.Println(`Component ReleaseVersion: V0.1(2021-04)`)
}

// Flag를 파싱합니다.
func ParseFlag(flagfield string) ([]string, error) {
	var result []string
	fields := strings.Split(flagfield, ",")
	for _, field := range fields {
		result = append(result, field)
	}
	return result, nil
}
