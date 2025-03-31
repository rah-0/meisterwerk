package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
	"github.com/rah-0/nabu"

	"github.com/rah-0/meisterwerk/model"
	"github.com/rah-0/meisterwerk/util"
)

const (
	EndpointLanguageInsert = "translations.language.insert"
	EndpointLanguageUpdate = "translations.language.update"
	EndpointLanguageDelete = "translations.language.delete"
	EndpointLanguageGet    = "translations.language.get"
	EndpointLanguageList   = "translations.language.list"

	EndpointLanguageKeyInsert     = "translations.language_key.insert"
	EndpointLanguageKeyUpdate     = "translations.language_key.update"
	EndpointLanguageKeyDelete     = "translations.language_key.delete"
	EndpointLanguageKeyGet        = "translations.language_key.get"
	EndpointLanguageKeyGetByValue = "translations.language_key.get_by_value"
	EndpointLanguageKeyList       = "translations.language_key.list"

	EndpointLanguageValueInsert = "translations.language_value.insert"
	EndpointLanguageValueUpdate = "translations.language_value.update"
	EndpointLanguageValueDelete = "translations.language_value.delete"
	EndpointLanguageValueGet    = "translations.language_value.get"
	EndpointLanguageValueList   = "translations.language_value.list"
)

var (
	languageStore      = NewLanguageStore()
	languageValueStore = NewLanguageValueStore()
	languageKeyStore   = NewLanguageKeyStore()
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := start(ctx); err != nil {
		nabu.FromError(err).Log()
	}
}

func start(ctx context.Context) error {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nabu.FromError(err).WithArgs(nats.DefaultURL).Log()
	}
	nabu.FromMessage("Connected to NATS").WithArgs(nats.DefaultURL).Log()

	// Register handlers
	if err = registerLanguageHandlers(nc); err != nil {
		return nabu.FromError(err).WithArgs(nats.DefaultURL).Log()
	}
	if err = registerLanguageKeyHandlers(nc); err != nil {
		return nabu.FromError(err).WithArgs(nats.DefaultURL).Log()
	}
	if err = registerLanguageValueHandlers(nc); err != nil {
		return nabu.FromError(err).WithArgs(nats.DefaultURL).Log()
	}

	// Block until context is cancelled
	<-ctx.Done()
	nabu.FromMessage("Shutting down NATS translation service").Log()
	return nc.Drain()
}

func registerLanguageHandlers(nc *nats.Conn) error {
	if err := util.NatsBindHandler(nc, EndpointLanguageInsert, func(lang model.Language) (any, error) {
		return nil, languageStore.Insert(lang)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageUpdate, func(lang model.Language) (any, error) {
		return nil, languageStore.Update(lang.Uuid, lang)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageDelete, func(req model.Language) (any, error) {
		return nil, languageStore.Delete(req.Uuid)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageGet, func(req model.Language) (any, error) {
		return languageStore.Get(req.Uuid)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageList, func(_ any) (any, error) {
		return languageStore.List(), nil
	}); err != nil {
		return err
	}

	return nil
}

func registerLanguageKeyHandlers(nc *nats.Conn) error {
	if err := util.NatsBindHandler(nc, EndpointLanguageKeyInsert, func(key model.LanguageKey) (any, error) {
		return nil, languageKeyStore.Insert(key)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageKeyUpdate, func(key model.LanguageKey) (any, error) {
		return nil, languageKeyStore.Update(key.Uuid, key)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageKeyDelete, func(req model.LanguageKey) (any, error) {
		return nil, languageKeyStore.Delete(req.Uuid)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageKeyGet, func(req model.LanguageKey) (any, error) {
		return languageKeyStore.Get(req.Uuid)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageKeyGetByValue, func(req model.LanguageKey) (any, error) {
		return languageKeyStore.GetByValue(req.Value)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageKeyList, func(_ any) (any, error) {
		return languageKeyStore.List(), nil
	}); err != nil {
		return err
	}

	return nil
}

func registerLanguageValueHandlers(nc *nats.Conn) error {
	if err := util.NatsBindHandler(nc, EndpointLanguageValueInsert, func(val model.LanguageValue) (any, error) {
		return nil, languageValueStore.Insert(val)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageValueUpdate, func(val model.LanguageValue) (any, error) {
		return nil, languageValueStore.Update(val.Uuid, val)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageValueDelete, func(req model.LanguageValue) (any, error) {
		return nil, languageValueStore.Delete(req.Uuid)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageValueGet, func(req model.LanguageValue) (any, error) {
		return languageValueStore.Get(req.Uuid)
	}); err != nil {
		return err
	}

	if err := util.NatsBindHandler(nc, EndpointLanguageValueList, func(_ any) (any, error) {
		return languageValueStore.List()
	}); err != nil {
		return err
	}

	return nil
}
