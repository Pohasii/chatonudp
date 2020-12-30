package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var secret string = "qwerty"

func main() {

	sAddr := "192.168.0.52:55442"
	serverAddr, err := net.ResolveUDPAddr("udp", sAddr)
	listener, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %serverAddr\n", listener.RemoteAddr().String())
	defer listener.Close()

	stop := true

	go reader(listener)

	for stop {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		text = strings.TrimSuffix(text, "\n")
		if text == "1000" {
			for i:= 0; i < 1000; i++ {
				data := []byte(secret+"Message number: "+strconv.Itoa(i))
				_, err = listener.Write(data)
			}
		}
		data := []byte(secret+text)

		_, err = listener.Write(data)

		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("off UDP client!")
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func reader(listener *net.UDPConn) {

	buffer := make([]byte, 10240)
	for {
		n, _, err := listener.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	}
}