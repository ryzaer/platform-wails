package stock

import (
	"errors"
	"fmt"
	"strings"

	"app-platform/internal/database"
)

type SaveRequest struct {
	ID         int64  `json:"id"`
	IDMaterial int64  `json:"id_material"`
	CodeBranch string `json:"code_branch"`
	IDSize     int64  `json:"id_size"`
	Quantity   int64  `json:"quantity"`
	SizePrice  int64  `json:"size_price"`
	FinalPrice int64  `json:"final_price"`
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

	material, err := db.Query(`SELECT id FROM material WHERE id = :id AND status = 1 LIMIT 1`).Fetch("id", req.IDMaterial)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return 0, errors.New("material not found")
		}
		return 0, err
	}
	materialID, ok := toInt64(material["id"])
	if !ok || materialID <= 0 {
		return 0, errors.New("invalid material id")
	}

	id, err := db.Query(`
        INSERT INTO material_criteria
            (id_material, code_branch, quantity, id_size, size_price, final_price, status, create_at, created_by)
        VALUES
            (:id_material, :code_branch, :quantity, :id_size, :size_price, :final_price, 1, NOW(), :created_by)
    `).LastInsertID(
		"id_material", materialID,
		"code_branch", req.CodeBranch,
		"quantity", req.Quantity,
		"id_size", req.IDSize,
		"size_price", req.SizePrice,
		"final_price", req.FinalPrice,
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
		SELECT id FROM material_criteria
		WHERE id = :id AND status = 1
		LIMIT 1
	`).Fetch("id", req.ID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return errors.New("stock not found")
		}
		return err
	}

	if _, err := db.Query(`
		SELECT id FROM material
		WHERE id = :id AND status = 1
		LIMIT 1
	`).Fetch("id", req.IDMaterial); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return errors.New("material not found")
		}
		return err
	}

	_, err := db.Query(`
        UPDATE material_criteria
        SET id_material = :id_material,
            code_branch = :code_branch,
            quantity = :quantity,
            id_size = :id_size,
            size_price = :size_price,
            final_price = :final_price
        WHERE id = :id
    `).RowsAffected(
		"id_material", req.IDMaterial,
		"code_branch", req.CodeBranch,
		"quantity", req.Quantity,
		"id_size", req.IDSize,
		"size_price", req.SizePrice,
		"final_price", req.FinalPrice,
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
	sql := "UPDATE material_criteria SET status = :status WHERE id IN (" + strings.Join(placeholders, ", ") + ")"
	return db.Query(sql).RowsAffected(args...)
}

func normalize(req *SaveRequest) {
	req.CodeBranch = strings.ToUpper(strings.TrimSpace(req.CodeBranch))
	req.CreatedBy = strings.ToUpper(strings.TrimSpace(req.CreatedBy))
}

func validate(req SaveRequest) error {
	if req.IDMaterial <= 0 {
		return errors.New("id_material is required")
	}
	if len(req.CodeBranch) == 0 || len(req.CodeBranch) > 3 {
		return errors.New("code_branch must be 1-3 characters")
	}
	if req.IDSize <= 0 {
		return errors.New("id_size is required")
	}
	if req.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	if req.SizePrice < 0 {
		return errors.New("size_price cannot be negative")
	}
	if req.FinalPrice < 0 {
		return errors.New("final_price cannot be negative")
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
