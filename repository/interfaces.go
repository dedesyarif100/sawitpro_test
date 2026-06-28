// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	CreateEstate(ctx context.Context, input CreateEstateInput) error
	CreateTree(ctx context.Context, input CreateTreeInput) error
	EstateExists(ctx context.Context, estateID string) (bool, error)
	TreeExists(ctx context.Context, estateID string, x, y int) (bool, error)
	GetTreeStats(ctx context.Context, estateID string) (TreeStats, error)
	GetDronePlanSummary(ctx context.Context, estateID string) (DronePlanSummary, error)
}
