package db

import (
	"context"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
	"github.com/sirupsen/logrus"
)

func (db *Db) SelectAboutUser(ctx context.Context, userName string) ([]*core.AboutUser, error) {
	result := make([]*core.AboutUser, 0) // such initialization is necessary so that the JSON response does not contain nil

	sqlRequest := `SELECT u.id, u.username, p.first_name, p.last_name, p.city, d.school 
		FROM user u
		INNER JOIN user_profile p ON u.id = p.user_id
		INNER JOIN user_data d ON u.id = d.user_id`

	if len(userName) > 0 {
		sqlRequest += " WHERE u.username = ?;"
		if err := db.SelectContext(ctx, &result, sqlRequest, userName); err != nil {
			logrus.Errorf("SQL=%s", sqlRequest)

			return nil, err
		}
	} else {
		sqlRequest += ";"
		if err := db.SelectContext(ctx, &result, sqlRequest); err != nil {
			return nil, err
		}
	}

	return result, nil
}
