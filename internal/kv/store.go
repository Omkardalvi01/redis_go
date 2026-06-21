package kv

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"
)

// Store owns the in-memory key/value map and the lock that protects it.
// Keeping the mutex next to the data makes the concurrency contract obvious.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Delete(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	deleted := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok {
			delete(s.data, key)
			deleted++
		}
	}
	return deleted
}

func (s *Store) Exists(keys ...string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok {
			count++
		}
	}
	return count
}

func (s *Store) Rename(oldKey, newKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[oldKey]
	if !ok {
		return "", fmt.Errorf("(error) no such key")
	}

	delete(s.data, oldKey)
	s.data[newKey] = value
	return value, nil
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Store) Keys(pattern string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Exact lookup is the fast path and preserves the old behavior for non-wildcard input.
	if !strings.ContainsAny(pattern, "*?[]") {
		if _, ok := s.data[pattern]; ok {
			return []string{pattern}, nil
		}
		return nil, nil
	}

	matched := make([]string, 0)
	for key := range s.data {
		ok, err := path.Match(pattern, key)
		if err != nil {
			return nil, fmt.Errorf("(error) invalid key pattern %q: %w", pattern, err)
		}
		if ok {
			matched = append(matched, key)
		}
	}
	sort.Strings(matched)
	return matched, nil
}
