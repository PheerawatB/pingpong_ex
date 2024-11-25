package player

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
)

var countMatch uint = 0
var mongoClient *mongo.Client
var logMatch string

type MatchLog struct {
	MatchID  uint      `json:"match_id" bson:"match_id"`
	MatchLog string    `json:"match_log" bson:"match_log"`
	Time     time.Time `json:"time" bson:"time"`
}

func PlayerService() {
	var err error
	// MongoDB connection
	mongoURI := "mongodb://localhost:27017"
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer mongoClient.Disconnect(context.TODO())
	countMatch, _ = getLastMatchID() // Get last match id

	app := fiber.New()

	// Define routes for the Player service
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Player Service")
	})

	app.Post("/new-match", func(c *fiber.Ctx) error {
		_ = newMatch()
		logMatchResultToMongoDB(countMatch, logMatch)
		return c.SendString(logMatch)
	})

	app.Get("/match", func(c *fiber.Ctx) error {
		listMatch, _ := getAllMatches()
		return c.JSON(listMatch)
	})

	// Listen on port 8888
	app.Listen(":8888")
}

func newMatch() string {
	logMatch = ""
	countMatch++
	logToCSV(fmt.Sprintf("------------------ New Match %d ------------------", countMatch))
	logToCSV("Player A 🧍🏻 & Player B 🧍🏻 on the court")
	time.Sleep(1 * time.Second)

	ach := make(chan uint)
	bch := make(chan uint)
	done := make(chan string) // Channel to signal the end of the match

	// Start the game by sending an initial power to Player A
	go func() {
		initialPower := uint(rand.Intn(51)) + 50
		ach <- initialPower
	}()

	closeCh := false
	// Main game loop
	go func() {
		for {
			if closeCh {
				break // Game over, exit the loop
			}

			select {
			case power, ok := <-ach: // Player A's turn
				if !ok {
					closeCh = true
					continue
				}
				go func() {
					time.Sleep(1 * time.Second)
					if !Player(power, bch, "Player A", "Player B") { // Player B's turn
						done <- "Player A" // Player B wins
						return             // Exit the goroutine
					}
				}()

			case power, ok := <-bch: // Player B's turn
				if !ok {
					closeCh = true
					continue
				}
				go func() {
					time.Sleep(1 * time.Second)
					if !Player(power, ach, "Player B", "Player A") { // Player A's turn
						done <- "Player B" // Player A wins
						return             // Exit the goroutine
					}
				}()
			}
		}
	}()

	// Block until we receive a winner from the "done" channel
	winner := <-done

	// Close channels to stop further game operations
	close(ach)
	close(bch)
	logToCSV(fmt.Sprintf("[Alert] %s wins!", winner))
	logToCSV("------------------- Game Over -------------------")

	return winner
}

func Player(power uint, wakeCh chan uint, name string, opponent string) bool {
	url := fmt.Sprintf("http://localhost:8889/ping-power?power=%d&name=%s", power, name)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer res.Body.Close()
	time.Sleep(1 * time.Second)

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return false
		}
		bodyInt, err := strconv.Atoi(string(body))
		if err != nil {
			fmt.Println("Error converting response to integer:", err)
			return false
		}

		newPower := uint(rand.Intn(51)) + 50
		time.Sleep(1 * time.Second)

		if newPower > uint(bodyInt) {
			wakeCh <- newPower
			logToCSV(fmt.Sprintf("[%s] 🏓 💥 {%d} ========== [%d] ==========> 🏓 [%s] ", name, power, bodyInt, opponent))
			return true
		} else {
			logToCSV(fmt.Sprintf("[%s] 🏓 💥 {%d} ========== [%d] ==========> 💀 [%d] [%s] ", name, power, bodyInt, newPower, opponent))
			return false
		}
	}
	return false
}

func logToCSV(message string) {

	fmt.Println(message)
	logDir := "./logs"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating log directory:", err)
		return
	}

	// Create the log file name with the match name and current date
	date := time.Now().Format("20060102_15")
	logFileName := fmt.Sprintf("%s_%s.csv", "match", date)
	logFilePath := filepath.Join(logDir, logFileName)

	// Open or create the CSV file for appending
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening/creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the log message and timestamp to the CSV file
	record := []string{
		time.Now().Format(time.RFC3339), // Timestamp
		message,                         // Log message
	}
	logMatch += time.Now().Format(time.RFC3339) + ":" + message + "\n"

	// Write the record to the CSV
	err = writer.Write(record)
	if err != nil {
		fmt.Println("Error writing to CSV file:", err)
	}
}

func logMatchResultToMongoDB(matchID uint, logMatch string) {
	collection := mongoClient.Database("match_results").Collection("results")

	matchResult := MatchLog{
		MatchID:  matchID,
		MatchLog: logMatch,
		Time:     time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), matchResult)
	if err != nil {
		fmt.Println("Error inserting match result to MongoDB:", err)
	}
}

func getLastMatchID() (uint, error) {
	collection := mongoClient.Database("match_results").Collection("results")

	var lastMatch MatchLog
	err := collection.FindOne(
		context.TODO(),
		bson.D{},
		options.FindOne().SetSort(bson.D{{Key: "match_id", Value: -1}}),
	).Decode(&lastMatch)
	if err != nil {
		return 0, err
	}

	return lastMatch.MatchID, nil
}

func getAllMatches() ([]MatchLog, error) {
	collection := mongoClient.Database("match_results").Collection("results")

	// Create a slice to hold the match logs
	var matches []MatchLog

	// Find all documents in the collection
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each match log into the slice
	for cursor.Next(context.TODO()) {
		var match MatchLog
		if err := cursor.Decode(&match); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	// Check for errors that may have occurred during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}
