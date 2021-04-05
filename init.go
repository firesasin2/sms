package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagHelp     bool
	flagVersion  bool
	flagPid      int
	flagPName    string
	flagInterval int
)

func init() {
	flag.BoolVar(&flagHelp, "h", false, "도움말")
	flag.BoolVar(&flagVersion, "v", false, "버전")
	//flag.IntVar(&flagPid, "pid", 1, "프로세스 아이디")
	flag.StringVar(&flagPName, "pname", "nginx", "프로세스 이름")
	flag.IntVar(&flagInterval, "i", 20, "수집 주기")
	flag.Parse()

	if flagHelp {
		PrintHelp()
		os.Exit(0)
	}
	if flagVersion {
		PrintVersion()
		os.Exit(0)
	}
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
	fmt.Println(`  -pname`)
	fmt.Println(`  프로세스 아이디`)
	fmt.Println(`  -i`)
	fmt.Println(`  수집 주기`)
}

// version info
func PrintVersion() {
	fmt.Println(`Component Name: sms V0.1`)
	fmt.Println(`Component ReleaseVersion: V0.1(2021-04)`)
}
