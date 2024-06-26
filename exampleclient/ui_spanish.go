package exampleclient

import (
	"fmt"

	"github.com/marianogappa/truco/truco"
)

func getLastActionString(rs renderState) string {
	if rs.lastActionLog == nil {
		if rs.roundNumber == 1 {
			return "¡Empezó el juego!"
		}
		return "¡Empezó la mano!"
	}

	return getActionString(*rs.lastActionLog, rs.you)
}

func getActionString(log truco.ActionLog, playerID int) string {
	lastAction, _ := truco.DeserializeAction(log.Action)

	said := "dijiste"
	revealed := "tiraste"
	who := "Vos"
	if playerID != log.PlayerID {
		who = "Elle"
		said = "dijo"
		revealed = "tiró"
	}

	var what string
	switch lastAction.GetName() {
	case truco.REVEAL_CARD:
		action := lastAction.(*truco.ActionRevealCard)
		what = fmt.Sprintf("%v la carta %v", revealed, getCardString(action.Card))
	case truco.SAY_ENVIDO:
		what = fmt.Sprintf("%v envido", said)
	case truco.SAY_REAL_ENVIDO:
		what = fmt.Sprintf("%v real envido", said)
	case truco.SAY_FALTA_ENVIDO:
		what = fmt.Sprintf("%v falta envido!", said)
	case truco.SAY_ENVIDO_QUIERO:
		action := lastAction.(*truco.ActionSayEnvidoQuiero)
		what = fmt.Sprintf("%v quiero con %d", said, action.Score)
	case truco.SAY_ENVIDO_NO_QUIERO:
		what = fmt.Sprintf("%v no quiero", said)
	case truco.SAY_TRUCO:
		what = fmt.Sprintf("%v truco", said)
	case truco.SAY_TRUCO_QUIERO:
		what = fmt.Sprintf("%v quiero", said)
	case truco.SAY_TRUCO_NO_QUIERO:
		what = fmt.Sprintf("%v no quiero", said)
	case truco.SAY_QUIERO_RETRUCO:
		what = fmt.Sprintf("%v quiero retruco", said)
	case truco.SAY_QUIERO_VALE_CUATRO:
		what = fmt.Sprintf("%v quiero vale cuatro", said)
	case truco.SAY_SON_BUENAS:
		what = fmt.Sprintf("%v son buenas", said)
	case truco.SAY_SON_MEJORES:
		action := lastAction.(*truco.ActionSaySonMejores)
		what = fmt.Sprintf("%v %d son mejores", said, action.Score)
	case truco.SAY_ME_VOY_AL_MAZO:
		what = fmt.Sprintf("%v me voy al mazo", said)
	default:
		what = "???"
	}

	return fmt.Sprintf("%v %v\n", who, what)
}

func spanishScore(score int) string {
	if score == 1 {
		return "1 mala"
	}
	if score < 15 {
		return fmt.Sprintf("%d malas", score)
	}
	if score == 15 {
		return "entraste"
	}
	return fmt.Sprintf("%d buenas", score-14)
}

func spanishAction(action truco.Action) string {
	switch action.GetName() {
	case truco.REVEAL_CARD:
		_action := action.(*truco.ActionRevealCard)
		return getCardString(_action.Card)
	case truco.SAY_ENVIDO:
		return "envido"
	case truco.SAY_REAL_ENVIDO:
		return "real envido"
	case truco.SAY_FALTA_ENVIDO:
		return "falta envido"
	case truco.SAY_ENVIDO_QUIERO:
		return "quiero"
	case truco.SAY_ENVIDO_NO_QUIERO:
		return "no quiero"
	case truco.SAY_TRUCO:
		return "truco"
	case truco.SAY_TRUCO_QUIERO:
		return "quiero"
	case truco.SAY_TRUCO_NO_QUIERO:
		return "no quiero"
	case truco.SAY_QUIERO_RETRUCO:
		return "quiero retruco"
	case truco.SAY_QUIERO_VALE_CUATRO:
		return "quiero vale cuatro"
	case truco.SAY_SON_BUENAS:
		return "son buenas"
	case truco.SAY_SON_MEJORES:
		_action := action.(*truco.ActionSaySonMejores)
		return fmt.Sprintf("%v son mejores", _action.Score)
	case truco.SAY_ME_VOY_AL_MAZO:
		return "me voy al mazo"
	default:
		return "???"
	}
}
