package protocol

import (
	"encoding/json"
)

type RequestType int

const (
	TypeGetHistoryReq RequestType = iota
	TypeNewMessageReq
	TypeSetAuthorReq
	TypeBadReq
)

type Request struct {
	ReqId   int             `json:"req_id"`
	Type    RequestType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type GetHistoryReq struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type NewMessageReq struct {
	Text string `json:"message"`
}

type SetAuthorReq struct {
	Author string `json:"author"`
}

func MakeRequest(t RequestType, resp any) Request {
	payload, _ := json.Marshal(resp)

	return Request{Type: t, Payload: payload}
}
