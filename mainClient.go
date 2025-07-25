package main

import (
	"encoding/json"
	"os"
	"strings"
)

func init() {
	if len(os.Args) < 2 {
		// Q: что за os.Args?
		panic("failed: len(os.Args < 2)")
	}
	var (
		addrStr = ""
		userNewStr = ""
		userLoadStr = ""
	)
	// Q: Зачем создавать 2 var, если можно всё поставить в один?
	var (
		addrExist = false
		userNewExist = false
		userLoadExist = false
	)
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case strings.HasPrefix(arg, "-loadaddr:"):
			// Q: Что тут сейчас произошло?
			addrStr = strings.Replace(arg, "-loadaddr:", "", 1)
			addrExist = true
		case strings.HasPrefix(arg, "-newuser:"):
			userNewStr = strings.Replace(arg, "-newuser:", "", 1)
			userNewExist = true
		case strings.HasPrefix(arg, "-loaduser:"):
			userLoadStr = strings.Replace(arg, "-loaduser:", "", 1)
			userLoadExist = true
		}
	}
	if !(userNewExist || userLoadExist || !addrExist) {
		panic("failed: !(userNewExist || userLoadExist || !addrExist)")
	}

	err := json.Unmarshal([]byte(readFile(addrStr), &Addresses))
	if err != nil {
		panic("failed: load addresses")
	}
	if len(Addresses) == 0 {
		panic("failed: len(Addresses) == 0")
	}
	if userNewExist {
		User = userNew(userNewStr)
	}
	if userLoadExist {
		user = userLoad(userLoadStr)
	}
	if User == nil {
		panic("failed: load user")
	}
}