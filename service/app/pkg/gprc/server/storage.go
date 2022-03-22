package grpcserver

import (
	"context"
	"github.com/vladazn/dhq/proto/gen/go/proto/storage"
	"github.com/vladazn/dhq/service/app/domain"
	"github.com/vladazn/dhq/service/app/service"
)

const defaultUserId = 1

type StorageServer struct {
	storage.UnimplementedStorageServer
	services *service.Services
}

func newStorageServer(services *service.Services) *StorageServer {
	return &StorageServer{services: services}
}

func (s *StorageServer) Create(ctx context.Context, req *storage.CreateRequest) (*storage.
	SuccessResponse, error) {

	success, err := s.services.Storage.Create(ctx, defaultUserId, &domain.Answer{
		Key:   req.Key,
		Value: req.Value,
	})

	if err != nil {
		return &storage.SuccessResponse{
			Error: &storage.Error{Msg: err.Error()},
		}, nil
	}

	return &storage.SuccessResponse{
		Result: &storage.SuccessResult{Success: success},
	}, nil
}

func (s *StorageServer) Update(ctx context.Context, req *storage.UpdateRequest) (*storage.
	SuccessResponse, error) {

	success, err := s.services.Storage.Update(ctx, defaultUserId, &domain.Answer{
		Key:   req.Key,
		Value: req.Value,
	})

	if err != nil {
		return &storage.SuccessResponse{
			Error: &storage.Error{Msg: err.Error()},
		}, nil
	}

	return &storage.SuccessResponse{
		Result: &storage.SuccessResult{Success: success},
	}, nil
}

func (s *StorageServer) Delete(ctx context.Context, req *storage.DeleteRequest) (*storage.
	SuccessResponse, error) {

	success, err := s.services.Storage.Delete(ctx, defaultUserId, req.Key)

	if err != nil {
		return &storage.SuccessResponse{
			Error: &storage.Error{Msg: err.Error()},
		}, nil
	}

	return &storage.SuccessResponse{
		Result: &storage.SuccessResult{Success: success},
	}, nil
}

func (s *StorageServer) Get(ctx context.Context, req *storage.GetRequest) (*storage.GetResponse, error) {

	data, err := s.services.Storage.Get(ctx, defaultUserId, req.Key)

	if err != nil {
		return &storage.GetResponse{
			Error: &storage.Error{Msg: err.Error()},
		}, nil
	}

	return &storage.GetResponse{
		Result: &storage.GetResult{
			Data: &storage.Answer{
				Key:   data.Key,
				Value: data.Value,
			},
		},
	}, nil
}

func (s *StorageServer) History(ctx context.Context, req *storage.HistoryRequest) (*storage.
	HistoryResponse, error) {

	data, err := s.services.Storage.History(ctx, defaultUserId, req.Key)

	if err != nil {
		return &storage.HistoryResponse{
			Error: &storage.Error{Msg: err.Error()},
		}, nil
	}

	arr := make([]*storage.Action, len(data))

	for i, d := range arr {
		arr[i] = &storage.Action{
			Event: d.Event,
		}
		if d.Data != nil {
			arr[i].Data = &storage.Answer{
				Key:   d.Data.Key,
				Value: d.Data.Value,
			}
		}
	}

	return &storage.HistoryResponse{
		Result: &storage.HistoryResult{
			Data: arr,
		},
	}, nil
}
