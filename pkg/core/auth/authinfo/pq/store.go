// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pq

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo"
	"github.com/skygeario/skygear-server/pkg/core/db"
	dbPq "github.com/skygeario/skygear-server/pkg/core/db/pq"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/skydb"
)

type authInfoStore struct {
	sqlBuilder  db.SQLBuilder
	sqlExecutor db.SQLExecutor
	logger      *logrus.Entry
}

func newAuthInfoStore(builder db.SQLBuilder, executor db.SQLExecutor, loggerFactory logging.Factory) *authInfoStore {
	return &authInfoStore{
		sqlBuilder:  builder,
		sqlExecutor: executor,
		logger:      loggerFactory.NewLogger("auth-info"),
	}
}

func NewAuthInfoStore(builder db.SQLBuilder, executor db.SQLExecutor, loggerFactory logging.Factory) authinfo.Store {
	return newAuthInfoStore(builder, executor, loggerFactory)
}

func (s authInfoStore) CreateAuth(authinfo *authinfo.AuthInfo) (err error) {
	var (
		lastSeenAt     *time.Time
		lastLoginAt    *time.Time
		disabledReason *string
		disabledExpiry *time.Time
		verifyInfo     dbPq.JSONMapBooleanValue
	)
	lastSeenAt = authinfo.LastSeenAt
	if lastSeenAt != nil && lastSeenAt.IsZero() {
		lastSeenAt = nil
	}
	lastLoginAt = authinfo.LastLoginAt
	if lastLoginAt != nil && lastLoginAt.IsZero() {
		lastLoginAt = nil
	}
	disabledReason = &authinfo.DisabledMessage
	if *disabledReason == "" {
		disabledReason = nil
	}
	disabledExpiry = authinfo.DisabledExpiry
	if disabledExpiry != nil && disabledExpiry.IsZero() {
		disabledExpiry = nil
	}

	verifyInfo = authinfo.VerifyInfo

	builder := s.sqlBuilder.Tenant().
		Insert(s.sqlBuilder.FullTableName("user")).
		Columns(
			"id",
			"last_seen_at",
			"last_login_at",
			"disabled",
			"disabled_message",
			"disabled_expiry",
			"verified",
			"verify_info",
		).
		Values(
			authinfo.ID,
			lastSeenAt,
			lastLoginAt,
			authinfo.Disabled,
			disabledReason,
			disabledExpiry,
			authinfo.Verified,
			verifyInfo,
		)

	_, err = s.sqlExecutor.ExecWith(builder)
	if db.IsUniqueViolated(err) {
		return skydb.ErrUserDuplicated
	}

	return err
}

// UpdateAuth updates an existing AuthInfo matched by the ID field.
// nolint: gocyclo
func (s authInfoStore) UpdateAuth(authinfo *authinfo.AuthInfo) (err error) {
	var (
		lastSeenAt     *time.Time
		lastLoginAt    *time.Time
		disabledReason *string
		disabledExpiry *time.Time
		verifyInfo     dbPq.JSONMapBooleanValue
	)
	lastSeenAt = authinfo.LastSeenAt
	if lastSeenAt != nil && lastSeenAt.IsZero() {
		lastSeenAt = nil
	}
	lastLoginAt = authinfo.LastLoginAt
	if lastLoginAt != nil && lastLoginAt.IsZero() {
		lastLoginAt = nil
	}
	disabledReason = &authinfo.DisabledMessage
	if *disabledReason == "" {
		disabledReason = nil
	}
	disabledExpiry = authinfo.DisabledExpiry
	if disabledExpiry != nil && disabledExpiry.IsZero() {
		disabledExpiry = nil
	}

	verifyInfo = authinfo.VerifyInfo

	builder := s.sqlBuilder.Tenant().
		Update(s.sqlBuilder.FullTableName("user")).
		Set("last_seen_at", lastSeenAt).
		Set("last_login_at", lastLoginAt).
		Set("disabled", authinfo.Disabled).
		Set("disabled_message", disabledReason).
		Set("disabled_expiry", disabledExpiry).
		Set("verified", authinfo.Verified).
		Set("verify_info", verifyInfo).
		Where("id = ?", authinfo.ID)

	result, err := s.sqlExecutor.ExecWith(builder)
	if err != nil {
		if db.IsUniqueViolated(err) {
			return skydb.ErrUserDuplicated
		}
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return skydb.ErrUserNotFound
	} else if rowsAffected > 1 {
		panic(fmt.Errorf("want 1 rows updated, got %v", rowsAffected))
	}

	return nil
}

func (s authInfoStore) baseUserBuilder() db.SelectBuilder {
	return s.sqlBuilder.Tenant().
		Select(
			"id",
			"last_seen_at",
			"last_login_at",
			"disabled",
			"disabled_message",
			"disabled_expiry",
			"verified",
			"verify_info",
		).
		From(s.sqlBuilder.FullTableName("user"))
}

func (s authInfoStore) doScanAuth(authinfo *authinfo.AuthInfo, scanner sq.RowScanner) error {
	var (
		id             string
		lastSeenAt     pq.NullTime
		lastLoginAt    pq.NullTime
		disabled       bool
		disabledReason sql.NullString
		disabledExpiry pq.NullTime
		verified       bool
		verifyInfo     dbPq.NullJSONMapBoolean
	)

	err := scanner.Scan(
		&id,
		&lastSeenAt,
		&lastLoginAt,
		&disabled,
		&disabledReason,
		&disabledExpiry,
		&verified,
		&verifyInfo,
	)
	if err == sql.ErrNoRows {
		return skydb.ErrUserNotFound
	} else if err != nil {
		return err
	}

	authinfo.ID = id
	if lastSeenAt.Valid {
		authinfo.LastSeenAt = &lastSeenAt.Time
	} else {
		authinfo.LastSeenAt = nil
	}
	if lastLoginAt.Valid {
		authinfo.LastLoginAt = &lastLoginAt.Time
	} else {
		authinfo.LastLoginAt = nil
	}
	authinfo.Disabled = disabled
	if disabled {
		if disabledReason.Valid {
			authinfo.DisabledMessage = disabledReason.String
		} else {
			authinfo.DisabledMessage = ""
		}
		if disabledExpiry.Valid {
			expiry := disabledExpiry.Time.UTC()
			authinfo.DisabledExpiry = &expiry
		} else {
			authinfo.DisabledExpiry = nil
		}
	} else {
		authinfo.DisabledMessage = ""
		authinfo.DisabledExpiry = nil
	}

	authinfo.Verified = verified
	authinfo.VerifyInfo = verifyInfo.JSON

	return nil
}

func (s authInfoStore) GetAuth(id string, authinfo *authinfo.AuthInfo) error {
	builder := s.baseUserBuilder().Where("id = ?", id)
	scanner := s.sqlExecutor.QueryRowWith(builder)
	return s.doScanAuth(authinfo, scanner)
}

func (s authInfoStore) DeleteAuth(id string) error {
	builder := s.sqlBuilder.Tenant().
		Delete(s.sqlBuilder.FullTableName("user")).
		Where("id = ?", id)

	result, err := s.sqlExecutor.ExecWith(builder)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return skydb.ErrUserNotFound
	} else if rowsAffected > 1 {
		panic(fmt.Errorf("want 1 rows deleted, got %v", rowsAffected))
	}

	return nil
}

// this ensures that our structure conform to certain interfaces.
var (
	_ authinfo.Store = &authInfoStore{}
)
