/*
 * DO NOT CALL ANY OF THESE FUNCTIONS DIRECTLY.
 * They should only be used by handlers.
 * TODO: Add additional wrapper around these functions for additional layer of vetting
 */

package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"github.com/binhonglee/GlobeTrotte/src/turbine/logger"
	structs "github.com/binhonglee/GlobeTrotte/src/turbine/structs"
	wings "github.com/binhonglee/GlobeTrotte/src/turbine/wings"

	"github.com/lib/pq"
)

// NewUserDB - Adding new user to the database.
func NewUserDB(newUser structs.IStructs) int {
	user, ok := newUser.(*wings.NewUser)
	if !ok {
		logger.Print("User add failed since interface passed in is not a NewUser.")
		return -1
	}

	return addNewUser(*user)
}

// GetUserDB - Retrieve user information from database with ID.
func GetUserDB(id int) structs.IStructs {
	newUser := getUserWithID(id)
	return &newUser
}

// GetUserPasswordHashDB - Retreives and return the password hash of the user account.
func GetUserPasswordHashDB(user wings.NewUser) string {
	return getUserWithEmail(user.Email).Password
}

// GetUserWithEmailDB - Retrieve user information from database with their email.
func GetUserWithEmailDB(user wings.NewUser) wings.User {
	return getUserWithID(getUserWithEmail(user.Email).ID)
}

// UpdateUserDB - Update user information back into the database.
func UpdateUserDB(updatedUser structs.IStructs) bool {
	user, ok := updatedUser.(*wings.User)
	if !ok {
		logger.Print("User update failed since interface passed in is not a user.")
		return false
	}

	return updateUser(*user)
}

// DeleteUserDB - Delete user from the database.
func DeleteUserDB(existingUser structs.IStructs) bool {
	user, ok := existingUser.(*wings.User)
	if !ok {
		logger.Print("User deletion failed since interface passed in is not a trip.")
		return false
	}

	existingUser = GetUserDB(user.GetID())

	if existingUser.GetID() == -1 {
		return false
	}

	//TODO: More testing to make sure this is the same user

	return deleteUserWithID(existingUser.GetID())
}

func addNewUser(newUser wings.NewUser) int {
	sqlStatement := `
		INSERT INTO users (name, email, password, bio, time_created)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	id := 0
	err := db.QueryRow(
		sqlStatement,
		newUser.Name,
		newUser.Email,
		newUser.Password,
		"",
		time.Now(),
	).Scan(&id)
	if err != nil {
		logger.Print(err)
		return -1
	}
	logger.Print("New record ID is: ", id)
	return id
}

func getUserWithID(id int) wings.User {
	var user wings.User
	var sqlInt64 []sql.NullInt64
	sqlStatement := `
		SELECT id, name, email, bio, time_created, trips
		FROM users WHERE id=$1;`
	switch err := db.QueryRow(sqlStatement, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Bio,
		&user.TimeCreated,
		pq.Array(&sqlInt64),
	); err {
	case sql.ErrNoRows:
		logger.Print("User not found.")
		user.ID = -1
	default:
		logger.Print(err)
	}

	user.Trips = []int{}
	for _, trip := range sqlInt64 {
		if trip.Valid {
			user.Trips = append(user.Trips, int(trip.Int64))
		}
	}

	return user
}

func getUserWithEmail(hashedPassword string) wings.NewUser {
	var user wings.NewUser
	sqlStatement := `
		SELECT id, password
		FROM users WHERE email=$1;`
	switch err := db.QueryRow(sqlStatement, hashedPassword).Scan(
		&user.ID,
		&user.Password,
	); err {
	case sql.ErrNoRows:
		logger.Print("User not found.")
		user.ID = -1
	default:
		logger.Print(err)
	}

	return user
}

func updateUser(updatedUser wings.User) bool {
	existingUser := GetUserDB(updatedUser.GetID())
	if existingUser.GetID() != updatedUser.GetID() {
		logger.Print("Existing User is not found. Aborting update.")
		logger.Print(
			"Given ID is "+strconv.Itoa(updatedUser.GetID()),
			" but found ID is "+strconv.Itoa(existingUser.GetID()),
			" instead.",
		)
		return false
	}

	sqlStatement := `
		UPDATE users
		SET name = $2,
		email = $3,
		bio = $4,
		trips = $5
		WHERE id = $1;`

	_, err := db.Exec(
		sqlStatement,
		updatedUser.ID,
		updatedUser.Name,
		updatedUser.Email,
		updatedUser.Bio,
		pq.Array(updatedUser.Trips),
	)

	if err != nil {
		logger.Print("Failed to update user.")
		logger.Print(err)
		return false
	}

	return true
}

func deleteUserWithID(id int) bool {
	sqlStatement := `
		DELETE FROM users
		WHERE id = $1;`
	if _, err := db.Exec(sqlStatement, id); err != nil {
		logger.Print(err)
		return false
	}
	logger.Print("User ID ", id, " deleted")
	return true
}
