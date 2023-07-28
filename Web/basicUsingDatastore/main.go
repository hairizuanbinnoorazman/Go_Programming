package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gofrs/uuid"
)

// Require the following environment variables
// DATASTORE_PROJECT_ID: test
// DATASTORE_EMULATOR_HOST: db:8081

var entityName = "visitors"

func main() {
	dsClient, err := datastore.NewClient(context.TODO(), "")
	if err != nil {
		log.Println("unable to connect to datastore")
		panic("failed")
	}

	log.Println("instantiation completed")

	createVisitor(dsClient, entityName, "Tom")
	createVisitor(dsClient, entityName, "Jerry")
	v := listVisitors(dsClient, entityName)
	prettyOutput, _ := json.MarshalIndent(v, "", "    ")
	fmt.Println(string(prettyOutput))

	fmt.Printf("\n\nBegin delete of first item and adding Harry:\n")
	deleteVisitor(dsClient, entityName, v[0].ID)
	createVisitor(dsClient, entityName, "Harry")
	newV := listVisitors(dsClient, entityName)
	prettyOutput2, _ := json.MarshalIndent(newV, "", "    ")
	fmt.Println(string(prettyOutput2))
}

type visitor struct {
	ID              string
	Name            string
	DatetimeEntered time.Time
}

func createVisitor(client *datastore.Client, entityName, name string) {
	log.Println("start createVisitor")
	defer log.Println("end createVisitor")

	uid, _ := uuid.NewV4()
	newVisitor := visitor{
		ID:              uid.String(),
		Name:            name,
		DatetimeEntered: time.Now(),
	}
	key := datastore.NameKey(entityName, newVisitor.ID, nil)
	_, err := client.Put(context.TODO(), key, &newVisitor)
	if err != nil {
		log.Printf("unable to save data. err: %v", err)
	}
}

func deleteVisitor(client *datastore.Client, entityName, id string) {
	log.Println("start deleteVisitor")
	defer log.Println("end deleteVisitor")

	k := datastore.NameKey(entityName, id, nil)
	err := client.Delete(context.TODO(), k)
	if err != nil {
		log.Printf("unable to delete visitor. err: %v\n", err)
	}
}

func listVisitors(client *datastore.Client, entityName string) []visitor {
	log.Println("start listVisitors")
	defer log.Println("end listVisitors")

	visitors := []visitor{}
	q := datastore.NewQuery(entityName)
	keys, err := client.GetAll(context.TODO(), q, &visitors)
	if err != nil {
		log.Printf("unable to list visitors. err: %v\n", err)
		// Will be empty if error-ed out
		return visitors
	}
	for k, _ := range visitors {
		visitors[k].ID = keys[k].Name
	}
	return visitors
}
