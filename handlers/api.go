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

const MAX_TRIALS = 3

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

func getData(c *fiber.Ctx) types.GameResponse {
	_wordData, err := session.Get("word_data", c)
	if err != nil || _wordData == nil {
		_wordData = newData(c)
	}
	fmt.Println(_wordData, "err")
	wordData := _wordData.(types.WordData)

	fmt.Println("Word data: ", wordData)

	var gameData types.GameData
	_gameData, err2 := session.Get("game_data", c)

	if err2 != nil || _gameData == nil {
		gameData = types.GameData{
			Length:     int16(len(wordData.Letters)),
			MaxTrials:  MAX_TRIALS,
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

func GameHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": map[string]any{
			"code":    200,
			"message": "OK",
		},
		"data": getData(c),
	})
}

func findWord(word string, order *models.Word) error {
	database.Database.Db.Find(&order, "word = ?", word)
	if order.Id == 0 {
		return errors.New("word does not exist")
	}
	return nil
}

func GameGuessHandler(c *fiber.Ctx) error {
	payload := struct {
		Guess string `json:"guess"`
	}{}

	gameData := getData(c)

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	// TODO: something works wrong here
	var word models.Word
	err := findWord(payload.Guess, &word)
	if err != nil {
		trials := gameData.GameData.Trials + 1
		endGame := trials == MAX_TRIALS
		var message string
		if MAX_TRIALS-trials > 1 {
			message = fmt.Sprintf("No problem! You have %v more trials", trials)
		} else if MAX_TRIALS-trials == 1 {
			message = "No problem! You have 1 more trial"
		} else {
			message = "You lost! The word was "
			_ = session.Delete("game_data", c)
		}

		_ = session.Delete("word_data", c)

		return c.JSON(types.GuessResponse{
			IsValid: false,
			Trials:  trials,
			Endgame: endGame,
			Message: types.GuessMessage{
				Type: "error",
				Text: message,
			},
		})
	}

	// TODO: check if the sent string has only the letters that the required word has
	_ = session.Delete("word_data", c)
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
