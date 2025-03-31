package main

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rah-0/nabu"

	"github.com/rah-0/meisterwerk/model"
	"github.com/rah-0/meisterwerk/util"
)

var testCtx context.Context
var cancel context.CancelFunc

var natsClientConn *nats.Conn

func TestMain(m *testing.M) {
	util.TestMainWrapper(util.TestConfig{
		M: m,
		LoadResources: func() error {
			testCtx, cancel = context.WithCancel(context.Background())
			go func() {
				if err := start(testCtx); err != nil {
					nabu.FromError(err).Log()
				}
			}()
			time.Sleep(100 * time.Millisecond) // give NATS handlers time to register

			var err error
			natsClientConn, err = nats.Connect(nats.DefaultURL)
			return err
		},
		UnloadResources: func() error {
			natsClientConn.Close()
			cancel()
			time.Sleep(100 * time.Millisecond) // wait for shutdown
			return nil
		},
	})
}

func TestLanguageInsertAndGet(t *testing.T) {
	lang := model.Language{
		Uuid:        uuid.NewString(),
		Prefix:      "en-US",
		Lang:        "English",
		Title:       "English",
		Img:         "/static/img/flags/us.png",
		MonthsShort: "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec",
	}

	// INSERT
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("nats insert request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var insertResp util.NatsResponse
	if err = model.Decode(&insertResp); err != nil {
		t.Fatalf("decode insert response failed: %v", err)
	}
	if insertResp.Error != "" {
		t.Fatalf("insert returned error: %s", insertResp.Error)
	}

	// GET
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode get failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageGet, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("nats get request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var getResp util.NatsResponse
	if err := model.Decode(&getResp); err != nil {
		t.Fatalf("decode get response failed: %v", err)
	}
	if getResp.Error != "" {
		t.Fatalf("get returned error: %s", getResp.Error)
	}

	// Decode payload into Language
	payload, ok := getResp.Data.(model.Language)
	if !ok {
		t.Fatalf("unexpected response payload type: %T", getResp.Data)
	}

	if payload.Uuid != lang.Uuid || payload.Lang != lang.Lang {
		t.Errorf("mismatch: got %+v, want %+v", payload, lang)
	}
}

func TestLanguageUpdate(t *testing.T) {
	lang := model.Language{
		Uuid:        uuid.NewString(),
		Prefix:      "en-US",
		Lang:        "Original",
		Title:       "English",
		Img:         "/static/img/flags/us.png",
		MonthsShort: "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec",
	}

	// Insert original
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var insertResp util.NatsResponse
	if err := model.Decode(&insertResp); err != nil || insertResp.Error != "" {
		t.Fatalf("insert failed: %v | %s", err, insertResp.Error)
	}

	// Modify and update
	lang.Lang = "Updated"
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode update failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageUpdate, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var updateResp util.NatsResponse
	if err := model.Decode(&updateResp); err != nil || updateResp.Error != "" {
		t.Fatalf("update failed: %v | %s", err, updateResp.Error)
	}

	// Get and verify
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode get failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageGet, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var getResp util.NatsResponse
	if err := model.Decode(&getResp); err != nil || getResp.Error != "" {
		t.Fatalf("get failed: %v | %s", err, getResp.Error)
	}
	updated, ok := getResp.Data.(model.Language)
	if !ok {
		t.Fatalf("unexpected type: %T", getResp.Data)
	}
	if updated.Lang != "Updated" {
		t.Errorf("expected Lang = Updated, got %s", updated.Lang)
	}
}

func TestLanguageDelete(t *testing.T) {
	lang := model.Language{
		Uuid:        uuid.NewString(),
		Prefix:      "en-US",
		Lang:        "ToDelete",
		Title:       "English",
		Img:         "/static/img/flags/us.png",
		MonthsShort: "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec",
	}

	// INSERT
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var insertResp util.NatsResponse
	if err := model.Decode(&insertResp); err != nil || insertResp.Error != "" || insertResp.Status != 200 {
		t.Fatalf("insert failed: %v | %s", err, insertResp.Error)
	}

	// DELETE
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode delete failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageDelete, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var deleteResp util.NatsResponse
	if err := model.Decode(&deleteResp); err != nil || deleteResp.Error != "" || deleteResp.Status != 200 {
		t.Fatalf("delete failed: %v | %s", err, deleteResp.Error)
	}

	// GET (should fail)
	model.BufferReset()
	if err := model.Encode(lang); err != nil {
		t.Fatalf("encode get after delete failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageGet, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("get after delete request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var getResp util.NatsResponse
	if err := model.Decode(&getResp); err != nil {
		t.Fatalf("decode get after delete failed: %v", err)
	}
	if getResp.Status == 200 {
		t.Error("expected error when getting deleted language, got success")
	}
}

func TestLanguageList(t *testing.T) {
	// Insert 2 languages
	for i := 0; i < 2; i++ {
		lang := model.Language{
			Uuid:        uuid.NewString(),
			Prefix:      "en-US",
			Lang:        "Lang" + uuid.NewString()[0:4],
			Title:       "English",
			Img:         "/static/img/flags/us.png",
			MonthsShort: "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec",
		}

		model.BufferReset()
		if err := model.Encode(lang); err != nil {
			t.Fatalf("encode insert #%d failed: %v", i+1, err)
		}

		respMsg, err := natsClientConn.Request(EndpointLanguageInsert, model.GetBytes(), time.Second)
		if err != nil {
			t.Fatalf("insert request #%d failed: %v", i+1, err)
		}

		model.SetBytes(respMsg.Data)
		var insertResp util.NatsResponse
		if err := model.Decode(&insertResp); err != nil {
			t.Fatalf("decode insert response #%d failed: %v", i+1, err)
		}
		if insertResp.Status != 200 || insertResp.Error != "" {
			t.Fatalf("insert #%d failed: status=%d error=%s", i+1, insertResp.Status, insertResp.Error)
		}
	}

	respMsg, err := natsClientConn.Request(EndpointLanguageList, nil, time.Second)
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}

	model.SetBytes(respMsg.Data)
	var listResp util.NatsResponse
	if err := model.Decode(&listResp); err != nil {
		t.Fatalf("decode list response failed: %v", err)
	}
	if listResp.Status != 200 || listResp.Error != "" {
		t.Fatalf("list returned error: status=%d error=%s", listResp.Status, listResp.Error)
	}

	// Decode and verify list
	langs, ok := listResp.Data.([]model.Language)
	if !ok {
		t.Fatalf("unexpected type for list response: %T", listResp.Data)
	}
	if len(langs) < 2 {
		t.Errorf("expected at least 2 languages, got %d", len(langs))
	}
}

func TestLanguageKeyInsertAndGet(t *testing.T) {
	key := model.LanguageKey{
		Uuid:  uuid.NewString(),
		Value: "welcome_message",
	}

	// INSERT
	model.BufferReset()
	if err := model.Encode(key); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageKeyInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var insertResp util.NatsResponse
	if err := model.Decode(&insertResp); err != nil || insertResp.Status != 200 || insertResp.Error != "" {
		t.Fatalf("insert failed: %v | %s", err, insertResp.Error)
	}

	// GET
	model.BufferReset()
	if err := model.Encode(key); err != nil {
		t.Fatalf("encode get failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageKeyGet, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var getResp util.NatsResponse
	if err := model.Decode(&getResp); err != nil || getResp.Status != 200 || getResp.Error != "" {
		t.Fatalf("get failed: %v | %s", err, getResp.Error)
	}
	got, ok := getResp.Data.(model.LanguageKey)
	if !ok {
		t.Fatalf("unexpected type from get: %T", getResp.Data)
	}
	if got.Uuid != key.Uuid || got.Value != key.Value {
		t.Errorf("get mismatch: got %+v, want %+v", got, key)
	}
}

func TestLanguageKeyUpdate(t *testing.T) {
	key := model.LanguageKey{
		Uuid:  uuid.NewString(),
		Value: "product.name",
	}

	// Insert
	model.BufferReset()
	if err := model.Encode(key); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageKeyInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var insertResp util.NatsResponse
	if err := model.Decode(&insertResp); err != nil || insertResp.Status != 200 || insertResp.Error != "" {
		t.Fatalf("insert failed: %v | %s", err, insertResp.Error)
	}

	// Update
	key.Value = "product.title"
	model.BufferReset()
	if err := model.Encode(key); err != nil {
		t.Fatalf("encode update failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageKeyUpdate, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var updateResp util.NatsResponse
	if err := model.Decode(&updateResp); err != nil || updateResp.Status != 200 || updateResp.Error != "" {
		t.Fatalf("update failed: %v | %s", err, updateResp.Error)
	}

	// Verify
	model.BufferReset()
	if err := model.Encode(key); err != nil {
		t.Fatalf("encode get failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageKeyGet, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var getResp util.NatsResponse
	if err := model.Decode(&getResp); err != nil || getResp.Status != 200 || getResp.Error != "" {
		t.Fatalf("get after update failed: %v | %s", err, getResp.Error)
	}
	updated, ok := getResp.Data.(model.LanguageKey)
	if !ok {
		t.Fatalf("unexpected type: %T", getResp.Data)
	}
	if updated.Value != "product.title" {
		t.Errorf("expected updated value, got %s", updated.Value)
	}
}

func TestLanguageKeyGetByValue(t *testing.T) {
	key := model.LanguageKey{
		Uuid:  uuid.NewString(),
		Value: "footer.message",
	}

	// Insert
	model.BufferReset()
	if err := model.Encode(key); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageKeyInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}

	// Get by value
	model.BufferReset()
	if err := model.Encode(key); err != nil {
		t.Fatalf("encode get_by_value failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageKeyGetByValue, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("get_by_value request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var resp util.NatsResponse
	if err = model.Decode(&resp); err != nil || resp.Status != 200 || resp.Error != "" {
		t.Fatalf("get_by_value failed: %v | %s", err, resp.Error)
	}
	byValue, ok := resp.Data.(model.LanguageKey)
	if !ok {
		t.Fatalf("unexpected type from get_by_value: %T", resp.Data)
	}
	if byValue.Uuid != key.Uuid {
		t.Errorf("get_by_value returned wrong key: got %s, want %s", byValue.Uuid, key.Uuid)
	}
}

func TestLanguageValue_InsertAndGet(t *testing.T) {
	value := model.LanguageValue{
		Uuid:            uuid.NewString(),
		UuidLanguage:    uuid.NewString(),
		UuidLanguageKey: uuid.NewString(),
		Value:           "Hallo Welt",
	}

	// INSERT
	model.BufferReset()
	if err := model.Encode(value); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageValueInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var insertResp util.NatsResponse
	if err := model.Decode(&insertResp); err != nil || insertResp.Status != 200 || insertResp.Error != "" {
		t.Fatalf("insert failed: %v | %s", err, insertResp.Error)
	}

	// GET
	model.BufferReset()
	if err := model.Encode(value); err != nil {
		t.Fatalf("encode get failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageValueGet, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var getResp util.NatsResponse
	if err := model.Decode(&getResp); err != nil || getResp.Status != 200 || getResp.Error != "" {
		t.Fatalf("get failed: %v | %s", err, getResp.Error)
	}
	got, ok := getResp.Data.(model.LanguageValue)
	if !ok {
		t.Fatalf("unexpected type: %T", getResp.Data)
	}
	if got.Value != value.Value {
		t.Errorf("mismatch: got %s, want %s", got.Value, value.Value)
	}
}

func TestLanguageValue_Update(t *testing.T) {
	value := model.LanguageValue{
		Uuid:            uuid.NewString(),
		UuidLanguage:    uuid.NewString(),
		UuidLanguageKey: uuid.NewString(),
		Value:           "Before",
	}

	// Insert first
	model.BufferReset()
	if err := model.Encode(value); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageValueInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}

	// Update
	value.Value = "After"
	model.BufferReset()
	if err := model.Encode(value); err != nil {
		t.Fatalf("encode update failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageValueUpdate, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var updateResp util.NatsResponse
	if err := model.Decode(&updateResp); err != nil || updateResp.Status != 200 || updateResp.Error != "" {
		t.Fatalf("update failed: %v | %s", err, updateResp.Error)
	}
}

func TestLanguageValue_Delete(t *testing.T) {
	value := model.LanguageValue{
		Uuid:            uuid.NewString(),
		UuidLanguage:    uuid.NewString(),
		UuidLanguageKey: uuid.NewString(),
		Value:           "DeleteMe",
	}

	// Insert
	model.BufferReset()
	if err := model.Encode(value); err != nil {
		t.Fatalf("encode insert failed: %v", err)
	}
	respMsg, err := natsClientConn.Request(EndpointLanguageValueInsert, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("insert request failed: %v", err)
	}

	// Delete
	model.BufferReset()
	if err := model.Encode(value); err != nil {
		t.Fatalf("encode delete failed: %v", err)
	}
	respMsg, err = natsClientConn.Request(EndpointLanguageValueDelete, model.GetBytes(), time.Second)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var deleteResp util.NatsResponse
	if err := model.Decode(&deleteResp); err != nil || deleteResp.Status != 200 || deleteResp.Error != "" {
		t.Fatalf("delete failed: %v | %s", err, deleteResp.Error)
	}
}

func TestLanguageValue_List(t *testing.T) {
	for i := 0; i < 2; i++ {
		value := model.LanguageValue{
			Uuid:            uuid.NewString(),
			UuidLanguage:    uuid.NewString(),
			UuidLanguageKey: uuid.NewString(),
			Value:           "val",
		}
		model.BufferReset()
		if err := model.Encode(value); err != nil {
			t.Fatalf("encode insert #%d failed: %v", i+1, err)
		}
		respMsg, err := natsClientConn.Request(EndpointLanguageValueInsert, model.GetBytes(), time.Second)
		if err != nil {
			t.Fatalf("insert request #%d failed: %v", i+1, err)
		}
		model.SetBytes(respMsg.Data)
		var insertResp util.NatsResponse
		if err := model.Decode(&insertResp); err != nil || insertResp.Status != 200 || insertResp.Error != "" {
			t.Fatalf("insert #%d failed: %v | %s", i+1, err, insertResp.Error)
		}
	}

	// List
	respMsg, err := natsClientConn.Request(EndpointLanguageValueList, nil, time.Second)
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}
	model.SetBytes(respMsg.Data)
	var listResp util.NatsResponse
	if err := model.Decode(&listResp); err != nil {
		t.Fatalf("decode list response failed: %v", err)
	}
	if listResp.Status != 200 || listResp.Error != "" {
		t.Fatalf("list failed: %v | %s", err, listResp.Error)
	}
	vals, ok := listResp.Data.([]model.LanguageValue)
	if !ok {
		t.Fatalf("unexpected type for list response: %T", listResp.Data)
	}
	if len(vals) < 2 {
		t.Errorf("expected at least 2 values, got %d", len(vals))
	}
}
