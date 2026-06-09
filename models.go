package review_sdk

type Review struct {
	ID        string `json:"id"`
	UserID    int    `json:"userId"`
	CityID    int    `json:"cityId"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"createdAt"`
}

type NewReview struct {
	UserID  int    `json:"userId"`
	CityID  int    `json:"cityId"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
