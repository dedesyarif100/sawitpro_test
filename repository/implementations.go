package repository

import "context"

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateEstate(ctx context.Context, input CreateEstateInput) error {
	_, err := r.Db.ExecContext(ctx, "INSERT INTO estates (name, length_plot, width_plot) VALUES ($1, $2, $3)", input.ID, input.Length, input.Width)
	return err
}
