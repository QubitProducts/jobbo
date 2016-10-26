package main

import (
	"context"
	pb "github.com/QubitProducts/jobbo/proto"
)

type PersistentStore interface {
	// Game is hard
	//	GetFrameworkLock(ctx context.Context, id string, lost func()) error

	GetFrameworkID(ctx context.Context) (string, error)
	SetFrameworkID(ctx context.Context, id string) error

	SetInstance(ctx context.Context, inst *pb.Instance) error
	ListInstances(ctx context.Context) ([]*pb.Instance, error)
	RemoveInstance(ctx context.Context, inst *pb.Instance) error
}
