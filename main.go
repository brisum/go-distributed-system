package main

import (
	"context"
	"distributes_system/lib/datastorage"
	eventsourcing "distributes_system/lib/event_sourcing"
	accountDomain "distributes_system/project/virtual_pay_network/domain/account/domain"
	accountDomainCommand "distributes_system/project/virtual_pay_network/domain/account/domain/command"
	"fmt"
	pgx "github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connection, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close(context.Background())

	store := eventsourcing.NewEventStore(connection)

	accountUuid := uuid.Must(uuid.FromString("65198e5e-f881-4d6e-ac98-502f2e3b9170"))
	accountAggregate := accountDomain.NewAccountAggregate(accountUuid)

	createAccountCommand := accountDomainCommand.NewCreateAccountCommand("Alex", "Dev")
	accountAggregate.ProcessCreateAccountCommand(*createAccountCommand)

	store.Save(&ctx, accountAggregate)

	storage := datastorage.NewEmptyDataStorage()
	storage.UnmarshalJSON("{\"balance\":{\"bonus\":0,\"cash\":0},\"firstName\":\"Alex\",\"lastName\":\"Dev\"}")

	//fmt.Printf("%+v\n", event)

	//stream.AppendEvent(event_store.NewEvent(
	//	"AccountCreated", "{balance: 0}", "", "",
	//))
	//stream.AppendEvent(event_store.NewEvent(
	//	"BalanceUpdated", "{balance: 20}", "", "",
	//))
	//stream.AppendEvent(event_store.NewEvent(
	//	"BalanceUpdated", "{balance: 30}", "", "",
	//))
	//err = store.Append(&ctx, stream)

	//stream, err = store.Load(&ctx, stream)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to generate uuid: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//for _, event := range stream.GetEvents() {
	//	fmt.Printf("%+v\n", event)
	//}
}
