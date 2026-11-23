package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	ws "github.com/gorilla/websocket"

	"fast-chat/internal/chat"
	"fast-chat/internal/protocol"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var app = struct {
	Mu              sync.Mutex
	ChatServ        chat.Service
	AuthorIdCounter int
}{ChatServ: chat.MakeService(), AuthorIdCounter: 1}

func handleRequest(ctx context.Context, req protocol.Request) ([]byte, error) {
	var resp protocol.Response

	switch req.Type {
	case protocol.TypeGetHistoryReq:
		var ghReq protocol.GetHistoryReq
		err := json.Unmarshal(req.Payload, &ghReq)
		if err != nil {
			resp = protocol.MakeErrorResp(err)
		}
		resp = handleGetHistoryReq(ctx, ghReq)
	case protocol.TypeNewMessageReq:
		var nmReq protocol.NewMessageReq
		err := json.Unmarshal(req.Payload, &nmReq)
		if err != nil {
			resp = protocol.MakeErrorResp(err)
		}
		resp = handleNewMessageReq(ctx, nmReq)
	default:
		resp = protocol.MakeErrorResp(errors.New("Unhandled request"))
	}

	data, err := json.Marshal(resp)
	return data, err
}

func handleWSTextMessage(ctx context.Context, payload []byte) ([]byte, error) {
	var req protocol.Request
	err := json.Unmarshal(payload, &req)
	var data []byte

	if err != nil {
		data, err = json.Marshal(protocol.MakeErrorResp(err))
	} else {
		data, err = handleRequest(ctx, req)
	}

	return data, err
}

func handleWSMessage(ctx context.Context, messageType int, payload []byte) (respMessageType int, respPayload []byte, err error) {
	switch messageType {
	case ws.TextMessage:
		respPayload, err = handleWSTextMessage(ctx, payload)
		return ws.TextMessage, respPayload, err
	default:
		return ws.TextMessage, nil, fmt.Errorf("Unsupported messageType: %c", messageType)
	}
}

func handleConn(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	app.Mu.Lock()
	ctx := context.WithValue(context.Background(), "author_id", app.AuthorIdCounter)
	app.AuthorIdCounter++
	app.Mu.Unlock()

	log.Printf("New connection (author_id: %d)", ctx.Value("author_id").(int))

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		respMessageType, respP, err := handleWSMessage(ctx, messageType, p)
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(respMessageType, respP); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", handleConn)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
