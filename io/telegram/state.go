package telegram

import (
	"context"
	"encoding/json"
	"financial/infrastructure/cache"
)

type StateManager struct {
	cache cache.ICache
}

type State struct {
	Uid  string
	Data interface{}
}

func (s StateManager) Set(ctx context.Context, state State) error {
	j, _ := json.Marshal(state)

	return s.cache.Add(ctx, s.stateKey(state.Uid), string(j), 0)
}

func (s StateManager) Get(ctx context.Context, uid string) (*State, error) {
	data, err := s.cache.Get(ctx, s.stateKey(uid))
	if err != nil {
		return nil, err
	}
	var state State
	json.Unmarshal([]byte(data), &state)
	return &state, nil
}

func (s StateManager) Delete(ctx context.Context, uid string) error {
	return s.cache.Delete(ctx, s.stateKey(uid))
}

func (s StateManager) stateKey(key string) string {
	return "states" + key
}
