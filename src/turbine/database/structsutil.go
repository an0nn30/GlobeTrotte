package database

import (
	"strconv"

	wings "github.com/binhonglee/GlobeTrotte/src/turbine/wings"
)

func daysToIDArray(days wings.Days) []int {
	var idArray []int = make([]int, len(days))
	for index := range days {
		idArray[index] = days[index].GetID()
	}
	return idArray
}

func placesToIDArray(places wings.Places) []int {
	var idArray []int = make([]int, len(places))
	for index := range places {
		idArray[index] = places[index].GetID()
	}
	return idArray
}

func cityEnumArrayToIDs(cities []wings.City) []int {
	var idArray []int = make([]int, len(cities))
	for index := range cities {
		idArray[index] = int(cities[index])
	}
	return idArray
}

func cityIDsToEnumArray(cityIDs []int64) []wings.City {
	var cityArray []wings.City = make([]wings.City, len(cityIDs))
	for index := range cityIDs {
		cityArray[index] = wings.City(cityIDs[index])
	}
	return cityArray
}

func tripsToString(tripIDs []int) string {
	var toReturn = ""

	for _, trip := range tripIDs {
		toReturn += strconv.Itoa(trip)
	}

	return toReturn
}
