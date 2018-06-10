package gql

type StatRanks struct {
	Strength     []*Stats `json:"strength"`
	Intelligence []*Stats `json:"intelligence"`
	Special      []*Stats `json:"special"`
}
