package logic_test

import (
	"challenge/internal/api"
	"challenge/internal/handler"
	"challenge/internal/logic"
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

func TestIsMutant(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ATGCGAAA", "CAGTGCAA", "TTATGTAA", "AGAAGGAA", "CCCCTAAA", "TCACTGAA", "CCCCTAAA", "TCACTGAA"}

	// Inserto un dna mutante
	isMutant := logic.IsMutant(dna, 4, 1)

	// Asserts
	assert.Equal(t, true, isMutant)
}

func TestAddMutant(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ATGCGAAA", "CAGTGCAA", "TTATGTAA", "AGAAGGAA", "CCCCTAAA", "TCACTGAA", "CCCCTAAA", "TCACTGAA"}

	// Inserto un dna mutante
	dnaType, err := handlerObj.LogicHandler.AddDNA(context.Background(), dna)

	// Asserts
	assert.Nil(t, err)
	assert.Equal(t, api.MUTANT, dnaType)
}

func TestAddHuman(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}

	// Inserto un dna mutante
	dnaType, postErr := handlerObj.LogicHandler.AddDNA(context.Background(), dna)
	stats, statsErr := handlerObj.LogicHandler.GetDnaCount(context.Background())

	// Asserts de codes
	assert.Nil(t, postErr)
	assert.Nil(t, statsErr)

	// Asserts de data
	assert.Equal(t, api.HUMAN, dnaType)
	assert.Equal(t, 1, stats.CountHumantDna)
	assert.Equal(t, 0, stats.CountMutantDna)
	assert.Equal(t, float64(0), stats.Ratio)
}

func TestAddHumanSmall(t *testing.T) {
	mr.FlushAll()

	dna := []string{"AT", "CA"}

	// Inserto un dna mutante
	dnaType, postErr := handlerObj.LogicHandler.AddDNA(context.Background(), dna)
	stats, statsErr := handlerObj.LogicHandler.GetDnaCount(context.Background())

	// Asserts de codes
	assert.Nil(t, postErr)
	assert.Nil(t, statsErr)

	// Asserts de data
	assert.Equal(t, api.HUMAN, dnaType)
	assert.Equal(t, 1, stats.CountHumantDna)
	assert.Equal(t, 0, stats.CountMutantDna)
	assert.Equal(t, float64(0), stats.Ratio)
}

func TestAddInvalid(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ERROR", "CAGTGCAA", "TTATGTAA", "AGAAGGAA", "CCCCTAAA", "TCACTGAA", "CCCCTAAA", "TCACTGAA"}

	// Inserto un dna mutante
	_, err := handlerObj.LogicHandler.AddDNA(context.Background(), dna)

	// Asserts
	assert.NotNil(t, err)
}

func TestAddTwice(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ATGCGAAA", "CAGTGCAA", "TTATGTAA", "AGAAGGAA", "CCCCTAAA", "TCACTGAA", "CCCCTAAA", "TCACTGAA"}

	// Inserto un dna mutante
	dnaType, postErr := handlerObj.LogicHandler.AddDNA(context.Background(), dna)
	stats, statsErr := handlerObj.LogicHandler.GetDnaCount(context.Background())

	// Asserts de codes
	assert.Nil(t, postErr)
	assert.Nil(t, statsErr)

	// Asserts de data
	assert.Equal(t, api.MUTANT, dnaType)
	assert.Equal(t, 0, stats.CountHumantDna)
	assert.Equal(t, 1, stats.CountMutantDna)
	assert.Equal(t, float64(-1), stats.Ratio)

	dnaType, postErr = handlerObj.LogicHandler.AddDNA(context.Background(), dna)
	stats, statsErr = handlerObj.LogicHandler.GetDnaCount(context.Background())

	// Asserts de codes
	assert.Nil(t, postErr)
	assert.Nil(t, statsErr)

	// Asserts de data
	assert.Equal(t, api.MUTANT, dnaType)
	assert.Equal(t, 0, stats.CountHumantDna)
	assert.Equal(t, 1, stats.CountMutantDna)
	assert.Equal(t, float64(-1), stats.Ratio)
}
