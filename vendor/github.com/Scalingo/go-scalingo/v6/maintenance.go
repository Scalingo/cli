package scalingo

type MaintenanceWindow struct {
	WeekdayUTC      int `json:"weekday_utc"`
	StartingHourUTC int `json:"starting_hour_utc"`
	DurationInHour  int `json:"duration_in_hour"`
}
