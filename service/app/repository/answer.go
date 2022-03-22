package repository

import (
	"context"
	"database/sql"
	"github.com/vladazn/dhq/service/app/domain"
	"github.com/vladazn/dhq/service/app/pkg/mariadb"
)

type AnswersModel struct {
	md mariadb.MariaDB
}

func newAnswersRepo(md mariadb.MariaDB) *AnswersModel {
	return &AnswersModel{
		md,
	}
}

func (a *AnswersModel) Migrate(ctx context.Context) error {

	q1 := `drop schema if exists dhq;`
	q2 := `create schema dhq collate utf8_general_ci;`
	q3 := `use dhq;`
	q4 := `
	create table answers
(
    id           int auto_increment,
    answer_key   varchar(255) null,
    answer_value varchar(255) null,
    user_id      int          not null,
    event        varchar(255) not null,
    constraint answers_pk
        primary key (id)
);
`

	_, err := a.md.GetClient().ExecContext(ctx, q1)
	if err != nil {
		return err
	}

	_, err = a.md.GetClient().ExecContext(ctx, q2)
	if err != nil {
		return err
	}

	_, err = a.md.GetClient().ExecContext(ctx, q3)
	if err != nil {
		return err
	}

	_, err = a.md.GetClient().ExecContext(ctx, q4)

	return err
}

func (a *AnswersModel) GetHistory(ctx context.Context, userId int, key string) ([]domain.Action,
	error) {
	client := a.md.GetClient()

	q := `select 
			answer_key, answer_value, event 
		from answers 
		where user_id = ? and answer_key = ?`

	r, err := client.QueryContext(ctx, q, userId, key)
	if err != nil {
		return nil, err
	}

	arr := make([]domain.Action, 0, 10)

	for r.Next() {

		key, value := &sql.NullString{}, &sql.NullString{}
		var event string

		var err = r.Scan(&key, &value, &event)
		if err != nil {
			return nil, err
		}

		v := domain.Action{
			Event: event,
		}

		if key.Valid {
			v.Data = &domain.Answer{
				Key:   key.String,
				Value: value.String,
			}
		}

		arr = append(arr, v)

	}

	return arr, nil
}

func (a *AnswersModel) GetLastAction(ctx context.Context, userId int,
	key string) (*domain.Action, error) {
	client := a.md.GetClient()

	q := `select 
			answer_key, answer_value, event 
		from answers
		where user_id = ? and answer_key = ?
		order by id desc limit 1;`

	r, err := client.QueryContext(ctx, q, userId, key)
	if err != nil {
		return nil, err
	}

	var response *domain.Action

	for r.Next() {
		key, value := &sql.NullString{}, &sql.NullString{}
		var event string

		var err = r.Scan(&key, &value, &event)
		if err != nil {
			return nil, err
		}

		response = &domain.Action{
			Event: event,
		}

		if key.String != "" {
			response.Data = &domain.Answer{
				Key:   key.String,
				Value: value.String,
			}
		}
	}

	return response, nil
}

func (a *AnswersModel) AddAction(ctx context.Context, userId int, action *domain.Action) error {
	client := a.md.GetClient()

	q := `insert into answers (user_id, answer_key, answer_value, event) 
			values (?,?,?,?);`

	var key, value string
	if action.Data != nil {
		key = action.Data.Key
		value = action.Data.Value
	}

	_, err := client.ExecContext(ctx, q, userId, key, value, action.Event)
	if err != nil {
		return err
	}

	return nil
}
