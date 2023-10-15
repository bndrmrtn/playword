package session

import (
	"encoding/gob"
	"github.com/bndrmrtn/playword/handlers/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

// This stores all of your app's sessions
var store *session.Store = session.New()

// Get tries to get a session by key, or returns an error
func Get(key string, c *fiber.Ctx) (interface{}, error) {
	sess, err := store.Get(c)
	if err != nil {
		return "", err
	}

	return sess.Get(key), nil
}

// Set tries to set a session by key and value, or returns an error
func Set(key string, value interface{}, c *fiber.Ctx) error {
	gob.Register(types.GameData{})
	gob.Register(types.WordData{})

	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	sess.Set(key, value)

	sess.SetExpiry(time.Hour * 3)

	if err2 := sess.Save(); err2 != nil {
		panic(err2.Error())
		return err2
	}

	return nil
}

// Delete tries to delete a session by key, or returns an error
func Delete(key string, c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	sess.Delete(key)

	if err2 := sess.Save(); err2 != nil {
		return err2
	}

	return nil
}
