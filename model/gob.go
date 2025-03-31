package model

import (
	"bytes"
	"encoding/gob"
	"sync"
)

func init() {
	preloadGob(&Language{})
	preloadGob(&LanguageKey{})
	preloadGob(&LanguageValue{})
}

var (
	mu      sync.Mutex
	buffer  = new(bytes.Buffer)
	encoder = gob.NewEncoder(buffer)
	decoder = gob.NewDecoder(buffer)
)

func Encode(a any) error {
	mu.Lock()
	defer mu.Unlock()
	if err := encoder.Encode(a); err != nil {
		return err
	}
	return nil
}

func Decode(a any) error {
	mu.Lock()
	defer mu.Unlock()
	return decoder.Decode(a)
}

func BufferReset() {
	mu.Lock()
	defer mu.Unlock()
	buffer.Reset()
}

func preloadGob(x any) {
	gob.Register(x)
	if err := Encode(x); err != nil {
		panic("failed to encode type metadata: " + err.Error())
	}
	if err := Decode(x); err != nil {
		panic("failed to decode type metadata: " + err.Error())
	}
	BufferReset()
}
