// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type CreateEstateInput struct {
	ID     string
	Name   string
	Length int
	Width  int
}

type CreateTreeInput struct {
	ID       string
	EstateID string
	X        int
	Y        int
	Height   int
}

type TreeStats struct {
	Count  int64
	Max    int64
	Min    int64
	Median int64
}

type DronePlanSummary struct {
	Distance           int64
	HorizontalDistance int64
	VerticalDistance   int64
}
