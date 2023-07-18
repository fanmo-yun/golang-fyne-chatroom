package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
)

var clientMap = make(map[string]*net.TCPConn, 20)
var lock sync.Mutex

func process(conn *net.TCPConn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		var buffer [1024]byte
		size, readerr := reader.Read(buffer[:])
		if readerr != nil {
			log.Println(readerr)
			break
		}
		if err := broadcast(buffer[:], size); err != nil {
			log.Println(err)
			break
		}
	}
	deleteuser(conn.RemoteAddr().String())
}

func broadcast(byteMsg []byte, size int) error {
	for _, conn := range clientMap {
		if _, writeerr := conn.Write(byteMsg[:size]); writeerr != nil {
			return writeerr
		}
	}
	return nil
}

func deleteuser(username string) {
	lock.Lock()
	delete(clientMap, username)
	lock.Unlock()
}

func main() {
	tcpAdr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8085")
	server, err := net.ListenTCP("tcp", tcpAdr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer server.Close()

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("coming: ", conn.RemoteAddr().String())
		lock.Lock()
		clientMap[conn.RemoteAddr().String()] = conn
		lock.Unlock()
		go process(conn)
	}
}
