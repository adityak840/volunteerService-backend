package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Volunteer struct {
	VolunteerID   string `json:"volunteerId,omitempty" bson:"volunteerId,omitempty"`
	VolunteerName string `json:"volunteerName,omitempty" bson:"volunteerName,omitempty"`
}

type Todo struct {
	ID               string      `json:"id,omitempty" bson:"_id,omitempty"`
	Task             string      `json:"task,omitempty" bson:"task,omitempty"`
	Description      string      `json:"description,omitempty" bson:"description,omitempty"`
	OrganisationName string      `json:"orgName,omitempty" bson:"orgName,omitempty"`
	VolunteerType    string      `json:"volType,omitempty" bson:"volType,omitempty"`
	OrganisationType string      `json:"orgType,omitempty" bson:"orgType,omitempty"`
	Completed        bool        `json:"completed" bson:"completed"`
	Time             time.Time   `json:"time,omitempty" bson:"time,omitempty"`
	Volunteer        []Volunteer `json:"volunteer,omitempty" bson:"volunteer,omitempty"` // Nested Volunteer struct
}

var client *mongo.Client

// New is used to initialize the mongo client for the Todo struct
func New(mongo *mongo.Client) Todo {
	client = mongo
	return Todo{}
}

// returnCollectionPointer returns a pointer to the 'todos' collection
func returnCollectionPointer(collection string) *mongo.Collection {
	return client.Database("volunteerService-backend-db").Collection(collection)
}

// GetAllTodos returns all the todos from the db
func (t *Todo) GetAllTodos() ([]Todo, error) {
	collection := returnCollectionPointer("todos")
	var todos []Todo

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodoById returns a single todo based on its ID
func (t *Todo) GetTodoById(id string) (Todo, error) {
	collection := returnCollectionPointer("todos")
	var todo Todo

	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Todo{}, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": mongoID}).Decode(&todo)
	if err != nil {
		log.Println(err)
		return Todo{}, err
	}

	return todo, nil
}

// InsertTodo creates a new todo in the collection
func (t *Todo) InsertTodo(entry Todo) error {
	collection := returnCollectionPointer("todos")

	// If the Time is not set in the request, set it to the current time
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	// Ensure the Volunteer field is not carrying over from previous operations
	entry.Volunteer = nil

	// Insert the entire 'entry' object as it contains all fields
	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Println("Error inserting todo:", err)
		return err
	}

	return nil
}

func (t *Todo) UpdateTodo(id string, entry Todo) (*mongo.UpdateResult, error) {
	collection := returnCollectionPointer("todos")
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	log.Println(entry)

	update := bson.D{
		{"$set", bson.D{
			{"task", entry.Task},
			{"completed", entry.Completed},
		}},
		{"$push", bson.D{
			{"volunteer", bson.M{"$each": entry.Volunteer}}, // Append the new volunteers to the array
		}},
	}

	res, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": mongoID},
		update,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}

// DeleteTodo deletes a todo by its ID
func (t *Todo) DeleteTodo(id string) error {
	collection := returnCollectionPointer("todos")
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = collection.DeleteOne(
		context.Background(),
		bson.M{"_id": mongoID},
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetTodosByOrg retrieves todos filtered by OrganisationName
func (t *Todo) GetTodosByOrg(orgName string) ([]Todo, error) {
	collection := returnCollectionPointer("todos")

	// Build filter to search by OrganisationName
	filter := bson.M{"orgName": orgName}

	var todos []Todo

	// Query the database
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Error finding todos by orgName:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the results and append to todos slice
	for cursor.Next(context.TODO()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Println("Error decoding todo:", err)
			continue
		}
		todos = append(todos, todo)
	}

	// Check if there was any error while iterating the cursor
	if err := cursor.Err(); err != nil {
		log.Println("Error with cursor:", err)
		return nil, err
	}

	return todos, nil
}

// GetTodosByVolType retrieves todos filtered by VolunteerType
func (t *Todo) GetTodosByVolType(volType string) ([]Todo, error) {
	collection := returnCollectionPointer("todos")

	// Build filter to search by VolunteerType
	filter := bson.M{"volType": volType}

	var todos []Todo

	// Query the database
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Error finding todos by volType:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the results and append to todos slice
	for cursor.Next(context.TODO()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Println("Error decoding todo:", err)
			continue
		}
		todos = append(todos, todo)
	}

	// Check if there was any error while iterating the cursor
	if err := cursor.Err(); err != nil {
		log.Println("Error with cursor:", err)
		return nil, err
	}

	return todos, nil
}
