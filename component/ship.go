package component

// ShipList is the list of ships returned by the game.
type ShipList struct {
	Ships []Ship `yaml:"ships"`
}

// Ship is the response from the API about ships.
type Ship struct {
	// Cargo []ShipCargo `yaml:"cargo"`
	Class          string `yaml:"class"`
	FlightPlanId   string `yaml:"flightPlanId"`
	Id             string `yaml:"id"`
	Location       string `yaml:"location"`
	Manufacturer   string `yaml:"manufacturer"`
	MaxCargo       int    `yaml:"maxCargo"`
	Plating        int    `yaml:"plating"`
	SpaceAvailable int    `yaml:"spaceAvailable"`
	Speed          int    `yaml:"speed"`
	Type           string `yaml:"type"`
	Weapons        int    `yaml:"weapons"`
	X              int    `yaml:"x"`
	Y              int    `yaml:"y"`
}
