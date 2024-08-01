package truco

import "fmt"

const (
	SAY_FLOR                = "say_flor"
	SAY_CON_FLOR_ME_ACHICO  = "say_con_flor_me_achico"
	SAY_CONTRAFLOR          = "say_contraflor"
	SAY_CONTRAFLOR_AL_RESTO = "say_contraflor_al_resto"
	SAY_CON_FLOR_QUIERO     = "say_con_flor_quiero"
	SAY_FLOR_SCORE          = "say_flor_score"
	SAY_FLOR_SON_BUENAS     = "say_flor_son_buenas"
	SAY_FLOR_SON_MEJORES    = "say_flor_son_mejores"
	REVEAL_FLOR_SCORE       = "reveal_flor_score"
)

type ActionSayFlor struct {
	act
	QuieroCost int `json:"quiero_cost"`
}
type ActionSayConFlorMeAchico struct {
	act
	Cost int `json:"cost"`
}
type ActionSayContraflor struct {
	act
	QuieroCost int `json:"quiero_cost"`
}
type ActionSayContraflorAlResto struct {
	act
	QuieroCost int `json:"quiero_cost"`
}
type ActionSayConFlorQuiero struct {
	act
	Cost int `json:"cost"`
}
type ActionSayFlorScore struct {
	act
	Score int `json:"score"`
}
type ActionSayFlorSonBuenas struct {
	act
}
type ActionSayFlorSonMejores struct {
	act
	Score int `json:"score"`
}
type ActionRevealFlorScore struct {
	act
	Score int `json:"score"`
}

func (a ActionSayFlor) IsPossible(g GameState) bool {
	return g.anyFlorActionIsPossible(&a) && len(_deserializeCurrentRoundActionsByPlayerID(a.PlayerID, g)) == 0
}

func (a ActionSayContraflor) IsPossible(g GameState) bool {
	return g.anyFlorActionIsPossible(&a)
}

func (a ActionSayContraflorAlResto) IsPossible(g GameState) bool {
	return g.anyFlorActionIsPossible(&a)
}

func (a ActionSayConFlorMeAchico) IsPossible(g GameState) bool {
	return g.anyFlorActionIsPossible(&a)
}

func (a ActionSayConFlorQuiero) IsPossible(g GameState) bool {
	return g.anyFlorActionIsPossible(&a)
}

func (a ActionSayFlorScore) IsPossible(g GameState) bool {
	return g.anyFlorActionIsPossible(&a)
}

func (a ActionSayFlorSonBuenas) IsPossible(g GameState) bool {
	var (
		myScore    = g.Players[a.PlayerID].Hand.FlorScore()
		theirScore = g.Players[g.OpponentOf(a.PlayerID)].Hand.FlorScore()
		iAmMano    = g.RoundTurnPlayerID == a.PlayerID
	)
	return g.anyFlorActionIsPossible(a) && (myScore < theirScore || (myScore == theirScore && !iAmMano))
}

func (a ActionSayFlorSonMejores) IsPossible(g GameState) bool {
	var (
		myScore    = g.Players[a.PlayerID].Hand.FlorScore()
		theirScore = g.Players[g.OpponentOf(a.PlayerID)].Hand.FlorScore()
		iAmMano    = g.RoundTurnPlayerID == a.PlayerID
	)
	return g.anyFlorActionIsPossible(&a) && (myScore > theirScore || (myScore == theirScore && iAmMano))
}

func (a ActionRevealFlorScore) IsPossible(g GameState) bool {
	if !g.RuleIsFlorEnabled {
		return false
	}
	if !g.Players[a.GetPlayerID()].Hand.HasFlor() {
		return false
	}
	roundLog := g.RoundsLog[g.RoundNumber]
	if roundLog.FlorWinnerPlayerID != a.PlayerID {
		return false
	}
	if g.FlorSequence.FlorPointsAwarded {
		return false
	}
	if !g.IsRoundFinished && g.Players[a.PlayerID].Score+roundLog.FlorPoints < g.RuleMaxPoints {
		return false
	}
	return len(g.Players[a.PlayerID].Hand.Revealed) != 3
}

func (g GameState) anyFlorActionIsPossible(a Action) bool {
	if !g.RuleIsFlorEnabled {
		return false
	}
	if !g.Players[a.GetPlayerID()].Hand.HasFlor() {
		return false
	}
	if g.IsRoundFinished {
		return false
	}
	// For any flor action except "say_flor" && "reveal_flor_score", both players must have flor
	if a.GetName() != SAY_FLOR && !g.Players[g.OpponentOf(a.GetPlayerID())].Hand.HasFlor() {
		return false
	}
	return g.FlorSequence.CanAddStep(a.GetName())
}

func (a ActionSayFlor) Run(g *GameState) error {
	if err := g.anyFlorActionRun(&a); err != nil {
		return err
	}
	if !g.Players[g.OpponentOf(a.PlayerID)].Hand.HasFlor() {
		g.FlorSequence.IsSinglePlayerFlor = true
		err := finalizeFlorSequence(a.PlayerID, g)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a ActionSayContraflor) Run(g *GameState) error {
	if err := g.anyFlorActionRun(&a); err != nil {
		return err
	}
	return nil
}

func (a ActionSayContraflorAlResto) Run(g *GameState) error {
	if err := g.anyFlorActionRun(&a); err != nil {
		return err
	}
	return nil
}

func (a ActionSayConFlorMeAchico) Run(g *GameState) error {
	if err := g.anyFlorActionRun(&a); err != nil {
		return err
	}
	return finalizeFlorSequence(g.OpponentOf(a.PlayerID), g)
}

func (a ActionSayConFlorQuiero) Run(g *GameState) error {
	if err := g.anyFlorActionRun(&a); err != nil {
		return err
	}
	return nil
}

func (a ActionSayFlorScore) Run(g *GameState) error {
	if err := g.anyFlorActionRun(&a); err != nil {
		return err
	}
	return nil
}

func (a ActionSayFlorSonBuenas) Run(g *GameState) error {
	if err := g.anyFlorActionRun(a); err != nil {
		return err
	}
	return finalizeFlorSequence(g.OpponentOf(a.PlayerID), g)
}

func (a ActionSayFlorSonMejores) Run(g *GameState) error {
	if err := g.anyFlorActionRun(&a); err != nil {
		return err
	}
	return finalizeFlorSequence(a.PlayerID, g)
}

func (a ActionRevealFlorScore) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.IsEnvidoFinished = true
	for len(g.Players[a.PlayerID].Hand.Unrevealed) > 0 {
		_ = g.Players[a.PlayerID].Hand.RevealCard(g.Players[a.PlayerID].Hand.Unrevealed[0])
	}
	if ok, err := g.tryAwardFlorPoints(); !ok {
		return fmt.Errorf("cannot award flor points: %w", err)
	}
	return nil
}

func (g *GameState) anyFlorActionRun(a Action) error {
	if !a.IsPossible(*g) {
		return errActionNotPossible
	}
	g.IsEnvidoFinished = true
	g.FlorSequence.AddStep(a.GetName())
	return nil
}

func finalizeFlorSequence(winnerPlayerID int, g *GameState) error {
	cost, err := g.FlorSequence.Cost(g.RuleMaxPoints, g.Players[winnerPlayerID].Score, g.Players[g.OpponentOf(winnerPlayerID)].Score, false)
	if err != nil {
		return err
	}
	g.RoundsLog[g.RoundNumber].FlorWinnerPlayerID = winnerPlayerID
	g.RoundsLog[g.RoundNumber].FlorPoints = cost
	_, _ = g.tryAwardFlorPoints()
	return nil
}

func (g *GameState) canAwardFlorPoints() (bool, error) {
	if !g.RuleIsFlorEnabled {
		return false, fmt.Errorf("flor is not enabled")
	}
	wonBy := g.RoundsLog[g.RoundNumber].FlorWinnerPlayerID
	if wonBy == -1 {
		return false, fmt.Errorf("no player won the flor")
	}
	if g.FlorSequence.FlorPointsAwarded {
		return false, fmt.Errorf("flor points already awarded")
	}
	// If the flor sequence was finished with "say_con_flor_me_achico", then
	// the points can be awarded immediately even though the sequence wasn't accepted
	// and the cards weren't revealed.
	if !g.FlorSequence.IsEmpty() && g.FlorSequence.IsFinished() && g.FlorSequence.Sequence[len(g.FlorSequence.Sequence)-1] == SAY_CON_FLOR_ME_ACHICO {
		return true, nil
	}
	if !g.FlorSequence.WasAccepted() {
		return false, fmt.Errorf("flor sequence was not accepted")
	}
	if len(g.Players[wonBy].Hand.Revealed) != 3 {
		return false, fmt.Errorf("player %v has not revealed all cards", wonBy)
	}
	return true, nil
}

func (g *GameState) tryAwardFlorPoints() (bool, error) {
	if ok, err := g.canAwardFlorPoints(); !ok {
		return false, err
	}
	wonBy := g.RoundsLog[g.RoundNumber].FlorWinnerPlayerID
	score := g.RoundsLog[g.RoundNumber].FlorPoints
	g.Players[wonBy].Score += score
	g.FlorSequence.FlorPointsAwarded = true
	return true, nil
}

func (a *ActionSayFlor) Enrich(g GameState) {
	var (
		seq, err0        = g.FlorSequence.WithStep(SAY_FLOR)
		youScore         = g.Players[a.GetPlayerID()].Score
		theirScore       = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		quieroSeq, err   = seq.WithStep(SAY_CON_FLOR_QUIERO)
		quieroCost, err2 = quieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err0 != nil {
		return
	}
	if err != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err))
	}
	if err2 != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err2))
	}
	a.QuieroCost = quieroCost
}
func (a *ActionSayContraflor) Enrich(g GameState) {
	var (
		seq, err0        = g.FlorSequence.WithStep(SAY_CONTRAFLOR)
		youScore         = g.Players[a.GetPlayerID()].Score
		theirScore       = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		quieroSeq, err   = seq.WithStep(SAY_CON_FLOR_QUIERO)
		quieroCost, err2 = quieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err0 != nil {
		return
	}
	if err != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err))
	}
	if err2 != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err2))
	}
	a.QuieroCost = quieroCost
}
func (a *ActionSayContraflorAlResto) Enrich(g GameState) {
	var (
		seq, err0        = g.FlorSequence.WithStep(SAY_CONTRAFLOR_AL_RESTO)
		youScore         = g.Players[a.GetPlayerID()].Score
		theirScore       = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		quieroSeq, err   = seq.WithStep(SAY_CON_FLOR_QUIERO)
		quieroCost, err2 = quieroSeq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err0 != nil {
		return
	}
	if err != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err))
	}
	if err2 != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err2))
	}
	a.QuieroCost = quieroCost
}
func (a *ActionSayConFlorMeAchico) Enrich(g GameState) {
	var (
		seq, err0         = g.FlorSequence.WithStep(SAY_CON_FLOR_ME_ACHICO)
		youScore          = g.Players[a.GetPlayerID()].Score
		theirScore        = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		noQuieroCost, err = seq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err0 != nil {
		return
	}
	if err != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err))
	}
	a.Cost = noQuieroCost
}
func (a *ActionSayConFlorQuiero) Enrich(g GameState) {
	var (
		seq, err0       = g.FlorSequence.WithStep(SAY_CON_FLOR_QUIERO)
		youScore        = g.Players[a.GetPlayerID()].Score
		theirScore      = g.Players[g.OpponentOf(a.GetPlayerID())].Score
		quieroCost, err = seq.Cost(g.RuleMaxPoints, youScore, theirScore, true)
	)
	if err0 != nil {
		return
	}
	if err != nil {
		panic(fmt.Errorf("when enriching %v: %w", a.GetName(), err))
	}
	a.Cost = quieroCost
}
func (a *ActionSayFlorScore) Enrich(g GameState) {
	a.Score = g.Players[a.PlayerID].Hand.FlorScore()
}
func (a *ActionSayFlorSonMejores) Enrich(g GameState) {
	a.Score = g.Players[a.PlayerID].Hand.FlorScore()
}
func (a *ActionRevealFlorScore) Enrich(g GameState) {
	a.Score = g.Players[a.PlayerID].Hand.FlorScore()
}

func (a ActionSayFlor) YieldsTurn(g GameState) bool {
	// If the opponent has flor too, it yields turn
	if g.Players[g.OpponentOf(a.PlayerID)].Hand.HasFlor() {
		return true
	}
	// If they don't, the turn should go to "mano", except if the opponent just said "truco"
	actionsOpponent := _deserializeCurrentRoundActionsByPlayerID(g.OpponentOf(a.PlayerID), g)
	if len(actionsOpponent) > 0 && actionsOpponent[len(actionsOpponent)-1].GetName() == SAY_TRUCO {
		return false
	}
	return g.RoundTurnPlayerID != a.PlayerID
}

func (a ActionSayFlorSonBuenas) YieldsTurn(g GameState) bool {
	// In son_buenas/son_mejores/no_quiero, the turn should go to whoever started the sequence
	return a.PlayerID != g.FlorSequence.StartingPlayerID
}

func (a ActionSayFlorSonMejores) YieldsTurn(g GameState) bool {
	// In son_buenas/son_mejores/no_quiero, the turn should go to whoever started the sequence
	// Unless the game should end due to the points won by this action.
	if g.Players[a.PlayerID].Score+g.RoundsLog[g.RoundNumber].FlorPoints >= g.RuleMaxPoints {
		return false
	}
	return a.PlayerID != g.FlorSequence.StartingPlayerID
}

func (a ActionSayConFlorMeAchico) YieldsTurn(g GameState) bool {
	// In son_buenas/son_mejores/no_quiero, the turn should go to whoever started the sequence
	return a.PlayerID != g.FlorSequence.StartingPlayerID
}

func (a ActionRevealFlorScore) YieldsTurn(g GameState) bool {
	// this action doesn't change turn because the round is finished at this point
	// and the current player must confirm round finished right after this action
	return false
}

func (a ActionSayConFlorQuiero) YieldsTurn(g GameState) bool {
	// In flor_quiero, the next turn should go to whoever has to reveal the score.
	// This should always be the "mano" player.
	return a.PlayerID != g.RoundTurnPlayerID
}

func (a ActionSayFlor) GetPriority() int {
	return 1
}

func (a ActionSayConFlorMeAchico) GetPriority() int {
	return 1
}

func (a ActionSayContraflor) GetPriority() int {
	return 1
}

func (a ActionSayContraflorAlResto) GetPriority() int {
	return 1
}

func (a ActionSayConFlorQuiero) GetPriority() int {
	return 1
}

func (a ActionSayFlorScore) GetPriority() int {
	return 1
}

func (a ActionSayFlorSonBuenas) GetPriority() int {
	return 1
}

func (a ActionSayFlorSonMejores) GetPriority() int {
	return 1
}

func (a ActionRevealFlorScore) GetPriority() int {
	return 2 // Because it's higher than confirming round finished
}
