package interfaces

import "context"

type StartStopper interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}
