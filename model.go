package main

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

// AddPerson implements the function behind POST /people method
func (s *APIServer) AddPerson(u *Person) error {
	if err := s.storage.QueryRow(
		"INSERT INTO people (survived, pclass, name, sex, age, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uuid",
		u.Survived, u.PassengerClass, u.Name,
		u.Sex, u.Age, u.SiblingsOrSpousesAboard,
		u.ParentsOrChildrenAboard, u.Fare,
	).Scan(&u.UUID); err != nil {
		return err
	}
	return nil
}

// GetPerson implements the function behind GET /people/{uuid} method
func (s *APIServer) GetPerson(uuid string) (*Person, error) {
	u := &Person{UUID: uuid}
	var age int
	if err := s.storage.QueryRow("SELECT survived, pclass, name, sex, age, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare FROM people WHERE uuid = $1", uuid).Scan(
		&u.Survived, &u.PassengerClass, &u.Name,
		&u.Sex, &age, &u.SiblingsOrSpousesAboard,
		&u.ParentsOrChildrenAboard, &u.Fare,
	); err != nil {
		return nil, err
	}
	u.Age = int(age)
	return u, nil
}

// GetPeople implements the function behind GET /people method
func (s *APIServer) GetPeople() (*[]Person, error) {
	people := &[]Person{}
	rows, err := s.storage.Query("SELECT uuid, survived, pclass, name, sex, age, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare FROM people")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		u := Person{}
		var age float32
		err = rows.Scan(&u.UUID, &u.Survived, &u.PassengerClass, &u.Name, &u.Sex, &age, &u.SiblingsOrSpousesAboard, &u.ParentsOrChildrenAboard, &u.Fare)
		if err != nil {
			return nil, err
		}
		u.Age = int(age)
		*people = append(*people, u)
	}
	return people, nil
}

// UpdatePerson implements the function behind PUT /people/{uuid} method
func (s *APIServer) UpdatePerson(u *Person) error {
	if err := s.storage.QueryRow(
		"UPDATE people SET (survived, pclass, name, sex, age, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE uuid = $9",
		u.Survived, u.PassengerClass, u.Name,
		u.Sex, u.Age, u.SiblingsOrSpousesAboard,
		u.ParentsOrChildrenAboard, u.Fare, u.UUID).Scan(); err != nil {
		return err
	}
	return nil
}
