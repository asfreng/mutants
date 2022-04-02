package model

import "fmt"

func getDnaKey(dna string) string {
	return fmt.Sprintf("dna:%s", dna)
}

func getDnaStatsKey() string {
	return "stats"
}