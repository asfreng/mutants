package model

import (
	"challenge/internal/util"
)

type Model struct {
	db     util.Pool
	Config util.Config
}

func NewModel(db util.Pool, config util.Config) *Model {
	return &Model{db: db, Config: config}
}
