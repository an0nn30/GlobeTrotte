package router

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	structs "github.com/binhonglee/GlobeTrotte/src/turbine/structs"
	wings "github.com/binhonglee/GlobeTrotte/src/turbine/wings"
)

var addedTrip wings.Trip
var addedUser wings.User

func TestAddInvalidEmailUser(t *testing.T) {
	var newUser = wings.NewUser{
		Email:    "wronglyformattedstring",
		Password: "shouldReplaceThisWithRand",
	}

	var returned *wings.User
	addTest("/user", t, &newUser, &returned, http.StatusNotAcceptable)
}

func TestAddUser(t *testing.T) {
	var newUser = wings.NewUser{
		Email:    "routertest@test.com",
		Password: "shouldReplaceThisWithRand",
	}

	var returned *wings.User
	addTest("/user", t, &newUser, &returned, http.StatusCreated)
	addedUser = *returned

	if returned.Email != newUser.Email {
		t.Errorf(
			"Sent Email is %v but returned Email is %v.",
			newUser.Email,
			returned.Email,
		)
	}
}

func TestGetUser(t *testing.T) {
	var returned *wings.User
	getTest(
		"/user/"+strconv.Itoa(addedUser.GetID()),
		t,
		&returned,
		http.StatusOK,
	)

	if _, diff := structs.CompareUser(addedUser, *returned); len(diff) > 0 {
		t.Errorf(
			"The followings fields does not match:\n%v",
			strings.Join(diff, ", "),
		)
	}
}

func TestAddTrip(t *testing.T) {
	var newTrip = wings.Trip{
		UserID: addedUser.ID,
		Name:   "TestUser",
		Cities: []wings.City{wings.SanFranciscoCAUS},
		Days: []wings.Day{
			{
				DayOf:  0,
				Places: []wings.Place{},
			},
		},
		Description: "Description",
	}

	var returned *wings.Trip
	addTest("/trip", t, &newTrip, &returned, http.StatusCreated)
	addedTrip = *returned

	if _, diff := removeIDFromArray(structs.CompareTrips(newTrip, *returned)); len(diff) > 1 {
		t.Errorf(
			"The followings fields does not match:\n%v",
			strings.Join(diff, ", "),
		)
	}
}

func TestGetTrip(t *testing.T) {
	var returned *wings.Trip
	getTest("/trip/"+strconv.Itoa(addedTrip.GetID()), t, &returned, http.StatusOK)

	if _, diff := structs.CompareTrips(addedTrip, *returned); len(diff) > 0 {
		t.Errorf(
			"The followings fields does not match:\n%v",
			strings.Join(diff, ", "),
		)
	}
}

func TestGetNonExistentTrip(t *testing.T) {
	var returned *wings.Trip
	getTest("/trip/"+strconv.Itoa(-1), t, &returned, http.StatusNotFound)
}

func TestUpdateTrip(t *testing.T) {
	var returned *wings.Trip
	addedTrip.Description = "Updated description"

	updateTest(
		"/trip/"+strconv.Itoa(addedTrip.ID), t,
		&addedTrip, http.StatusAccepted,
	)
	getTest(
		"/trip/"+strconv.Itoa(addedTrip.ID), t,
		&returned, http.StatusOK,
	)

	if _, diff := structs.CompareTrips(addedTrip, *returned); len(diff) > 0 {
		t.Errorf(
			"The followings fields does not match:\n%v",
			strings.Join(diff, ", "),
		)
	}
}

func TestUpdateUser(t *testing.T) {
	var returned *wings.User
	addedUser.Name = "My NewName"

	updateTest(
		"/user/"+strconv.Itoa(addedUser.ID), t,
		&addedUser, http.StatusAccepted,
	)
	getTest(
		"/user/"+strconv.Itoa(addedUser.ID), t,
		&returned, http.StatusOK,
	)

	if _, diff := structs.CompareUser(addedUser, *returned); len(diff) > 0 {
		t.Errorf(
			"The followings fields does not match:\n%v",
			strings.Join(diff, ", "),
		)
	}
}

func TestDeleteTrip(t *testing.T) {
	deleteTest("/trip/"+strconv.Itoa(addedTrip.GetID()), t, http.StatusOK)
}

func TestDeleteNonExistentTrip(t *testing.T) {
	deleteTest("/trip/"+strconv.Itoa(-1), t, http.StatusNotFound)
}

func TestDeleteUser(t *testing.T) {
	deleteTest("/user/"+strconv.Itoa(addedUser.GetID()), t, http.StatusOK)
}

func TestDeleteNonExistentUser(t *testing.T) {
	deleteTest("/user/"+strconv.Itoa(-1), t, http.StatusNotFound)
}

func removeIDFromArray(status bool, arr []string) (bool, []string) {
	ret := []string{}
	for _, s := range arr {
		if s != "ID" {
			ret = append(ret, s)
		}
	}
	return status, ret
}
