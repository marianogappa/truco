package newbot

import (
	"encoding/json"

	"github.com/marianogappa/truco/truco"
)

func _deserializeActions(as []json.RawMessage) []truco.Action {
	_as := []truco.Action{}
	for _, a := range as {
		_a, _ := truco.DeserializeAction(a)
		_as = append(_as, _a)
	}
	return _as
}

func isPossibleAll(st state, actionNames ...string) bool {
	for _, actionName := range actionNames {
		if _, ok := st["possibleActionNameSet"].(map[string]struct{})[actionName]; !ok {
			return false
		}
	}
	return true
}

// func isPossibleAny(st state, actionNames ...string) bool {
// 	for _, actionName := range actionNames {
// 		if _, ok := st["possibleActionNameSet"].(map[string]struct{})[actionName]; ok {
// 			return true
// 		}
// 	}
// 	return false
// }

func getAction(st state, actionName string) truco.Action {
	return st["possibleActions"].(map[string]truco.Action)[actionName]
}

func aggresiveness(st state) string {
	return st["aggresiveness"].(string)
}

func florScore(st state) int {
	return st["florScore"].(int)
}

func envidoScore(st state) int {
	return st["envidoScore"].(int)
}

func pointsToLose(st state) int {
	return st["pointsToLose"].(int)
}
