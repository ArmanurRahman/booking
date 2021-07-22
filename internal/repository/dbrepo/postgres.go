package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/ArmanurRahman/booking/internal/models"
	"golang.org/x/crypto/bcrypt"
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
			where $1 <= end_date and $2 >= start_date and room_id = $3`

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

func (m *postgressDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	sql := `select * from users where id=$1`

	row := m.DB.QueryRowContext(ctx, sql, id)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *postgressDBRepo) UpdateUserById(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `update user set first_name=$1, last_name=$2, email=$3, access_level=$4, update_at=$5
			where id=&6
	`
	_, err := m.DB.ExecContext(ctx, sql,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgressDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	sql := `select id, password from users where email=$1`
	row := m.DB.QueryRowContext(ctx, sql, email)

	err := row.Scan(
		&id,
		&hashedPassword,
	)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgressDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	sql := `select r.id, r.first_name, r.last_name, r.email, r.phone,	
			r.start_date, r.end_date, r.create_at, r.update_at, r.process,
			rm.id, rm.room_name 
			from reservations r left join rooms rm on r.room_id=rm.id`
	rows, err := m.DB.QueryContext(ctx, sql)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Process,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

func (m *postgressDBRepo) NewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	sql := `select r.id, r.first_name, r.last_name, r.email, r.phone,	
			r.start_date, r.end_date, r.create_at, r.update_at, r.process,
			rm.id, rm.room_name 
			from reservations r left join rooms rm on r.room_id=rm.id
			where process=0`
	rows, err := m.DB.QueryContext(ctx, sql)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Process,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

func (m *postgressDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	sql := `select r.id, r.first_name, r.last_name, r.email, r.phone,	
		r.start_date, r.end_date, r.create_at, r.update_at, r.process,
		rm.id, rm.room_name 
		from reservations r left join rooms rm on r.room_id=rm.id
		where r.id = $1`

	row := m.DB.QueryRowContext(ctx, sql, id)

	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Process,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)

	if err != nil {
		return reservation, err
	}
	return reservation, nil
}

func (m *postgressDBRepo) UpdateReservationById(reservation models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, update_at=$5
			where id=$6
	`
	_, err := m.DB.ExecContext(ctx, sql,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		time.Now(),
		reservation.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgressDBRepo) DeleteReservationById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `delete from reservations 
			where id=$1
	`
	_, err := m.DB.ExecContext(ctx, sql, id)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgressDBRepo) UpdateProcessedForReservation(process, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `update reservations set process=$1 
			where id=$2
	`
	_, err := m.DB.ExecContext(ctx, sql, process, id)

	if err != nil {
		return err
	}
	return nil
}
