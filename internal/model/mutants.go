package model

import (
	"context"
	"fmt"
	"strconv"

	"challenge/internal/api"

	"github.com/mediocregopher/radix/v3"
)

func (model *Model) GetDna(ctx context.Context, dnaHash string) (dnaType string, ok bool) {
	fmt.Printf("Model: GetDna.\n")

	key := getDnaKey(dnaHash)

	if err := model.db.Do(radix.Cmd(&dnaType, "GET", key)); err != nil {
		fmt.Printf("Error getting dnaType\n")
		return "", false
	}

	if dnaType != "" {
		return dnaType, true
	}

	return "", false
}

func (model *Model) SaveDna(ctx context.Context, dnaHash string, dnaType string) api.Error {
	fmt.Printf("Model: addDna.\n")

	key := getDnaKey(dnaHash)

	err := model.db.Do(radix.WithConn(key, func(conn radix.Conn) error {
		// Watch sobre la key que voy a agregar
		if err := conn.Do(radix.Cmd(nil, "WATCH", key)); err != nil {
			fmt.Printf("Error watching\n")
			return &api.ErrorInternalServerError{InternalErrorDescription: err.Error()}
		}

		// Obtengo el tipo de dna de la base (si existe)
		var dnaTypeDB string
		if err := conn.Do(radix.Cmd(&dnaTypeDB, "GET", key)); err != nil {
			fmt.Printf("Error getting DNA\n")
			return err
		}

		// No incremento el valor si ya lo agregaron
		if dnaTypeDB != "" {
			return nil
		}

		fmt.Printf("DNA type DB: %v\n", dnaTypeDB)

		// Begin transaction
		if err := conn.Do(radix.Cmd(nil, "MULTI")); err != nil {
			fmt.Printf("Error in MULTI\n")
			return err
		}

		// Agrego la key
		if err := conn.Do(radix.Cmd(nil, "SET", key, dnaType)); err != nil {
			fmt.Printf("Error setting DNA\n")
			return err
		}

		// Incremento el contador de stats
		dnaStatsKey := getDnaStatsKey()
		if err := conn.Do(radix.Cmd(nil, "HINCRBY", dnaStatsKey, dnaType, "1")); err != nil {
			fmt.Printf("Error setting DNA\n")
			return err
		}

		// Finish transaction
		var result []string
		if err := conn.Do(radix.Cmd(&result, "EXEC")); err != nil {
			fmt.Printf("Error while sending EXEC to Redis\n")
			return err
		}

		return nil
	}))

	if err != nil {
		fmt.Printf("Error while adding new dna to Redis. err: %v", err)
		return &api.ErrorInternalServerError{InternalErrorDescription: err.Error()}
	}

	return nil
}

func (model *Model) GetDnaCount(ctx context.Context) (mutantCount int, humanCount int, err api.Error) {
	fmt.Printf("Model: getDnaCount.\n")

	key := getDnaStatsKey()

	dnaStats := make(map[string]string)
	if err := model.db.Do(radix.Cmd(&dnaStats, "HGETALL", key)); err != nil {
		fmt.Printf("Error getting dnaStats\n")
		return 0, 0, &api.ErrorInternalServerError{InternalErrorDescription: err.Error()}
	}

	// Mutant count
	mutantCountStr, ok := dnaStats[api.MUTANT]
	if ok {
		var parseErr error
		mutantCount, parseErr = strconv.Atoi(mutantCountStr)
		if parseErr != nil {
			mutantCount = 0
		}
	} else {
		mutantCount = 0
	}

	// Human count
	humanCountStr, ok := dnaStats[api.HUMAN]
	if ok {
		var parseErr error
		humanCount, parseErr = strconv.Atoi(humanCountStr)
		if parseErr != nil {
			humanCount = 0
		}
	} else {
		humanCount = 0
	}

	return mutantCount, humanCount, nil
}
