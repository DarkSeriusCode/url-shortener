package storage

import (
	"database/sql"
	"net/url"

	"github.com/DarkSeriusCode/url-shortener/internal/config"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(cfg *config.ServerConfig) (*PostgresStorage, error) {
	dbCfg := pq.Config{
		Host: cfg.Database.Host,
		Port: cfg.Database.Port,
		User: cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Dbname,
		SSLMode: "disable",
	}

	conn, err := pq.NewConnectorConfig(dbCfg)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	s := PostgresStorage{
		db: db,
	}

	if err = s.initTables(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *PostgresStorage) initTables() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		url TEXT NOT NULL,
		url_id TEXT NOT NULL,
		author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		token_hash TEXT NOT NULL UNIQUE,
		owner_id BIGINT NOT NULL REFERENCES users(id),
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	return err
}

func (s *PostgresStorage) Close() {
	s.db.Close()
}

func (s *PostgresStorage) NewUser(name string, password string) error {
	_, err := s.db.Exec(`INSERT INTO users (name, password) VALUES ($1, $2)`, name, password)
	return err
}

func (s *PostgresStorage) SaveUrl(originURL url.URL, urlID string) error {
	_, err := s.db.Exec(`INSERT INTO urls (url, url_id, author_id) VALUES ($1, $2, 0)`,
		originURL.String(), urlID)
	return err
}

func (s *PostgresStorage) GetUrl(urlID string) (*ShortenedURL, error) {
	var shortenedURL ShortenedURL
	var rawURL string
	err := s.db.QueryRow(`SELECT * FROM urls WHERE url_id = $1`, urlID).Scan(
		&shortenedURL.ID, &rawURL, &shortenedURL.UrlID,
	)
	if err != nil {
		return nil, err
	}
	shortenedURL.OriginURL, err = url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	shortenedURL.Author = nil
	return &shortenedURL, nil
}
