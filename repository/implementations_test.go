package repository

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateEstate(t *testing.T) {
	t.Run("inserts estate successfully with generated name", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("expected no error creating sqlmock, got %v", err)
		}
		defer db.Close()

		repo := &Repository{Db: db}
		input := CreateEstateInput{
			ID:     "estate-1",
			Name:   "",
			Length: 100,
			Width:  50,
		}

		mock.ExpectExec("INSERT INTO estates").
			WithArgs(
				input.ID,
				sqlmock.AnyArg(), // generated Estate-xxxxxxxx
				input.Length,
				input.Width,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateEstate(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expected all expectations to be met, got %v", err)
		}
	})

	t.Run("use provided estate name", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		mock.ExpectExec("INSERT INTO estates").
			WithArgs(
				"id-1",
				"My Estate",
				10,
				20,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateEstate(context.Background(), CreateEstateInput{
			ID:     "id-1",
			Name:   "My Estate",
			Length: 10,
			Width:  20,
		})

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("inserts estate successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("expected no error creating sqlmock, got %v", err)
		}
		defer db.Close()

		repo := &Repository{Db: db}
		input := CreateEstateInput{ID: "estate-1", Name: "Estate One", Length: 10, Width: 20}

		mock.ExpectExec("INSERT INTO estates").WithArgs(input.ID, input.Name, input.Length, input.Width).WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateEstate(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expected all expectations to be met, got %v", err)
		}
	})
}

func TestCreateTree(t *testing.T) {
	t.Run("returns error when insert fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("expected no error creating sqlmock, got %v", err)
		}
		defer db.Close()

		repo := &Repository{Db: db}

		input := CreateTreeInput{
			ID:       "tree-1",
			EstateID: "estate-1",
			X:        1,
			Y:        2,
			Height:   5,
		}

		mock.ExpectExec("INSERT INTO trees").
			WithArgs(input.ID, input.EstateID, input.X, input.Y, input.Height).
			WillReturnError(errors.New("insert error"))

		err = repo.CreateTree(context.Background(), input)
		if err == nil {
			t.Fatal("expected error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("inserts tree successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("expected no error creating sqlmock, got %v", err)
		}
		defer db.Close()

		repo := &Repository{Db: db}
		input := CreateTreeInput{ID: "tree-1", EstateID: "estate-1", X: 1, Y: 2, Height: 5}

		mock.ExpectExec("INSERT INTO trees").WithArgs(input.ID, input.EstateID, input.X, input.Y, input.Height).WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateTree(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expected all expectations to be met, got %v", err)
		}
	})
}

func TestEstateExists(t *testing.T) {
	t.Run("estate exists", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{"exists"}).
			AddRow(true)

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("estate-1").
			WillReturnRows(rows)

		exists, err := repo.EstateExists(context.Background(), "estate-1")
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Fatal("expected true")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("estate-1").
			WillReturnError(errors.New("db error"))

		_, err := repo.EstateExists(context.Background(), "estate-1")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTreeExists(t *testing.T) {
	t.Run("tree exists", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{"exists"}).
			AddRow(true)

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("estate-1", 1, 2).
			WillReturnRows(rows)

		exists, err := repo.TreeExists(context.Background(), "estate-1", 1, 2)
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Fatal("expected true")
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("estate-1", 1, 2).
			WillReturnError(errors.New("db error"))

		_, err := repo.TreeExists(context.Background(), "estate-1", 1, 2)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestGetTreeStats(t *testing.T) {
	t.Run("returns empty stats", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("expected no error creating sqlmock, got %v", err)
		}
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{"height"})

		mock.ExpectQuery("SELECT height FROM trees").
			WithArgs("estate-1").
			WillReturnRows(rows)

		stats, err := repo.GetTreeStats(context.Background(), "estate-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if stats.Count != 0 {
			t.Fatalf("expected count 0, got %d", stats.Count)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("returns stats for odd heights", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{"height"}).
			AddRow(int64(10)).
			AddRow(int64(20)).
			AddRow(int64(30))

		mock.ExpectQuery("SELECT height FROM trees").
			WithArgs("estate-1").
			WillReturnRows(rows)

		stats, err := repo.GetTreeStats(context.Background(), "estate-1")
		if err != nil {
			t.Fatal(err)
		}

		if stats.Count != 3 ||
			stats.Min != 10 ||
			stats.Max != 30 ||
			stats.Median != 20 {
			t.Fatalf("unexpected stats: %+v", stats)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("returns stats for even heights", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{"height"}).
			AddRow(int64(10)).
			AddRow(int64(20)).
			AddRow(int64(30)).
			AddRow(int64(40))

		mock.ExpectQuery("SELECT height FROM trees").
			WithArgs("estate-1").
			WillReturnRows(rows)

		stats, err := repo.GetTreeStats(context.Background(), "estate-1")
		if err != nil {
			t.Fatal(err)
		}

		if stats.Count != 4 ||
			stats.Min != 10 ||
			stats.Max != 40 ||
			stats.Median != 25 {
			t.Fatalf("unexpected stats: %+v", stats)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		mock.ExpectQuery("SELECT height FROM trees").
			WithArgs("estate-1").
			WillReturnError(errors.New("db error"))

		_, err := repo.GetTreeStats(context.Background(), "estate-1")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns scan error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{"height"}).
			AddRow("invalid")

		mock.ExpectQuery("SELECT height FROM trees").
			WithArgs("estate-1").
			WillReturnRows(rows)

		_, err := repo.GetTreeStats(context.Background(), "estate-1")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns rows error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{"height"}).
			AddRow(int64(10)).
			RowError(0, errors.New("row error"))

		mock.ExpectQuery("SELECT height FROM trees").
			WithArgs("estate-1").
			WillReturnRows(rows)

		_, err := repo.GetTreeStats(context.Background(), "estate-1")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestGetDronePlanSummary(t *testing.T) {
	t.Run("returns summary successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("expected no error creating sqlmock, got %v", err)
		}
		defer db.Close()

		repo := &Repository{Db: db}

		rows := sqlmock.NewRows([]string{
			"total_distance",
			"total_horizontal_distance",
			"total_vertical_distance",
		}).AddRow(int64(100), int64(60), int64(40))

		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(total_distance\\), 0\\)").
			WithArgs("estate-1").
			WillReturnRows(rows)

		summary, err := repo.GetDronePlanSummary(context.Background(), "estate-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if summary.Distance != 100 ||
			summary.HorizontalDistance != 60 ||
			summary.VerticalDistance != 40 {
			t.Fatalf("unexpected summary: %+v", summary)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("expected no error creating sqlmock, got %v", err)
		}
		defer db.Close()

		repo := &Repository{Db: db}

		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(total_distance\\), 0\\)").
			WithArgs("estate-1").
			WillReturnError(errors.New("db error"))

		_, err = repo.GetDronePlanSummary(context.Background(), "estate-1")
		if err == nil {
			t.Fatal("expected error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetTestById(t *testing.T) {
	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		mock.ExpectQuery("SELECT name FROM test WHERE id = \\$1").
			WithArgs("1").
			WillReturnError(errors.New("database error"))

		_, err = repo.GetTestById(context.Background(), GetTestByIdInput{
			Id: "1",
		})

		require.Error(t, err)
		require.EqualError(t, err, "database error")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("John")

		mock.ExpectQuery("SELECT name FROM test WHERE id = \\$1").
			WithArgs("1").
			WillReturnRows(rows)

		output, err := repo.GetTestById(context.Background(), GetTestByIdInput{
			Id: "1",
		})

		require.NoError(t, err)
		require.Equal(t, "John", output.Name)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
