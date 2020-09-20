package main

import (
	"fmt"
	"sync"

	"tcpudp-server/src/comm/mylog"
	"tcpudp-server/src/server"
)

func main() {
	mylog.Init("./log/s.log", "debug", 21600)
	var wg sync.WaitGroup

	ts := server.NewTCP("localhost:9999")
	ts.OnNewClient(func(c *server.TCPClient) {
		mylog.LogDebug("New client connected: %s", c.Conn().RemoteAddr().String())
	})
	ts.OnNewMessage(func(c *server.TCPClient, m string) {
		fmt.Printf("recieve message: %s", m)
		if err := c.Send(m); err != nil {
			mylog.LogError("send message failed: %v", err)
		} else {
			mylog.LogDebug("send message success: %s", m)
		}
	})
	go func() {
		wg.Add(1)
		ts.Listen()
	}()

	us := server.NewUDP("localhost:6666", 1024)
	go func() {
		wg.Add(1)
		for {
			pkg := <-us.ChRead
			us.ChWrite <- pkg
		}
	}()

	wg.Wait()
}
