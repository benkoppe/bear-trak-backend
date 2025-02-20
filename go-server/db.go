package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func connectToDbPool(ctx context.Context) (*pgxpool.Pool, error) {
	time.Sleep(3 * time.Second) // delay for db startup

	config := getDbConfig()

	pool, err := pgxpool.New(ctx, config.connStr())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("pgx"); err != nil {
		return nil, fmt.Errorf("unable to set dialect: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("unable to run migrations: %v", err)
	}

	return pool, nil
}

type dbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

func getDbConfig() dbConfig {
	password := ""
	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		passwordFile := os.Getenv("POSTGRES_PASSWORD_FILE")

		data, _ := os.ReadFile(passwordFile)
		password = strings.TrimSpace(string(data))
	}

	return dbConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: password,
		DbName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}

func (c *dbConfig) connStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DbName, c.SSLMode)
}
