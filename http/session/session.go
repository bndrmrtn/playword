package session

import (
	"encoding/gob"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

// This stores all of your app's sessions
var store *session.Store = session.New()

func Get(key string, c *fiber.Ctx) (interface{}, error) {
	sess, err := store.Get(c)
	if err != nil {
		return "", err
	}

	return sess.Get(key), nil
}

func Set(key string, value interface{}, c *fiber.Ctx) error {
	gob.Register(map[string]interface{}{})

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
