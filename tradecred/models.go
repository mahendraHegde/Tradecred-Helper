package tradecred

type Deal struct {
	ID         string
	Attributes struct {
		Name               string
		Code               string
		Type               string
		State              string
		Days               float64 `json:"financier_credit_period_in_days,omitempty"`
		Rate               float64 `json:"net_irr_to_be_given_to_investor_percentage,omitempty"`
		MinAmount          float64 `json:"minimum_purchase_amount,omitempty"`
		MinAmountSecondary float64 `json:"secondary_consideration_amount,omitempty"`
		DaysSecondary      float64 `json:"secondary_investment_tenure,omitempty"`
		RateSecondary      float64 `json:"effective_irr,omitempty"`
	}
	Relationships struct {
		DealTransaction struct {
			Data struct {
				ID string `json:",omitempty"`
			} `json:",omitempty"`
		} `json:"deal_transaction,omitempty"`
	} `json:",omitempty"`
}

type LiquidationReq struct {
	Data     []Deal
	Included []Deal
}
