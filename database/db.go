package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Database struct {
	// connection
	conn *sql.DB
}

func sanitizeString(input string) (result string) {
	result = input
	result = strings.Replace(result, "'", "''", -1)
	return
}

func (database *Database) execQuery(query string) {
	_, err := database.conn.Exec(query)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (database *Database) Connect(fileName string) error {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	database.conn = db

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" global_vars(name TEXT NOT NULL PRIMARY KEY" +
		",integer_value INTEGER" +
		",string_value TEXT" +
		")")

	// the same user in two different chats is treated as two different users
	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" users(messenger_id INTEGER NOT NULL" +
		",chat_id INTEGER NOT NULL" +
		",score INTEGER NOT NULL" +
		",name TEXT NOT NULL" +
		",PRIMARY KEY (messenger_id, chat_id)" +
		")")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" bets(id INTEGER NOT NULL PRIMARY KEY" +
		",chat_id INTEGER NOT NULL" +
		",bet_description TEXT NOT NULL" +
		")")

	return nil
}

func (database *Database) Disconnect() {
	database.conn.Close()
	database.conn = nil
}

func (database *Database) IsConnectionOpened() bool {
	return database.conn != nil
}

func (database *Database) createUniqueRecord(table string, values string) int64 {
	var err error
	if len(values) == 0 {
		_, err = database.conn.Exec(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES ", table))
	} else {
		_, err = database.conn.Exec(fmt.Sprintf("INSERT INTO %s VALUES (%s)", table, values))
	}

	if err != nil {
		log.Fatal(err.Error())
		return -1
	}

	rows, err := database.conn.Query(fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT 1", table))

	if err != nil {
		log.Fatal(err.Error())
		return -1
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err.Error())
			return -1
		}

		return id
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal("No record created")
	return -1
}

func (database *Database) GetDatabaseVersion() (version string) {
	rows, err := database.conn.Query("SELECT string_value FROM global_vars WHERE name='version'")

	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	var safeVersion string

	if rows.Next() {
		err := rows.Scan(&safeVersion)
		if err != nil {
			log.Fatal(err.Error())
		}
		version = strings.Replace(safeVersion, "_", ".", -1)
	} else {
		// that means it's a new clean database
		version = latestVersion
	}

	return
}

func (database *Database) SetDatabaseVersion(version string) {
	database.execQuery("DELETE FROM global_vars WHERE name='version'")

	safeVersion := sanitizeString(version)
	database.execQuery(fmt.Sprintf("INSERT INTO global_vars (name, string_value) VALUES ('version', '%s')", safeVersion))
}

func (database *Database) GetUserName(chatId int64, messengerUserId int64) (name string) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT name FROM users WHERE chat_id=%d AND messenger_id=%d",
		chatId,
		messengerUserId,
	))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	return
}

func (database *Database) UpdateUser(chatId int64, messengerUserId int64, name string) {
	sanitizedName := sanitizeString(name)

	database.execQuery(fmt.Sprintf(
		"INSERT OR IGNORE INTO users(messenger_id, chat_id, name, score) VALUES (%d, %d, '%s', 0);" +
		"UPDATE users SET name='%s' WHERE messenger_id=%d",
		messengerUserId,
		chatId,
		sanitizedName,
		sanitizedName,
		messengerUserId,
	))
}

func (database *Database) AddBet(chatId int64, messengerUserId int64, timeHours int, text string) {

}

func (database *Database) GetBetText(betId int64) (text string) {
	return
}

func (database *Database) GetActiveBets() (betIds []int64) {
	return
}

func (database *Database) GetBetRequirements(betId int64) (time int64) {
	return
}
