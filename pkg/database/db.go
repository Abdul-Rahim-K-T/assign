package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Debug: Log environment variables to ensure they're loaded correctly
	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))
	log.Println("DB_USER:", os.Getenv("DB_USER"))
	log.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))

	// Create a PostgreSQL connection string
	connStr := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	// Retry mechanism for database connection (useful in Docker environments)
	var db *sql.DB
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error opening DB connection (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		// Test the connection with Ping
		err = db.Ping()
		if err != nil {
			log.Printf("Error pinging DB (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
		} else {
			log.Println("Database connected successfully!")
			break
		}
	}

	if err != nil {
		log.Fatal("Failed to connect to the database after 5 attempts: ", err)
	}

	DB = db
	log.Println("Database connected")

	// Create table if they don't exist
	createTables()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL, -- Ensure email column is created
		address TEXT,
		user_type VARCHAR(50) NOT NULL,
		password_hash TEXT NOT NULL,
		profile_headline TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS profiles (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		resume_file TEXT,
		skills TEXT,
		education TEXT,
		experience TEXT,
		phone VARCHAR(15),
		name VARCHAR(255),
		email VARCHAR(255)
	);

	-- Create jobs table
	CREATE TABLE IF NOT EXISTS jobs (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		posted_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		total_applications INT DEFAULT 0,
		company_name VARCHAR(255) NOT NULL,
		posted_by_id INT NOT NULL,
		FOREIGN KEY (posted_by_id) REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS applications (
    	id SERIAL PRIMARY KEY,
    	user_id INT NOT NULL,
    	job_id INT NOT NULL,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    	FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
	);
	-- Add email column if not exists
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='email') THEN
			ALTER TABLE users ADD COLUMN email VARCHAR(255) UNIQUE NOT NULL;
		END IF;
	END $$;
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	log.Println("Tables created successfully (if not already existing)")
}
