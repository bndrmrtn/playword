package handlers

import (
	"errors"
	"fmt"
	"github.com/bndrmrtn/playword/database"
	"github.com/bndrmrtn/playword/handlers/types"
	"github.com/bndrmrtn/playword/helpers"
	"github.com/bndrmrtn/playword/http/session"
	"github.com/bndrmrtn/playword/models"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

// MaxTrials The maximum trials to guess the word
const MaxTrials = 3

// newData Returns a new word from the database and saves it to the session
func newData(c *fiber.Ctx) types.WordData {
	var randomWord models.Word
	database.Database.Db.Order("RAND()").First(&randomWord)

	if randomWord.Id == 0 {
		panic("Failed to find a random word")
	}
	letters := strings.Split(randomWord.Word, "")
	wordData := types.WordData{
		WordId:  randomWord.Id,
		Letters: helpers.ShuffleStringSlice(letters),
	}

	err := session.Set("word_data", wordData, c)
	if err != nil {
		panic("Failed to set session")
	}

	return wordData
}

// getData Creates a new game or returns the existing one
func getData(c *fiber.Ctx) types.GameResponse {
	_wordData, err := session.Get("word_data", c)
	if err != nil || _wordData == nil {
		_wordData = newData(c)
	}
	wordData := _wordData.(types.WordData)

	var gameData types.GameData
	_gameData, err2 := session.Get("game_data", c)

	if err2 != nil || _gameData == nil {
		gameData = types.GameData{
			Length:     int16(len(wordData.Letters)),
			MaxTrials:  MaxTrials,
			Trials:     0,
			Guessed:    0,
			Expiration: time.Now().Add(time.Minute * 10),
		}
		_ = session.Set("game_data", gameData, c)
	} else {
		gameData = _gameData.(types.GameData)
	}

	return types.GameResponse{
		WordData: wordData,
		GameData: gameData,
	}
}

// GameHandler Returns the current game data to the frontend application
func GameHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": map[string]any{
			"code":    200,
			"message": "OK",
		},
		"data": getData(c),
	})
}

// findWord Search for a word in the database (returns an error if it doesn't exist)
func findWord(word string, order *models.Word) error {
	database.Database.Db.Find(&order, "word = ?", word)
	if order.Id == 0 {
		return errors.New("word does not exist")
	}
	return nil
}

// TODO: make it work ;)
func expirationChecker(time time.Time) error {
	return nil
}

// GameGuessHandler Handles the game main logic, checks if the received word is in the database, and manages the trials, then returns a response to the frontend application
func GameGuessHandler(c *fiber.Ctx) error {
	payload := struct {
		Guess string `json:"guess"`
	}{}

	gameData := getData(c)

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	var word models.Word
	database.Database.Db.Find(&word, "id = ?", gameData.WordData.WordId)
	if word.Id == 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	var guessWord models.Word
	err := findWord(payload.Guess, &guessWord)
	if err != nil || !helpers.CharsMatch(word.Word, guessWord.Word) {
		gameData.GameData.Trials++
		_ = session.Set("game_data", gameData.GameData, c)
		endGame := gameData.GameData.Trials == MaxTrials
		var message string
		if MaxTrials-gameData.GameData.Trials > 1 {
			message = fmt.Sprintf("No problem! You have %v more trials", MaxTrials-gameData.GameData.Trials)
		} else if MaxTrials-gameData.GameData.Trials == 1 {
			message = "No problem! You have 1 more trial"
		} else {
			message = "You lost! The word was \"" + word.Word + "\""
			_ = session.Delete("game_data", c)
		}
		_ = session.Delete("word_data", c)

		return c.JSON(types.GuessResponse{
			IsValid: false,
			Trials:  gameData.GameData.Trials,
			Endgame: endGame,
			Message: types.GuessMessage{
				Type: "error",
				Text: message,
			},
		})
	}

	_ = session.Delete("word_data", c)
	gameData.GameData.Trials = 0
	_ = session.Set("game_data", gameData.GameData, c)

	return c.JSON(types.GuessResponse{
		IsValid: true,
		Trials:  gameData.GameData.Trials,
		Endgame: false,
		Message: types.GuessMessage{
			Type: "success",
			Text: "Congratulation you\nguessed the word! :D",
		},
	})
}
