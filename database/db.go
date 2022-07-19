package database

import (
	"context"
	"fmt"
	"log"
	"os"

	c "github.com/Brawdunoir/dionysos-server/constants"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/arangodb/go-driver"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetDatabase returns a database instance
func GetDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(createDSN()), createConfig())
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// nolint:errcheck
	db.AutoMigrate(&models.User{}, &models.Room{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	return db
}

// createDSN creates a DSN string from the environment variables to connect to the database.
func createDSN() string {
	username, found := os.LookupEnv("POSTGRES_USER")
	if !found {
		log.Fatal("POSTGRES_USER environment variable not found")
	}
	password, found := os.LookupEnv("POSTGRES_PASSWORD")
	if !found {
		log.Fatal("POSTGRES_PASSWORD environment variable not found")
	}
	host, found := os.LookupEnv("POSTGRES_HOST")
	if !found {
		log.Fatal("POSTGRES_HOST environment variable not found")
	}
	port, found := os.LookupEnv("POSTGRES_PORT")
	if !found {
		log.Fatal("POSTGRES_PORT environment variable not found")
	}
	dbname, found := os.LookupEnv("POSTGRES_DB")
	if !found {
		log.Fatal("POSTGRES_DB environment variable not found")
	}

	return "host=" + host + " port=" + port + " user=" + username + " password=" + password + " dbname=" + dbname
}

// createConfig creates a Gorm config depending on the environment variables.
func createConfig() *gorm.Config {
	env, found := os.LookupEnv("ENVIRONMENT")
	if !found {
		log.Println("ENVIRONMENT variable not found")
		log.Println("Possible values are : " + c.ENVIRONMENT_TESTING + ", " + c.ENVIRONMENT_DEVELOPMENT + ", " + c.ENVIRONMENT_PRODUCTION)
		return &gorm.Config{}
	}

	switch env {
	case c.ENVIRONMENT_TESTING:
		return &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	case c.ENVIRONMENT_DEVELOPMENT:
		return &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	case c.ENVIRONMENT_PRODUCTION:
		return &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}
	default:
		log.Println("ENVIRONMENT variable not valid, using default config")
		log.Println("Possible values are : " + c.ENVIRONMENT_TESTING + ", " + c.ENVIRONMENT_DEVELOPMENT + ", " + c.ENVIRONMENT_PRODUCTION)
		return &gorm.Config{}
	}
}

// GetGraph returns a graph instance
func GetGraph(db driver.Database, graphName string) (graph driver.Graph) {

	graphExists, err := db.GraphExists(context.TODO(), graphName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if graphExists {
		fmt.Printf("%s graph exists already\n", graphName)
		graph, err = db.Graph(context.TODO(), graphName)
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		graph = SetupGraph(db, graphName, cols)
	}

	return graph
}

// SetupGraph creates the edgeDefinition and the corresponding graph
func SetupGraph(db driver.Database, graphName string, cols []string) driver.Graph {

	var edgeDefinition driver.EdgeDefinition

	edgeDefinition.Collection = EdgeCollection

	// define a set of collections where an edge is going out...
	edgeDefinition.From = []string{UsersCollection}

	// repeat this for the collections where an edge is going into
	edgeDefinition.To = cols

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	graph, err := db.CreateGraphV2(context.TODO(), graphName, &options)
	if err != nil {
		fmt.Printf("Failed to create graph: %v", err)
	} else {
		fmt.Printf("Created graph '%s' in database '%s'\n", graphName, db.Name())
	}

	return graph
}
