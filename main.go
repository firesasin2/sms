package main

import "log"

func main() {

	Pid := 1
	if flagPid > 0 {
		Pid = flagPid
	}

	p, err := NewProcess(Pid)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(p)
}
