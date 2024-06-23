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

	TRUCO_COST_NOT_READY = -1
)

var (
	validTrucoSequenceCosts = map[string]int{
		SAY_TRUCO: COST_NOT_READY,
		fmt.Sprintf("%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO):                                                                                              COST_NOT_READY,
		fmt.Sprintf("%s,%s,%s", SAY_TRUCO, SAY_QUIERO_RETRUCO, SAY_QUIERO_VALE_CUATRO):                                                                   COST_NOT_READY,
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

func (ts TrucoSequence) Cost() (int, error) {
	if !ts.isValid() {
		return COST_NOT_READY, errInvalidTrucoSequence
	}
	if !ts.IsFinished() {
		return COST_NOT_READY, errUnfinishedTrucoSequence
	}
	return validTrucoSequenceCosts[ts.String()], nil
}

var (
	errInvalidTrucoSequence    = errors.New("invalid truco sequence")
	errUnfinishedTrucoSequence = errors.New("unfinished truco sequence")
)
