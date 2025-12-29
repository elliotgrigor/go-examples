package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"
)

type SessionData struct {
	Subject string
	Iat     string
	Exp     string
}

type SessionStore struct {
	Store map[string][]byte
	Mu    sync.RWMutex
}

func (s *SessionStore) Get(key string) (SessionData, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	session, ok := s.Store[key]
	if !ok {
		return SessionData{}, errors.New("Session doesn't exist")
	}
	buf := bytes.NewBuffer(session)
	decoder := gob.NewDecoder(buf)
	var data SessionData
	if err := decoder.Decode(&data); err != nil {
		return SessionData{}, err
	}
	return data, nil
}

func (s *SessionStore) Set(key string, session SessionData) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(session)
	s.Store[key] = buf.Bytes()
}

func (s *SessionStore) Delete(key string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Store, key)
}
