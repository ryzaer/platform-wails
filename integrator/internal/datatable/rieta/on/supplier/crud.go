package supplier

import (
	"errors"
	"fmt"
	"strings"

	"app-platform/internal/database"
)

type SaveRequest struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Contact   string  `json:"contact"`
	Email     string  `json:"email"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Notes     string  `json:"notes"`
	Rating    float64 `json:"rating"`
	CreatedBy string  `json:"created_by"`
}

type IDsRequest struct {
	IDs []int64 `json:"ids"`
}

func Create(db *database.Connect, req SaveRequest) (int64, error) {
	normalize(&req)

	if err := validate(req); err != nil {
		return 0, err
	}

	return db.Query(`
		INSERT INTO suppliers (
			name,
			contact,
			email,
			address,
			latitude,
			longitude,
			notes,
			rating,
			status,
			created_by,
			created_at
		)
		VALUES (
			:name,
			:contact,
			:email,
			:address,
			:latitude,
			:longitude,
			:notes,
			:rating,
			1,
			:created_by,
			NOW()
		)
	`).LastInsertID(
		"name", req.Name,
		"contact", req.Contact,
		"email", req.Email,
		"address", req.Address,
		"latitude", req.Latitude,
		"longitude", req.Longitude,
		"notes", req.Notes,
		"rating", req.Rating,
		"created_by", req.CreatedBy,
	)
}

func Update(db *database.Connect, req SaveRequest) error {
	normalize(&req)

	if req.ID <= 0 {
		return errors.New("id is required")
	}

	if err := validate(req); err != nil {
		return err
	}

	if _, err := db.Query(`
		SELECT id
		FROM suppliers
		WHERE id = :id
		LIMIT 1
	`).Fetch("id", req.ID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return errors.New("supplier not found")
		}
		return err
	}

	_, err := db.Query(`
		UPDATE suppliers
		SET name = :name,
			contact = :contact,
			email = :email,
			address = :address,
			latitude = :latitude,
			longitude = :longitude,
			notes = :notes,
			rating = :rating
		WHERE id = :id
	`).RowsAffected(
		"name", req.Name,
		"contact", req.Contact,
		"email", req.Email,
		"address", req.Address,
		"latitude", req.Latitude,
		"longitude", req.Longitude,
		"notes", req.Notes,
		"rating", req.Rating,
		"id", req.ID,
	)

	return err
}

func Delete(db *database.Connect, req IDsRequest) (int64, error) {
	if len(req.IDs) == 0 {
		return 0, errors.New("ids is required")
	}

	placeholders := make([]string, len(req.IDs))
	args := make([]any, 0, len(req.IDs))

	for i, id := range req.IDs {
		if id <= 0 {
			return 0, errors.New("invalid id")
		}

		name := fmt.Sprintf("id_%d", i)

		placeholders[i] = ":" + name
		args = append(args, name, id)
	}

	sql := `
		DELETE FROM suppliers
		WHERE id IN (` + strings.Join(placeholders, ", ") + `)
	`

	return db.Query(sql).RowsAffected(args...)
}

func SoftDelete(db *database.Connect, req IDsRequest) (int64, error) {
	return setStatus(db, req.IDs, 0)
}

func normalize(req *SaveRequest) {
	req.Name = strings.TrimSpace(req.Name)
	req.Contact = strings.TrimSpace(req.Contact)
	req.Email = strings.TrimSpace(req.Email)
	req.Address = strings.TrimSpace(req.Address)
	req.Notes = strings.TrimSpace(req.Notes)
}

func validate(req SaveRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if len(req.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}

	if len(req.Email) > 150 {
		return errors.New("email must be at most 150 characters")
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		return errors.New("latitude tidak valid")
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		return errors.New("longitude tidak valid")
	}

	if req.Rating < 0 || req.Rating > 5 {
		return errors.New("rating harus antara 0 dan 5")
	}

	return nil
}

func Restore(db *database.Connect, req IDsRequest) (int64, error) {
	return setStatus(db, req.IDs, 1)
}

func setStatus(db *database.Connect, ids []int64, status int) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("ids is required")
	}

	placeholders := make([]string, len(ids))
	args := make([]any, 0, len(ids)*2)

	for i, id := range ids {
		if id <= 0 {
			return 0, errors.New("invalid id")
		}

		name := fmt.Sprintf("id_%d", i)

		placeholders[i] = ":" + name
		args = append(args, name, id)
	}

	args = append(args, "status", status)

	sql := `
		UPDATE suppliers
		SET status = :status
		WHERE id IN (` + strings.Join(placeholders, ", ") + `)
	`

	return db.Query(sql).RowsAffected(args...)
}
