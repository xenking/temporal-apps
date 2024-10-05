package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/xenking/temporal-apps/pkg/currency"
)

type SQLite struct {
	db *sql.DB
}

func (mw *SQLite) Get(key string) (interface{}, error) {
	var value []byte
	if err := sqlite.Get(mw.db, "SELECT value FROM cache WHERE key = ?", &value, key); err != nil {
		return nil, err
	}
	var rates []*currency.Rate
	if err := json.Unmarshal(value, &rates); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return rates, nil
}

func (mw *SQLite) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	if _, err := sqlite.Exec(mw.db, "INSERT OR REPLACE INTO cache (key, value) VALUES (?, ?)", key, data); err != nil {
		return fmt.Errorf("sqlite.Exec: %w", err)
	}
	return nil
}
