package repository

import (
	"time"

	"github.com/ArmanurRahman/booking/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsetIntoRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserById(id int) (models.User, error)
	UpdateUserById(user models.User) error
	Authenticate(email, testPassword string) (int, string, error)
}
