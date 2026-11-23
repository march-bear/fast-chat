package main

import (
	"context"
	"fast-chat/internal/protocol"
)

func handleGetHistoryReq(ctx context.Context, req protocol.GetHistoryReq) protocol.Response {
	messages, err := app.ChatServ.GetHistory(req.From, req.To)
	if err != nil {
		return protocol.MakeErrorResp(err)
	}

	hRep := protocol.HistoryResp{From: req.From, To: req.From + len(messages), Messages: messages}
	return protocol.MakeResponse(protocol.TypeHistoryResp, hRep)
}

func handleNewMessageReq(ctx context.Context, req protocol.NewMessageReq) protocol.Response {
	message, err := app.ChatServ.AddMessage(req.Text, ctx.Value("author_id").(int))
	if err != nil {
		return protocol.MakeErrorResp(err)
	}

	mRep := protocol.MessageOkResp{Message: message}
	return protocol.MakeResponse(protocol.TypeMessageOkResp, mRep)
}
