package handlers

import (
	"fmt"
	"github.com/bndrmrtn/playword/database"
	"github.com/bndrmrtn/playword/helpers"
	"github.com/bndrmrtn/playword/http/session"
	"github.com/bndrmrtn/playword/models"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

func newData(c *fiber.Ctx) map[string][]string {
	var randomWord models.Word
	database.Database.Db.Order("RAND()").First(&randomWord)

	if randomWord.Id == 0 {
		panic("Failed to find a random word")
	}
	letters := strings.Split(randomWord.Word, "")
	err := session.Set("word_data", map[string]interface{}{
		"word_id": randomWord.Id,
		"letters": letters,
	}, c)
	if err != nil {
		panic("Failed to set session")
	}

	return map[string][]string{
		"letters": letters,
	}
}

func getData(c *fiber.Ctx) interface{} {
	wordData, err := session.Get("word_data", c)
	if err != nil || wordData == nil {
		wordData = newData(c)
	}

	fmt.Println("Word data: ", wordData)
	mappedWordData := wordData.(map[string]interface{})

	gameData, err2 := session.Get("game_data", c)
	if err2 != nil {
		gameData = map[string]interface{}{
			"len":        len(mappedWordData["letters"].([]string)),
			"maxtrials":  3,
			"trials":     0,
			"guessed":    0,
			"expiration": time.Now().Add(time.Minute * 10),
		}
		_ = session.Set("game_data", gameData, c)
	}

	mappedGameData := gameData.(map[string]interface{})

	return helpers.MergeMaps(map[string]interface{}{
		"letters": mappedWordData["letters"],
	}, mappedGameData)
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
