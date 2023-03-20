package db

import "context"

func (db *Db) SelectIDByAPIKey(ctx context.Context, apyKey string) (uint, error) {
	sqlRequest := "SELECT id FROM auth WHERE `api-key` = ?;"

	var result uint

	err := db.QueryRowxContext(ctx, sqlRequest, apyKey).Scan(&result)

	return result, err
}
