package types

type Place struct {
	ID       int      `csv:"-" json:"id"`
	Name     string   `csv:"Name" json:"name"`
	Address  string   `csv:"Address" json:"address"`
	Phone    string   `csv:"Phone" json:"phone"`
	Location GeoPoint `csv:"-" json:"location"`
}

type Response struct { 
	Name 		string   `json:"name"`
	Total 		int 	`json:"total"`
	Places 		[]Place    `json:"places"`
	Previous    int
	Next		int
	Last		int
}

type User struct {
	Username		string
	Password		string 
}

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Result struct {
	Name    string `csv:"Name" json:"name"`
	Address string `csv:"Address" json:"address"`
	Phone   string `csv:"Phone" json:"phone"`
}
