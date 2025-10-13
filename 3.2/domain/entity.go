package domain

import "time"

type Link struct {
	Line string `json:"link"`
}

type ShortLink struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	ClickCount  int       `json:"click_count"`
}

type Visit struct {
	ID          int64     `json:"id"`
	ShortLinkID int64     `json:"short_link_id"`
	VisitedAt   time.Time `json:"visited_at"`
	UserAgent   string    `json:"user_agent"`
	IPAddress   string    `json:"ip_address"`
	DeviceType  string    `json:"device_type,omitempty"`
}

type VisitStats struct {
	TotalVisits     int            `json:"total_visits"`
	UniqueIPs       int            `json:"unique_ips"`
	DailyActivity   map[string]int `json:"daily_activity"`
	MonthlyActivity map[string]int `json:"monthly_activity"`
	DeviceStats     map[string]int `json:"device_stats"`
	Visits          []Visit        `json:"visits,omitempty"`
}
