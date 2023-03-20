package db

import (
	"context"
)

func (db *Db) exec(ctx context.Context, request string) (err error) {
	_, err = db.ExecContext(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
