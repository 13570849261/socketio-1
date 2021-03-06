package engine_test

import (
	"log"
	"net/http"
	"time"

	"github.com/zyxar/socketio/engine"
)

func ExampleDial() {
	c, err := engine.Dial("ws://localhost:8080/engine.io/", nil, engine.WebsocketTransport)
	if err != nil {
		log.Printf("dial err=%s", err)
		return
	}
	defer c.Close()
	log.Printf("id=%s\n", c.Sid())
}

func ExampleServer() {
	server, _ := engine.NewServer(time.Second*5, time.Second*5, func(so *engine.Socket) {
		log.Println("connect", so.RemoteAddr())
	})
	server.On(engine.EventMessage, engine.Callback(func(so *engine.Socket, typ engine.MessageType, data []byte) {
		switch typ {
		case engine.MessageTypeString:
			log.Printf("txt: %s\n", data)
		case engine.MessageTypeBinary:
			log.Printf("bin: %x\n", data)
		default:
			log.Printf("???: %x\n", data)
		}
	}))
	server.On(engine.EventPing, engine.Callback(func(so *engine.Socket, _ engine.MessageType, _ []byte) {
		log.Printf("socket ping\n")
	}))
	server.On(engine.EventClose, engine.Callback(func(so *engine.Socket, _ engine.MessageType, _ []byte) {
		log.Printf("socket close\n")
	}))
	server.On(engine.EventUpgrade, engine.Callback(func(so *engine.Socket, _ engine.MessageType, _ []byte) {
		log.Printf("socket upgrade\n")
	}))
	defer server.Close()
	log.Fatalln(http.ListenAndServe("localhost:8081", server))
}
