package data

import (
	effx_api "github.com/effxhq/effx-api/generated/go"
)

type Data struct {
	Service *effx_api.ServicePayload
	Event   *effx_api.EventPayload
}
