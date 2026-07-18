package noaa

type StationsResponse struct {
	Stations []Station `json:"stations"`
}

type Station struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	State    string  `json:"state"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Timezone string  `json:"timezone"`
}
