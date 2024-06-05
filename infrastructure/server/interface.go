package server

import "context"

type Server interface {
	Start() error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
}
