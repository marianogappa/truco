package truco

type act struct {
	Name string `json:"name"`
}

func (a act) GetName() string {
	return a.Name
}

func (a act) YieldsTurn(g GameState) bool {
	return true
}
