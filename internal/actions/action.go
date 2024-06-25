package action

import "context"

type Action interface {
	Run(ctx context.Context) error
}
