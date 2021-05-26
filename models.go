package main

type Cowin struct {
	Centers []Center `json:"centers"`
}

type Center struct {
	CenterID      int     `json:"center_id"`
	Name          string  `json:"name"`
	NameL         string  `json:"name_l"`
	Address       string  `json:"address"`
	AddressL      string  `json:"address_l"`
	StateName     string  `json:"state_name"`
	StateNameL    string  `json:"state_name_l"`
	DistrictName  string  `json:"district_name"`
	DistrictNameL string  `json:"district_name_l"`
	BlockName     string  `json:"block_name"`
	BlockNameL    string  `json:"block_name_l"`
	Pincode       int     `json:"pincode"`
	Lat           float64 `json:"lat"`
	Long          float64 `json:"long"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	FeeType       string  `json:"fee_type"`
	VaccineFees   []struct {
		Vaccine string `json:"vaccine"`
		Fee     string `json:"fee"`
	} `json:"vaccine_fees"`
	Sessions []Session `json:"sessions"`
}

type Session struct {
	SessionID              string   `json:"session_id"`
	Date                   string   `json:"date"`
	AvailableCapacity      int      `json:"available_capacity"`
	AvailableCapacityDose1 int      `json:"available_capacity_dose1"`
	AvailableCapacityDose2 int      `json:"available_capacity_dose2"`
	MinAgeLimit            int      `json:"min_age_limit"`
	Vaccine                string   `json:"vaccine"`
	Slots                  []string `json:"slots"`
}
