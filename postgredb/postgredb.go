package postgredb

import (
	"billing-api/model"
	"github.com/go-pg/pg"
	"os"
)

func New() *pg.DB {

	db := pg.Connect(&pg.Options{
		//Addr: getenv("DATABASE_SERVICE_HOST", "localhost") + ":" + getenv("DATABASE_SERVICE_DB", "5432"),
		Addr:     "db:5432",
		Database: getenv("DATABASE_SERVICE_DB", "app"),
		User:     getenv("DATABASE_SERVICE_USER", "root"),
		Password: getenv("DATABASE_SERVICE_PASSWORD", ""),
	})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	return db
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&model.Wallet{}, &model.Transaction{}} {
		db.DropTable(model, nil) //TODO warn!
		err := db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func TruncateTables(db *pg.DB) {
	db.Exec(`TRUNCATE TABLE wallets;`)
	db.Exec(`TRUNCATE TABLE transactions;`)
}
