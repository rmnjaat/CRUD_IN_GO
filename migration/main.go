package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	CATALOG_DB_NAME := os.Getenv("CATALOG_DB_NAME")

	// run migration for catalog db

	DATABASE_DRIVER := "postgres"
	CATALOG_DATABASE_SOURCE := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, CATALOG_DB_NAME)
	log.Println("Catalog database source: ", CATALOG_DATABASE_SOURCE)
	log.Println("Running migration for catalog db...")
	goose.SetDialect(DATABASE_DRIVER)

	db, err := sql.Open(DATABASE_DRIVER, CATALOG_DATABASE_SOURCE)
	if err != nil {
		log.Fatal("Error connecting to catalog database: ", err)
	}
	defer db.Close()
	db = run_migration_catalog_db(db)

	// fetch all tenant db name from catalog db

	rows := get_db_names(db)

	run_tenant_migration(rows, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DATABASE_DRIVER)

}

func run_tenant_migration(rows *sql.Rows, DB_HOST string, DB_PORT string, DB_USER string, DB_PASSWORD string, DATABASE_DRIVER string) {
	for rows.Next() {
		var tenantDbName string
		err := rows.Scan(&tenantDbName)
		if err != nil {
			log.Fatal("Error scanning tenant db name: ", err)
		}

		log.Printf("Processing tenant database: %s", tenantDbName)

		// Check if database exists, create if not
		err = ensureDatabaseExists(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, tenantDbName)
		if err != nil {
			log.Printf("Error ensuring database exists for %s: %v", tenantDbName, err)
			continue
		}

		// Connect to tenant database
		tenantDatabaseSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, tenantDbName)
		tenantDb, err := sql.Open(DATABASE_DRIVER, tenantDatabaseSource)
		if err != nil {
			log.Printf("Error connecting to tenant database %s: %v", tenantDbName, err)
			continue
		}
		defer tenantDb.Close()

		err = tenantDb.Ping()
		if err != nil {
			log.Printf("Error pinging tenant database %s: %v", tenantDbName, err)
			continue
		}
		log.Printf("Successfully connected to tenant database: %s", tenantDbName)

		// Run migrations
		err = goose.Up(tenantDb, "tenant_db_migration")
		if err != nil {
			log.Printf("Error running migration for tenant db %s: %v", tenantDbName, err)
			continue
		}
		log.Printf("Migration completed for tenant db: %s", tenantDbName)
	}
}

func get_db_names(db *sql.DB) *sql.Rows {
	rows, err := db.Query("SELECT database_name FROM public.tenants")
	if err != nil {
		log.Fatal("Error fetching tenant db names: ", err)
	}
	defer rows.Close()
	return rows
}

func run_migration_catalog_db(db *sql.DB) *sql.DB {
	err := db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database (ping failed): ", err)
	}
	log.Println("Successfully connected to catalog database")

	err = goose.Up(db, "catalog_db_migrations") // for running migrations
	if err != nil {
		log.Fatal("Error running migration: catalog db", err)
	}
	log.Println("Migration completed for catalog db successfully")
	return db
}

// ensureDatabaseExists checks if database exists and creates it if not
func ensureDatabaseExists(host, port, user, password, dbName string) error {

	defaultDbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)
	defaultDb, err := sql.Open("postgres", defaultDbSource)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer defaultDb.Close()

	err = defaultDb.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping postgres database: %w", err)
	}

	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = defaultDb.QueryRow(checkQuery, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Printf("Database %s does not exist, creating it...", dbName)
		createQuery := fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
		_, err = defaultDb.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database %s: %w", dbName, err)
		}
		log.Printf("Database %s created successfully", dbName)
	} else {
		log.Printf("Database %s already exists", dbName)
	}

	return nil
}
