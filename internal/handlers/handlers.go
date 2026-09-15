package handlers

import (
	"context"

	"github.com/patrickisgreat/pb-go-api-starter/internal/api"
	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

type Quotes struct {
	ServerContext *appserver.ServerContext
}

func (h *Quotes) GetQuote(ctx context.Context, request *proto.GetQuoteRequest) (*proto.GetQuoteResponse, error) {
	return api.GetQuoteHandler(ctx, h.ServerContext, request)
}

func (h *Quotes) GetQuoteFlaky(ctx context.Context, request *proto.GetQuoteFlakyRequest) (*proto.GetQuoteFlakyResponse, error) {
	return api.GetQuoteFlakyHandler(ctx, h.ServerContext, request)
}

func (h *Quotes) GetQuoteLate(ctx context.Context, request *proto.GetQuoteLateRequest) (*proto.GetQuoteLateResponse, error) {
	return api.GetQuoteLateHandler(ctx, h.ServerContext, request)
}
