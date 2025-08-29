package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("order not found")

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, url string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MinConns = 1
	cfg.MaxConns = 8
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Store{Pool: pool}, nil
}

func (s *Store) Close() { s.Pool.Close() }

func (s *Store) GetOrderRawJSON(ctx context.Context, orderUID string) (json.RawMessage, error) {
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT raw_json FROM orders WHERE order_uid=$1`, orderUID).Scan(&raw)
	if err != nil {
		return nil, ErrNotFound
	}
	return json.RawMessage(raw), nil
}

func (s *Store) LoadRecentRaw(ctx context.Context, limit int) (map[string]json.RawMessage, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT order_uid, raw_json
		FROM orders
		ORDER BY updated_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]json.RawMessage, limit)
	for rows.Next() {
		var id string
		var raw []byte
		if err := rows.Scan(&id, &raw); err != nil {
			return nil, err
		}
		out[id] = json.RawMessage(raw)
	}
	return out, rows.Err()
}
