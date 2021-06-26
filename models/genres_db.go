package models

import (
	"context"
	"time"
)

func (m *DBModel) GetAllGenres() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, genre_name, created_at, updated_at FROM genres`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gs []*Genre
	for rows.Next() {
		var g Genre
		err := rows.Scan(
			&g.ID,
			&g.GenreName,
			&g.CreatedAt,
			&g.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		gs = append(gs, &g)
	}

	return gs, nil
}
