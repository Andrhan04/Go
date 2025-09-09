package main

import (
	"database/sql"
)

// Create - Add a new cat
func CreateCat(cat *Cat) (int64, error) {
	result, err := db.Exec("INSERT INTO cat (name, age) VALUES (?, ?)",
		cat.Name, cat.Age)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Read - Get all cats
func GetAllCats() ([]Cat, error) {
	rows, err := db.Query("SELECT id, name, age FROM cat")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []Cat
	for rows.Next() {
		var cat Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Age)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	return cats, nil
}

// Read - Get cat by ID
func GetCatByID(id int) (Cat, error) {
	var cat Cat
	err := db.QueryRow("SELECT id, name, age FROM cat WHERE id = ?", id).
		Scan(&cat.ID, &cat.Name, &cat.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return Cat{}, nil
		}
		return Cat{}, err
	}

	return cat, nil
}

// Update - Update a cat
func UpdateCat(cat *Cat) error {
	_, err := db.Exec("UPDATE cat SET name = ?, age = ? WHERE id = ?",
		cat.Name, cat.Age, cat.ID)
	return err
}

// Delete - Delete a cat
func DeleteCat(id int) error {
	_, err := db.Exec("DELETE FROM cat WHERE id = ?", id)
	return err
}
