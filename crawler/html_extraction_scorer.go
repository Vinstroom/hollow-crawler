package crawler

import (
	"sort"
	"sync"
	"time"
)

// Score represents the score of an HTML extraction helper.
type Score struct {
	Helper   HTMLExtractionHelper
	Score    float64
	LastUsed time.Time
	Mutex    sync.Mutex
}

// HTMLExtractionScorer manages and scores HTML extraction helpers.
type HTMLExtractionScorer struct {
	helpers  []Score
	lastUse  map[string]time.Time
	scorer   func(score, oldScore float64) float64
	threshold float64
}

// NewHTMLExtractionScorer creates a new HTMLExtractionScorer with the given scorer function and threshold.
func NewHTMLExtractionScorer(scorer func(score, oldScore float64) float64, threshold float64) *HTMLExtractionScorer {
	return &HTMLExtractionScorer{
		scorer:   scorer,
		threshold: threshold,
		lastUse:  make(map[string]time.Time),
	}
}

// AddHelper adds a new HTML extraction helper to the scorer.
func (s *HTMLExtractionScorer) AddHelper(h HTMLExtractionHelper) {
	s.helpers = append(s.helpers, Score{Helper: h})
}

// ScoreHelpers calculates the score for each HTML extraction helper based on its usage and performance.
func (s *HTMLExtractionScorer) ScoreHelpers() {
	for _, h := range s.helpers {
		s.lastUse[h.Helper.Name()] = time.Now()
	//
	// ... Implement your scoring logic here using the scorer function and threshold ...
	//
		h.Score = s.scorer(h.Score, h.LastUsed.Sub(s.lastUse[h.Helper.Name()]).Seconds())
	}
}

// GetBestHelper returns the best performing HTML extraction helper based on its score.
func (s *HTMLExtractionScorer) GetBestHelper() HTMLExtractionHelper {
	sort.Slice(s.helpers, func(i, j int) bool {
		return s.helpers[i].Score > s.helpers[j].Score
	})
	if len(s.helpers) > 0 {
		return s.helpers[0].Helper
	}
	return nil
}
