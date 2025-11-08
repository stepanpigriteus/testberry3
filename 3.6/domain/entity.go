package domain

import "time"

type Item struct {
	ID          string    `json:"id" db:"id"`
	Type        string    `json:"type" db:"type"`
	Category    string    `json:"category" db:"category"`
	Amount      float64   `json:"amount" db:"amount"`
	Date        time.Time `json:"date" db:"date"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type AnalyticsResult struct {
	Count        int              `json:"count"`
	Sum          float64          `json:"sum"`
	Avg          float64          `json:"avg"`
	Median       float64          `json:"median"`
	Percentile90 float64          `json:"percentile_90"`
	Grouped      []GroupedMetrics `json:"grouped,omitempty"`
}

type GroupedMetrics struct {
	Group string  `json:"group"`
	Sum   float64 `json:"sum"`
}
