package http

import (
	"context"

	transport "multienv/agent/client/http"
	storageclient "multienv/agent/storage/client"
	storagehttp "multienv/agent/storage/client/http"
	"multienv/common/codex/model"
	"multienv/common/codex/operations"
	storagemodel "multienv/common/storage/model"
	storageops "multienv/common/storage/operations"
)

const sessionsKey = "@codex"

type httpClient struct {
	storage storageclient.Client
}

func New(config transport.Config) httpClient {
	return httpClient{storage: storagehttp.New(config)}
}

func (c httpClient) SetSessions(
	ctx context.Context,
	req operations.SetSessionsRequest,
) (*operations.SetSessionsResponse, error) {
	at := storagemodel.Location{Key: sessionsKey, Path: []string{req.Hostname}}
	if err := storageclient.SetEncode(ctx, c.storage, at, req.Sessions); err != nil {
		return nil, err
	}

	return &operations.SetSessionsResponse{}, nil
}

func (c httpClient) ListAll(
	ctx context.Context,
	req operations.ListAllRequest,
) (*operations.ListAllResponse, error) {
	at := storagemodel.Location{Key: sessionsKey}
	sessions, _, err := storageclient.GetDecode[model.SessionsByHost](ctx, c.storage, at)
	if err != nil {
		return nil, err
	}

	if sessions == nil {
		sessions = make(model.SessionsByHost)
	}

	return &operations.ListAllResponse{Sessions: sessions}, nil
}

func (c httpClient) SetSession(
	ctx context.Context,
	req operations.SetSessionRequest,
) (*operations.SetSessionResponse, error) {
	at := storagemodel.Location{Key: sessionsKey, Path: []string{req.Hostname, req.Session.Id}}
	if err := storageclient.SetEncode(ctx, c.storage, at, req.Session); err != nil {
		return nil, err
	}

	return &operations.SetSessionResponse{}, nil
}

func (c httpClient) DeleteSession(
	ctx context.Context,
	req operations.DeleteSessionRequest,
) (*operations.DeleteSessionResponse, error) {
	at := storagemodel.Location{Key: sessionsKey, Path: []string{req.Hostname, req.SessionId}}
	if _, err := c.storage.Delete(ctx, storageops.DeleteRequest{Location: at}); err != nil {
		return nil, err
	}

	return &operations.DeleteSessionResponse{}, nil
}
