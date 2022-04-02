package handler

import (
	"challenge/internal/logic"
	"challenge/internal/model"
)

type Handler struct {
	LogicHandler *logic.LogicHandler
}

// NewFraudcheckHandler .
func NewHandler(model *model.Model) *Handler {
	LogicHandler := logic.NewLogicHandler(model)

	return &Handler{LogicHandler: LogicHandler}
}
