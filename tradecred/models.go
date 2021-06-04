package tradecred

type Deal struct {
	Attributes struct {
		Code      string
		Type      string
		State     string
		Days      float64 `json:"financier_credit_period_in_days"`
		Rate      float64 `json:"net_irr_to_be_given_to_investor_percentage"`
		MinAmount float64 `json:"minimum_purchase_amount"`
	}
}
