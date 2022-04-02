package logic

import (
	"challenge/internal/model"
)

// Logic handler
type LogicHandler struct {
	Model *model.Model
}

// NewFraudcheckHandler .
func NewLogicHandler(model *model.Model) *LogicHandler {
	logicHandler := LogicHandler{Model: model}

	return &logicHandler
}
