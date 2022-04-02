package handler

import (
	"challenge/internal/api"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Handler) MutantsPost(c *gin.Context) {
	fmt.Printf("Handle POST Mutants\n")

	var mutantsRequest api.MutantPostRequest
	if err := c.ShouldBindJSON(&mutantsRequest); err != nil {
		fmt.Printf("Error reading expected JSON body: %s\n", err)

		e := &api.ErrorBadRequest{}
		c.JSON(e.GetCode(), e)
		return
	}

	dnaType, err := p.LogicHandler.AddDNA(c, mutantsRequest.Dna)
	if err != nil {
		fmt.Printf("Error adding dna: %s\n", err.Error())
		c.JSON(err.GetCode(), err)
		return
	}

	fmt.Printf("Response DNA Type: %s\n", dnaType)

	if (dnaType == api.MUTANT) {
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusForbidden)
	}
}
