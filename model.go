package main

//Model

// People

// Person ...
type Person struct {
	UUID                    string  `json:"uuid"`
	Survived                bool    `json:"survived"`
	PassengerClass          int     `json:"passengerClass"`
	Name                    string  `json:"name"`
	Sex                     string  `json:"sex"`
	Age                     int     `json:"age"`
	SiblingsOrSpousesAboard int     `json:"siblingsOrSpousesAboard"`
	ParentsOrChildrenAboard int     `json:"parentsOrChildrenAboard"`
	Fare                    float32 `json:"fare"`
}

// PersonRepository resource

// Get a person data
func (s *APIServer) Get(u *Person) (*Person, error) {
	if err := s.storage.QueryRow(
		"INSERT INTO people (survived, passengerClass, name, sex, age, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uuid",
		u.Survived, u.Name, u.PassengerClass,
		u.Name, u.Sex, u.Age, u.SiblingsOrSpousesAboard,
		u.ParentsOrChildrenAboard, u.Fare,
	).Scan(&u.UUID); err != nil {
		return nil, err
	}
	return u, nil

}

// FindByUUID ...
func (s *APIServer) FindByUUID(uuid string) (*Person, error) {
	u := &Person{UUID: uuid}
	if err := s.storage.QueryRow("SELECT survived, pclass, name, sex, age, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare FROM people WHERE uuid = $1", uuid).Scan(
		&u.Survived, &u.PassengerClass, &u.Name,
		&u.Sex, &u.Age, &u.SiblingsOrSpousesAboard,
		&u.ParentsOrChildrenAboard, &u.Fare,
	); err != nil {
		return nil, err
	}

	return u, nil
}
