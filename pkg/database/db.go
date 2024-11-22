package database

import (
	"database/sql"
	"log"
	"os"

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

	// Create a PostgreSQL connection string
	connStr := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	// Open the connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
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
		phone VARCHAR(15)
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
