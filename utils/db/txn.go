package db

import (
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TxnFunc func(txn boil.Transactor) error

func WithTxn(fn TxnFunc) (err error) {
	txn, err := boil.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			txn.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			txn.Rollback()
		} else {
			// all good, commit
			err = txn.Commit()
		}
	}()

	err = fn(txn)
	return err
}
