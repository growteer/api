package entities

import "time"

type Profile struct {
	DID          string
	FirstName    string
	LastName     string
	DateOfBirth  time.Time
	PrimaryEmail string
	Location     Location
	Website      string
	PersonalGoal string
	About        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Location struct {
	Country    string
	PostalCode string
	City       string
}
