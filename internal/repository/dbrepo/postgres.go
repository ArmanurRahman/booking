package dbrepo

import (
	"context"
	"time"

	"github.com/ArmanurRahman/booking/internal/models"
)

func (m *postgressDBRepo) AllUsers() bool {
	return true
}

func (m *postgressDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	sql := `insert into reservations (first_name, last_name, email, phone, start_date, end_date,
		room_id, create_at, update_at) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id `

	err := m.DB.QueryRowContext(ctx, sql,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (m *postgressDBRepo) InsetIntoRoomRestriction(res models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id,
	create_at, update_at)
		values($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := m.DB.ExecContext(ctx, sql,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ResevationID,
		res.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgressDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	sql := `select count(id) from room_restrictions 
			where $1 < end_date and $2 > start_date and room_id = $3`

	row := m.DB.QueryRowContext(ctx, sql, start, end, roomId)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil

}

func (m *postgressDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `select r.id, r.room_name from rooms r where r.id not in 
	(select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

	var rooms []models.Room
	rows, err := m.DB.QueryContext(ctx, sql, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)

		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (m *postgressDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	sql := `select id, room_name, create_at, update_at from rooms where id = $1`
	row := m.DB.QueryRowContext(ctx, sql, id)

	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil

}