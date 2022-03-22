package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladazn/dhq/service/app/domain"
	"github.com/vladazn/dhq/service/app/pkg/cache"
	"github.com/vladazn/dhq/service/app/pkg/redis"
	"github.com/vladazn/dhq/service/app/repository"
)

const (
	historyCachePrefix = "history"
	actionCachePrefix  = "action"
	defaultUserId      = 1
)

type StorageService struct {
	repo         *repository.Repositories
	historyCache cache.Cache
	actionCache  cache.Cache
}

func newStorageService(repo *repository.Repositories, r redis.Redis) *StorageService {
	return &StorageService{
		repo:         repo,
		historyCache: cache.NewCacheService(r.GetClient(), historyCachePrefix),
		actionCache:  cache.NewCacheService(r.GetClient(), actionCachePrefix),
	}
}

func (s *StorageService) Create(ctx context.Context, userId int, answer *domain.Answer) (bool,
	error) {

	key := answer.Key

	action, err := s.getLastAction(ctx, userId, key)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	if action != nil {
		return false, errors.New("answer already exists")
	}

	err = s.clearCache(ctx, userId, key)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	lastAction := &domain.Action{
		Event: "create",
		Data:  answer,
	}

	err = s.addAction(ctx, userId, key, lastAction)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	return true, nil
}

func (s *StorageService) Delete(ctx context.Context, userId int, key string) (bool, error) {

	action, err := s.getLastAction(ctx, userId, key)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	if action == nil {
		return false, errors.New("answer does not exist")
	}

	err = s.clearCache(ctx, userId, key)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	lastAction := &domain.Action{
		Event: "delete",
	}

	err = s.addAction(ctx, userId, key, lastAction)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	return true, nil
}

func (s *StorageService) Update(ctx context.Context, userId int, answer *domain.Answer) (bool,
	error) {

	key := answer.Key

	action, err := s.getLastAction(ctx, userId, key)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	if action == nil {
		return false, errors.New("answer does not exist")
	}

	err = s.clearCache(ctx, userId, key)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	lastAction := &domain.Action{
		Event: "update",
		Data:  answer,
	}

	err = s.addAction(ctx, userId, key, lastAction)
	if err != nil {
		return false, errors.New("something went wrong")
	}

	return true, nil
}

func (s *StorageService) Get(ctx context.Context, userId int, key string) (*domain.Answer, error) {
	action, err := s.getLastAction(ctx, userId, key)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	if action == nil {
		return nil, errors.New("does not exist")
	}

	return action, nil
}

func (s *StorageService) History(ctx context.Context, userId int, key string) ([]domain.Action,
	error) {
	cKey := fmt.Sprintf("%v:%v", userId, key)

	v, ex, err := s.historyCache.Get(ctx, cKey)

	if err != nil {
		return nil, err
	}

	if !ex {
		history, err := s.repo.Answers.GetHistory(ctx, userId, key)
		if err != nil {
			return nil, err
		}

		b, _ := json.Marshal(history)

		err = s.historyCache.Set(ctx, cKey, string(b))
		if err != nil {
			fmt.Println(err)
		}

		return history, nil
	}

	r := make([]domain.Action, 0)

	err = json.Unmarshal([]byte(v), &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *StorageService) addAction(ctx context.Context, userId int, key string,
	lastAction *domain.Action) error {

	err := s.repo.Answers.AddAction(ctx, userId, lastAction)
	if err != nil {
		return err
	}

	err = s.saveLastActionCache(ctx, userId, lastAction, key)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (s *StorageService) getLastAction(ctx context.Context, userId int,
	key string) (*domain.Answer, error) {

	cKey := fmt.Sprintf("%v:%v", userId, key)

	v, ex, err := s.actionCache.Get(ctx, cKey)

	if err != nil {
		return nil, err
	}

	if !ex {
		action, err := s.repo.Answers.GetLastAction(ctx, userId, key)
		if err != nil {
			return nil, err
		}

		if action == nil {
			return nil, nil
		}

		b, _ := json.Marshal(action)

		err = s.actionCache.Set(ctx, cKey, string(b))
		if err != nil {
			fmt.Println(err)
		}

		return action.Data, nil
	}

	r := &domain.Action{}

	err = json.Unmarshal([]byte(v), r)
	if err != nil {
		return nil, err
	}

	return r.Data, nil
}

func (s *StorageService) clearCache(ctx context.Context, userId int, key string) error {
	cKey := fmt.Sprintf("%v:%v", userId, key)
	err := s.actionCache.Clear(ctx, cKey)
	if err != nil {
		return err
	}

	err = s.historyCache.Clear(ctx, cKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) saveLastActionCache(ctx context.Context, userId int,
	action *domain.Action, key string) error {

	cKey := fmt.Sprintf("%v:%v", userId, key)

	b, _ := json.Marshal(action)

	err := s.actionCache.Set(ctx, cKey, string(b))
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
