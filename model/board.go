package model

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kdakean/kdakean/db"
	"github.com/kdakean/kdakean/errors"
	"github.com/kdakean/kdakean/model/types"
)

type Board struct {
	Id          uint             `json:"id"`
	Name        string           `json:"name" valid:"required,maxlen(150)"`
	Description types.NullString `json:"desc" valid:"maxlen(500)"`
	UserId      uint             `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

type CreateBoardJSON struct {
	Name        string `json:"name" valid:"required,maxlen(150)"`
	Description string `json:"desc" valid:"maxlen(500)"`
	UserId      uint   `json:"user_id" valid:"required"`
}

type UpdateBoardJSON struct {
	Name        string `json:"name" valid:"required,maxlen(150)"`
	Description string `json:"desc" valid:"maxlen(500)"`
}

func CreateBoard(b *Board) (*Board, error) {
	createdAt := time.Now().UTC()
	columns := []string{
		"name", "description", "user_id",
		"created_at", "updated_at",
	}
	values := []interface{}{
		b.Name, b.Description, b.UserId,
		createdAt, createdAt,
	}
	sql, args, err := db.SQ.Insert("boards").Columns(columns...).
		Values(values...).Suffix("RETURNING id").ToSql()
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	board := Board{
		Name:        b.Name,
		Description: b.Description,
		UserId:      b.UserId,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	if err := db.DBX.QueryRowx(sql, args...).Scan(&board.Id); err != nil {
		if isReferencesError(err) {
			return nil, err
		} else {
			return nil, errors.ErrInternalServer
		}
	}
	return &board, nil
}

func UpdateBoard(b *Board) (*Board, error) {
	var board Board
	sql, args, err := db.SQ.Select("*").From("boards").
		Where(sq.Eq{"id": b.Id}).ToSql()
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if err := db.DBX.QueryRowx(sql, args...).StructScan(&board); err != nil {
		return nil, errors.ErrRecordNotFound
	}

	setMap := make(map[string]interface{})
	if b.Name != board.Name {
		setMap["name"] = b.Name
		board.Name = b.Name
	}
	if b.Description.Valid && b.Description.String != board.Description.String {
		setMap["description"] = b.Description.String
		board.Description = b.Description
	}

	if len(setMap) != 0 {
		updatedAt := time.Now().UTC()
		setMap["updated_at"] = updatedAt
		board.UpdatedAt = updatedAt

		sql, args, err := db.SQ.Update("boards").SetMap(setMap).
			Where(sq.Eq{"id": b.Id}).ToSql()
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		if _, err := db.DBX.Exec(sql, args...); err != nil {
			return nil, errors.ErrInternalServer
		}
	}

	return &board, nil
}
