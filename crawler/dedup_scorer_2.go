package crawler

import (
	"sort"
	"time"
)

// DeduplicationScore represents the score for a given item in the deduplication queue.
type DeduplicationScore struct {
	Item      interface{}
	Timestamp time.Time
	Score     float64
}

// DeduplicationScorer defines the interface for scoring items in the deduplication queue.
type DeduplicationScorer interface {
	Score(item interface{}, timestamp time.Time) float64
}

// HtmlScoreScorer implements DeduplicationScorer using HTML tags and content for scoring.
type HtmlScoreScorer struct{}

// Score calculates the score for the given item and timestamp.
func (s HtmlScoreScorer) Score(item interface{}, timestamp time.Time) float64 {
	// Assume item is *HtmlResolution
	htmlRes := item.(*HtmlResolution)
	// example scoring function: return the inverse of the number of unique tags and the length of the content
	return 1.0 / (float64(len(htmlRes.Tags)) + float64(len(htmlRes.Content)))
}

// ScoreItems returns a slice of DeduplicationScorer for the given items and timestamps.
func ScoreItems(items []interface{}, timestamps []time.Time, scorers ...DeduplicationScorer) []DeduplicationScore {
	scores := make([]DeduplicationScore, len(items))
	for i, item := range items {
		score := 0.0
		for _, scorer := range scorers {
			score += scorer.Score(item, timestamps[i])
		}
		scores[i] = DeduplicationScore{
			Item:      item,
			Timestamp: timestamps[i],
			Score:     score / float64(len(scorers)),
		}
	}
	return scores
}

// SortScores sorts the scores slice by score in ascending order.
func SortScores(scores []DeduplicationScore) {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score < scores[j].Score
	})
}
