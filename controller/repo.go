package controller

import "time"

type Repo struct {
	Username string        `json:"username"`
	Repos    int           `json:"repos"`
	Cached   bool          `json:"cached"`
	Duration time.Duration `json:"duration"`
	Faster   int           `json:"faster"`
}
