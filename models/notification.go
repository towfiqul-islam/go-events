package models

import (
	"time"

	"example.com/rest-api/db"
)

type Notification struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	EventID   int64     `json:"event_id"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // "upcoming_event", "event_reminder", etc.
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Notification) Save() error {
	query := `
		INSERT INTO notifications (user_id, event_id, message, type, is_read, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(n.UserID, n.EventID, n.Message, n.Type, n.IsRead, n.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	n.ID = id
	return nil
}

func GetNotificationsByUserID(userID int64) ([]Notification, error) {
	query := `SELECT id, user_id, event_id, message, type, is_read, created_at 
			  FROM notifications WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.EventID,
			&notification.Message, &notification.Type, &notification.IsRead, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func MarkNotificationAsRead(notificationID int64) error {
	query := `UPDATE notifications SET is_read = true WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(notificationID)
	return err
}

// GetUpcomingEventsForNotification finds events that are within the next 24 hours
// and returns registered users for those events
func GetUpcomingEventsForNotification() ([]struct {
	EventID   int64
	EventName string
	DateTime  time.Time
	UserID    int64
}, error) {
	query := `
		SELECT e.id, e.name, e.dateTime, er.user_id
		FROM events e
		INNER JOIN events_registry er ON e.id = er.event_id
		WHERE e.dateTime BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 24 HOUR)
		AND NOT EXISTS (
			SELECT 1 FROM notifications n 
			WHERE n.event_id = e.id 
			AND n.user_id = er.user_id 
			AND n.type = 'upcoming_event'
			AND DATE(n.created_at) = CURDATE()
		)
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		EventID   int64
		EventName string
		DateTime  time.Time
		UserID    int64
	}

	for rows.Next() {
		var result struct {
			EventID   int64
			EventName string
			DateTime  time.Time
			UserID    int64
		}

		err := rows.Scan(&result.EventID, &result.EventName, &result.DateTime, &result.UserID)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}
