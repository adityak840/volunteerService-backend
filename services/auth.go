package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// User struct for storing user data
type User struct {
	ID               string `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName        string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email            string `json:"email,omitempty" bson:"email,omitempty"`
	Password         string `json:"password,omitempty" bson:"password,omitempty"`
	ContactNumber    string `json:"contactNo,omitempty" bson:"contactNo,omitempty"`
	UserType         string `json:"userType,omitempty" bson:"userType,omitempty"`
	VolunteerType    string `json:"volType,omitempty" bson:"volType,omitempty"`
	OrganisationName string `json:"orgName" bson:"orgName"`
}

// JWT secret key used for signing JWT tokens (should be stored in environment variable in production)
var jwtSecret = []byte("your-secret-key") // Replace with a secure key in production

// NewAuthService initializes the Mongo client for authentication
func NewAuthService(mongo *mongo.Client) {
	client = mongo
}

// Signup handles user registration by hashing the password and saving user data
func Signup(user User) (string, error) {
	collection := returnCollectionPointer("users")

	// Check if the email already exists
	emailResult := collection.FindOne(context.TODO(), bson.M{"email": user.Email})
	if emailResult.Err() == nil {
		log.Println("Email already exists")
		return "", fmt.Errorf("email already exists")
	} else if emailResult.Err() != mongo.ErrNoDocuments {
		log.Println("Error checking for email existence:", emailResult.Err())
		return "", emailResult.Err()
	}

	// Check if the contact number already exists
	contactResult := collection.FindOne(context.TODO(), bson.M{"contactNo": user.ContactNumber})
	if contactResult.Err() == nil {
		log.Println("Contact number already exists")
		return "", fmt.Errorf("contact number already exists")
	} else if contactResult.Err() != mongo.ErrNoDocuments {
		log.Println("Error checking for contact number existence:", contactResult.Err())
		return "", contactResult.Err()
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}
	user.Password = string(hashedPassword)

	// Insert the user into the MongoDB "users" collection
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Error inserting user:", err)
		return "", err
	}

	// Return userType as JSON response, no need to send token during signup
	return user.UserType, nil
}

// Login handles user login by verifying the password and returning a JWT token
func Login(w http.ResponseWriter, email, password string) (string, error) {
	// Find the user in the MongoDB "users" collection
	collection := returnCollectionPointer("users")
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Println("Error finding user:", err)
		return "", fmt.Errorf("invalid credentials")
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Invalid credentials")
		return "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := generateJWT(user.Email)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return "", err
	}

	// Set the token as an HTTPOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true, // Make the cookie accessible only through HTTP requests (can't be accessed via JavaScript)
		Secure:   true, // Should be true if you're using HTTPS
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24), // Token expires in 24 hours
	})

	// Return the userType as JSON response
	response := struct {
		UserType      string `json:"userType"`
		ID            string `json:"ID"`
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		ContactNumber string `json:"contactNo"`
		Email         string `json:"email"`
		OrgName       string `json:"orgName"`
	}{
		UserType:      user.UserType,
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		ContactNumber: user.ContactNumber,
		Email:         user.Email,
		OrgName:       user.OrganisationName,
	}

	// Send the userType as part of the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding JSON response:", err)
		return "", err
	}

	return token, nil
}

// generateJWT generates a JWT token for the user
func generateJWT(email string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create the token with claims and the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and return it
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return signedToken, nil
}

// GetUsersByID retrieves user details by multiple IDs
func GetUsersByID(ids []string) ([]User, error) {
	collection := returnCollectionPointer("users")
	var users []User
	// Convert string slice to BSON array of ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Invalid ID format:", id)
			return nil, fmt.Errorf("invalid ID format: %s", id)
		}
		objectIDs = append(objectIDs, objectID)
	}
	// Query to find all users with the given IDs
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := collection.Find(context.TODO(), filter, options.Find())
	if err != nil {
		log.Println("Error finding users:", err)
		return nil, err
	}
	// Iterate through the cursor and decode each user
	for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			log.Println("Error decoding user:", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}
	// Close the cursor once finished
	cursor.Close(context.TODO())
	return users, nil
}
