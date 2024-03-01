package copy

import (
	"database/sql"
	"github.com/lib/pq"
	"go-sql/introspection"
)

// Item Columns and Values must be in the same order
type Item interface {
	GetValues(columns []string) []interface{}
	SetUpdatedAt()
}

func Copy[T Item](model T, db *sql.DB, tableName string, items []Item) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	columns := introspection.GetStructColumns(model)
	copyTx, err := tx.Prepare(pq.CopyIn(tableName, columns...))
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	errsCount := 0

	for _, item := range items {
		// force updated_at to be set at now
		item.SetUpdatedAt()
		_, err := copyTx.Exec(item.GetValues(columns)...)
		if err != nil {
			errsCount++
		}
	}
	err = copyTx.Close()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
