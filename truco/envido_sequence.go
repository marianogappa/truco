package truco

import (
	"errors"
	"fmt"
	"strings"
)

const (
	SAY_ENVIDO           = "say_envido"
	SAY_REAL_ENVIDO      = "say_real_envido"
	SAY_FALTA_ENVIDO     = "say_falta_envido"
	SAY_ENVIDO_QUIERO    = "say_envido_quiero"
	SAY_ENVIDO_NO_QUIERO = "say_envido_no_quiero"
	SAY_SON_BUENAS       = "say_son_buenas"
	SAY_SON_MEJORES      = "say_son_mejores"
	SAY_ME_VOY_AL_MAZO   = "say_me_voy_al_mazo"
	REVEAL_CARD          = "reveal_card"

	COST_NOT_READY    = -1
	COST_FALTA_ENVIDO = -2
)

var (
	validEnvidoSequenceCosts = map[string]int{
		SAY_ENVIDO:       COST_NOT_READY,
		SAY_REAL_ENVIDO:  COST_NOT_READY,
		SAY_FALTA_ENVIDO: COST_NOT_READY,
		fmt.Sprintf("%s,%s", SAY_ENVIDO, SAY_ENVIDO):                                                                                    COST_NOT_READY,
		fmt.Sprintf("%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO):                                                                               COST_NOT_READY,
		fmt.Sprintf("%s,%s", SAY_ENVIDO, SAY_FALTA_ENVIDO):                                                                              COST_NOT_READY,
		fmt.Sprintf("%s,%s", SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO):                                                                         COST_NOT_READY,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO):                                                                COST_NOT_READY,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO):                                                          COST_NOT_READY,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO):                                           COST_NOT_READY,
		fmt.Sprintf("%s,%s", SAY_ENVIDO, SAY_ENVIDO_QUIERO):                                                                             2,
		fmt.Sprintf("%s,%s", SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO):                                                                        3,
		fmt.Sprintf("%s,%s", SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                                                       COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_QUIERO):                                                              4,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO):                                                         5,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                                        COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s", SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                                   COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO):                                          7,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                    COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                     COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                                                         2,
		fmt.Sprintf("%s,%s,%s", SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                                                    3,
		fmt.Sprintf("%s,%s,%s", SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                                                   COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                                          4,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                                     5,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                                    COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s", SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                               COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                      7,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES):                COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_MEJORES): COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                                                          2,
		fmt.Sprintf("%s,%s,%s", SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                                                     3,
		fmt.Sprintf("%s,%s,%s", SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                                                    COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                                           4,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                                      5,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                                     COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s", SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                                COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                       7,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):                 COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_SON_BUENAS):  COST_FALTA_ENVIDO,
		fmt.Sprintf("%s,%s", SAY_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                                          1,
		fmt.Sprintf("%s,%s", SAY_REAL_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                                     1,
		fmt.Sprintf("%s,%s", SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                                    1,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                           2,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                      2,
		fmt.Sprintf("%s,%s,%s", SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                     2,
		fmt.Sprintf("%s,%s,%s", SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                3,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                       4,
		fmt.Sprintf("%s,%s,%s,%s", SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                 5,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                  7,
	}
)

type EnvidoSequence struct {
	Sequence []string `json:"sequence"`

	// This is necessary because when son_buenas/son_mejores/no_quiero is said,
	// the turn goes to whoever started the envido sequence (i.e. affects YieldsTurn)
	StartingPlayerID int `json:"startingPlayerID"`
}

func (es EnvidoSequence) String() string {
	return strings.Join(es.Sequence, ",")
}

func (es EnvidoSequence) IsEmpty() bool {
	return len(es.Sequence) == 0
}

func (es EnvidoSequence) isValid() bool {
	_, ok := validEnvidoSequenceCosts[es.String()]
	return ok
}

func (es *EnvidoSequence) CanAddStep(step string) bool {
	es.Sequence = append(es.Sequence, step)
	isValid := es.isValid()
	es.Sequence = es.Sequence[:len(es.Sequence)-1]
	return isValid
}

func (es *EnvidoSequence) AddStep(step string) bool {
	if !es.CanAddStep(step) {
		return false
	}
	es.Sequence = append(es.Sequence, step)
	return true
}

func (es *EnvidoSequence) IsFinished() bool {
	if len(es.Sequence) == 0 {
		return false
	}
	last := es.Sequence[len(es.Sequence)-1]
	return last == SAY_SON_BUENAS || last == SAY_SON_MEJORES || last == SAY_ENVIDO_NO_QUIERO
}

func (es EnvidoSequence) Cost(currentPlayerScore int, otherPlayerScore int) (int, error) {
	if !es.isValid() {
		return COST_NOT_READY, errInvalidEnvidoSequence
	}
	if !es.IsFinished() {
		return COST_NOT_READY, errUnfinishedEnvidoSequence
	}
	cost := validEnvidoSequenceCosts[es.String()]
	if cost == COST_FALTA_ENVIDO {
		return es.calculateFaltaEnvidoCost(currentPlayerScore, otherPlayerScore), nil
	}
	return cost, nil
}

func (es EnvidoSequence) calculateFaltaEnvidoCost(meScore int, youScore int) int {
	if meScore < 15 && youScore < 15 {
		return 15 - meScore
	}
	return 30 - max(meScore, youScore)
}

var (
	errInvalidEnvidoSequence    = errors.New("invalid envido sequence")
	errUnfinishedEnvidoSequence = errors.New("unfinished envido sequence")
)
