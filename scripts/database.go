package scripts

import (
	"database/sql"
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
	statement.Exec()

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
