package oob

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lib/pq"

	"github.com/authgear/authgear-server/pkg/lib/authn"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator"
	"github.com/authgear/authgear-server/pkg/lib/infra/db"
)

type Store struct {
	SQLBuilder  db.SQLBuilder
	SQLExecutor db.SQLExecutor
}

func (s *Store) selectQuery() db.SelectBuilder {
	return s.SQLBuilder.Tenant().
		Select(
			"a.id",
			"a.user_id",
			"a.created_at",
			"a.updated_at",
			"a.tag",
			"ao.channel",
			"ao.phone",
			"ao.email",
		).
		From(s.SQLBuilder.FullTableName("authenticator"), "a").
		Join(s.SQLBuilder.FullTableName("authenticator_oob"), "ao", "a.id = ao.id")
}

func (s *Store) scan(scn db.Scanner) (*Authenticator, error) {
	a := &Authenticator{}
	var tag []byte

	err := scn.Scan(
		&a.ID,
		&a.UserID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&tag,
		&a.Channel,
		&a.Phone,
		&a.Email,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, authenticator.ErrAuthenticatorNotFound
	} else if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(tag, &a.Tag); err != nil {
		return nil, err
	}

	return a, nil
}

func (s *Store) Get(userID string, id string) (*Authenticator, error) {
	q := s.selectQuery().Where("a.user_id = ? AND a.id = ?", userID, id)

	row, err := s.SQLExecutor.QueryRowWith(q)
	if err != nil {
		return nil, err
	}

	return s.scan(row)
}

func (s *Store) GetMany(ids []string) ([]*Authenticator, error) {
	builder := s.selectQuery().Where("a.id = ANY (?)", pq.Array(ids))

	rows, err := s.SQLExecutor.QueryWith(builder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []*Authenticator
	for rows.Next() {
		a, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		as = append(as, a)
	}

	return as, nil
}

func (s *Store) List(userID string) ([]*Authenticator, error) {
	q := s.selectQuery().Where("a.user_id = ?", userID)

	rows, err := s.SQLExecutor.QueryWith(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authenticators []*Authenticator
	for rows.Next() {
		a, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		authenticators = append(authenticators, a)
	}

	return authenticators, nil
}

func (s *Store) Delete(id string) error {
	q := s.SQLBuilder.Tenant().
		Delete(s.SQLBuilder.FullTableName("authenticator_oob")).
		Where("id = ?", id)
	_, err := s.SQLExecutor.ExecWith(q)
	if err != nil {
		return err
	}

	q = s.SQLBuilder.Tenant().
		Delete(s.SQLBuilder.FullTableName("authenticator")).
		Where("id = ?", id)
	_, err = s.SQLExecutor.ExecWith(q)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Create(a *Authenticator) error {
	tag, err := json.Marshal(a.Tag)
	if err != nil {
		return err
	}

	q := s.SQLBuilder.Tenant().
		Insert(s.SQLBuilder.FullTableName("authenticator")).
		Columns(
			"id",
			"type",
			"user_id",
			"created_at",
			"updated_at",
			"tag",
		).
		Values(
			a.ID,
			authn.AuthenticatorTypeOOB,
			a.UserID,
			a.CreatedAt,
			a.UpdatedAt,
			tag,
		)
	_, err = s.SQLExecutor.ExecWith(q)
	if err != nil {
		return err
	}

	q = s.SQLBuilder.Tenant().
		Insert(s.SQLBuilder.FullTableName("authenticator_oob")).
		Columns(
			"id",
			"channel",
			"phone",
			"email",
		).
		Values(
			a.ID,
			a.Channel,
			a.Phone,
			a.Email,
		)
	_, err = s.SQLExecutor.ExecWith(q)
	if err != nil {
		return err
	}

	return nil
}
