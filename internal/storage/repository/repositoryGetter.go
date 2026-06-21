package repository

import (
	"database/sql"
	"rent/internal/storage/db"
)

var apartmentRepository *ApartmentRepository = nil
var bookingRepository *BookingRepository = nil
var reviewRepository *ReviewRepository = nil
var userRepository *UserRepository = nil
var apartmentImageRepository *ApartmentImageRepository = nil
var avatarRepository *AvatarRepository = nil

func getDb() *sql.DB {
	d, err := db.GetSingletonDB()

	if err != nil {
		panic(err)
	}
	return d
}

func GetApartmentRepository() *ApartmentRepository {
	if apartmentRepository == nil {
		apartmentRepository = &ApartmentRepository{Db: getDb()}
	}

	return apartmentRepository
}

func GetBookingRepository() *BookingRepository {
	if bookingRepository == nil {
		bookingRepository = &BookingRepository{Db: getDb()}
	}

	return bookingRepository
}

func GetReviewRepository() *ReviewRepository {
	if reviewRepository == nil {
		reviewRepository = &ReviewRepository{Db: getDb()}
	}

	return reviewRepository
}

func GetUserRepository() *UserRepository {
	if userRepository == nil {
		userRepository = &UserRepository{Db: getDb()}
	}

	return userRepository
}

func GetApartmentImageRepository() *ApartmentImageRepository {
	if apartmentImageRepository == nil {
		apartmentImageRepository = &ApartmentImageRepository{Db: getDb()}
	}

	return apartmentImageRepository
}

func GetAvatarRepository() *AvatarRepository {
	if avatarRepository == nil {
		avatarRepository = &AvatarRepository{Db: getDb()}
	}

	return avatarRepository
}