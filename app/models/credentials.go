package models

import (
	"amplifier/app/db"
	"amplifier/app/helpers"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const (
	createCredentialSQL      = `insert into credentials (app, url, username, password, created_at) values ($1, $2, $3, $4, $5) returning id`
	selectCredentialSQL      = `select id, app, url, username, password, created_at, updated_at from credentials`
	selectCredentialSQLByID  = selectCredentialSQL + ` where id=$1`
	selectCredentialSQLByApp = selectCredentialSQL + ` where app=$1`
	countCredentialSQL       = `select count(id) from credentials`
	updateCredentialSQL      = `update credentials set (app, url, username, password, updated_at) = ($1, $2, $3, $4, $5) where id=$6`
)

type (
	Credential struct {
		SequentialIdentifier
		App      string
		Url      string
		Username string
		Password string
		Timestamps
	}
)

func (m *Credential) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	m.Timestamps.Touch()

	var err error
	if m.IsNew() {
		err = db.QueryRowContext(
			ctx,
			createCredentialSQL,
			m.App,
			m.Url,
			m.Username,
			m.Password,
			m.Timestamps.CreatedAt,
		).Scan(&m.ID)
		return err
	}
	_, err = db.ExecContext(
		ctx,
		updateCredentialSQL,
		m.App,
		m.Url,
		m.Username,
		m.Password,
		m.Timestamps.UpdatedAt,
		m.ID,
	)
	return err
}

func (m *Credential) All(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) ([]*Credential, error) {
	credentials := make([]*Credential, 0)

	query, args := m.buildQuery(
		selectCredentialSQL,
		filter,
	)
	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return credentials, err
	}

	for rows.Next() {
		var credential Credential
		err = rows.Scan(
			&credential.ID,
			&credential.App,
			&credential.Url,
			&credential.Username,
			&credential.Password,
			&credential.CreatedAt,
			&credential.UpdatedAt,
		)
		if err != nil {
			return credentials, err
		}
		credentials = append(credentials, &credential)
	}

	return credentials, err
}

func (m *Credential) ByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*Credential, error) {
	var credential Credential
	row := db.QueryRowContext(ctx, selectCredentialSQLByID, id)
	err := m.scan(row, &credential)
	return &credential, err
}

func (m *Credential) ByApp(
	ctx context.Context,
	db db.SQLOperations,
	app string,
) (*Credential, error) {
	var credential Credential
	row := db.QueryRowContext(ctx, selectCredentialSQLByApp, app)
	err := m.scan(row, &credential)
	return &credential, err
}

func (m *Credential) scan(
	row *sql.Row,
	credential *Credential,
) error {
	return row.Scan(
		&credential.ID,
		&credential.App,
		&credential.Url,
		&credential.Username,
		&credential.Password,
		&credential.CreatedAt,
		&credential.UpdatedAt,
	)
}

func (m *Credential) Count(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) (int, error) {
	query, args := m.buildQuery(
		countCredentialSQL,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (m *Credential) buildQuery(
	query string,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"app", "username"}
		for _, col := range columns {
			search := fmt.Sprintf(" (lower(%s) like '%%' || $%d || '%%')", col, placeholder.Touch())
			likeStmt = append(likeStmt, search)
			args = append(args, filter.Term)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(likeStmt, " or")))
	}

	if len(conditions) > 0 {
		query += " where" + strings.Join(conditions, " and")
	}

	if filter.Per > 0 && filter.Page > 0 {
		query += fmt.Sprintf(" order by id desc limit $%d offset $%d", placeholder.Touch(), placeholder.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}
