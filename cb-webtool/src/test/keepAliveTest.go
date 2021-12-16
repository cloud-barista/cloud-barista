package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	var LadyBugURL = os.Getenv("LADYBUG_URL")
	// var TumbleBugURL = os.Getenv("TUMBLE_URL")
	//LadyBugURL = "52.197.30.208:8000"
	LadyBugURL = "52.197.30.209:8000"
	// var TumbleBugURL = os.Getenv("TUMBLE_URL")
	fmt.Println(LadyBugURL)
	// fmt.Println(TumbleBugURL)

	// conn, _ := net.Dial("tcp", "52.197.30.208:1323")
	//conn, _ := net.Dial("tcp", LadyBugURL+"/healthy")
	conn, connErr := net.Dial("tcp", LadyBugURL)
	if connErr != nil {
		fmt.Println(connErr)
		return
	}

	// conn, _ := net.Dial("tcp", TumbleBugURL)
	//
	err := conn.(*net.TCPConn).SetKeepAlive(true) // conn 이 nil 이면 panic으로 빠짐.
	if err != nil {
		fmt.Println(err)
		return
	}

	err = conn.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	notify := make(chan error)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				notify <- err
				if io.EOF == err {
					close(notify)
					return
				}
			}
			if n > 0 {
				fmt.Println("unexpected data : %s", buf[:n])
			}
		}
	}()

	for {
		select {
		case err := <-notify:
			fmt.Println("Connection dropped message", err)
			if err == io.EOF {
				fmt.Println("connection to server was closed")
				return
			}
			break
		case <-time.After(time.Second * 1):
			fmt.Println("timeout 1 , still alive")
		}
	}
}
