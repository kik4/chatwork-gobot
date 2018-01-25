package env

import (
	"log"

	"github.com/joho/godotenv"
)

// ChatWorkToken is gobot's token
var ChatWorkToken string

// RoomID is ID of general chat room
var RoomID string

// Load loads .env
func Load() {
	// load env
	envMap, err := godotenv.Read("my.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ChatWorkToken = envMap["ChatWorkToken"]
	if len(ChatWorkToken) < 1 {
		panic("ChatWorkToken is not found in .env file")
	}
	RoomID = envMap["RoomID"]
	if len(RoomID) < 1 {
		panic("RoomID is not found in .env file")
	}
}
