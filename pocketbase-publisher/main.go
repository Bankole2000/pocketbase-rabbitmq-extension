package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PBEvent struct {
	Action     string         `json:"action"`
	Collection string         `json:"collection"`
	Record     *models.Record `json:"record"`
}

type ServiceEvent struct {
	EventType string  `json:"type"`
	Data      PBEvent `json:"data"`
	Origin    string  `json:"origin"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func sendEventMessage(eventType string, action string, record *models.Record) {
	message := PBEvent{Action: strings.ToUpper(action), Collection: strings.ToUpper(record.Collection().Name), Record: record}
	eventMessage := ServiceEvent{EventType: strings.ToUpper(eventType), Data: message, Origin: "POCKETBASE"}
	jsonBytes, error := json.Marshal(eventMessage)
	failOnError(error, "Failed to convert message to JsonBytes")
	rmqUrl := goDotEnvVariable("RABBITMQ_URL")
	exchange := goDotEnvVariable("RABBITMQ_EXCHANGE")
	conn, err := amqp.Dial(rmqUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"gistable", // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonBytes,
			// Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %v => %v", message.Collection, message.Action)
}

func main() {

	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})

	app.OnBeforeBootstrap().Add(func(e *core.BootstrapEvent) error {
		return nil
	})

	app.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		log.Printf("MODEL %v => CREATED", e.Model.TableName())
		return nil
	})

	app.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
		log.Printf("MODEL %v => UPDATED", e.Model.TableName())
		return nil
	})

	app.OnModelAfterDelete().Add(func(e *core.ModelEvent) error {
		log.Printf("MODEL %v => DELETED", e.Model.TableName())
		return nil
	})

	app.OnCollectionAfterCreateRequest().Add(func(e *core.CollectionCreateEvent) error {
		log.Printf("COLLECTION %v => CREATED", e.Collection.Name)
		return nil
	})

	app.OnCollectionAfterUpdateRequest().Add(func(e *core.CollectionUpdateEvent) error {
		log.Printf("COLLECTION %v => UPDATED", e.Collection.Name)
		return nil
	})

	app.OnCollectionAfterDeleteRequest().Add(func(e *core.CollectionDeleteEvent) error {
		log.Printf("COLLECTION %v => DELETED", e.Collection.Name)
		return nil
	})

	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		sendEventMessage("records", "created", e.Record)
		return nil
	})
	// fires for every collection
	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		sendEventMessage("records", "updated", e.Record)
		return nil
	})

	// fires for every collection
	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		sendEventMessage("records", "deleted", e.Record)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
