package db

import "errors"

type DBBooking struct {
	Id        int64
	HotelName string
	Price     string
	Currency  string
	ClientId  int64
}

func (db *BookingDb) GetAllBookings() ([]DBBooking, error) {
	res, err := db.db.Query("select id, hotel_name, price, currency, client_id from booking")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	bookings := make([]DBBooking, 0)
	for res.Next() {
		var book DBBooking
		err = res.Scan(&book.Id, &book.HotelName, &book.Price, &book.Currency, &book.ClientId)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, book)
	}

	return bookings, nil
}

func (db *BookingDb) GetBooking(id int64) (DBBooking, error) {
	res, err := db.db.Query("select id, hotel_name, price, currency, client_id from booking where id = $1", id)
	if err != nil {
		return DBBooking{}, err
	}
	defer res.Close()

	if !res.Next() {
		return DBBooking{}, errors.New("booking not found")
	}

	var result DBBooking
	err = res.Scan(&result.Id, &result.HotelName, &result.Price, &result.Currency, &result.ClientId)
	if err != nil {
		return DBBooking{}, err
	}
	return result, nil
}

func (db *BookingDb) RegisterBooking(hotelName, price, currency string, clientId int64) (int64, error) {
	var id int64
	err := db.db.QueryRow("insert into booking(hotel_name, price, currency, client_id) values ($1, $2, $3, $4) returning id",
		hotelName, price, currency, clientId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
