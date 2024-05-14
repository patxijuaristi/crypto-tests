package scripts

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() *sql.DB {
	database, _ := sql.Open("sqlite3", "././crypto_tests.db")

	statement, _ := database.Prepare(`
		CREATE TABLE IF NOT EXISTS crypto_tests_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			test_name TEXT,
			algorithm_name TEXT,
			datetime DATETIME,
			duration INT,
			alloc INT,
			total_alloc INT,
			sys INT,
			num_gc INT,
			hash TEXT
		);
	`)
	statement2, _ := database.Prepare(`
		CREATE TABLE IF NOT EXISTS key_sig_sizes_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			algorithm_name TEXT,
			private_key INT,
			public_key INT,
			signature INT
		);
	`)
	statement.Exec()
	statement2.Exec()

	if IsSizesTableEmpty(database) {
		GenerateSizesDB(database)
	}

	return database
}

func InsertDB(database *sql.DB, testName string, currentTime time.Time, result map[string]interface{}) {

	statement, _ := database.Prepare(`
		INSERT INTO crypto_tests_data (test_name, algorithm_name, datetime,
										duration, alloc, total_alloc, sys, num_gc)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	statement.Exec(testName, result["algorithm"].(string), currentTime, result["execution_time"].(int), result["alloc_kb"].(int), result["total_alloc_kb"].(int), result["sys_kb"].(int), result["num_gc"].(int))
}

func InsertDBHash(database *sql.DB, testName string, currentTime time.Time, result map[string]interface{}, hash string) {

	statement, _ := database.Prepare(`
		INSERT INTO crypto_tests_data (test_name, algorithm_name, datetime,
										duration, alloc, total_alloc, sys, num_gc, hash)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	statement.Exec(testName, result["algorithm"].(string), currentTime, result["execution_time"].(int), result["alloc_kb"].(int), result["total_alloc_kb"].(int), result["sys_kb"].(int), result["num_gc"].(int), hash)
}

func IsSizesTableEmpty(database *sql.DB) bool {
	var count int
	err := database.QueryRow("SELECT COUNT(*) FROM key_sig_sizes_data").Scan(&count)
	if err != nil {
		fmt.Println("Error checking table:", err)
		return true // Assuming an error means the table is empty
	}
	return count == 0
}

func InsertSizesDB(database *sql.DB, result map[string]interface{}) {

	statement, _ := database.Prepare(`
		INSERT INTO key_sig_sizes_data (algorithm_name, private_key, public_key, signature)
				VALUES (?, ?, ?, ?)
	`)
	statement.Exec(result["algorithm"].(string), result["private_key"].(int), result["public_key"].(int), result["signature"].(int))
}
