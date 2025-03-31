package main

import (
	"testing"

	"github.com/google/uuid"

	"github.com/rah-0/meisterwerk/model"
)

func TestLanguageKeyStore_InsertAndGet(t *testing.T) {
	store := NewLanguageKeyStore()

	id := uuid.NewString()
	val := "welcome_message"
	key := model.LanguageKey{
		Uuid:  id,
		Value: val,
	}

	if err := store.Insert(key); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	got, ok := store.Get(id)
	if !ok {
		t.Fatal("Get failed: key not found")
	}
	if got.Uuid != id || got.Value != val {
		t.Errorf("Get returned wrong key: got %+v, want %+v", got, key)
	}
}

func TestLanguageKeyStore_InsertDuplicateUuid(t *testing.T) {
	store := NewLanguageKeyStore()

	id := uuid.NewString()
	key := model.LanguageKey{Uuid: id, Value: "foo"}

	if err := store.Insert(key); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if err := store.Insert(key); err == nil {
		t.Error("Expected error on duplicate UUID insert, got nil")
	}
}

func TestLanguageKeyStore_InsertDuplicateValue(t *testing.T) {
	store := NewLanguageKeyStore()

	key1 := model.LanguageKey{Uuid: uuid.NewString(), Value: "shared"}
	key2 := model.LanguageKey{Uuid: uuid.NewString(), Value: "shared"}

	if err := store.Insert(key1); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if err := store.Insert(key2); err == nil {
		t.Error("Expected error on duplicate Value insert, got nil")
	}
}

func TestLanguageKeyStore_GetByValue(t *testing.T) {
	store := NewLanguageKeyStore()

	val := "product.title"
	key := model.LanguageKey{Uuid: uuid.NewString(), Value: val}

	if err := store.Insert(key); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	got, ok := store.GetByValue(val)
	if !ok {
		t.Fatal("GetByValue failed: key not found")
	}
	if got.Value != val {
		t.Errorf("Expected value %s, got %s", val, got.Value)
	}
}

func TestLanguageKeyStore_List(t *testing.T) {
	store := NewLanguageKeyStore()

	for i := 0; i < 3; i++ {
		store.Insert(model.LanguageKey{
			Uuid:  uuid.NewString(),
			Value: "key_" + uuid.NewString(),
		})
	}

	keys := store.List()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
}

func TestLanguageKeyStore_Update(t *testing.T) {
	store := NewLanguageKeyStore()

	id := uuid.NewString()
	key := model.LanguageKey{Uuid: id, Value: "before"}
	if err := store.Insert(key); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	updated := model.LanguageKey{Uuid: id, Value: "after"}
	if err := store.Update(id, updated); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ := store.Get(id)
	if got.Value != "after" {
		t.Errorf("Expected value 'after', got %s", got.Value)
	}
}

func TestLanguageKeyStore_Update_DuplicateValue(t *testing.T) {
	store := NewLanguageKeyStore()

	id1 := uuid.NewString()
	id2 := uuid.NewString()
	store.Insert(model.LanguageKey{Uuid: id1, Value: "val1"})
	store.Insert(model.LanguageKey{Uuid: id2, Value: "val2"})

	// Attempt to update id2 to use val1 → should fail
	err := store.Update(id2, model.LanguageKey{Uuid: id2, Value: "val1"})
	if err == nil {
		t.Error("Expected error on value conflict during update")
	}
}

func TestLanguageKeyStore_Delete(t *testing.T) {
	store := NewLanguageKeyStore()

	id := uuid.NewString()
	key := model.LanguageKey{Uuid: id, Value: "deletable"}
	store.Insert(key)

	if err := store.Delete(id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, ok := store.Get(id); ok {
		t.Error("Expected key to be deleted")
	}

	if _, ok := store.GetByValue("deletable"); ok {
		t.Error("Expected value mapping to be removed")
	}
}

func TestLanguageKeyStore_Delete_NotFound(t *testing.T) {
	store := NewLanguageKeyStore()

	err := store.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error on delete of nonexistent key")
	}
}
