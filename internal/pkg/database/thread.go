package database

import (
	"forum-db/internal/models"
	"github.com/jackc/pgx"
	"strconv"
	"sync"
	"time"
)

var x sync.Mutex

// /thread/{slug_or_id}/create
func CreateThreadPost(check string, posts []models.Post) ([]models.Post, int) {

	thread := models.Thread{}

	query := `SELECT id, forum
					FROM threads
					WHERE lower(slug) = lower($1)`

	row := dbPool.QueryRow(query, check)
	err := row.Scan(&thread.Id, &thread.Forum)

	if err != nil {
		return posts, NotFound
	}

	times := time.Now()

	if len(posts) == 0 {
		return posts, OK
	}


	tx, err := dbPool.Begin()
	query = `INSERT INTO posts (author, post, created_at, forum,  isEdited, parent, thread, path) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			RETURNING  id;`
	ins, _ := tx.Prepare("insert", query)


	result := []models.Post{}
	for _, p := range posts {

		pr := p
		pr.Forum = thread.Forum
		pr.Thread = thread.Id

		err := tx.QueryRow(ins.Name, p.Author, p.Message, times, thread.Forum, false, p.Parent, thread.Id, []int{}).Scan(&pr.Id)

		pr.CreatedAt = times.Format(time.RFC3339)

		if err != nil {
			tx.Rollback()
			if pqError, ok := err.(pgx.PgError); ok {
				switch pqError.Code {
				case "23503":
					return []models.Post{}, NotFound
				case "23505":
					return []models.Post{}, ForumConflict
				default:
					return []models.Post{}, ForumConflict
				}
			}
		}
		result = append(result, pr)
		info.Post++
	}

	tx.Commit()
	return result, OK
}

func GetThreadBySlugOrId(check string, thread models.Thread) (models.Thread, int) {
	var row *pgx.Row

	if value, err := strconv.Atoi(check); err != nil {
		thread.Slug = check
		query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE lower(slug) = lower($1)`
		row = dbPool.QueryRow(query, thread.Slug)

	} else {
		query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE id = $1`
		row = dbPool.QueryRow(query, value)
	}

	err := row.Scan(&thread.Id, &thread.Author, &thread.Message, &thread.Title,
		&thread.CreatedAt, &thread.Forum, &thread.Slug, &thread.Votes)

	if err != nil {
		return thread, NotFound
	}

	return thread, OK
}

// thread/{slug_or_id}/posts
func GetThreadsPosts(limit, since, desc, sort, check string) (models.Posts, int) {

	var row *pgx.Rows
	ps := models.Posts{}
	//TODO: получить только id
	thread, status := GetThreadBySlug(check, models.Thread{})
	if status == NotFound {
		return ps, NotFound
	}

	switch sort {
	case "flat":
		row = getFlat(thread.Id, since, limit, desc)

	case "tree":
		row = getTree(thread.Id, since, limit, desc)

	case "parent_tree":
		row = getParentTree(thread.Id, since, limit, desc)

	default:
		row = getFlat(thread.Id, since, limit, desc)
	}

	defer row.Close()
	for row.Next() {


		pr := models.Post{}
		times := time.Time{}
		err := row.Scan(&pr.Id, &pr.Author, &pr.Message, &times, &pr.Forum, &pr.IsEdited, &pr.Parent, &pr.Thread)
		pr.CreatedAt = times.Format(time.RFC3339)
		if err != nil {
		}
		ps = append(ps, pr)
	}

	return ps, OK
}

///thread/{slug_or_id}/vote
func ThreadVote(check string, vote models.Vote) (models.Thread, int) {

	thread, status := GetSlugID(check, models.Thread{})
	if status == NotFound {
		return thread, NotFound
	}

	query := `INSERT INTO VOTES (author, vote, thread) VALUES ($1, $2, $3)`
	_, err := dbPool.Exec(query, vote.NickName, vote.Voice, thread.Id)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case "23503":
				return thread, NotFound
			case "23505":
				upd := `UPDATE votes SET vote =  $1 WHERE author = $2 AND thread = $3`
				dbPool.Exec(upd, vote.Voice, vote.NickName, thread.Id)
			}
		}
	}

	thread, status = GetThreadByID(thread.Id, models.Thread{})
	return thread, OK
}

// /thread/{slug_or_id}/details
func ThreadUpdate(check string, thread models.Thread) (models.Thread, int) {
	t, status := GetThreadBySlugOrId(check, thread)
	if status == NotFound {
		return thread, NotFound
	}

	if thread.Message != "" {
		t.Message = thread.Message
	}

	if thread.Title != "" {
		t.Title = thread.Title
	}

	query := `UPDATE threads
	SET message=$1, title=$2
	WHERE id = $3
	RETURNING id, author, message, title, created_at, forum, slug, votes`

	row := dbPool.QueryRow(query, t.Message, t.Title, t.Id)
	res := models.Thread{}

	err := row.Scan(&res.Id, &res.Author, &res.Message, &res.Title,
		&res.CreatedAt, &res.Forum, &res.Slug, &res.Votes)
	if err == nil {
	}

	return res, OK
}

func ThreadVoteID(check int, vote models.Vote) (models.Thread, int) {

	   query := `INSERT INTO VOTES (author, vote, thread) VALUES ($1, $2, $3)`
		_, err := dbPool.Exec(query, vote.NickName, vote.Voice, check)

		if err != nil {
			if pqError, ok := err.(pgx.PgError); ok {
				switch pqError.Code {
				case "23503":
					return models.Thread{}, NotFound
				case "23505":
					upd := `UPDATE votes SET vote =  $1 WHERE author = $2 AND thread = $3`
					dbPool.Exec(upd, vote.Voice, vote.NickName, check)
				}
			}
		}

		thread, _ := GetThreadByID(check, models.Thread{})

		return thread, OK
	}


func CreateThreadPostID(id int, posts []models.Post) ([]models.Post, int) {
	thread := models.Thread{}

	query := `SELECT id, forum
					FROM threads
					WHERE id = $1`

	row := dbPool.QueryRow(query, id)
	err := row.Scan(&thread.Id, &thread.Forum)

	if err != nil {
		return posts, NotFound
	}

	times := time.Now()

	if len(posts) == 0 {
		return posts, OK
	}

	tx, err := dbPool.Begin()
	query = `INSERT INTO posts (author, post, created_at, forum,  isEdited, parent, thread, path) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			RETURNING  id;`
	ins, _ := tx.Prepare("insert", query)


	result := []models.Post{}
	for _, p := range posts {

		pr := p
		pr.Forum = thread.Forum
		pr.Thread = thread.Id

		err := tx.QueryRow(ins.Name, p.Author, p.Message, times, thread.Forum, false, p.Parent, thread.Id, []int{}).Scan(&pr.Id)

		pr.CreatedAt = times.Format(time.RFC3339)

		if err != nil {
			tx.Rollback()
			if pqError, ok := err.(pgx.PgError); ok {
				switch pqError.Code {
				case "23503":
					return []models.Post{}, NotFound
				case "23505":
					return []models.Post{}, ForumConflict
				default:
					return []models.Post{}, ForumConflict
				}
			}
		}
		result = append(result, pr)
		info.Post++
	}

	tx.Commit()
	return result, OK
}

func GetThreadBySlug(check string, thread models.Thread) (models.Thread, int) {
	var row *pgx.Row

	query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE lower(slug) = lower($1)`
	row = dbPool.QueryRow(query, check)

	err := row.Scan(&thread.Id, &thread.Author, &thread.Message, &thread.Title,
		&thread.CreatedAt, &thread.Forum, &thread.Slug, &thread.Votes)

	if err != nil {
		return thread, NotFound
	}

	return thread, OK
}

func GetThreadByID(id int, thread models.Thread) (models.Thread, int) {
	var row *pgx.Row

	query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE id = $1`
	row = dbPool.QueryRow(query, id)

	err := row.Scan(&thread.Id, &thread.Author, &thread.Message, &thread.Title,
		&thread.CreatedAt, &thread.Forum, &thread.Slug, &thread.Votes)

	if err != nil {
		return thread, NotFound
	}

	return thread, OK
}


// /thread/{slug_or_id}/details
func ThreadUpdateID(id int, thread models.Thread) (models.Thread, int) {
	t := models.Thread{}
	var row *pgx.Row

	query := `SELECT id, author, message, title, created_at, forum, slug, votes
					FROM threads
					WHERE id = $1`
	row = dbPool.QueryRow(query, id)

	err := row.Scan(&t.Id, &t.Author, &t.Message, &t.Title,
		&t.CreatedAt, &t.Forum, &t.Slug, &t.Votes)



	if err != nil {
		return thread, NotFound
	}


	if thread.Message != "" {
		t.Message = thread.Message
	}

	if thread.Title != "" {
		t.Title = thread.Title
	}

	query = `UPDATE threads
	SET message=$1, title=$2
	WHERE id = $3`

	_, err = dbPool.Exec(query, t.Message, t.Title, t.Id)

	if err == nil {
	}


	return t, OK
}

func GetThreadsPostsID(limit, since, desc, sort string, id int) (models.Posts, int) {

	var row *pgx.Rows
	ps := models.Posts{}
	thread := models.Thread{}

	query := `SELECT id	FROM threads WHERE id = $1`
	rows := dbPool.QueryRow(query, id)

	err := rows.Scan(&thread.Id)

	if err != nil {
		return ps, NotFound
	}

	switch sort {
	case "flat":
		row = getFlat(thread.Id, since, limit, desc)

	case "tree":
		row = getTree(thread.Id, since, limit, desc)

	case "parent_tree":
		row = getParentTree(thread.Id, since, limit, desc)

	default:
		row = getFlat(thread.Id, since, limit, desc)
	}

	defer row.Close()
	for row.Next() {

		pr := models.Post{}
		times := time.Time{}

		err := row.Scan(&pr.Id, &pr.Author, &pr.Message, &times, &pr.Forum, &pr.IsEdited, &pr.Parent, &pr.Thread)
		pr.CreatedAt = times.Format(time.RFC3339)
		if err != nil {
		}
		ps = append(ps, pr)
	}

	return ps, OK
}

func GetSlugID(check string, thread models.Thread) (models.Thread, int) {
	var row *pgx.Row

	query := `SELECT id FROM threads WHERE lower(slug) = lower($1)`
	row = dbPool.QueryRow(query, check)

	err := row.Scan(&thread.Id)

	if err != nil {
		return thread, NotFound
	}

	return thread, OK
}
