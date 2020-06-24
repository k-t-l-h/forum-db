package database

import (
	"forum-db/internal/models"
	"github.com/jackc/pgx"
	"log"
)

var slugs map[string]string

//DONE
///forum/create

func CreateForum(forum models.Forum) ([]models.Forum, int) {

	query := ` INSERT INTO forums (title, author, slug, posts, threads) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING slug;`


	user, errs := CheckUser(models.User{NickName: forum.User})
	if errs != OK {
		return []models.Forum{}, UserNotFound
	}

	results := []models.Forum{}
	res := forum
	res.User = user.NickName

	row := dbPool.QueryRow(query, forum.Title, user.NickName, forum.Slug, 0, 0)
	err := row.Scan(&res.Slug)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case "23505":
				result, _ := GetForumOnConflict(forum)
				return result, ForumConflict
			case "23503":
				return results, UserNotFound
			default:
				result, _ := GetForumOnConflict(forum)
				return result, ForumConflict
			}
		}
	}

	info.Forum = info.Forum + 1

	results = append(results, res)
	//InsertForumUser(forum.User, forum.Slug)
	return results, OK
}

func GetForumOnConflict(f models.Forum) ([]models.Forum, int) {
	forums := []models.Forum{}
	query := `SELECT title, author, slug, posts, threads
				FROM forums 
				WHERE slug = $1 or title = $2;`

	row, _ := dbPool.Query(query, f.Slug, f.Title)
	defer row.Close()

	for row.Next() {
		forum := models.Forum{}
		err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
		if err != nil {
			return forums, NotFound
		}
		forums = append(forums, forum)
	}
	return forums, OK
}

///forum/{slug}/create
func CreateSlug(thread models.Thread) ([]models.Thread, int) {

	f, status := ForumCheck(models.Forum{Slug: thread.Forum})
	if status == NotFound {
		return nil, UserNotFound
	}

	thread.Forum = f.Slug

	t := thread

	if thread.Slug != "" {
		thread, status := CheckSlug(thread)
		if status == OK {
			th, _ := GetThreadBySlug(thread.Slug, t)
			return []models.Thread{th}, ForumConflict
		}
	}

	query := `INSERT INTO threads (author, message, title, created_at, forum, slug, votes)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING
				id
				`

	row := dbPool.QueryRow(query, thread.Author, thread.Message, thread.Title,
		thread.CreatedAt, thread.Forum, thread.Slug, 0)

	err := row.Scan(&t.Id)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case "23505":
				return []models.Thread{t}, ForumConflict
			case "23503":
				return []models.Thread{}, UserNotFound
			default:
				log.Print(pqError.Code)
				return []models.Thread{}, UserNotFound
			}
		}
	}
	info.Thread = info.Thread + 1
	return []models.Thread{t}, OK
}

///forum/{slug}/details
func GetForumBySlag(forum models.Forum) (models.Forum, int) {
	query := `SELECT title, author, slug, posts, threads
				FROM forums 
				WHERE slug = $1;`

	row := dbPool.QueryRow(query, forum.Slug)

	err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)

	if err != nil {
		return forum, NotFound
	}

	return forum, OK
}

func ForumCheck(forum models.Forum) (models.Forum, int) {
	query := `SELECT slug FROM forums 
				WHERE slug = $1;`

	row := dbPool.QueryRow(query, forum.Slug)

	err := row.Scan(&forum.Slug)

	if err != nil {
		return forum, NotFound
	}

	return forum, OK
}

///forum/{slug}/users
func GetForumUsers(forum models.Forum, limit, since, desc string) ([]models.User, int) {
	us := []models.User{}
	var row *pgx.Rows
	var err error

	if err == nil {
	}
	forum, state := ForumCheck(forum)
	if state != OK {
		return us, NotFound
	}

	query := ``

	if limit == "" && since == "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY lower(forum_users.nickname) DESC`
		} else {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY lower(forum_users.nickname)  ASC`
		}

		row, err = dbPool.Query(query, forum.Slug)
	}

	if limit != "" && since == "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY lower(forum_users.nickname)  DESC LIMIT $2`
		} else {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		ORDER BY lower(forum_users.nickname)  ASC LIMIT $2`
		}

		row, err = dbPool.Query(query, forum.Slug, limit)
	}

	if limit == "" && since != "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND lower(forum_users.nickname) < lower($2)
		ORDER BY lower(forum_users.nickname)  DESC  `
		} else {
			query = `SELECT email,
                      fullname,
                     users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND lower(forum_users.nickname) > lower($2)
		ORDER BY lower(forum_users.nickname)  ASC`
		}

		row, err = dbPool.Query(query, forum.Slug, since)
	}

	if limit != "" && since != "" {
		if desc == "true" {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND lower(forum_users.nickname) < lower($2)
		ORDER BY lower(forum_users.nickname)  DESC  LIMIT $3`
		} else {
			query = `SELECT email,
                      fullname,
                      users.nickname,
                      about
		FROM users JOIN forum_users ON users.nickname = forum_users.nickname WHERE forum = $1
		AND lower(forum_users.nickname) > lower($2)
		ORDER BY lower(forum_users.nickname) ASC LIMIT $3`
		}

		row, err = dbPool.Query(query, forum.Slug, since, limit)
		log.Print(err)
	}

	defer row.Close()

	for row.Next() {
		a := models.User{}
		log.Print(row.Scan(&a.Email, &a.FullName, &a.NickName, &a.About))
		us = append(us, a)
	}
	return us, OK
}

///forum/{slug}/threads
func GetForumThreads(t models.Thread, limit, since, desc string) ([]models.Thread, int) {

	th := []models.Thread{}
	var row *pgx.Rows
	var err error

	query := ``

	if limit == "" && since == "" {
		if desc == "" || desc == "false" {
			query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1)
						ORDER BY created_at ASC`

		} else {
			query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1)
						ORDER BY created_at DESC`
		}
		row, err = dbPool.Query(query, t.Forum)
	} else {

		if limit != "" && since == "" {
			if desc == "" || desc == "false" {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1) 
						ORDER BY created_at ASC  LIMIT $2`

			} else {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1) 
						ORDER BY created_at DESC  LIMIT $2`
			}

			row, err = dbPool.Query(query, t.Forum, limit)
		}

		if since != "" && limit == "" {
			if desc == "" || desc == "false" {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1) AND created_at >= $2
						ORDER BY created_at ASC `
			} else {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1) AND created_at <= $2
						ORDER BY created_at DESC `
			}

			row, err = dbPool.Query(query, t.Forum, since)
		}

		if since != "" && limit != "" {

			if desc == "" || desc == "false" {

				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1) AND created_at >= $2
						ORDER BY created_at ASC LIMIT $3`
			} else {
				query = `SELECT id, slug, author, created_at, forum, title, message, votes
						FROM threads
						WHERE lower(forum) = lower($1) AND created_at <= $2
						ORDER BY created_at DESC LIMIT $3`
			}
			row, err = dbPool.Query(query, t.Forum, since, limit)
		}
	}
	defer row.Close()
	for row.Next() {
		t := models.Thread{}
		err = row.Scan(&t.Id, &t.Slug, &t.Author, &t.CreatedAt, &t.Forum, &t.Title, &t.Message, &t.Votes)

		th = append(th, t)
	}
	if err == nil {

	}

	if len(th) == 0 {
		_, status := GetForumBySlag(models.Forum{Slug: t.Forum})
		if status != OK {
			return th, NotFound
		}
		return th, OK
	}

	return th, OK
}

func CheckSlug(thread models.Thread) (models.Thread, int) {
	query := `SELECT slug, author
				FROM threads 
				WHERE lower(slug) = lower($1);`

	row := dbPool.QueryRow(query, thread.Slug)

	err := row.Scan(&thread.Slug, &thread.Author)

	if err != nil {
		return thread, NotFound
	}

	return thread, OK
}
