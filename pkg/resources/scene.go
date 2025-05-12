package resources

type Action struct {
	Target Reference `json:"target"`
	Action struct {
		On               LightOn          `json:"on"`
		Dimming          Dimming          `json:"dimming"`
		ColorTemperature ColorTemperature `json:"color_temperature"`
	}
}

type Palette struct{}

type Scene struct {
	ID       string   `json:"id"`
	IDV1     string   `json:"id_v1"`
	Actions  []Action `json:"actions"`
	Palette  Palette  `json:"palette"`
	Recall   struct{} `json:"recall"`
	Metadata struct {
		Name  string    `json:"name"`
		Image Reference `json:"image"`
	}
	// Can either be room or zone.
	Group       Reference `json:"group"`
	Speed       float64   `json:"speed"`
	AutoDynamic bool      `json:"auto_dynamic"`
	Status      struct {
		Active     string `json:"active"`
		LastRecall string `json:"last_recall"`
	}
	Type string `json:"type"`
}
