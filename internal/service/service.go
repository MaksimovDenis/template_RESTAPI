package service

import "context"

type ServerService interface {
	Probe(ctx context.Context) error
}
