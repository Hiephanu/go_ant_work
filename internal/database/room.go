package database

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
)

type Room struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	HostId       string    `json:"host_id"`
	Participants []string  `json:"participants"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsPrivate    bool      `json:"is_private"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (r *service) CreateRoom(room *Room) (string, error) {
	query := `INSERT INTO rooms (id, name, host_id, participants, start_time, end_time, is_private, password, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(query, room.Id, room.Name, room.HostId, room.Participants, room.StartTime, room.EndTime, room.IsPrivate, room.Password, room.CreatedAt, room.UpdatedAt)
	if err != nil {
		return "", err
	}

	return room.Id, nil
}

func (r *service) FindRoomById(roomId string) (*Room, error) {
	query := `SELECT id, name, host_id, participants, start_time, end_time, is_private, password, created_at, updated_at 
	          FROM rooms WHERE id = $1`
	var room Room
	row := r.db.QueryRow(query, roomId)
	err := row.Scan(&room.Id, &room.Name, &room.HostId, &room.Participants, &room.StartTime, &room.EndTime, &room.IsPrivate, &room.Password, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	return &room, nil
}

func (r *service) FindAllRooms(page int64, perPage int64) ([]Room, error) {
	offset := (page - 1) * perPage
	query := `SELECT id, name, host_id, participants, start_time, end_time, is_private, password, created_at, updated_at
	          FROM rooms
	          LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.Id, &room.Name, &room.HostId, &room.Participants, &room.StartTime, &room.EndTime, &room.IsPrivate, &room.Password, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *service) UpdateRoom(room *Room) (*Room, error) {
	room.UpdatedAt = time.Now()

	query := `UPDATE rooms 
	          SET name = $1, host_id = $2, participants = $3, start_time = $4, end_time = $5, is_private = $6, password = $7, updated_at = $8
	          WHERE id = $9`

	_, err := r.db.Exec(query, room.Name, room.HostId, room.Participants, room.StartTime, room.EndTime, room.IsPrivate, room.Password, room.UpdatedAt, room.Id)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *service) DeleteRoom(roomId string) (string, error) {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.Exec(query, roomId)
	if err != nil {
		return "", err
	}

	return roomId, nil
}

func (r *service) FindRoomByHostId(hostId string) ([]Room, error) {
	query := `SELECT id, name, host_id, participants, start_time, end_time, is_private, password, created_at, updated_at
	          FROM rooms WHERE host_id = $1`
	rows, err := r.db.Query(query, hostId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.Id, &room.Name, &room.HostId, &room.Participants, &room.StartTime, &room.EndTime, &room.IsPrivate, &room.Password, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
