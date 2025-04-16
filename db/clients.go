package db

import "errors"

type DBClient struct {
	Id       int64
	FullName string
}

func (db *BookingDb) GetAllClients() ([]DBClient, error) {
	res, err := db.db.Query("select id, fullName from client")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	result := make([]DBClient, 0)
	for res.Next() {
		var cli DBClient
		err = res.Scan(&cli.Id, &cli.FullName)
		if err != nil {
			return nil, err
		}
		result = append(result, cli)
	}

	return result, nil
}

func (db *BookingDb) ClientInsert(fullName string) error {
	_, err := db.db.Exec("insert into client(fullName) values ($1)", fullName)
	return err
}

func (db *BookingDb) GetClient(id int64) (DBClient, error) {
	res, err := db.db.Query("select id, fullName from client where id = $1", id)
	if err != nil {
		return DBClient{}, err
	}
	defer res.Close()

	if !res.Next() {
		return DBClient{}, errors.New("client not found")
	}
	var result DBClient
	err = res.Scan(&result.Id, &result.FullName)
	if err != nil {
		return DBClient{}, err
	}
	return result, nil
}
