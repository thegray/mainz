package store

import (
	"log"
	"pbk-main/pkg/db"
	"pbk-main/pkg/model/dto"
	"pbk-main/pkg/store/query"
	"pbk-main/pkg/util/dbutil"
)

type AccountStore struct {
	dbImpl *db.DBImpl
}

func NewAccountStore(dbImpl *db.DBImpl) *AccountStore {
	st := AccountStore{dbImpl: dbImpl}
	return &st
}

func (as *AccountStore) GetAccountById(id int) (*dto.Account, error) {
	conn, err := as.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][Account] Can't open db connection. %+v \n", err)
		return nil, err
	}
	defer conn.Close()

	found := &dto.Account{}
	stmt, err := conn.Prepare(query.GetAccountById)
	if err != nil {
		return nil, err
	}
	res := stmt.QueryRow(id)
	err = res.Scan(
		&found.ID,
		&found.Username,
		&found.EmailPhonenum,
		&found.Platform,
		&found.Details,
		&found.CreatedAt,
		&found.UpdatedAt,
		&found.CredID,
	)
	if err != nil {
		log.Printf("[Store][Account] Error scan data to dto %+v \n", err.Error())
		return nil, err
	}

	return found, nil
}

func (as *AccountStore) InsertNewAccount(s *dto.Account, catIds []int) (id int64, err error) {
	conn, err := as.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][Account][Insert] Can't open db connection. %+v \n", err)
		return -1, err
	}
	defer conn.Close()

	var lastTaskId int64
	err = dbutil.WithTransaction(conn, func(tx dbutil.Transaction) error {

		stmt, err := tx.Prepare(query.InsertAccount)
		if err != nil {
			return err
		}
		res, err := stmt.Exec(
			s.Username,
			s.EmailPhonenum,
			s.Platform,
			s.Details,
			s.CredID,
		)
		if err != nil {
			return err
		}
		lastTaskId, err = res.LastInsertId()
		if err != nil {
			return err
		}

		stmt, err = tx.Prepare(query.InsertAccToCategories)
		if err != nil {
			return err
		}
		for i := range catIds {
			_, err = stmt.Exec(catIds[i], lastTaskId)
			if err != nil {
				return err
			}
		}

		stmt, err = tx.Prepare(query.InsertAccountHistory)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			s.Username,
			s.EmailPhonenum,
			s.Platform,
			s.Details,
			lastTaskId,
			0,
			)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return -1, err
	}

	return lastTaskId, nil
}

func (as *AccountStore) InsertNewSafetyBox(sb *dto.SafetyBox) error {
	conn, err := as.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][Account][Insert SB] Can't open db connection. %+v \n", err)
		return err
	}
	defer conn.Close()

	stmt, err := conn.Prepare(query.InsertSafetyBox)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sb.SecretPass, sb.SecretInfo, sb.AccID)
	if err != nil {
		return err
	}

	return nil
}

func (as *AccountStore) UpdateAccount(id int, s *dto.Account) error {
	conn, err := as.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][Account][Update] Can't open db connection. %+v \n", err)
		return err
	}
	defer conn.Close()

	stmt, err := conn.Prepare(query.UpdateAccount)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.Username, s.EmailPhonenum, s.Platform, s.Details, s.CredID, id)
	if err != nil {
		return err
	}
	return nil
}

func (as *AccountStore) UpdateSecret(id int, sp string, si string) error {
	conn, err := as.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][Account][UpdSecret] Can't open db connection. %+v \n", err)
		return err
	}
	defer conn.Close()

	stmt, err := conn.Prepare(query.UpdateSecret)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(sp, si, id)
	if err != nil {
		return err
	}
	return nil
}

func (as *AccountStore) DeleteAccount(id int) (int64, error) {
	conn, err := as.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][Account][Delete] Can't open db connection. %+v \n", err)
		return -1, err
	}
	defer conn.Close()

	var rowsDeleted int64
	err = dbutil.WithTransaction(conn, func(tx dbutil.Transaction) error {
		stmt, err := tx.Prepare(query.DeleteSBToCategory)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(id)
		if err != nil {
			return err
		}

		stmt, err = tx.Prepare(query.DeleteAccount)
		if err != nil {
			return err
		}
		res, err := stmt.Exec(id)
		if err != nil {
			return err
		}
		rowsDeleted, err = res.RowsAffected()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return -1, err
	}
	return rowsDeleted, nil
}
