package enrichmodels

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Nationalize struct {
	Count     int       `json:"count"`
	Name      string    `json:"name"`
	Countries []Country `json:"country"`
}
