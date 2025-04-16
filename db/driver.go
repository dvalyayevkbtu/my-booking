package db

import (
	"database/sql"
	"dvalyayevkbtu/my-booking/config"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type BookingDb struct {
	db *sql.DB
}

func InitDatabase(conf config.DBConfig) (*BookingDb, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	result := &BookingDb{db}
	result.migrate()
	return result, nil
}

func (b *BookingDb) migrate() error {
	err := b.migrateClient()
	if err != nil {
		return err
	}
	err = b.migrateBook()
	if err != nil {
		return err
	}
	return b.migratePayment()
}

func (b *BookingDb) migrateClient() error {
	_, err := b.db.Exec(`
		create table if not exists client(
			id bigserial not null,
			fullName varchar(512) not null,

			constraint client_pk primary key (id)
		)
	`)
	if err != nil {
		return err
	}
	log.Debug("Clients migrated!")
	return nil
}

func (b *BookingDb) migrateBook() error {
	_, err := b.db.Exec(`
		create table if not exists booking(
			id bigserial not null,
			hotel_name varchar(512) not null,
			price varchar(256) not null,
			currency char(3) not null,
			client_id bigint not null,

			constraint booking_pk primary key (id),
			constraint booking_fk foreign key (client_id) references client(id)
		)
	`)
	if err != nil {
		return err
	}
	log.Debug("Booking migrated!")
	return nil
}

func (b *BookingDb) migratePayment() error {
	_, err := b.db.Exec(`
		create table if not exists payment(
			id bigserial not null,
			booking_id bigint not null,
			state varchar(16) not null,

			constraint payment_pk primary key (id),
			constraint payment_fk foreign key (booking_id) references booking(id)
		)
	`)
	if err != nil {
		return err
	}
	log.Debug("Payment migrated!")
	return nil
}
