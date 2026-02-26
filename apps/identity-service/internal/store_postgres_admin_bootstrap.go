package internal

import "fmt"

func (s *PostgresStore) UpsertAdminUser(email, passwordHash string) error {
	_, err := s.db.Exec(
		`insert into identity.admin_users (email, password_hash)
		 values ($1, $2)
		 on conflict (email) do update set
		   password_hash = excluded.password_hash`,
		email,
		passwordHash,
	)
	if err != nil {
		return fmt.Errorf("upsert admin user: %w", err)
	}
	return nil
}
