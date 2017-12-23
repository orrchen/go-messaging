package lib

import (
	"log"
	"net"
	"sync"
	"github.com/satori/go.uuid"
	"errors"
	"fmt"
)



// TCP Server
type TcpServer struct {
	address                  string // Address to open connection: localhost:9999
	connLock sync.RWMutex
	connections map[string]*Client
	callbacks Callbacks
	listener net.Listener
}

func (s *TcpServer) onConnectionEvent(c *Client,eventType ConnectionEventType, e error ) {
	switch eventType {
	case CONNECTION_EVENT_TYPE_NEW_CONNECTION:
		s.connLock.Lock()
		u1 := uuid.NewV4()
		uidString := u1.String()
		c.Uid = uidString
		s.connections[uidString] = c
		s.connLock.Unlock()
		//log.Println(eventType , " ,  uid:", c.Uid, " , ip: ", c.Conn().RemoteAddr().String())
		if s.callbacks.OnNewConnection != nil {
			s.callbacks.OnNewConnection(uidString)
		}
	case CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED, CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR:
		//log.Println(eventType , " ,  uid:", c.Uid, " , ip: ", c.Conn().RemoteAddr().String(), " , error: ", e.Error())
		s.connLock.Lock()
		delete(s.connections,c.Uid)
		s.connLock.Unlock()
		if s.callbacks.OnConnectionTerminated!=nil {
			s.callbacks.OnConnectionTerminated(c.Uid)
		}
	}
}

func (s *TcpServer) onDataEvent(c *Client, data []byte) {
	//log.Println("onDataEvent, ", c.Conn().RemoteAddr().String(), " data: " , string(data))
	if s.callbacks.OnDataReceived!=nil {
		s.callbacks.OnDataReceived(c.Uid, data)
	}
}


// Start network Server
func (s *TcpServer) Listen() {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP Server.: " , err)
	}
	for {
		conn, _ := s.listener.Accept()
		client := NewClient(conn,s.onConnectionEvent,s.onDataEvent)
		s.onConnectionEvent(client, CONNECTION_EVENT_TYPE_NEW_CONNECTION,nil)
		go client.listen()

	}
}

// Creates new tcp Server instance
func NewServer(address string, callbacks Callbacks ) *TcpServer {
	log.Println("Creating Server with address", address)
	s := &TcpServer{
		address: address,
		callbacks: callbacks,
	}
	s.connections = make(map[string]*Client)
	return s
}

func (s *TcpServer) SendDataByClientId(clientUid string, data []byte) error{
	if s.connections[clientUid]!=nil {
		return s.connections[clientUid].Send(data)
	} else {
		return errors.New(fmt.Sprint("no connection with uid ", clientUid))
	}

	return nil
}

func (s *TcpServer) SendDataByDeviceUid(deviceUid string, data []byte) error{
	for k := range s.connections {
		if s.connections[k].DeviceUid == deviceUid {
			return s.connections[k].Send(data)
		}
	}
	return errors.New(fmt.Sprint("no connection with deviceUid ", deviceUid))
}

func (s *TcpServer) Close(){
	log.Println("TcpServer.Close()")
	log.Println("s.connections length: " , len(s.connections))
	for k := range s.connections {
		fmt.Printf("key[%s]\n", k)
		s.connections[k].Close()
	}
	s.listener.Close()
}

func (s *TcpServer) SetDeviceUidToClient(clientUid string, deviceUid string) error{
	if s.connections[clientUid]!=nil {
		s.connections[clientUid].DeviceUid = deviceUid
		return nil
	} else {
		return errors.New(fmt.Sprint("no connection with uid ", clientUid))
	}
}


