package logic

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"strings"
)

type ResultNode struct {
	Character        string
	LeftStreak       int
	UpperLeftStreak  int
	UpperStreak      int
	UpperRightStreak int
}

type Direction int64

const (
	Left Direction = iota
	UpperLeft
	Upper
	UpperRight
)

func getDnaHash(dna []string) string {
	dnaStr := ""
	for _, row := range dna {
		dnaStr += row
	}

	md5Hash := md5.Sum([]byte(dnaStr))

	md5HashSlice := md5Hash[:]

	return base64.RawStdEncoding.EncodeToString(md5HashSlice)
}

func getUpdatedStreak(result [][]ResultNode, i, j int, direction Direction, character string, count, maxStreak int) (newStreak int, newCount int) {
	if i < 0 || j < 0 || j >= len(result) {
		return 1, count
	}

	node := result[i][j]
	adjacentCharacter := node.Character

	var streak int

	switch direction {
	case Left:
		streak = node.LeftStreak
	case UpperLeft:
		streak = node.UpperLeftStreak
	case Upper:
		streak = node.UpperStreak
	case UpperRight:
		streak = node.UpperRightStreak
	}

	if character == adjacentCharacter {
		streak++

		if streak%maxStreak == 0 {
			count++
		}
	} else {
		streak = 1
	}

	return streak, count
}

func IsMutant(dna []string, maxStreak int, maxCount int) bool {
	size := len(dna)

	if size < maxStreak {
		return false
	}

	count := 0

	result := make([][]ResultNode, size)

	for i := range result {
		result[i] = make([]ResultNode, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			currentCharacter := dna[i][j : j+1]

			var leftStreak, upperLeftStreak, upperStreak, upperRightStreak int

			// Left
			leftStreak, count = getUpdatedStreak(result, i, j-1, Left, currentCharacter, count, maxStreak)
			if count > maxCount {
				return true
			}

			// UpperLeft
			upperLeftStreak, count = getUpdatedStreak(result, i-1, j-1, UpperLeft, currentCharacter, count, maxStreak)
			if count > maxCount {
				return true
			}

			// Upper
			upperStreak, count = getUpdatedStreak(result, i-1, j, Upper, currentCharacter, count, maxStreak)
			if count > maxCount {
				return true
			}

			// UpperRight
			upperRightStreak, count = getUpdatedStreak(result, i-1, j+1, UpperRight, currentCharacter, count, maxStreak)
			if count > maxCount {
				return true
			}

			result[i][j] = ResultNode{
				Character:        currentCharacter,
				LeftStreak:       leftStreak,
				UpperLeftStreak:  upperLeftStreak,
				UpperStreak:      upperStreak,
				UpperRightStreak: upperRightStreak,
			}
		}
	}

	return false
}

func validateDna(dna []string) ([]string, error) {
	size := len(dna)
	validDna := make([]string, size)

	for i, row := range dna {
		upperCaseRow := strings.ToUpper(row)
		if !isValidRow(upperCaseRow) || len(upperCaseRow) != size {
			return nil, errors.New("INVALID ROW")
		}

		validDna[i] = upperCaseRow
	}

	return validDna, nil
}

func isValidLetter(c rune) bool {
	return c == 'A' || c == 'T' || c == 'C' || c == 'G'
}

func isValidRow(s string) bool {
	for _, c := range s {
		if !isValidLetter(c) {
			return false
		}
	}
	return true
}
