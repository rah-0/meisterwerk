package model

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"sync"
)

func init() {
	PreloadGob(Language{})
	PreloadGob([]Language{})
	PreloadGob(LanguageKey{})
	PreloadGob([]LanguageKey{})
	PreloadGob(LanguageValue{})
	PreloadGob([]LanguageValue{})
}

var (
	mu      sync.Mutex
	buffer  = new(bytes.Buffer)
	encoder = gob.NewEncoder(buffer)
	decoder = gob.NewDecoder(buffer)
)

func SetBytes(b []byte) {
	mu.Lock()
	defer mu.Unlock()
	buffer.Reset()
	buffer.Write(b)
}

func GetBytes() []byte {
	mu.Lock()
	defer mu.Unlock()
	return buffer.Bytes()
}

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

func PreloadGob(x any) {
	gob.Register(x)

	if err := Encode(x); err != nil {
		panic("failed to encode type metadata: " + err.Error())
	}

	decoded := reflect.New(reflect.TypeOf(x)).Interface()
	if err := Decode(decoded); err != nil {
		panic("failed to decode type metadata: " + err.Error())
	}

	BufferReset()
}
