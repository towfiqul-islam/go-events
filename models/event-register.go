package models

import "example.com/rest-api/db"

type EventRegister struct {
	ID int64
	EventID int64
	UserID int64
}

func (ER *EventRegister) Register() error {
	query := `INSERT INTO events_registry (event_id, user_id) VALUES (?, ?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ER.EventID, ER.UserID)

	if err != nil {
		return err
	}

	return nil
}

func (ER *EventRegister) Cancel() error {
	query := `DELETE FROM events_registry WHERE event_id = ? AND user_id = ?`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ER.EventID, ER.UserID)

	if err != nil {
		return err
	}

	return nil
}