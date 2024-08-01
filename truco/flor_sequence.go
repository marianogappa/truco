package truco

import (
	"errors"
	"fmt"
	"strings"
)

const (
	COST_CONTRAFLOR_AL_RESTO = -2
)

var (
	// TODD: add son_buenas, son_mejores & say_flor_score
	validFlorSequenceCosts = map[string]int{
		SAY_FLOR:                              3,
		_s(SAY_FLOR, SAY_CONTRAFLOR):          COST_NOT_READY,
		_s(SAY_FLOR, SAY_CONTRAFLOR_AL_RESTO): COST_NOT_READY,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CONTRAFLOR_AL_RESTO): COST_NOT_READY,

		// All "me achico"
		_s(SAY_FLOR, SAY_CON_FLOR_ME_ACHICO):                                          3,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CON_FLOR_ME_ACHICO):                          4,
		_s(SAY_FLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_ME_ACHICO):                 4,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_ME_ACHICO): 6,

		// "quiero" to "flor"
		_s(SAY_FLOR, SAY_CON_FLOR_QUIERO):                                       4,
		_s(SAY_FLOR, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE):                       4,
		_s(SAY_FLOR, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_BUENAS):  4,
		_s(SAY_FLOR, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_MEJORES): 4,

		// "quiero" to "contraflor"
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CON_FLOR_QUIERO):                                       6,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE):                       6,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_BUENAS):  6,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_MEJORES): 6,

		// "quiero" to "contraflor" => "contraflor al resto"
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO):                                       COST_CONTRAFLOR_AL_RESTO,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE):                       COST_CONTRAFLOR_AL_RESTO,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_BUENAS):  COST_CONTRAFLOR_AL_RESTO,
		_s(SAY_FLOR, SAY_CONTRAFLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_MEJORES): COST_CONTRAFLOR_AL_RESTO,

		// "quiero" to "contraflor al resto"
		_s(SAY_FLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO):                                       COST_CONTRAFLOR_AL_RESTO,
		_s(SAY_FLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE):                       COST_CONTRAFLOR_AL_RESTO,
		_s(SAY_FLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_BUENAS):  COST_CONTRAFLOR_AL_RESTO,
		_s(SAY_FLOR, SAY_CONTRAFLOR_AL_RESTO, SAY_CON_FLOR_QUIERO, SAY_FLOR_SCORE, SAY_FLOR_SON_MEJORES): COST_CONTRAFLOR_AL_RESTO,
	}
)

func _s(ss ...string) string {
	return strings.Join(ss, ",")
}

type FlorSequence struct {
	Sequence []string `json:"sequence"`

	// IsSinglePlayerFlor is necessary because when only one player has a flor,
	// there's no way to know if .IsFinished() is true when there's only a SAY_FLOR step.
	IsSinglePlayerFlor bool `json:"isSinglePlayerFlor"`

	// This is necessary because when son_buenas/son_mejores/no_quiero is said,
	// the turn goes to whoever started the sequence (i.e. affects YieldsTurn)
	StartingPlayerID int `json:"startingPlayerID"`

	// FlorPointsAwarded is used to determine if the points have already been awarded.
	//
	// When a flor is won, points are not automatically awarded. They are when the
	// winning score is revealed. This could be at the time it is won, but normally
	// it is revealed later.
	FlorPointsAwarded bool `json:"florPointsAwarded"`
}

func (es FlorSequence) String() string {
	return strings.Join(es.Sequence, ",")
}

func (es FlorSequence) IsEmpty() bool {
	return len(es.Sequence) == 0
}

func (es FlorSequence) isValid() bool {
	_, ok := validFlorSequenceCosts[es.String()]
	return ok
}

func (es *FlorSequence) CanAddStep(step string) bool {
	es.Sequence = append(es.Sequence, step)
	isValid := es.isValid()
	es.Sequence = es.Sequence[:len(es.Sequence)-1]
	return isValid
}

func (es *FlorSequence) AddStep(step string) bool {
	if !es.CanAddStep(step) {
		return false
	}
	es.Sequence = append(es.Sequence, step)
	return true
}

func (es *FlorSequence) IsFinished() bool {
	if len(es.Sequence) == 0 {
		return false
	}
	last := es.Sequence[len(es.Sequence)-1]
	return last == SAY_FLOR_SON_BUENAS || last == SAY_FLOR_SON_MEJORES || last == SAY_CON_FLOR_ME_ACHICO || (last == SAY_FLOR && es.IsSinglePlayerFlor)
}

func (es FlorSequence) Cost(maxPoints, winnerPlayerScore, loserPlayerScore int, forHint bool) (int, error) {
	if !es.isValid() {
		return COST_NOT_READY, fmt.Errorf("%w: [%v]", errInvalidFlorSequence, strings.Join(es.Sequence, ","))
	}
	cost := validFlorSequenceCosts[es.String()]
	if cost == COST_CONTRAFLOR_AL_RESTO {
		return calculateFaltaEnvidoCost(maxPoints, winnerPlayerScore, loserPlayerScore), nil
	}
	// If this calculation is for enriching an action, we don't care if it's finished.
	// If it's for assigning cost, then it must be finished.
	if forHint {
		return cost, nil
	}
	if !es.IsFinished() {
		return cost, fmt.Errorf("%w: %v", errUnfinishedFlorSequence, strings.Join(es.Sequence, ","))
	}
	return cost, nil
}

func (es FlorSequence) WasAccepted() bool {
	for _, step := range es.Sequence {
		if step == SAY_CON_FLOR_QUIERO || (step == SAY_FLOR && es.IsSinglePlayerFlor) {
			return true
		}
	}
	return false
}

func (es FlorSequence) Clone() *FlorSequence {
	return &FlorSequence{
		Sequence:          append([]string{}, es.Sequence...),
		StartingPlayerID:  es.StartingPlayerID,
		FlorPointsAwarded: es.FlorPointsAwarded,
	}
}

func (es FlorSequence) WithStep(step string) (FlorSequence, error) {
	if !es.CanAddStep(step) {
		return es, fmt.Errorf("%w: [%v]", errInvalidFlorSequence, _s(es.Sequence...))
	}
	newEs := es.Clone()
	newEs.AddStep(step)
	return *newEs, nil
}

var (
	errInvalidFlorSequence    = errors.New("invalid flor sequence")
	errUnfinishedFlorSequence = errors.New("unfinished flor sequence")
)
