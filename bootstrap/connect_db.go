package bootstrap

import (
	"database/sql"
	"fmt"
)

func ConnectDB(c *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DbHost, c.DbPort, c.DbLogin, c.DbPassword, c.DbName)

	db, err := sql.Open(c.DbScheme, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
