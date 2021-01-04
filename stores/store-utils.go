package stores

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type TransactionFunction func(*sqlx.Tx)

/*
Transaction Wrapper.
DB Calls within the transaction should panic on fail.  i.e. use MustExec vs Exec.
*/
func transaction(db *sqlx.DB, fn TransactionFunction) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		log.Printf("Unable to start transaction: %s\n", err)
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			err = r.(error)
			txerr := tx.Rollback()
			if txerr != nil {
				log.Printf("Unable to rollback from transaction: %s", txerr)
			}
		} else {
			txerr := tx.Commit()
			if txerr != nil {
				log.Printf("Unable to commit transaction: %s", txerr)
			}
		}
	}()
	fn(tx)
	return err
}
