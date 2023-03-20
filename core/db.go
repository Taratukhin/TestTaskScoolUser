package core

import "context"

type AboutUser struct {
	ID        uint   `db:"id" json:"id"`
	UserName  string `db:"username" json:"username"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	City      string `db:"city" json:"city"`
	School    string `db:"school" json:"school"`
}

type Db interface {
	SelectIDByAPIKey(ctx context.Context, apiKey string) (uint, error)
	SelectAboutUser(ctx context.Context, userName string) ([]*AboutUser, error)
}
