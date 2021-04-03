package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagHelp    bool
	flagVersion bool
	flagPid     int
)

func init() {
	flag.BoolVar(&flagHelp, "h", false, "도움말")
	flag.BoolVar(&flagVersion, "v", false, "버전")
	flag.IntVar(&flagPid, "p", 0, "프로세스아이디")
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
}

// version info
func PrintVersion() {
	fmt.Println(`Component Name: sms V0.1`)
	fmt.Println(`Component ReleaseVersion: V0.1(2021-04)`)
}
