package core

type GetDealsQS struct {
	Days      int     `form:"days"`
	MaxAmount float64 `form:"maxAmount"`
	Rate      float64 `form:"rate"`
}
type Credentials struct {
	Email, Password string
}
type Headers struct {
	ApiKey string `header:"x-api-key"`
}
