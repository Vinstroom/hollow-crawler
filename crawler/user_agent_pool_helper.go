package crawler

import (
	"math/rand"
	"time"
)

// GetRandomUserAgent returns a random user agent from the provided pool.
// It uses the system time for seeding the random number generator to ensure different
// user agents are chosen across multiple crawler instances and over time.
func GetRandomUserAgent(pool []string) string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	return pool[rand.Intn(len(pool))]
}
