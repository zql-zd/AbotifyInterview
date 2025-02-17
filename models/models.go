package models

import "time"

type Advertiser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Budget    float64   `json:"budget"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Advertisement struct {
	ID           int64     `json:"id"`
	AdvertiserID int64     `json:"advertiser_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CPMCount     int64     `json:"cpm_count"` // Number of impressions
	CPCCount     int64     `json:"cpc_count"` // Number of clicks
	CPMRate      float64   `json:"cpm_rate"`  // Cost per thousand impressions
	CPCRate      float64   `json:"cpc_rate"`  // Cost per click
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
