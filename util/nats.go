package util

import (
	"github.com/nats-io/nats.go"
	"github.com/rah-0/nabu"

	"github.com/rah-0/meisterwerk/model"
)

func init() {
	model.PreloadGob(NatsResponse{})
}

type NatsResponse struct {
	Status int
	Error  string
	Data   any
}

func NatsRespondWith(nc *nats.Msg, payload any, err error) {
	model.BufferReset()

	resp := NatsResponse{}

	if err != nil {
		resp.Status = 500
		resp.Error = err.Error()
	} else {
		resp.Status = 200
		resp.Data = payload
	}

	if err = model.Encode(resp); err != nil {
		resp.Status = 500
		resp.Error = err.Error()
	}

	if err := nc.Respond(model.GetBytes()); err != nil {
		nabu.FromError(err).Log()
	}
}

func NatsBindHandler[T any](nc *nats.Conn, subject string, handler func(req T) (any, error)) error {
	_, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		var req T

		if len(msg.Data) > 0 {
			model.SetBytes(msg.Data)
			if err := model.Decode(&req); err != nil {
				NatsRespondWith(msg, nil, err)
				return
			}
		}

		resp, err := handler(req)
		NatsRespondWith(msg, resp, err)
	})
	return err
}
