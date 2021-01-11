package statistics

// Statable is something that we can apply these functions on. Nababa!
type Statable interface {
	GetMean() (float64, error)
	GetMedian() (float64, error)
	GetStd() (float64, error)
	GetPercentile(p int) (float64, error)
	GetCount() (int, error)
}
