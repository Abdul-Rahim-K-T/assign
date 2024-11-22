package models

import "time"

// User represents the application's users (Applicants or Admins)
type User struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Address         string    `json:"address"`
	UserType        string    `json:"user_type"` // Applicant/Admin
	PasswordHash    string    `json:"password_hash"`
	ProfileHeadline string    `json:"profile_headline"`
	CreatedAt       time.Time `json:"created_at"`
	Profile         Profile   `json:"profile"` // Embedded Profile for Applicants
}

type Profile struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"` // fk refrencing user
	ResumeFile string `json:"resume_file"`
	Skills     string `json:"skills"`
	Education  string `json:"education"`
	Experience string `json:"experience"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}
