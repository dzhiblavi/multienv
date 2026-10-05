package dispatcher

import (
	"context"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, payload string) (any, error)
}
