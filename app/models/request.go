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
	createRequestSQL     = `insert into requests (app, multi, number, message, times, created_at) values ($1, $2, $3, $4, $5, $6) returning id`
	selectRequestSQL     = `select id, app, multi, number, message, times, created_at from requests`
	selectRequestSQLByID = selectCredentialSQL + ` where id=$1`
	countRequestSQL      = `select count(id) from requests`
)

type (
	Request struct {
		SequentialIdentifier
		App     string
		Multi   bool
		Number  int
		Message string
		Times   int
		Timestamps
	}
)

func (m *Request) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	m.Timestamps.Touch()
	err := db.QueryRowContext(
		ctx,
		createRequestSQL,
		m.App,
		m.Multi,
		m.Number,
		m.Message,
		m.Times,
		m.Timestamps.CreatedAt,
	).Scan(&m.ID)
	return err
}

func (r *Request) All(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) ([]*Request, error) {
	requests := make([]*Request, 0)

	query, args := r.buildQuery(
		selectRequestSQL,
		filter,
	)
	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return requests, err
	}

	for rows.Next() {
		var request Request
		err = rows.Scan(
			&request.ID,
			&request.App,
			&request.Multi,
			&request.Number,
			&request.Message,
			&request.Times,
			&request.CreatedAt,
		)
		if err != nil {
			return requests, err
		}
		requests = append(requests, &request)
	}

	return requests, err
}

func (r *Request) ByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*Request, error) {
	var request Request
	row := db.QueryRowContext(ctx, selectRequestSQLByID, id)
	err := r.scan(row, &request)
	return &request, err
}

func (*Request) scan(
	row *sql.Row,
	request *Request,
) error {
	return row.Scan(
		&request.ID,
		&request.App,
		&request.Multi,
		&request.Message,
		&request.Times,
		&request.CreatedAt,
	)
}

func (r *Request) Count(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) (int, error) {
	query, args := r.buildQuery(
		countRequestSQL,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (r *Request) buildQuery(
	query string,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"app", "message"}
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
