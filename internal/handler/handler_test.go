package handler_test

import (
	"bytes"
	"challenge/internal/api"
	"challenge/internal/handler"
	"challenge/internal/model"
	"challenge/internal/routers"
	"challenge/internal/util"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func PostDna(router *gin.Engine, dna []string) (int) {
	DnaReq := api.MutantPostRequest{
		Dna: dna,
	}

	body, _ := json.Marshal(&DnaReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/mutant", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	return  w.Code
}

func PostDnaObject(router *gin.Engine, obj interface{}) (int) {
	body, _ := json.Marshal(&obj)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/mutant", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	return  w.Code
}

func GetStats(router *gin.Engine) (int, api.StatsResponse) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats", nil)

	router.ServeHTTP(w, req)

	statusCode := w.Code

	if statusCode != 200 {
		return statusCode, api.StatsResponse{}
	}

	JSONBody := w.Body
	responseBody := api.StatsResponse{}
	json.Unmarshal(JSONBody.Bytes(), &responseBody)

	return statusCode, responseBody
}

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

func TestPostMutant(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ATGCGAAA", "CAGTGCAA", "TTATGTAA", "AGAAGGAA", "CCCCTAAA", "TCACTGAA", "CCCCTAAA", "TCACTGAA"}

	// Inserto un dna mutante
	statusMutantCode := PostDna(router, dna)
	statusStatsCode, stats := GetStats(router)

	// Asserts de codes
	assert.Equal(t, 200, statusMutantCode)
	assert.Equal(t, 200, statusStatsCode)

	// Asserts de data
	assert.Equal(t, 0, stats.CountHumantDna)
	assert.Equal(t, 1, stats.CountMutantDna)
	assert.Equal(t, float64(-1), stats.Ratio)
}

func TestPostHuman(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}

	// Inserto un dna mutante
	statusMutantCode := PostDna(router, dna)
	statusStatsCode, stats := GetStats(router)

	// Asserts de codes
	assert.Equal(t, 403, statusMutantCode)
	assert.Equal(t, 200, statusStatsCode)

	// Asserts de data
	assert.Equal(t, 1, stats.CountHumantDna)
	assert.Equal(t, 0, stats.CountMutantDna)
	assert.Equal(t, float64(0), stats.Ratio)
}

func TestPostError(t *testing.T) {
	mr.FlushAll()

	dna := []string{"ERROR", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}

	// Inserto un dna mutante
	statusMutantCode := PostDna(router, dna)
	statusStatsCode, stats := GetStats(router)

	// Asserts de codes
	assert.Equal(t, 400, statusMutantCode)
	assert.Equal(t, 200, statusStatsCode)

	// Asserts de data
	assert.Equal(t, 0, stats.CountHumantDna)
	assert.Equal(t, 0, stats.CountMutantDna)
	assert.Equal(t, float64(-1), stats.Ratio)
}

func TestPostErrorBody(t *testing.T) {
	mr.FlushAll()

	// Inserto un dna mutante
	statusMutantCode := PostDnaObject(router, "ERROR")
	statusStatsCode, stats := GetStats(router)

	// Asserts de codes
	assert.Equal(t, 400, statusMutantCode)
	assert.Equal(t, 200, statusStatsCode)

	// Asserts de data
	assert.Equal(t, 0, stats.CountHumantDna)
	assert.Equal(t, 0, stats.CountMutantDna)
	assert.Equal(t, float64(-1), stats.Ratio)
}
