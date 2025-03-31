package main

import (
	"errors"
	"sync"
	"time"

	"github.com/rah-0/meisterwerk/model"
)

type LanguageKeyStore struct {
	mu      sync.RWMutex
	items   map[string]model.LanguageKey
	byValue map[string]string // map[Value]Uuid for uniqueness check
}

func NewLanguageKeyStore() *LanguageKeyStore {
	return &LanguageKeyStore{
		items:   make(map[string]model.LanguageKey),
		byValue: make(map[string]string),
	}
}

func (s *LanguageKeyStore) Insert(k model.LanguageKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[k.Uuid]; exists {
		return errors.New("key already exists")
	}
	if _, exists := s.byValue[k.Value]; exists {
		return errors.New("key value must be unique")
	}

	k.FirstInsert = time.Now().Truncate(time.Microsecond)
	s.items[k.Uuid] = k
	s.byValue[k.Value] = k.Uuid
	return nil
}

func (s *LanguageKeyStore) Get(uuid string) (model.LanguageKey, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	k, ok := s.items[uuid]
	return k, ok
}

func (s *LanguageKeyStore) GetByValue(value string) (model.LanguageKey, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uuid, ok := s.byValue[value]
	if !ok {
		return model.LanguageKey{}, false
	}
	k, exists := s.items[uuid]
	return k, exists
}

func (s *LanguageKeyStore) List() []model.LanguageKey {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.LanguageKey, 0, len(s.items))
	for _, k := range s.items {
		out = append(out, k)
	}
	return out
}

func (s *LanguageKeyStore) Update(uuid string, updated model.LanguageKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.items[uuid]
	if !exists {
		return errors.New("key not found")
	}

	if current.Value != updated.Value {
		if _, exists := s.byValue[updated.Value]; exists {
			return errors.New("key value must be unique")
		}
		delete(s.byValue, current.Value)
		s.byValue[updated.Value] = uuid
	}

	updated.FirstInsert = current.FirstInsert
	updated.LastUpdate = time.Now().Truncate(time.Microsecond)
	s.items[uuid] = updated
	return nil
}

func (s *LanguageKeyStore) Delete(uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k, exists := s.items[uuid]
	if !exists {
		return errors.New("key not found")
	}
	delete(s.items, uuid)
	delete(s.byValue, k.Value)
	return nil
}
