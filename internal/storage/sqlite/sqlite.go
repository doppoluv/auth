package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"

	"auth/internal/domain/model"
	"auth/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %v", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email, username string,
	passwordHash []byte,
) (int64, error) {
	// TODO: завернуть в транзакцию
	stmt, err := s.db.Prepare(
		"INSERT INTO users(username, email, password_hash) VALUES(?, ?, ?)",
	)
	if err != nil {
		return 0, fmt.Errorf("create prepared statement: %w", err)
	}

	res, err := stmt.ExecContext(ctx, username, email, passwordHash)
	if err != nil {
		// TODO: проверять на ошибку существующего юзера

		return 0, fmt.Errorf("execute prepared statement: %w", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert id: %w", err)
	}

	return userID, nil
}

func (s *Storage) GetUserByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {
	stmt, err := s.db.Prepare(
		"SELECT id, username, email, is_admin, password_hash FROM users WHERE username = ?",
	)
	if err != nil {
		return nil, fmt.Errorf("create prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, username)

	var user model.User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("scan row: %w", storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("scan row: %w", err)
	}

	return &user, nil
}

func (s *Storage) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	stmt, err := s.db.Prepare(
		"SELECT id, username, email, is_admin, password_hash FROM users WHERE email = ?",
	)
	if err != nil {
		return nil, fmt.Errorf("create prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user model.User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("scan row: %w", storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("scan row: %w", err)
	}

	return &user, nil
}

func (s *Storage) IsUserAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	stmt, err := s.db.Prepare(
		"SELECT is_admin FROM users WHERE id = ?",
	)
	if err != nil {
		return false, fmt.Errorf("create prepared statement: %w", err) // TODO: желательно не возвращать false
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("scan row: %w", storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("scan row: %w", err)
	}

	return isAdmin, nil
}
