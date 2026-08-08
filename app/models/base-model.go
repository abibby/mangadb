package models

import (
	"context"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/model"
	"github.com/google/uuid"
)

type BaseModel struct {
	ID  uuid.UUID `json:"id" db:"id,primary"`
	ctx context.Context
}

var _ model.Model = &BaseModel{}
var _ model.Contexter = &BaseModel{}

func (b *BaseModel) InDatabase() bool {
	return b.ID != uuid.UUID{}
}

func (m *BaseModel) Context() context.Context {
	return m.ctx
}

func (m *BaseModel) AfterLoad(ctx context.Context, tx database.DB) error {
	m.ctx = ctx
	return nil
}
func (m *BaseModel) BeforeSave(ctx context.Context, tx database.DB) error {
	if !m.InDatabase() {
		m.ID = uuid.New()
	}
	return nil
}
func (m *BaseModel) AfterSave(ctx context.Context, tx database.DB) error {
	m.ctx = ctx
	return nil
}
