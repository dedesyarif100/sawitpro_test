package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateEstate(ctx context.Context, input CreateEstateInput) error {
	name := input.Name
	if name == "" {
		name = fmt.Sprintf("Estate-%s", uuid.NewString()[:8])
	}
	_, err := r.Db.ExecContext(ctx, "INSERT INTO estates (id, name, length_plot, width_plot) VALUES ($1, $2, $3, $4)", input.ID, name, input.Length, input.Width)
	return err
}

func (r *Repository) CreateTree(ctx context.Context, input CreateTreeInput) error {
	_, err := r.Db.ExecContext(ctx, "INSERT INTO trees (id, estate_id, plot_x, plot_y, height) VALUES ($1, $2, $3, $4, $5)", input.ID, input.EstateID, input.X, input.Y, input.Height)
	return err
}

func (r *Repository) EstateExists(ctx context.Context, estateID string) (bool, error) {
	var exists bool
	err := r.Db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM estates WHERE id::text = $1)", estateID).Scan(&exists)
	return exists, err
}

func (r *Repository) TreeExists(ctx context.Context, estateID string, x, y int) (bool, error) {
	var exists bool
	err := r.Db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM trees WHERE estate_id::text = $1 AND plot_x = $2 AND plot_y = $3)", estateID, x, y).Scan(&exists)
	return exists, err
}

func (r *Repository) GetTreeStats(ctx context.Context, estateID string) (TreeStats, error) {
	var stats TreeStats
	rows, err := r.Db.QueryContext(ctx, "SELECT height FROM trees WHERE estate_id::text = $1 ORDER BY height", estateID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	var heights []int64
	for rows.Next() {
		var h int64
		if err := rows.Scan(&h); err != nil {
			return stats, err
		}
		heights = append(heights, h)
	}
	if err := rows.Err(); err != nil {
		return stats, err
	}

	stats.Count = int64(len(heights))
	if stats.Count == 0 {
		return stats, nil
	}

	stats.Min = heights[0]
	stats.Max = heights[0]
	for _, h := range heights[1:] {
		if h < stats.Min {
			stats.Min = h
		}
		if h > stats.Max {
			stats.Max = h
		}
	}

	mid := len(heights) / 2
	if len(heights)%2 == 0 {
		stats.Median = (heights[mid-1] + heights[mid]) / 2
	} else {
		stats.Median = heights[mid]
	}

	return stats, nil
}

func (r *Repository) GetDronePlanSummary(ctx context.Context, estateID string) (DronePlanSummary, error) {
	var summary DronePlanSummary
	err := r.Db.QueryRowContext(ctx, "SELECT COALESCE(SUM(total_distance), 0), COALESCE(SUM(total_horizontal_distance), 0), COALESCE(SUM(total_vertical_distance), 0) FROM drone_simulations WHERE estate_id::text = $1", estateID).Scan(&summary.Distance, &summary.HorizontalDistance, &summary.VerticalDistance)
	if err != nil {
		return summary, err
	}
	return summary, nil
}
