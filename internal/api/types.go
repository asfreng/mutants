package api

type EmptyResponse struct{}

type MutantPostRequest struct {
	Dna []string `json:"dna"`
}

type StatsResponse struct {
	CountMutantDna int     `json:"count_mutant_dna"`
	CountHumantDna int     `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}
