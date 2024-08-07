package truco

import (
	"errors"
	"fmt"
	"strings"
)

const (
	SAY_TRUCO              = "say_truco"
	SAY_QUIERO_RETRUCO     = "say_quiero_retruco"
	SAY_QUIERO_VALE_CUATRO = "say_quiero_vale_cuatro"
	SAY_TRUCO_QUIERO       = "say_truco_quiero"
	SAY_TRUCO_NO_QUIERO    = "say_truco_no_quiero"
)

var (
	validTrucoSequenceCosts = map[string]int{
		SAY_TRUCO: 1,
		fmt.Sprintf("%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO):                                                                                              2,
		fmt.Sprintf("%s,%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO, SAY_QUIERO_VALE_CUATRO):                                                                   3,
		fmt.Sprintf("%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO):                                                                                                2,
		fmt.Sprintf("%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO):                                                                         2,
		fmt.Sprintf("%s,%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO, SAY_TRUCO_QUIERO):                                                                         3,
		fmt.Sprintf("%s,%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_TRUCO_QUIERO):                                                    3,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_VALE_CUATRO):                         3,
		fmt.Sprintf("%s,%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO, SAY_TRUCO_NO_QUIERO):                                                                      2,
		fmt.Sprintf("%s,%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_TRUCO_NO_QUIERO):                                                 2,
		fmt.Sprintf("%s,%s,%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO, SAY_QUIERO_VALE_CUATRO, SAY_TRUCO_QUIERO):                                              4,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_QUIERO_VALE_CUATRO, SAY_TRUCO_QUIERO):                         4,
		fmt.Sprintf("%s,%s,%s,%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_VALE_CUATRO, SAY_TRUCO_QUIERO):    4,
		fmt.Sprintf("%s,%s", SAY_TRUCO, SAY_TRUCO_NO_QUIERO):                                                                                             1,
		fmt.Sprintf("%s,%s,%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO, SAY_QUIERO_VALE_CUATRO, SAY_TRUCO_NO_QUIERO):                                           3,
		fmt.Sprintf("%s,%s,%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_QUIERO_VALE_CUATRO, SAY_TRUCO_NO_QUIERO):                      3,
		fmt.Sprintf("%s,%s,%s,%s,%s,%s", SAY_TRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_RETRUCO, SAY_TRUCO_QUIERO, SAY_QUIERO_VALE_CUATRO, SAY_TRUCO_NO_QUIERO): 3,
	}
)

type TrucoSequence struct {
	// StartingPlayerID is the player ID that started the truco sub-sequence.
	//
	// It is used to determine "YieldsTurn" after a "quiero" action.
	//
	// Note the word "sub-sequence". There can be 0 to 3 truco sub-sequences in a round.
	//
	// Sub-sequences are separated by "quiero" actions.
	//
	// StartingPlayerID holds the player ID that started the _current_ sub-sequence.
	StartingPlayerID int `json:"starting_player_id"`

	// QuieroOwnerPlayerID is the player ID of the player who said "quiero" last in the truco
	// sequence. This is used to determine who can raise the stakes in the truco sequence.
	QuieroOwnerPlayerID int `json:"quiero_owner_player_id"`

	Sequence []string `json:"sequence"`
}

func (ts TrucoSequence) String() string {
	return strings.Join(ts.Sequence, ",")
}

func (ts TrucoSequence) IsEmpty() bool {
	return len(ts.Sequence) == 0
}

func (ts TrucoSequence) isValid() bool {
	_, ok := validTrucoSequenceCosts[ts.String()]
	return ok
}

func (ts *TrucoSequence) CanAddStep(step string) bool {
	ts.Sequence = append(ts.Sequence, step)
	isValid := ts.isValid()
	ts.Sequence = ts.Sequence[:len(ts.Sequence)-1]
	return isValid
}

func (ts *TrucoSequence) AddStep(step string) bool {
	if !ts.CanAddStep(step) {
		return false
	}
	ts.Sequence = append(ts.Sequence, step)
	return true
}

func (ts *TrucoSequence) IsFinished() bool {
	if len(ts.Sequence) == 0 {
		return false
	}
	last := ts.Sequence[len(ts.Sequence)-1]
	return last == SAY_TRUCO_QUIERO || last == SAY_TRUCO_NO_QUIERO
}

func (ts TrucoSequence) Cost() int {
	return validTrucoSequenceCosts[ts.String()]
}

func (ts TrucoSequence) IsSubsequenceStart() bool {
	// Subsequences are delimited by "quiero" actions.
	// It's necessary to store the playerID that started the current sub-sequence,
	// so that we can determine "YieldsTurn" after a "quiero" action.
	if len(ts.Sequence) == 0 {
		return false
	}
	if len(ts.Sequence) == 1 {
		return true
	}
	if ts.Sequence[len(ts.Sequence)-2] == SAY_TRUCO_QUIERO {
		return true
	}
	return false
}

func (es TrucoSequence) Clone() *TrucoSequence {
	return &TrucoSequence{
		Sequence:            append([]string{}, es.Sequence...),
		StartingPlayerID:    es.StartingPlayerID,
		QuieroOwnerPlayerID: es.QuieroOwnerPlayerID,
	}
}

func (es TrucoSequence) WithStep(step string) (TrucoSequence, error) {
	if !es.CanAddStep(step) {
		return es, errInvalidTrucoSequence
	}
	newEs := es.Clone()
	newEs.AddStep(step)
	return *newEs, nil
}

func (ts TrucoSequence) WasAccepted() bool {
	for _, step := range ts.Sequence {
		if step == SAY_TRUCO_QUIERO {
			return true
		}
	}
	return false
}

var (
	errInvalidTrucoSequence = errors.New("invalid truco sequence")
)
