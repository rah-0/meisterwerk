package main

import (
	"errors"
	"sync"
	"time"

	"github.com/rah-0/meisterwerk/model"
)

type LanguageStore struct {
	mu    sync.RWMutex
	items map[string]model.Language
}

func NewLanguageStore() *LanguageStore {
	return &LanguageStore{
		items: make(map[string]model.Language),
	}
}

func (s *LanguageStore) Insert(l model.Language) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[l.Uuid]; exists {
		return errors.New("language already exists")
	}

	l.FirstInsert = time.Now().Truncate(time.Microsecond)
	s.items[l.Uuid] = l
	return nil
}

func (s *LanguageStore) Get(uuid string) (model.Language, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	l, ok := s.items[uuid]
	if !ok {
		return model.Language{}, errors.New("language not found")
	}
	return l, nil
}

func (s *LanguageStore) List() []model.Language {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.Language, 0, len(s.items))
	for _, l := range s.items {
		out = append(out, l)
	}
	return out
}

func (s *LanguageStore) Update(uuid string, updated model.Language) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.items[uuid]
	if !exists {
		return errors.New("language not found")
	}

	updated.FirstInsert = current.FirstInsert // preserve insert timestamp
	updated.LastUpdate = time.Now().Truncate(time.Microsecond)
	s.items[uuid] = updated
	return nil
}

func (s *LanguageStore) Delete(uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[uuid]; !exists {
		return errors.New("language not found")
	}
	delete(s.items, uuid)
	return nil
}
