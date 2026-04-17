package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"inkflow/internal/domain/models"
	"inkflow/internal/storage"
)

func (s *Storage) SaveProfile(ctx context.Context, userID int64, email string, displayName string) error {
	const op = "storage.postgres.SaveProfile"

	query := `INSERT INTO profiles (user_id, email, display_name) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, query, userID, email, displayName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Profile(ctx context.Context, userID int64) (models.Profile, error) {
	const op = "storage.postgres.Profile"

	query := `SELECT user_id, email, display_name FROM profiles WHERE user_id = $1`

	var profile models.Profile

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.Email,
		&profile.DisplayName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Profile{}, fmt.Errorf("%s: %w", op, storage.ErrProfileNotFound)
		}

		return models.Profile{}, fmt.Errorf("%s: %w", op, err)
	}

	return profile, nil
}
