package models

import (
	"database/sql"
	"errors"
	"time"

	"example.com/rest_api/database"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (e *Event) Save() error {
	//稍后，保存到数据库中
	query := `
	INSERT INTO events(name,description,location,dateTime,user_id)
	VALUES(?,?,?,?,?)
	`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := database.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	} //一行一行查下来，直到为NULL

	// rows.Err() 和上面 err != nil 检查的不是同一个错误来源：
	// Query 的 err 检查的是"执行查询/建立连接"阶段的错误，必须在循环前判断；
	// rows.Err() 检查的是"遍历 rows 内部"（Next）阶段可能发生的错误，
	// 只能放在循环结束后判断，否则永远拿不到它。二者不重复，都要保留。
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := database.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *Event) Update() error {
	query := `
	UPDATE events
	SET name = ?,description = ?,location = ?,dateTime = ?
	WHERE id = ?
	`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID) //传值给sql

	return err

}

func (e *Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID) //传id值给删除函数
	return err
}

func (e *Event) Register(userId int64) error {
	IsRegistered, err := IsUserRegistered(e.ID, userId)
	if err != nil {
		return err
	}

	var ErrAlreadyRegistered = errors.New("user already registered for this event")
	if IsRegistered {
		return ErrAlreadyRegistered
	}

	query := "INSERT INTO registrations (event_id,user_id) VALUES (?,?)"
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e *Event) CanceRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId) //传id值给删除函数
	return err
}

func GetRegistrationsForUser(userId int64) ([]Event, error) {
	query := `
	SELECT e.id, e.name, e.description, e.location, e.dateTime, e.user_id
	FROM events e
	JOIN registrations r ON e.id = r.event_id
	WHERE r.user_id = ?
	`
	rows, err := database.DB.Query(query, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	} //一行一行查下来，直到为NULL

	// rows.Err() 和上面 err != nil 检查的不是同一个错误来源：
	// Query 的 err 检查的是"执行查询/建立连接"阶段的错误，必须在循环前判断；
	// rows.Err() 检查的是"遍历 rows 内部"（Next）阶段可能发生的错误，
	// 只能放在循环结束后判断，否则永远拿不到它。二者不重复，都要保留。
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil

}

type Attendee struct {
	ID    int64
	Email string
}

func GetAttendees(eventId int64) ([]Attendee, error) {
	query := "SELECT u.id, u.email FROM users u JOIN registrations r ON u.id = r.user_id WHERE r.event_id = ?"
	rows, err := database.DB.Query(query, eventId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []Attendee

	for rows.Next() {
		var attendee Attendee
		err := rows.Scan(&attendee.ID, &attendee.Email)
		if err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return attendees, nil
}

func IsUserRegistered(eventId, userId int64) (bool, error) {
	query := "SELECT COUNT(*) FROM registrations WHERE event_id = ? AND user_id = ?"
	row := database.DB.QueryRow(query, eventId, userId)

	var count int
	err := row.Scan(&count)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
