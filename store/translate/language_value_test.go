package main

import (
	"testing"

	"github.com/google/uuid"

	"github.com/rah-0/meisterwerk/model"
)

func TestLanguageValueStore_InsertAndGet(t *testing.T) {
	store := NewLanguageValueStore()

	id := uuid.NewString()
	langID := uuid.NewString()
	keyID := uuid.NewString()
	val := "Hallo"

	v := model.LanguageValue{
		Uuid:            id,
		UuidLanguage:    langID,
		UuidLanguageKey: keyID,
		Value:           val,
	}

	if err := store.Insert(v); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	got, err := store.Get(id)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Uuid != id || got.Value != val {
		t.Errorf("Retrieved value mismatch: got %+v, want %+v", got, v)
	}
}

func TestLanguageValueStore_InsertDuplicate(t *testing.T) {
	store := NewLanguageValueStore()

	id := uuid.NewString()
	v := model.LanguageValue{Uuid: id, Value: "Hello"}

	if err := store.Insert(v); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if err := store.Insert(v); err == nil {
		t.Error("Expected error on duplicate insert, got nil")
	}
}

func TestLanguageValueStore_Get_NotFound(t *testing.T) {
	store := NewLanguageValueStore()

	_, err := store.Get("not-found")
	if err == nil {
		t.Error("Expected error on Get for nonexistent value")
	}
}

func TestLanguageValueStore_List(t *testing.T) {
	store := NewLanguageValueStore()

	for i := 0; i < 3; i++ {
		store.Insert(model.LanguageValue{
			Uuid:            uuid.NewString(),
			UuidLanguage:    uuid.NewString(),
			UuidLanguageKey: uuid.NewString(),
			Value:           "val",
		})
	}

	list, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("Expected 3 items in List, got %d", len(list))
	}
}

func TestLanguageValueStore_List_Empty(t *testing.T) {
	store := NewLanguageValueStore()

	_, err := store.List()
	if err == nil {
		t.Error("Expected error when listing empty store")
	}
}

func TestLanguageValueStore_Update(t *testing.T) {
	store := NewLanguageValueStore()

	id := uuid.NewString()
	initial := model.LanguageValue{
		Uuid:            id,
		UuidLanguage:    "lang1",
		UuidLanguageKey: "key1",
		Value:           "Old",
	}

	if err := store.Insert(initial); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	updated := model.LanguageValue{
		Uuid:            id,
		UuidLanguage:    "lang1",
		UuidLanguageKey: "key1",
		Value:           "New",
	}

	if err := store.Update(id, updated); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ := store.Get(id)
	if got.Value != "New" {
		t.Errorf("Expected updated value to be 'New', got '%s'", got.Value)
	}
}

func TestLanguageValueStore_Update_NotFound(t *testing.T) {
	store := NewLanguageValueStore()

	err := store.Update("missing-id", model.LanguageValue{Uuid: "missing-id", Value: "Whatever"})
	if err == nil {
		t.Error("Expected error on Update for non-existent value")
	}
}

func TestLanguageValueStore_Delete(t *testing.T) {
	store := NewLanguageValueStore()

	id := uuid.NewString()
	store.Insert(model.LanguageValue{
		Uuid:            id,
		UuidLanguage:    "lang",
		UuidLanguageKey: "key",
		Value:           "DeleteMe",
	})

	if err := store.Delete(id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, err := store.Get(id); err == nil {
		t.Error("Expected error after deletion")
	}
}

func TestLanguageValueStore_Delete_NotFound(t *testing.T) {
	store := NewLanguageValueStore()

	err := store.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error when deleting nonexistent item")
	}
}
