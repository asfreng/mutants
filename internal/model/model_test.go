package model_test

import (
	"challenge/internal/api"
	"challenge/internal/handler"
	"challenge/internal/model"
	"challenge/internal/routers"
	"challenge/internal/util"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var mr *miniredis.Miniredis
var modelObj *model.Model
var handlerObj *handler.Handler

func TestMain(m *testing.M) {
	// Cargo la configuracion
	config, err := util.LoadConfig("../../configs", "challenge")
	if err != nil {
		fmt.Printf("Could not load the config file: %s\n", err.Error())
	}

	// Corro el modulo que mockea redis
	mr, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer mr.Close()

	// Conecto al mock de redis
	config.DatabaseConfig.Addresses = []string{mr.Addr()}
	config.DatabaseConfig.PoolSize = 10

	pool, err := util.ConnectToRedis(config.DatabaseConfig)
	if err != nil {
		fmt.Printf("Could not connect to Redis: %s\n", err.Error())
		return
	}
	defer util.CloseRedis(config.DatabaseConfig, pool)

	modelObj = model.NewModel(pool, config)
	handlerObj = handler.NewHandler(modelObj)
	router = routers.SetupRouter(handlerObj, config)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestGetEmptyStatsAndDna(t *testing.T) {
	mr.FlushAll()

	ctx := context.Background()

	dnaHash := "AAA"

	// Inserto un dna mutante
	_, getOk := handlerObj.LogicHandler.Model.GetDna(ctx, dnaHash)
	mutantCount, humanCount, statsErr := handlerObj.LogicHandler.Model.GetDnaCount(ctx)

	// Asserts
	assert.Equal(t, false, getOk)
	assert.Nil(t, statsErr)
	assert.Equal(t, 0, mutantCount)
	assert.Equal(t, 0, humanCount)
}

func TestSaveMutantDna(t *testing.T) {
	mr.FlushAll()

	ctx := context.Background()

	dnaHash := "AAA"

	// Inserto un dna mutante
	saveErr1 := handlerObj.LogicHandler.Model.SaveDna(ctx, dnaHash, api.MUTANT)
	saveErr2 := handlerObj.LogicHandler.Model.SaveDna(ctx, dnaHash, api.MUTANT)
	dnaGetType, getOk := handlerObj.LogicHandler.Model.GetDna(ctx, dnaHash)
	mutantCount, humanCount, statsErr := handlerObj.LogicHandler.Model.GetDnaCount(ctx)

	// Asserts
	assert.Nil(t, saveErr1)
	assert.Nil(t, saveErr2)
	assert.Equal(t, true, getOk)
	assert.Nil(t, statsErr)

	assert.Equal(t, api.MUTANT, dnaGetType)
	assert.Equal(t, 1, mutantCount)
	assert.Equal(t, 0, humanCount)
}

func TestSaveHumanDna(t *testing.T) {
	mr.FlushAll()

	ctx := context.Background()

	dnaHash1 := "AAA"
	dnaHash2 := "BBB"
	dnaHash3 := "CCC"

	// Inserto un dna mutante
	saveErr1 := handlerObj.LogicHandler.Model.SaveDna(ctx, dnaHash1, api.HUMAN)
	saveErr2 := handlerObj.LogicHandler.Model.SaveDna(ctx, dnaHash1, api.HUMAN)

	saveErr3 := handlerObj.LogicHandler.Model.SaveDna(ctx, dnaHash2, api.MUTANT)
	saveErr4 := handlerObj.LogicHandler.Model.SaveDna(ctx, dnaHash3, api.MUTANT)

	dnaGetType1, getOk1 := handlerObj.LogicHandler.Model.GetDna(ctx, dnaHash1)
	dnaGetType2, getOk2 := handlerObj.LogicHandler.Model.GetDna(ctx, dnaHash2)
	dnaGetType3, getOk3 := handlerObj.LogicHandler.Model.GetDna(ctx, dnaHash3)

	mutantCount, humanCount, statsErr := handlerObj.LogicHandler.Model.GetDnaCount(ctx)

	// Asserts
	assert.Nil(t, saveErr1)
	assert.Nil(t, saveErr2)
	assert.Nil(t, saveErr3)
	assert.Nil(t, saveErr4)
	assert.Equal(t, true, getOk1)
	assert.Equal(t, true, getOk2)
	assert.Equal(t, true, getOk3)
	assert.Nil(t, statsErr)

	assert.Equal(t, api.HUMAN, dnaGetType1)
	assert.Equal(t, api.MUTANT, dnaGetType2)
	assert.Equal(t, api.MUTANT, dnaGetType3)
	assert.Equal(t, 2, mutantCount)
	assert.Equal(t, 1, humanCount)
}