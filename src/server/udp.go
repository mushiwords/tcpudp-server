package server

import (
	"net"
	"runtime"

	"tcpudp-server/src/comm/mylog"
)

/************************************************************
 *  UDP Server
************************************************************/
const chanLen = 10

type (
	// Server main struct
	udpserver struct {
		conn      *net.UDPConn //
		frameSize uint16       //
		ChRead    chan *TPack  //
		ChWrite   chan *TPack  //
	}
	// TPack pack struct
	TPack struct {
		Pack []byte
		Addr *net.UDPAddr
	}
)

// New constructor of a new server
func NewUDP(addr string, frameSize uint16) *udpserver {
	var laddr *net.UDPAddr
	var err error
	if frameSize == 0 {
		frameSize = 65507
	}

	server := &udpserver{
		ChRead:    make(chan *TPack, chanLen),
		ChWrite:   make(chan *TPack, chanLen),
		frameSize: frameSize,
	}

	laddr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		mylog.LogError("ResolveUDPAddr server error: %w", err)
		return nil
	}
	server.conn, err = net.ListenUDP("udp", laddr)
	if err != nil {
		mylog.LogError("ListenUDP error: %w", err)
		return nil
	}

	go server.reader()
	go server.writer()

	// destroy action
	stopAllGorutines := func(t *udpserver) {
		close(t.ChWrite)
		t.conn.Close()
	}
	runtime.SetFinalizer(server, stopAllGorutines)
	return server
}

func (t *udpserver) reader() {
	defer func() {
		mylog.LogDebug("udp server reader exited ")
	}()

	var (
		err error
		n   int
	)

	for {
		pack := &TPack{}
		pack.Pack = make([]byte, t.frameSize)
		n, pack.Addr, err = t.conn.ReadFromUDP(pack.Pack)
		mylog.LogDebug("read from: %s", string(pack.Pack[:n]))
		if err != nil {
			mylog.LogError("read package error: %w", err)
			return
		} else {
			pack.Pack = pack.Pack[:n]
			t.ChRead <- pack
		}
	}
}

func (t *udpserver) writer() {
	defer func() {
		mylog.LogDebug("udp server write exited ")
	}()

	var (
		v   *TPack
		err error
	)
	for v = range t.ChWrite {
		_, err = t.conn.WriteToUDP(v.Pack, v.Addr)
		mylog.LogDebug("read from: %s", string(v.Pack))
		if err != nil {
			mylog.LogError("write package error: %w", err)
		}
	}
}
