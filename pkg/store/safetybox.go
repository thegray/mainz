package store

import (
	"log"
	"pbk-main/pkg/db"
	"pbk-main/pkg/model/dto"
	"pbk-main/pkg/store/query"
	"pbk-main/pkg/util/dbutil"
)

type SafetyBoxStore struct {
	dbImpl *db.DBImpl
}

func NewSafetyBoxStore(dbImpl *db.DBImpl) *SafetyBoxStore {
	st := SafetyBoxStore{dbImpl: dbImpl}
	return &st
}

func (sb *SafetyBoxStore) GetSafetyBoxById(id int) (*dto.SafetyBox, error) {
	conn, err := sb.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][SafetyBox] Can't open db connection. %+v \n", err)
		return nil, err
	}
	defer conn.Close()

	found := &dto.SafetyBox{}
	stmt, err := conn.Prepare(query.GetSafetyBoxById)
	if err != nil {
		return nil, err
	}
	res := stmt.QueryRow(id)
	err = res.Scan(
		&found.ID,
		&found.Username,
		&found.SecretPass,
		&found.Email,
		&found.Platform,
		&found.Details,
		&found.SecretInfo,
		&found.DateAdd,
		&found.DateModif,
		&found.CredID,
	)
	if err != nil {
		log.Printf("[Store][SafetyBox] Error scan data to dto %+v \n", err.Error())
		return nil, err
	}

	return found, nil
}

func (sb *SafetyBoxStore) InsertNewSafetyBox(s *dto.SafetyBox, catIds []int) (id int64, err error) {
	conn, err := sb.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][SafetyBox][Insert] Can't open db connection. %+v \n", err)
		return -1, err
	}
	defer conn.Close()

	var lastTaskId int64
	err = dbutil.WithTransaction(conn, func(tx dbutil.Transaction) error {

		stmt, err := tx.Prepare(query.InsertSafetyBox)
		if err != nil {
			return err
		}
		res, err := stmt.Exec(
			s.Username,
			s.SecretPass,
			s.Email,
			s.Platform,
			s.Details,
			s.SecretInfo,
			s.CredID,
		)
		if err != nil {
			return err
		}
		lastTaskId, err = res.LastInsertId()
		if err != nil {
			return err
		}

		stmt, err = tx.Prepare(query.InsertSBToCategory)
		if err != nil {
			return err
		}
		for i := range catIds {
			_, err = stmt.Exec(catIds[i], lastTaskId)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return -1, err
	}

	return lastTaskId, nil
}

func (sb *SafetyBoxStore) UpdateSafetyBox(id int, s *dto.SafetyBox) error {
	conn, err := sb.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][SafetyBox][Update] Can't open db connection. %+v \n", err)
		return err
	}
	defer conn.Close()

	stmt, err := conn.Prepare(query.UpdateSafetyBox)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.Username, s.Email, s.Platform, s.Details, s.CredID, id)
	if err != nil {
		return err
	}
	return nil
}

func (sb *SafetyBoxStore) UpdateSecret(id int, sp string, si string) error {
	conn, err := sb.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][SafetyBox][UpdSecret] Can't open db connection. %+v \n", err)
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

func (sb *SafetyBoxStore) DeleteSafetyBox(id int) (int64, error) {
	conn, err := sb.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][SafetyBox][Delete] Can't open db connection. %+v \n", err)
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

		stmt, err = tx.Prepare(query.DeleteSafetyBox)
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
