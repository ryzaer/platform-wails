package material

import (
	"errors"
	"fmt"
	"strings"

	"app-platform/internal/database"
)

type SaveRequest struct {
	ID         int64  `json:"id"`
	IDSupplier int64  `json:"id_supplier"`
	CodeQR     string `json:"code_qr"`
	ItemName   string `json:"item_name"`
	CreatedBy  string `json:"created_by"`
}

type IDsRequest struct {
	IDs []int64 `json:"ids"`
}

func Create(db *database.Connect, req SaveRequest) (int64, error) {
	normalize(&req)

	if err := validate(req); err != nil {
		return 0, err
	}

	if req.CreatedBy == "" {
		req.CreatedBy = "SYSTEM"
	}

	supplier, err := db.Query(`
		SELECT id
		FROM suppliers
		WHERE id = :id AND id > 0
		LIMIT 1
	`).Fetch("id", req.IDSupplier)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return 0, errors.New("supplier not found")
		}
		return 0, err
	}

	supplierID, ok := toInt64(supplier["id"])
	if !ok || supplierID <= 0 {
		return 0, errors.New("invalid supplier id")
	}

	id, err := db.Query(`
		INSERT INTO material
			(id_supplier, code_qr, item_name, status, created_by, created_at)
		VALUES
			(:id_supplier, :code_qr, :item_name, 1, :created_by, NOW())
	`).LastInsertID(
		"id_supplier", supplierID,
		"code_qr", req.CodeQR,
		"item_name", req.ItemName,
		"created_by", req.CreatedBy,
	)

	return id, err
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
		FROM material
		WHERE id = :id AND status = 1
		LIMIT 1
	`).Fetch("id", req.ID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return errors.New("material not found")
		}
		return err
	}

	if _, err := db.Query(`
		SELECT id
		FROM suppliers
		WHERE id = :id
		LIMIT 1
	`).Fetch("id", req.IDSupplier); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return errors.New("supplier not found")
		}
		return err
	}

	_, err := db.Query(`
		UPDATE material
		SET id_supplier = :id_supplier,
			code_qr = :code_qr,
			item_name = :item_name
		WHERE id = :id
	`).RowsAffected(
		"id_supplier", req.IDSupplier,
		"code_qr", req.CodeQR,
		"item_name", req.ItemName,
		"id", req.ID,
	)

	return err
}

func SoftDelete(db *database.Connect, req IDsRequest) (int64, error) {
	return setStatus(db, req.IDs, 0)
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
		UPDATE material
		SET status = :status
		WHERE id IN (` + strings.Join(placeholders, ", ") + `)
	`

	return db.Query(sql).RowsAffected(args...)
}

func normalize(req *SaveRequest) {
	req.CodeQR = strings.TrimSpace(req.CodeQR)
	req.ItemName = strings.TrimSpace(req.ItemName)
	req.CreatedBy = strings.ToUpper(strings.TrimSpace(req.CreatedBy))
}

func validate(req SaveRequest) error {
	if req.IDSupplier <= 0 {
		return errors.New("id_supplier is required")
	}

	if req.CodeQR == "" {
		return errors.New("code_qr is required")
	}

	if len(req.CodeQR) > 225 {
		return errors.New("code_qr must be at most 225 characters")
	}

	if req.ItemName == "" {
		return errors.New("item_name is required")
	}

	if len(req.ItemName) > 225 {
		return errors.New("item_name must be at most 225 characters")
	}

	if req.CreatedBy != "" && len(req.CreatedBy) > 7 {
		return errors.New("created_by must be at most 7 characters")
	}

	return nil
}

func toInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case int32:
		return int64(n), true
	case uint64:
		return int64(n), true
	case uint32:
		return int64(n), true
	case []byte:
		var x int64
		_, err := fmt.Sscan(string(n), &x)
		return x, err == nil
	default:
		return 0, false
	}
}
