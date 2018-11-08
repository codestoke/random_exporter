package collector

import "math/rand"

type RandomClient struct {
	r *rand.Rand
}

func NewClient(seed int64) *RandomClient {
	rc := RandomClient{r: rand.New(rand.NewSource(seed))}
	return &rc
}

func (rc *RandomClient) GetMetrics() map[string]float64 {
	result := map[string]float64{}

	result["int10"] = float64(rc.r.Intn(10))

	return result
}
