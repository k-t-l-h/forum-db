package database

import (
	"forum-db/internal/models"
	"github.com/jackc/pgx"
	"sync/atomic"
)

var users map[string]models.User

func CreateUser(user models.User) ([]models.User, int) {
	results := []models.User{}
	result := user
	query := `INSERT INTO users (email, fullname, nickname, about) 
			VALUES ($1, $2, $3, $4) RETURNING nickname`
	row := dbPool.QueryRow(query, user.Email, user.FullName, user.NickName, user.About)

	err := row.Scan(&result.NickName)
	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case "23505":
				us, _ := GetUserOnConflict(user)
				return us, ForumConflict
			default:
				us, _ := GetUserOnConflict(user)
				return us, ForumConflict
			}
		}
	}

	atomic.AddInt64(&info.User, 1)
	results = append(results, result)
	return results, OK
}

func GetUserOnConflict(user models.User) ([]models.User, int) {
	results := []models.User{}
	query := `SELECT email, fullname, nickname, about
	FROM users
	WHERE lower(email) = lower($1) or  lower(nickname) =  lower($2)`

	rows, _ := dbPool.Query(query, user.Email, user.NickName)
	defer rows.Close()

	for rows.Next() {
		result := models.User{}
		rows.Scan(&result.Email, &result.FullName, &result.NickName, &result.About)
		results = append(results, result)
	}
	return results, OK
}

//GET /user/{nickname}/profile
func GetUser(user models.User) (models.User, int) {
	result := models.User{}
	query := `SELECT email, fullname, nickname, about
	FROM users
	WHERE  lower(nickname) =  lower($1)`

	rows := dbPool.QueryRow(query, user.NickName)

	err := rows.Scan(&result.Email, &result.FullName, &result.NickName, &result.About)
	if err != nil {
		return result, NotFound
	}
	return result, OK
}

func CheckUser(user models.User) (models.User, int) {
	result := models.User{}
	query := `SELECT nickname
	FROM users
	WHERE  lower(nickname) =  lower($1)`

	rows := dbPool.QueryRow(query, user.NickName)

	err := rows.Scan(&result.NickName)
	if err != nil {
		return result, NotFound
	}
	return result, OK
}

//POST /user/{nickname}/profile
func UpdateUser(user models.User) (models.User, int) {

	us, status := GetUser(user)
	if status == NotFound {
		return us, NotFound
	}

	if user.FullName != "" {
		us.FullName = user.FullName
	}
	if user.Email != "" {
		us.Email = user.Email
	}
	if user.About != "" {
		us.About = user.About
	}

	query := `UPDATE users 
	SET fullname=$1, email=$2, about=$3 
	WHERE nickname = $4
	RETURNING nickname, fullname, about, email;`

	rows := dbPool.QueryRow(query, us.FullName, us.Email, us.About, us.NickName)
	err := rows.Scan(&us.NickName, &us.FullName, &us.About, &us.Email)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case "23505":
				return us, ForumConflict
			case "23503":
				return us, UserNotFound
			}
		}
	}

	return us, OK
}
