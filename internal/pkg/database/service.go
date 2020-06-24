package database

import "forum-db/internal/models"

// /service/clear
func Clear() error {
	query := `TRUNCATE TABLE users, forums, threads, post CASCADE;`
	dbPool.Exec(query)

	info = models.Status{
		Forum:  0,
		Post:   0,
		Thread: 0,
		User:   0,
	}

	ForumClearSlug()
	return nil
}

// /service/status
func Info() models.Status {
	return info
}
