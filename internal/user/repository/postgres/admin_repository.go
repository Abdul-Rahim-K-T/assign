package postgres

import (
	"database/sql"
	"fmt"
	"recruitment-management/internal/models"
)

// GetJobByTitle retrieves a job by its title
func (r *UserRepositoryImpl) GetJobByTitleAndUserID(title string, userID int) (models.Job, error) {
	var job models.Job
	query := `SELECT id, title, description, posted_by_id FROM jobs WHERE title = $1 AND posted_by_id = $2`
	row := r.DB.QueryRow(query, title, userID)

	err := row.Scan(&job.ID, &job.Title, &job.Description, &job.PostedByID)
	if err != nil {
		if err == sql.ErrNoRows {
			return job, nil // No job found, return an empty job struct
		}
		return job, fmt.Errorf("failed to fetch job: %v", err) // Some other error occurred
	}
	return job, nil
}

// CreateJob saves a new job to the database
func (r *UserRepositoryImpl) CreateJob(job models.Job) error {
	query := `INSERT INTO jobs (title, description, posted_on, total_applications, company_name, posted_by_id) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.Exec(query, job.Title, job.Description, job.PostedOn, job.TotalApplications, job.CompanyName, job.PostedByID)
	if err != nil {
		return fmt.Errorf("failed to create job: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetUserByID(userID int, user *models.User) error {
	query := `SELECT id, name ,email, address, profile_headline FROM users WHERE id = $1`

	// Query the database to retrieve user details
	row := r.DB.QueryRow(query, userID)

	// Scan the result into the user struct
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.ProfileHeadline)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with id %d not found", userID) // User not found
		}
		return fmt.Errorf("failed to fetch user: %v", err)
	}
	return nil
}

// GetJobDetaials retrieves job details and list of applicants by job ID
func (r *UserRepositoryImpl) GetJobDetails(jobID int) (models.Job, []models.User, error) {
	var job models.Job
	var applicants []models.User

	// Query to fetch job details
	jobQuery := `SELECT id, title, description, posted_on, total_applications, company_name, posted_by_id FROM jobs WHERE id = $1`
	err := r.DB.QueryRow(jobQuery, jobID).Scan(&job.ID, &job.Title, &job.Description, &job.PostedOn, &job.TotalApplications, &job.CompanyName, &job.PostedByID)
	if err != nil {
		if err == sql.ErrNoRows {
			return job, nil, fmt.Errorf("job with id %d not found", jobID)
		}
		return job, nil, fmt.Errorf("failed to fetch job details: %v", err)
	}

	// Query to fetch applicants for the job
	applicantQuery := `SELECT u.id, u.name, u.email, u.address, u.profile_headline FROM users u
	INNER JOIN applications a ON u.id = a.user_id WHERE a.job_id = $1`
	rows, err := r.DB.Query(applicantQuery, jobID)
	if err != nil {
		return job, nil, fmt.Errorf("failed to fetch applicants: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var applicant models.User
		err := rows.Scan(&applicant.ID, &applicant.Name, &applicant.Email, &applicant.Address, &applicant.ProfileHeadline)
		if err != nil {
			return job, nil, fmt.Errorf("failed to scan applicant: %v", err)
		}
		applicants = append(applicants, applicant)
	}

	if err = rows.Err(); err != nil {
		return job, nil, fmt.Errorf("rows error: %v", err)
	}

	return job, applicants, nil
}

// Getalllusers retrieves the list of all users from the databas
func (r *UserRepositoryImpl) GetAllUsers() ([]models.User, error) {
	var users []models.User

	// Query to fetch alll users
	userQuery := `SELECT id, name, email, address, profile_headlin FROM users`
	rows, err := r.DB.Query(userQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.ProfileHeadline)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// GetProfileByID fetches the applicant's profile details by ID
func (r *UserRepositoryImpl) GetProfileByID(userID int, user *models.User) error {
	query := `
		SELECT u.id, u.name, u.email, u.address, u.profile_headline, 
		       p.id AS profile_id, p.resume_file, p.skills, p.education, p.experience, p.name AS profile_name, 
		       p.email AS profile_email, p.phone AS profile_phone
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		WHERE u.id = $1
	`

	// Query the database to retrieve user and profile details
	row := r.DB.QueryRow(query, userID)

	// Scan the result into the user struct
	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Address, &user.ProfileHeadline,
		&user.Profile.ID, &user.Profile.ResumeFile, &user.Profile.Skills, &user.Profile.Education,
		&user.Profile.Experience, &user.Profile.Name, &user.Profile.Email, &user.Profile.Phone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with id %d not found", userID)
		}
		return fmt.Errorf("failed to fetch user: %v", err)
	}
	return nil
}
