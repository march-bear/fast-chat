package protocol

import (
	"encoding/json"
)

type ResponseType int

type Response struct {
	ReqId   int             `json:"req_id"`
	Type    ResponseType    `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	TypeHistoryResp ResponseType = iota
	TypeMessageOkResp
	TypeOkResp
	TypeErrorResp
	TypeBadResp
)

type HistoryResp struct {
	From     int           `json:"from"`
	To       int           `json:"to"`
	Messages []ChatMessage `json:"messages"`
}

type MessageOkResp struct {
	Message ChatMessage `json:"message"`
}

type OkResp struct{}

type ErrorResp struct {
	ErrorText string `json:"error"`
}

func MakeErrorResp(err error) Response {
	eResp := ErrorResp{ErrorText: err.Error()}
	payload, _ := json.Marshal(eResp)
	return Response{Type: TypeErrorResp, Payload: payload}
}

func MakeResponse(t ResponseType, resp any) Response {
	payload, err := json.Marshal(resp)
	if err != nil {
		return MakeErrorResp(err)
	}

	return Response{Type: t, Payload: payload}
}
