package truco

import (
	"errors"
	"strings"
)

const (
	SAY_ENVIDO             = "say_envido"
	SAY_REAL_ENVIDO        = "say_real_envido"
	SAY_FALTA_ENVIDO       = "say_falta_envido"
	SAY_ENVIDO_QUIERO      = "say_envido_quiero"
	SAY_ENVIDO_NO_QUIERO   = "say_envido_no_quiero"
	SAY_SON_BUENAS         = "say_son_buenas"
	SAY_SON_MEJORES        = "say_son_mejores"
	SAY_ME_VOY_AL_MAZO     = "say_me_voy_al_mazo"
	REVEAL_CARD            = "reveal_card"
	CONFIRM_ROUND_FINISHED = "confirm_round_finished"
	SAY_ENVIDO_SCORE       = "say_envido_score"
	REVEAL_ENVIDO_SCORE    = "reveal_envido_score"

	COST_NOT_READY    = -1
	COST_FALTA_ENVIDO = -2
)

var (
	validEnvidoSequenceCosts = map[string]int{
		SAY_ENVIDO:                                                                                                          COST_NOT_READY,
		SAY_REAL_ENVIDO:                                                                                                     COST_NOT_READY,
		SAY_FALTA_ENVIDO:                                                                                                    COST_NOT_READY,
		_s(SAY_ENVIDO, SAY_ENVIDO):                                                                                          COST_NOT_READY,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO):                                                                                     COST_NOT_READY,
		_s(SAY_ENVIDO, SAY_FALTA_ENVIDO):                                                                                    COST_NOT_READY,
		_s(SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO):                                                                               COST_NOT_READY,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO):                                                                         COST_NOT_READY,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO):                                                                   COST_NOT_READY,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO):                                                       COST_NOT_READY,
		_s(SAY_ENVIDO, SAY_ENVIDO_QUIERO):                                                                                   2,
		_s(SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO):                                                                              3,
		_s(SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                                                             COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_QUIERO):                                                                       4,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO):                                                                  5,
		_s(SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                                                 COST_FALTA_ENVIDO,
		_s(SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                                            COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO):                                                      7,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                                COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO):                                    COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                                                 2,
		_s(SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                                            3,
		_s(SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                                           COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                                     4,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                                5,
		_s(SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                               COST_FALTA_ENVIDO,
		_s(SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                          COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                                    7,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                              COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE):                  COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                                                2,
		_s(SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                                           3,
		_s(SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                                          COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                                    4,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                               5,
		_s(SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                              COST_FALTA_ENVIDO,
		_s(SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                         COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):                   7,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES):             COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_MEJORES): COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                                                 2,
		_s(SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                                            3,
		_s(SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                                           COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                                     4,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                                5,
		_s(SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                               COST_FALTA_ENVIDO,
		_s(SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                          COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):                    7,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):              COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_QUIERO, SAY_ENVIDO_SCORE, SAY_SON_BUENAS):  COST_FALTA_ENVIDO,
		_s(SAY_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                                                1,
		_s(SAY_REAL_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                                           1,
		_s(SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                                          1,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                                    2,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                               2,
		_s(SAY_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                              2,
		_s(SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                         3,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                                   4,
		_s(SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                             5,
		_s(SAY_ENVIDO, SAY_ENVIDO, SAY_REAL_ENVIDO, SAY_FALTA_ENVIDO, SAY_ENVIDO_NO_QUIERO):                                 7,
	}
)

type EnvidoSequence struct {
	Sequence []string `json:"sequence"`

	// This is necessary because when son_buenas/son_mejores/no_quiero is said,
	// the turn goes to whoever started the envido sequence (i.e. affects YieldsTurn)
	StartingPlayerID int `json:"startingPlayerID"`

	// EnvidoPointsAwarded is used to determine if the points have already been awarded.
	//
	// When an envido is won, points are not automatically awarded. They are when the
	// winning score is revealed. This could be at the time it is won, but normally
	// it is revealed later.
	EnvidoPointsAwarded bool `json:"envidoPointsAwarded"`
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

func (es EnvidoSequence) Cost(maxPoints, winnerPlayerScore, loserPlayerScore int, forHint bool) (int, error) {
	if !es.isValid() {
		return COST_NOT_READY, errInvalidEnvidoSequence
	}
	cost := validEnvidoSequenceCosts[es.String()]
	if cost == COST_FALTA_ENVIDO {
		return calculateFaltaEnvidoCost(maxPoints, winnerPlayerScore, loserPlayerScore), nil
	}
	// If this is a hint, return the cost as is, without checking if the sequence is finished.
	if forHint {
		return cost, nil
	}
	if !es.IsFinished() {
		return cost, errUnfinishedEnvidoSequence
	}
	return cost, nil
}

func (es EnvidoSequence) WasAccepted() bool {
	for _, step := range es.Sequence {
		if step == SAY_ENVIDO_QUIERO {
			return true
		}
	}
	return false
}

func (es EnvidoSequence) Clone() *EnvidoSequence {
	return &EnvidoSequence{
		Sequence:            append([]string{}, es.Sequence...),
		StartingPlayerID:    es.StartingPlayerID,
		EnvidoPointsAwarded: es.EnvidoPointsAwarded,
	}
}

func (es EnvidoSequence) WithStep(step string) (EnvidoSequence, error) {
	if !es.CanAddStep(step) {
		return es, errInvalidEnvidoSequence
	}
	newEs := es.Clone()
	newEs.AddStep(step)
	return *newEs, nil
}

func calculateFaltaEnvidoCost(maxPoints, winnerScore, loserScore int) int {
	// maxPoints is normally only 15 or 30, but if it's set to less then
	// use the same rule as for 15, but using maxPoints instead.
	if maxPoints < 15 {
		return maxPoints - winnerScore
	}
	if winnerScore < 15 && loserScore < 15 {
		return 15 - winnerScore
	}
	return 30 - max(winnerScore, loserScore)
}

var (
	errInvalidEnvidoSequence    = errors.New("invalid envido sequence")
	errUnfinishedEnvidoSequence = errors.New("unfinished envido sequence")
)
