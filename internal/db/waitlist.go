package db

import "database/sql"

func SaveWaitlist(db *sql.DB, chatID int64, fullname, email, course string) error {
	_, err := db.Exec(`
		INSERT INTO waitlist (chat_id, fullname, email, course)
		VALUES ($1, $2, $3, $4)
	`, chatID, fullname, email, course)

	return err
}
