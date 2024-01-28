package dbmodels

type UserOut struct {
	UserId     int    `json:"user_id"`
	FName      string `json:"fname"`
	SName      string `json:"sname"`
	PName      string `json:"pname"`
	GenderName string `json:"gender_name"`
	Age        uint8  `json:"age"`
	RegionCode string `json:"region_code"`
}

type UserIn struct {
	UserId   int    `json:"user_id"`
	FName    string `json:"fname"`
	SName    string `json:"sname"`
	PName    string `json:"pname"`
	GenderId int    `json:"gender_id"`
	Age      uint8  `json:"age"`
	RegionId int    `json:"region_code"`
}
