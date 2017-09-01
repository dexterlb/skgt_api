package common

// Stop contains information about a single stop
type Stop struct {
	ID                int
	Name              string
	InternationalName string
	CommunityName     string
	Description       string
	Latitude          float64
	Longitude         float64
}
