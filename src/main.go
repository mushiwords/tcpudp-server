package main

import (
	"fmt"
	"sync"

	"tcpudp-server/src/comm/mylog"
	"tcpudp-server/src/server"
)

func main() {
	if err := mylog.Init("./log/s.log", "debug", 21600);err != nil {
		fmt.Printf("init log failed: %w\n", err)
	}

	var wg sync.WaitGroup

	ts := server.NewTCP("localhost:9999")
	ts.OnNewClient(func(c *server.TCPClient) {
		mylog.LogDebug("New client connected: %s", c.Conn().RemoteAddr().String())
	})
	ts.OnNewMessage(func(c *server.TCPClient, m string) {
		mylog.LogDebug("recieve message: %s", m)
		if err := c.Send(m); err != nil {
			mylog.LogError("send message failed: %v", err)
		} else {
			mylog.LogDebug("send message success: %s", m)
		}
	})
	go func() {
		fmt.Printf("start tcp server [:9999] success \n")
		wg.Add(1)
		ts.Listen()
	}()

	us := server.NewUDP("localhost:6666", 1024)
	go func() {
		fmt.Printf("start udp server [:6666] success \n")

		wg.Add(1)
		for {
			pkg := <-us.ChRead
			us.ChWrite <- pkg
		}
	}()

	mylog.LogDebug("service started.")

	wg.Wait()
}
