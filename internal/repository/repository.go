package repository

import "github.com/ArmanurRahman/booking/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsetIntoRoomRestriction(res models.RoomRestriction) error
}
