package mock

import (
	"context"
	"fmt"
)

type Client struct{}

func (client *Client) SendObject(ctx context.Context, method, path string, snd, rcv any) (err error) {

	if method != "POST" {
		err = fmt.Errorf("mock does not like method: %s", method)
	}
	return
}
