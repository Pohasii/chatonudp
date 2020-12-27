package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Server Start")

	// init server
	sAddr := "192.168.0.52:12345" //localhost:0
	adr, err := net.ResolveUDPAddr("udp", sAddr) //192.168.0.52:12345
	if err != nil {
		log.Println(err)
	}

	Conn, err := net.ListenUDP("udp", adr)
	if err != nil {
		log.Println(err)
	}

	echoServer(Conn)
}

func echoServer(Conn *net.UDPConn) {
	defer Conn.Close()

	var sms = make([]byte, 512)

	for {
		size, caddr, err := Conn.ReadFromUDP(sms)
		if err != nil {
			log.Println(err)
		}

		if size > 0 {
			Conn.WriteTo(sms, caddr)
			fmt.Println("caddr: ", caddr, "size: ", size, "mess: ", string(sms[:size]))
		}
	}
}
