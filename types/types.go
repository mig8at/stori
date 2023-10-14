package types

type Process struct {
	TotalBalance        float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	TotalByMonth        map[string]int
}
