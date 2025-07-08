package network

import (
	"net"
	"time"
)

type Package struct
{
	Option int
	Data string
}

func Send(address string, pack *Package) *Package {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	conn.Write([]byte(SerializePackage(pack) + ENDBYTES))
	var res = new(Package)
	ch := make(chan bool)
	go func() {
		res = readPackage(conn)
		ch <- true
	}()
	select {
		case <-ch:
		case <-time.After(WAITTIME * time.Second):
	}
	return res
}