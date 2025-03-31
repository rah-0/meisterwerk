package main

import (
	"testing"

	"github.com/google/uuid"

	"github.com/rah-0/meisterwerk/model"
)

func TestLanguageStore_InsertAndGet(t *testing.T) {
	store := NewLanguageStore()

	id := uuid.NewString()
	lang := model.Language{
		Uuid:        id,
		Prefix:      "en-US",
		Lang:        "English",
		Title:       "English",
		Img:         "/static/img/flags/us.svg",
		MonthsShort: "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec",
	}

	if err := store.Insert(lang); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	got, ok := store.Get(id)
	if !ok {
		t.Fatal("Get failed: language not found")
	}

	if got.Uuid != lang.Uuid || got.Lang != lang.Lang {
		t.Errorf("Retrieved language does not match: got %+v, want %+v", got, lang)
	}
}

func TestLanguageStore_InsertDuplicate(t *testing.T) {
	store := NewLanguageStore()

	id := uuid.NewString()
	lang := model.Language{Uuid: id, Lang: "English"}

	if err := store.Insert(lang); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if err := store.Insert(lang); err == nil {
		t.Fatal("Expected error for duplicate insert, got nil")
	}
}

func TestLanguageStore_Get_NotFound(t *testing.T) {
	store := NewLanguageStore()

	_, ok := store.Get("nonexistent-id")
	if ok {
		t.Error("Expected not found result for Get")
	}
}

func TestLanguageStore_List(t *testing.T) {
	store := NewLanguageStore()

	ids := []string{uuid.NewString(), uuid.NewString(), uuid.NewString()}
	for _, id := range ids {
		store.Insert(model.Language{Uuid: id, Lang: "Test"})
	}

	langs := store.List()
	if len(langs) != 3 {
		t.Errorf("Expected 3 languages, got %d", len(langs))
	}
}

func TestLanguageStore_Update(t *testing.T) {
	store := NewLanguageStore()

	id := uuid.NewString()
	original := model.Language{Uuid: id, Lang: "Old"}
	if err := store.Insert(original); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	updated := model.Language{Uuid: id, Lang: "Updated"}
	if err := store.Update(id, updated); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ := store.Get(id)
	if got.Lang != "Updated" {
		t.Errorf("Update failed: got %s, want %s", got.Lang, "Updated")
	}
	if got.LastUpdate.Before(original.FirstInsert) {
		t.Error("Update should update LastUpdate timestamp")
	}
}

func TestLanguageStore_Update_NotFound(t *testing.T) {
	store := NewLanguageStore()

	err := store.Update("nonexistent", model.Language{Uuid: "nonexistent", Lang: "Doesn't Matter"})
	if err == nil {
		t.Error("Expected error for update on nonexistent item")
	}
}

func TestLanguageStore_Delete(t *testing.T) {
	store := NewLanguageStore()

	id := uuid.NewString()
	store.Insert(model.Language{Uuid: id, Lang: "DeleteMe"})

	if err := store.Delete(id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, ok := store.Get(id); ok {
		t.Error("Expected deleted language to be gone")
	}
}

func TestLanguageStore_Delete_NotFound(t *testing.T) {
	store := NewLanguageStore()

	err := store.Delete("not-there")
	if err == nil {
		t.Error("Expected error on deleting non-existent language")
	}
}
