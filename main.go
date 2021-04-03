package main

import "log"

func main() {
	if flagPid > 0 {
		log.Println(NewProcess(flagPid))
	} else {
		log.Println(NewProcess(1))
	}
}
