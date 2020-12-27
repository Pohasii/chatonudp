package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type message struct {
	userAddr *net.UDPAddr
	mes      []byte
}

type Messages chan message

type ValidUser map[string]time.Time // *net.UDPAddr

var secret string = "qwerty"

func reader(Conn *net.UDPConn, MessagesFromUser Messages, verifyUser ValidUser) {
	defer Conn.Close()

	var sms = make([]byte, 10240)

	ticker := time.Tick(5 * time.Millisecond)
	for range ticker {
		size, caddr, err := Conn.ReadFromUDP(sms)
		if err != nil {
			log.Println(err)
		}

		if validation(verifyUser, caddr, &sms, &size); size > 0 {
			MessagesFromUser <- message{caddr, sms[:size]}
			fmt.Println("caddr: ", caddr, "size: ", size, "mess: ", string(sms[:size]))
		}
	}
}

func validation(verifyUser ValidUser, userAddr *net.UDPAddr, envelope *[]byte, size *int) {

	// check user in verified Users or validation secret
	if _, ok := verifyUser[userAddr.String()]; ok || secret == string((*envelope)[:6]) {

		// update time last message
		verifyUser[userAddr.String()] = time.Now()

		// del secrets byte from message
		*envelope = (*envelope)[6:]
	} else {
		*size = 0
	}
}

func sender(Conn *net.UDPConn, MessagesToUser Messages) {
	for letter := range MessagesToUser {
		_, err := Conn.WriteTo(letter.mes, letter.userAddr)
		if err != nil {
			log.Println(err)
		}
	}
}

func handler(verifiedUser ValidUser, MessagesFromUser, MessagesToUser Messages) {
	for mess := range MessagesFromUser {
		test := verifiedUser
		for client := range test {
			if client != mess.userAddr.String() {
				addr, err := net.ResolveUDPAddr("udp", client)
				if err != nil {
					log.Println(err)
				}
				MessagesToUser <- message{addr, mess.mes}
			}
		}
	}
}

func main() {
	fmt.Println("Server Start")

	// init server
	sAddr := "192.168.0.52:12345"                //localhost:0
	adr, err := net.ResolveUDPAddr("udp", sAddr) //192.168.0.52:12345
	if err != nil {
		log.Println(err)
	}

	Conn, err := net.ListenUDP("udp", adr)
	if err != nil {
		log.Println(err)
	}

	Users := make(ValidUser, 0)
	MessagesFromUsers := make(Messages, 5000)
	MessagesToUsers := make(Messages, 5000)

	go handler(Users, MessagesFromUsers, MessagesToUsers)
	go sender(Conn, MessagesToUsers)
	reader(Conn, MessagesFromUsers, Users)
}
