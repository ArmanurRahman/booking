package dbrepo

import (
	"errors"
	"time"

	"github.com/ArmanurRahman/booking/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) InsetIntoRoomRestriction(res models.RoomRestriction) error {
	if res.RoomID == 1000 {
		return errors.New("some error")
	}

	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {

	return false, nil

}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("Some error")
	}
	return room, nil

}

func (m *testDBRepo) GetUserById(id int) (models.User, error) {

	var user models.User

	return user, nil
}

func (m *testDBRepo) UpdateUserById(user models.User) error {

	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	return 0, "", nil
}
