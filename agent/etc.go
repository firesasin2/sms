package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Account struct {
	Name              string
	Status            string
	Uid               string
	Gid               string
	PrimaryGroupName  string
	Gecos             string
	HomeDir           string
	LoginShell        string
	PwdChangeDate     string // Date of last change.
	PwdChangeMinDays  string // Minimum number of days between changes.
	PwdChangeMaxDays  string // Maximum number of days between changes.
	PwdChangeWarnDays string // Number of days to warn user to change the password.
	InactiveDay       string // Number of days the account may be inactive.
	ExpireDay         string // Number of days since 1970-01-01 until account expires.
}

// Account 이름을 가져옵니다.
func findAccountNameFromUid(uid int, Accounts []Account) (string, error) {

	suid := strconv.Itoa(uid)

	for _, ac := range Accounts {
		if ac.Uid == suid {
			return ac.Name, nil
		}
	}

	return "", fmt.Errorf("")
}

// /etc/passwd파일로부터 account정보를 만듭니다.
func GetAccounts() ([]Account, error) {
	acs := []Account{}

	if bs, err := ioutil.ReadFile("/etc/passwd"); err != nil {
		return acs, err
	} else {
		ns := strings.Split(string(bs), "\n")

		for i := 0; i < len(ns); i++ {
			if len(ns[i]) > 0 && !strings.HasPrefix(ns[i], "#") {
				if ac, err := makeAccount(ns[i]); err != nil {
					return acs, err
				} else {
					acs = append(acs, ac)
				}
			}
		}
	}

	return acs, nil
}

func makeAccount(line string) (Account, error) {
	var err error

	ac := Account{}
	field := strings.Split(line, ":")

	if len(field) == 7 {
		ac.Name = field[0]
		ac.Uid = field[2]
		ac.Gid = field[3]
		if ac.PrimaryGroupName, err = convertGidToName(ac.Gid); err != nil {
			log.Println(err)
		}

		ac.Gecos = field[4]
		ac.HomeDir = field[5]
		ac.LoginShell = field[6]

		if ac.PwdChangeDate, ac.PwdChangeMinDays, ac.PwdChangeMaxDays, ac.PwdChangeWarnDays, ac.ExpireDay, ac.InactiveDay, err = readShadow(field[0]); err != nil {
			return ac, err
		}
	}

	return ac, nil
}

func convertGidToName(gid string) (string, error) {

	bgroup, err := ioutil.ReadFile("/etc/group")
	if err != nil {
		return "", err
	}

	line := strings.Split(string(bgroup), "\n")
	for i := 0; i < len(line); i++ {
		if len(line[i]) > 0 && !strings.HasPrefix(line[i], "#") {
			if strings.Contains(line[i], ":x:"+gid+":") {
				sp := strings.Split(line[i], ":")
				if (len(sp) > 3) && (sp[2] == gid) {
					return sp[0], nil
				}
			}
		}
	}

	return "", fmt.Errorf("Can not find the groupname of gid %s", gid)
}

// /etc/shadow파일을 읽습니다.
func readShadow(name string) (string, string, string, string, string, string, error) {
	shadow, err := ioutil.ReadFile("/etc/shadow")
	if err != nil {
		return "", "", "", "", "", "", err
	}

	line := strings.Split(string(shadow), "\n")
	for i := 0; i < len(line); i++ {
		if len(line[i]) > 0 && strings.HasPrefix(line[i], name+":") {
			sp := strings.Split(line[i], ":")
			if len(sp) == 9 {
				return sp[2], sp[3], sp[4], sp[5], sp[6], sp[7], nil
			}
		}
	}

	return "", "", "", "", "", "", nil
}
