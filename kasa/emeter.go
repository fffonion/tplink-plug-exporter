package kasa

type KasaClientEmeterService struct {
	c *KasaClient
}

type GetRealtimeRequest struct {
}

type GetRealtimeResponse struct {
	// some fw/hw may use the following keys
	Current float64 `mapstructure:"current"` // unit is A
	Voltage float64 `mapstructure:"voltage"` // unit is V
	Power   float64 `mapstructure:"power"`   // unit is W
	Total   float64 `mapstructure:"total"`   // unit is kWh

	// some may use these
	CurrentmA float64 `mapstructure:"current_ma"`
	VoltagemV float64 `mapstructure:"voltage_mv"`
	PowermW   float64 `mapstructure:"power_mw"`
	TotalWh   float64 `mapstructure:"total_wh"`
}

func (r *GetRealtimeResponse) Normalize() {
	if r.TotalWh != -1 {
		r.Current = r.CurrentmA / 1000
		r.Voltage = r.VoltagemV / 1000
		r.Power = r.PowermW / 1000
		r.Total = r.TotalWh / 1000
	}
}

func (s *KasaClientEmeterService) GetRealtime() (*GetRealtimeResponse, error) {
	response := GetRealtimeResponse{
		TotalWh: -1,
	}
	err := s.c.RPC("emeter", "get_realtime", GetRealtimeRequest{}, &response)

	response.Normalize()

	return &response, err
}
