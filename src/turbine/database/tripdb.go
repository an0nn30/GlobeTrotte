/*
 * DO NOT CALL ANY OF THESE FUNCTIONS DIRECTLY.
 * They should only be used by handlers or inside the database class.
 *
 * TODO: Add additional wrapper around these functions for additional layer of vetting
 */

package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"github.com/binhonglee/GlobeTrotte/src/turbine/logger"
	structs "github.com/binhonglee/GlobeTrotte/src/turbine/structs"
	wings "github.com/binhonglee/GlobeTrotte/src/turbine/wings"

	"github.com/lib/pq"
)

// AddTripDB - Adding new trip into the database.
func AddTripDB(newTrip structs.IStructs) int {
	trip, ok := newTrip.(*wings.Trip)
	if !ok {
		logger.Print("Trip add failed since interface passed in is not a trip.")
		return -1
	}

	newTripID := addTrip(*trip)
	user := getUserWithID(trip.UserID)
	if user.ID == -1 {
		logger.Print("User adding the new trip is not found.")
		return failAddingTripToUser(newTripID)
	}

	user.Trips = append(user.Trips, newTripID)
	if ok = updateUser(user); !ok {
		logger.Print("Fail to add trip id to new user.")
		return failAddingTripToUser(newTripID)
	}

	return newTripID
}

// GetTripDB - Retrieve trip information from database with ID.
func GetTripDB(id int) structs.IStructs {
	var trip wings.Trip = fetchTrip(id)
	return &trip
}

// UpdateTripDB - Update trip information back into the database.
func UpdateTripDB(updatedTrip structs.IStructs) bool {
	trip, ok := updatedTrip.(*wings.Trip)
	if !ok {
		logger.Print("Trip update failed since interface passed in is not a trip.")
		return false
	}

	for index, day := range trip.Days {
		if updateDay(&day, true) {
			trip.Days[index] = day
		} else {
			logger.Print("Trip update failed since one of the day update failed.")
			return false
		}
	}

	return updateTrip(*trip)
}

// DeleteTripDB - Delete trip from the database.
func DeleteTripDB(existingTrip structs.IStructs) bool {
	trip, ok := existingTrip.(*wings.Trip)
	if !ok {
		logger.Print("Trip deletion failed since interface passed in is not a trip.")
		return false
	}

	existingTrip = GetTripDB(trip.GetID())

	if existingTrip.GetID() == -1 {
		return false
	}

	//TODO: More testing to make sure this is the same trip
	return deleteTripWithID(existingTrip.GetID())
}

func addTrip(newTrip wings.Trip) int {
	sqlStatement := `
		INSERT INTO trips (userid, name, cities, description, days, time_created, last_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	id := 0
	err := db.QueryRow(
		sqlStatement,
		newTrip.UserID,
		newTrip.Name,
		pq.Array(cityEnumArrayToIDs(newTrip.Cities)),
		newTrip.Description,
		pq.Array([]int{}),
		newTrip.TimeCreated,
		newTrip.LastUpdated,
	).Scan(&id)

	if len(newTrip.Days) > 0 {
		for index, day := range newTrip.Days {
			dayID := addDay(day)
			newTrip.Days[index].ID = dayID
		}

		newTrip.ID = id
		updateTrip(newTrip)
	}

	if err != nil {
		logger.Print(err)
		return -1
	}
	logger.Print("New trip ID is: ", id)
	return id
}

func addDay(newDay wings.Day) int {
	sqlStatement := `
		INSERT INTO days (id, trip_id, day_of, places)
		VALUES($1, $2, $3, $4)
		RETURNING id`
	id := 0
	err := db.QueryRow(
		sqlStatement,
		newDay.ID,
		newDay.TripID,
		newDay.DayOf,
		pq.Array(placesToIDArray(newDay.Places)),
	)

	if err != nil {
		logger.Print(err)
		return -1
	}
	logger.Print("New day ID is: ", id)
	return id
}

func fetchTrip(id int) wings.Trip {
	var trip wings.Trip
	var days []int
	var cities []int64
	sqlStatement := `
		SELECT id, userid, name, cities, description, days, time_created, last_updated
		FROM trips WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(
		&trip.ID,
		&trip.UserID,
		&trip.Name,
		pq.Array(&cities),
		&trip.Description,
		pq.Array(&days),
		&trip.TimeCreated,
		&trip.LastUpdated,
	); err {
	case sql.ErrNoRows:
		logger.Print("Trip not found.")
		trip.ID = -1
	default:
		logger.Print(err)
	}

	trip.Cities = cityIDsToEnumArray(cities)
	trip.Days = fetchDays(days)
	return trip
}

func fetchDays(ids []int) wings.Days {
	var days wings.Days = make([]wings.Day, len(ids))
	for index, id := range ids {
		days[index] = fetchDay(id)
	}
	return days
}

func fetchDay(id int) wings.Day {
	var day wings.Day
	var places []int
	sqlStatement := `
		SELECT id, trip_id, day_of, places
		FROM days WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(
		&day.ID,
		&day.TripID,
		&day.DayOf,
		pq.Array(&places),
	); err {
	case sql.ErrNoRows:
		logger.Print("Day not found.")
		day.ID = -1
	default:
		logger.Print(err)
	}
	day.Places = fetchPlaces(places)
	return day
}

func fetchPlaces(ids []int) wings.Places {
	var places wings.Places = make([]wings.Place, len(ids))

	for index, id := range ids {
		var place wings.Place
		sqlStatement := `
			SELECT id, label, url, description
			FROM days WHERE id=$1;`
		row := db.QueryRow(sqlStatement, id)
		switch err := row.Scan(
			&place.ID,
			&place.Label,
			&place.URL,
			&place.Description,
		); err {
		case sql.ErrNoRows:
			logger.Print("Place not found.")
			place.ID = -1
		default:
			logger.Print(err)
		}
		places[index] = place
	}

	return places
}

func updateTrip(updatedTrip wings.Trip) bool {
	existingTrip := GetTripDB(updatedTrip.GetID())
	if existingTrip.GetID() != updatedTrip.GetID() {
		logger.Print("Existing Trip is not found. Aborting update.")
		logger.Print("Given ID is " + strconv.Itoa(updatedTrip.GetID()) +
			" but found ID is " + strconv.Itoa(existingTrip.GetID()) +
			" instead.")
		return false
	}

	sqlStatement := `
		UPDATE trips
		SET name = $2,
		description = $3,
		cities = $4,
		days = $5,
		last_updated = $6
		WHERE id = $1;`

	_, err := db.Exec(
		sqlStatement,
		updatedTrip.ID,
		updatedTrip.Name,
		updatedTrip.Description,
		pq.Array(cityEnumArrayToIDs(updatedTrip.Cities)),
		pq.Array(daysToIDArray(updatedTrip.Days)),
		updatedTrip.LastUpdated,
	)

	if err != nil {
		logger.Print("Failed to update trip.")
		logger.Print(err)
		return false
	}

	return true
}

func updateDay(updatedDay *wings.Day, createOnNonExist bool) bool {
	existingDay := fetchDay(updatedDay.ID)

	if existingDay.ID != updatedDay.ID {
		logger.Print("Existing Day not found.")

		if createOnNonExist {
			updatedDay.ID = addDay(*updatedDay)
			return updatedDay.ID == -1
		} else {
			return false
		}
	}

	sqlStatement := `
		UPDATE days
		SET trip_id = $2,
		day_of = $3,
		places = $4
		WHERE id = $1;`

	_, err := db.Exec(
		sqlStatement,
		updatedDay.ID,
		updatedDay.TripID,
		updatedDay.DayOf,
		pq.Array(placesToIDArray(updatedDay.Places)),
	)

	if err != nil {
		logger.Print("Failed to update day.")
		logger.Print(err)
		return false
	}

	return true
}

func deleteTripWithID(id int) bool {
	sqlStatement := `
		DELETE FROM trips
		WHERE id = $1;`
	if _, err := db.Exec(sqlStatement, id); err != nil {
		logger.Print(err)
		return false
	}
	logger.Print("Trip ID ", id, " deleted")
	return true
}

func failAddingTripToUser(id int) int {
	deleteTripWithID(id)
	return -1
}
