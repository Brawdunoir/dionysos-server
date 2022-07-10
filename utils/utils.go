package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/arangodb/go-driver"
)

func GetUsersFromRoom(roomID string, db driver.Database, graph driver.Graph) (result []driver.DocumentID) {

	ctx := driver.WithQueryCount(context.Background())

	edgeCollection, _, err := graph.EdgeCollection(ctx, database.EdgeCollection)
	if err != nil {
		log.Fatalf("Failed to select edge collection: %v", err)
	}

	query := "FOR r IN rooms FILTER r._id == \"rooms/" + string(roomID) + "\" FOR v IN 1..1 INBOUND r " + edgeCollection.Name() + " RETURN v"
	fmt.Print(query)
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
	}
	defer cursor.Close()
	fmt.Printf("Query yields %d users\n", cursor.Count())

	for {
		var user models.User
		meta, err := cursor.ReadDocument(ctx, &user)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Printf("Failed to fetch query result: %v", err)
		}
		result = append(result, meta.ID)
	}

	return result
}

func IsUserConnectedToRoom(roomID string, user string, db driver.Database, graph driver.Graph) (found bool) {
	connnectedUsers := GetUsersFromRoom(roomID, db, graph)
	fmt.Println("Connected users: ", connnectedUsers)

	// Check if roomUpdate.Owner is in connnectedUsers
	found = false

	for _, v := range connnectedUsers {
		// check if the strings match
		if string(v) == user {
			found = true
			break
		}
	}
	return found
}
