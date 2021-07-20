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
