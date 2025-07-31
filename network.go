package main

import (
	"fmt"
	"strings"
	"time"

	nt "github.com/siyovush-hamidov/Blockchain/network"
)

const (
	TO_UPPER = iota + 1
	TO_LOWER
)

const (
	ADDRESS = ":8080"
)

func main() {
	var (
		res = new(nt.Package)
		msg = "Hello World!"
	)
	go nt.Listen(ADDRESS, handleServer)
	time.Sleep(500 * time.Millisecond)
	// отправляем Hello World!
	// получаем HELLO WORLD!
	res = nt.Send(ADDRESS, &nt.Package{
		Option: TO_UPPER,
		Data: msg,
	})
	fmt.Println(res.Data)

	// теперь отправим HELLO WORLD!
	// и получим hello world!
	res = nt.Send(ADDRESS, &nt.Package{
		Option: TO_LOWER,
		Data: msg,
	})
	fmt.Println(msg)
}

func handleServer(conn nt.Conn, pack *nt.Package) {
	nt.Handle(TO_UPPER, conn, pack, handleToUpper)
	nt.Handle(TO_LOWER, conn, pack, handleToLower)
}

func handleToUpper(pack *nt.Package) string {
	return strings.ToUpper(pack.Data)
}

func handleToLower(pack *nt.Package) string {
	return strings.ToLower(pack.Data)
}