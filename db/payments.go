package db

import "errors"

var (
	PaymentCreated         = "CREATED"
	PaymentPartiallyFilled = "PARTIALLY_FILLED"
	PaymentFulfilled       = "FULFILLED"
)

type DBPayment struct {
	Id        int64
	BookingId int64
	Status    string
}

func (db *BookingDb) CreatePayment(bookingId int64) (int64, error) {
	var id int64
	err := db.db.QueryRow("insert into payment (booking_id, state) values ($1, $2) returning id",
		bookingId, PaymentCreated).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *BookingDb) GetPayment(bookingId int64) (DBPayment, error) {
	res, err := db.db.Query("select id, booking_id, state from payment where booking_id = $1", bookingId)
	if err != nil {
		return DBPayment{}, err
	}
	defer res.Close()

	if !res.Next() {
		return DBPayment{}, errors.New("payment not found")
	}

	var result DBPayment
	err = res.Scan(&result.Id, &result.BookingId, &result.Status)
	if err != nil {
		return DBPayment{}, err
	}

	return result, nil
}

func (db *BookingDb) UpdatePayment(id int64, status string) error {
	_, err := db.db.Exec("update payment set state = $1 where id = $2", status, id)
	return err
}
