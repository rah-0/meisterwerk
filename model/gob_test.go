package model

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGobConcurrentEncodingDecoding(t *testing.T) {
	var goroutines = runtime.NumCPU()
	var wg sync.WaitGroup

	lang := Language{
		Uuid:   uuid.NewString(),
		Prefix: "en-US",
		Lang:   "English",
	}

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			if err := Encode(lang); err != nil {
				t.Errorf("goroutine %d: encode failed: %v", id, err)
				return
			}

			var decoded Language
			if err := Decode(&decoded); err != nil {
				t.Errorf("goroutine %d: decode failed: %v", id, err)
				return
			}

			if decoded.Prefix != lang.Prefix || decoded.Lang != lang.Lang {
				t.Errorf("goroutine %d: decoded mismatch: got %+v, want %+v", id, decoded, lang)
			}
		}(i)
	}

	wg.Wait()
}

func TestGobEncodingDecodingLanguage(t *testing.T) {
	id := uuid.NewString()
	now := time.Now().Truncate(time.Microsecond)
	prefix := "en-US"
	lang := "English"
	title := "English"
	img := "/static/img/flags/us.png"
	months := "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec"

	original := &Language{
		Uuid:        id,
		FirstInsert: now,
		LastUpdate:  now,
		Prefix:      prefix,
		Lang:        lang,
		Title:       title,
		Img:         img,
		MonthsShort: months,
	}

	if err := Encode(original); err != nil {
		t.Fatalf("encoding failed: %v", err)
	}

	var decoded Language
	if err := Decode(&decoded); err != nil {
		t.Fatalf("decoding failed: %v", err)
	}

	BufferReset()

	if decoded.Uuid != id {
		t.Errorf("Uuid mismatch: got %s, want %s", decoded.Uuid, id)
	}
	if !decoded.FirstInsert.Equal(now) {
		t.Errorf("FirstInsert mismatch: got %v, want %v", decoded.FirstInsert, now)
	}
	if !decoded.LastUpdate.Equal(now) {
		t.Errorf("LastUpdate mismatch: got %v, want %v", decoded.LastUpdate, now)
	}
	if decoded.Prefix != prefix {
		t.Errorf("Prefix mismatch: got %s, want %s", decoded.Prefix, prefix)
	}
	if decoded.Lang != lang {
		t.Errorf("Lang mismatch: got %s, want %s", decoded.Lang, lang)
	}
	if decoded.Title != title {
		t.Errorf("Title mismatch: got %s, want %s", decoded.Title, title)
	}
	if decoded.Img != img {
		t.Errorf("Img mismatch: got %s, want %s", decoded.Img, img)
	}
	if decoded.MonthsShort != months {
		t.Errorf("MonthsShort mismatch: got %s, want %s", decoded.MonthsShort, months)
	}
}

func TestGobEncodingDecodingLanguageKey(t *testing.T) {
	id := uuid.NewString()
	now := time.Now().Truncate(time.Microsecond)
	value := "quote.confirmation_title"

	original := &LanguageKey{
		Uuid:        id,
		FirstInsert: now,
		LastUpdate:  now,
		Value:       value,
	}

	if err := Encode(original); err != nil {
		t.Fatalf("encoding failed: %v", err)
	}

	var decoded LanguageKey
	if err := Decode(&decoded); err != nil {
		t.Fatalf("decoding failed: %v", err)
	}

	BufferReset()

	if decoded.Uuid != id {
		t.Errorf("Uuid mismatch: got %s, want %s", decoded.Uuid, id)
	}
	if !decoded.FirstInsert.Equal(now) {
		t.Errorf("FirstInsert mismatch: got %v, want %v", decoded.FirstInsert, now)
	}
	if !decoded.LastUpdate.Equal(now) {
		t.Errorf("LastUpdate mismatch: got %v, want %v", decoded.LastUpdate, now)
	}
	if decoded.Value != value {
		t.Errorf("Value mismatch: got %s, want %s", decoded.Value, value)
	}
}

func TestGobEncodingDecodingLanguageValue(t *testing.T) {
	id := uuid.NewString()
	now := time.Now().Truncate(time.Microsecond)
	langUUID := uuid.NewString()
	keyUUID := uuid.NewString()
	value := "Angebot bestätigt"

	original := &LanguageValue{
		Uuid:            id,
		FirstInsert:     now,
		LastUpdate:      now,
		UuidLanguage:    langUUID,
		UuidLanguageKey: keyUUID,
		Value:           value,
	}

	if err := Encode(original); err != nil {
		t.Fatalf("encoding failed: %v", err)
	}

	var decoded LanguageValue
	if err := Decode(&decoded); err != nil {
		t.Fatalf("decoding failed: %v", err)
	}

	BufferReset()

	if decoded.Uuid != id {
		t.Errorf("Uuid mismatch: got %s, want %s", decoded.Uuid, id)
	}
	if !decoded.FirstInsert.Equal(now) {
		t.Errorf("FirstInsert mismatch: got %v, want %v", decoded.FirstInsert, now)
	}
	if !decoded.LastUpdate.Equal(now) {
		t.Errorf("LastUpdate mismatch: got %v, want %v", decoded.LastUpdate, now)
	}
	if decoded.UuidLanguage != langUUID {
		t.Errorf("UuidLanguage mismatch: got %s, want %s", decoded.UuidLanguage, langUUID)
	}
	if decoded.UuidLanguageKey != keyUUID {
		t.Errorf("UuidLanguageKey mismatch: got %s, want %s", decoded.UuidLanguageKey, keyUUID)
	}
	if decoded.Value != value {
		t.Errorf("Value mismatch: got %s, want %s", decoded.Value, value)
	}
}
