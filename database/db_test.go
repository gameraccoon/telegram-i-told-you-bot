package database

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	testDbPath = "./testDb.db"
)

func dropDatabase(fileName string) {
	os.Remove(fileName)
}

func clearDb() {
	dropDatabase(testDbPath)
}

func connectDb(t *testing.T) *Database {
	assert := require.New(t)
	db := &Database{}

	err := db.Connect(testDbPath)
	if err != nil {
		assert.Fail("Problem with creation db connection:" + err.Error())
		return nil
	}
	return db
}

func createDbAndConnect(t *testing.T) *Database {
	clearDb()
	return connectDb(t)
}

func TestConnection(t *testing.T) {
	assert := require.New(t)
	dropDatabase(testDbPath)

	db := &Database{}

	assert.False(db.IsConnectionOpened())

	err := db.Connect(testDbPath)
	defer dropDatabase(testDbPath)
	if err != nil {
		assert.Fail("Problem with creation db connection:" + err.Error())
		return
	}

	assert.True(db.IsConnectionOpened())

	db.Disconnect()

	assert.False(db.IsConnectionOpened())
}

func TestDatabaseVersion(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}

	{
		version := db.GetDatabaseVersion()
		assert.Equal(latestVersion, version)
	}

	db.SetDatabaseVersion("1.0")

	{
		version := db.GetDatabaseVersion()
		assert.Equal("1.0", version)
	}

	db.SetDatabaseVersion("1.4")
	db.Disconnect()

	{
		db = connectDb(t)
		version := db.GetDatabaseVersion()
		assert.Equal("1.4", version)
		db.Disconnect()
	}
}

func TestUpdateUser(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId1 int64 = 321
	var chatId2 int64 = 123

	var userId1 int64 = 1234
	var userId2 int64 = 4321

	db.UpdateUser(chatId1, userId1, "test1")
	db.UpdateUser(chatId2, userId2, "test2")

	db.UpdateUser(chatId1, userId2, "test3")
	db.UpdateUser(chatId2, userId1, "test4")

	db.UpdateUser(chatId2, userId1, "test5")

	{
		assert.Equal("test5", db.GetUserName(chatId1, userId1))
		assert.Equal("test3", db.GetUserName(chatId1, userId2))
	}

	{
		assert.Equal("test5", db.GetUserName(chatId2, userId1))
		assert.Equal("test3", db.GetUserName(chatId2, userId2))
	}
}

func TestSanitizeString(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	testText := "text'test''test\"test\\"

	db.SetDatabaseVersion(testText)
	assert.Equal(testText, db.GetDatabaseVersion())
}

func TestAddBet(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId1 int64 = 321

	endTime := time.Now().Add(time.Hour)
	message := "test message"

	betId := db.AddBet(chatId1, endTime, message)

	{
		gotChatId, gotEndTime, gotMessage := db.GetBetData(betId)
		assert.Equal(chatId1, gotChatId)
		assert.True(endTime.Equal(gotEndTime))
		assert.Equal(message, gotMessage)
	}

	{
		activeBets := db.GetActiveBets()
		assert.Equal(1, len(activeBets))
		if len(activeBets) > 0 {
			assert.Equal(betId, activeBets[0])
		}
	}

	db.RemoveBet(betId)

	{
		activeBets := db.GetActiveBets()
		assert.Equal(0, len(activeBets))
	}
}
