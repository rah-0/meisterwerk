package main

import (
	"errors"
	"sync"
	"time"

	"github.com/rah-0/meisterwerk/model"
)

type LanguageValueStore struct {
	mu    sync.RWMutex
	items map[string]model.LanguageValue
}

func NewLanguageValueStore() *LanguageValueStore {
	return &LanguageValueStore{
		items: make(map[string]model.LanguageValue),
	}
}

func (s *LanguageValueStore) Insert(v model.LanguageValue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[v.Uuid]; exists {
		return errors.New("value already exists")
	}

	v.FirstInsert = time.Now().Truncate(time.Microsecond)
	s.items[v.Uuid] = v
	return nil
}

func (s *LanguageValueStore) Get(uuid string) (model.LanguageValue, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.items[uuid]
	return v, ok
}

func (s *LanguageValueStore) List() []model.LanguageValue {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.LanguageValue, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, v)
	}
	return out
}

func (s *LanguageValueStore) Update(uuid string, updated model.LanguageValue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.items[uuid]
	if !exists {
		return errors.New("value not found")
	}

	updated.FirstInsert = current.FirstInsert
	updated.LastUpdate = time.Now().Truncate(time.Microsecond)
	s.items[uuid] = updated
	return nil
}

func (s *LanguageValueStore) Delete(uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[uuid]; !exists {
		return errors.New("value not found")
	}
	delete(s.items, uuid)
	return nil
}
