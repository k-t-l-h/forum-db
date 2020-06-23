package database

import (
	"errors"
	"forum-db/internal/models"
	"github.com/jackc/pgx"
	"time"
)

var dbPool *pgx.ConnPool
var info models.Status

func Open() (err error) {
	port := 5432

	connConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     uint16(port),
		Database: "docker",
		User:     "docker",
		Password: "docker",
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: 200,
		AcquireTimeout: 20 * time.Second,
		AfterConnect:   nil,
	}

	//panic(poolConfig.Host)
	if dbPool != nil {
		return errors.New("connection pool was already init")
	}

	dbPool, err = pgx.NewConnPool(poolConfig)
	if err != nil {
		return errors.New("connection cannot be established")
	}

	info = models.Status{
		Forum:  0,
		Post:   0,
		Thread: 0,
		User:   0,
	}
	return nil
}

func Close() {
	dbPool.Close()
}
