package store

import "golang.org/x/net/context"

type GeoLocation struct {
	Latitude  string
	Longitude string
}

type Connector struct {
	Format      string
	Id          string
	MaxAmperage int32
	MaxVoltage  int32
	PowerType   string
	Standard    string
	LastUpdated string
}

type Evse struct {
	Connectors  []Connector
	EvseId      *string
	Status      string
	Uid         string
	LastUpdated string
}

type Location struct {
	Address     string
	City        string
	Coordinates GeoLocation
	Country     string
	Evses       *[]Evse
	Id          string
	LastUpdated string
	Name        string
	ParkingType string
	PostalCode  string
}

type LocationStore interface {
	SetLocation(ctx context.Context, location *Location) error
	LookupLocation(ctx context.Context, locationId string) (*Location, error)
	ListLocations(context context.Context, offset int, limit int) ([]*Location, error)
}
