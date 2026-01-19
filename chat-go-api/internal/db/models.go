package db

import "time"

// Plan representa a tabela plans
type Plan struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       string  `json:"price"` // numeric do postgres vem como string
	Active      bool    `json:"active"`
}

// Schedule representa a tabela class_schedules
type Schedule struct {
	ID        int64      `json:"id"`
	Weekday   string     `json:"weekday"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Modality  *string    `json:"modality"`
	Coach     *string    `json:"coach"`
	Active    bool       `json:"active"`
}
