package client

import "github.com/effxhq/effx-go/data"

type Client interface {
	Synchronize(object *data.Data) error
}
