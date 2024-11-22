package models

import "time"

type Job struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	PostedOn          time.Time `json:"posted_on"`
	TotalApplications int       `json:"total_applications"`
	CompanyName       string    `json:"company_name"`
	PostedByID        int       `json:"posted_by_id"` // Foreign key referencing User
	PostedBy          User      `json:"posted_by"`    // To load the user who posted the job
}

type Application struct {
	ID     int `json:"id" gorm:"primary_key"`
	UserID int `json:"user_id"`
	JobID  int `json:"job_id"`
}
