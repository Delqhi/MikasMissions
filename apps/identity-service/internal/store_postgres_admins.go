package internal

import (
	"database/sql"
	"fmt"
)

func (s *PostgresStore) FindAdminByEmail(email string) (AdminUser, bool, error) {
	var admin AdminUser
	err := s.db.QueryRow(
		`select id::text, email, coalesce(password_hash, '')
		 from identity.admin_users
		 where lower(email) = lower($1)`,
		email,
	).Scan(&admin.ID, &admin.Email, &admin.PasswordHash)
	if err == sql.ErrNoRows {
		return AdminUser{}, false, nil
	}
	if err != nil {
		return AdminUser{}, false, fmt.Errorf("find admin by email: %w", err)
	}
	return admin, true, nil
}

func (s *PostgresStore) UpdateAdminLastLogin(adminUserID string) error {
	_, err := s.db.Exec(
		`update identity.admin_users
		 set last_login_at = now()
		 where id::text = $1`,
		adminUserID,
	)
	if err != nil {
		return fmt.Errorf("update admin last login: %w", err)
	}
	return nil
}
