package logic

import (
	"context"
	"fmt"

	"challenge/internal/api"
)

func (logic *LogicHandler) AddDNA(ctx context.Context, dna []string) (string, api.Error) {
	fmt.Printf("Logic: addDNA.\n")

	validDna, err := validateDna(dna)
	if err != nil {
		errDesc := "Error: Invalid DNA format."
		fmt.Println(errDesc)
		return "", &api.ErrorBadRequest{InternalErrorDescription: errDesc}
	}

	config := logic.Model.Config
	maxCount := config.MaxCount
	maxStreak := config.MaxStreak
	dnaHash := getDnaHash(dna)

	dnaType, ok := logic.Model.GetDna(ctx, dnaHash)
	if ok {
		// La clave ya esta guardada en la base, no vuelvo a calcularla
		fmt.Printf("Returning DNA type from db.\n")
		return dnaType, nil
	}

	// Lammo a isMutant para obtener el tipo de ADN y lo guardo en base
	fmt.Printf("Saving new DNA type.\n")
	if IsMutant(validDna, maxStreak, maxCount) {
		dnaType = api.MUTANT
	} else {
		dnaType = api.HUMAN
	}

	apiErr := logic.Model.SaveDna(ctx, dnaHash, dnaType)
	if err != nil {
		return "", apiErr
	}

	return dnaType, nil
}

func (logic *LogicHandler) GetDnaCount(ctx context.Context) (*api.StatsResponse, api.Error) {
	fmt.Printf("Logic: GetDnaCount.\n")

	mutantCount, humanCount, err := logic.Model.GetDnaCount(ctx)
	if err != nil {
		return nil, err
	}

	var ratio float64

	if humanCount == 0 {
		ratio = -1
	} else {
		ratio = float64(mutantCount) / float64(humanCount)
	}

	stats := &api.StatsResponse{
		CountHumantDna: humanCount,
		CountMutantDna: mutantCount,
		Ratio:          ratio,
	}

	return stats, nil
}
