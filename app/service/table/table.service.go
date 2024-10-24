package table

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/rand"
)

type Table struct {
	ID   string
	Name string
	Date time.Time
}

func TableService() {
	app := fiber.New()

	// Define routes for the Table service
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Service Table!")
	})

	app.Get("/ping-power", func(c *fiber.Ctx) error {
		power := c.QueryInt("power")
		name := c.Query("name")
		if power <= 0 {
			return c.Status(fiber.StatusBadRequest).SendString("Power must be a positive integer")
		}
		if name == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Name is required")
		}
		lowPower := BallPowerTo(uint(power), name)
		return c.SendString(fmt.Sprintf("%d", lowPower))
	})

	// Listen on port 8889
	app.Listen(":8889")
}

func BallPowerTo(power uint, name string) uint {
	// Calculate a random low power between 70% and 90% of the original power
	lowPowerPercentage := 70 + rand.Intn(21) // Random number between 70 and 90
	lowPower := uint(float64(power) * float64(lowPowerPercentage) / 100.0)
	return lowPower
}
