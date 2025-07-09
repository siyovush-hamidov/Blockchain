package network

import (
	"encoding/json"
	"net"
	"strings"
	"time"
)

type Package struct
{
	Option int
	Data string
}

const (
	ENDBYTES = "\000\005\007\001\001\007\005\000"
)

const(
	WAITTIME = 5 // seconds
	DMAXSIZE = (2 << 20) // 2^20*2 = 2MiB
	BUFFSIZE = 4 << 10 // 2^10*4 = 4KiB
)

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

func SerializePackage(pack *Package) string{
	jsonData, err := json.MarshalIndent(*pack, "", "\t")
	if err != nil{
		return ""
	}
	return string(jsonData)
}

func readPackage(conn net.Conn) *Package {
	var(
		data string
		size = uint64(0)
		buffer = make([]byte, BUFFSIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil{
			return nil
		}
		size += uint64(length)
		if size > DMAXSIZE{
			return nil
		}
		data += string(buffer[:length])
		if strings.Contains(data, ENDBYTES) {
			data = strings.Split(data, ENDBYTES)[0]
			break
		}
	}
	return DeserializePackage(data)
}

func DeserializePackage(data string) *Package{
	var pack Package
	err := json.Unmarshal([]byte(data), &pack)
	if err != nil{
		return nil
	}
	return &pack
}