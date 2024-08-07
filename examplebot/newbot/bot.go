package newbot

import (
	"fmt"
	"os"

	"log"

	"github.com/marianogappa/truco/truco"
)

type Bot struct {
	orderedRules []rule
	st           state
	logger       Logger
}

func WithDefaultLogger(b *Bot) {
	b.logger = log.New(os.Stderr, "", log.LstdFlags)
}

func New(opts ...func(*Bot)) *Bot {
	// Rules organically form a DAG. Kahn flattens them into a linear order.
	// If this is not possible (i.e. it's not a DAG), it blows up.
	orderedRules, err := topologicalSortKahn(rules)
	if err != nil {
		panic(fmt.Errorf("couldn't sort rules: %w; bot is defective! please report this bug!", err))
	}

	b := &Bot{orderedRules: orderedRules, logger: NoOpLogger{}, st: state{}}
	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (m Bot) ChooseAction(gs truco.ClientGameState) truco.Action {
	// Trivial cases
	if len(gs.PossibleActions) == 0 {
		return nil
	}
	if len(gs.PossibleActions) == 1 {
		return _deserializeActions(gs.PossibleActions)[0]
	}

	// If trickier, run rules
	for _, r := range m.orderedRules {
		if !r.isApplicable(m.st, gs) {
			continue
		}
		res, err := r.run(m.st, gs)
		if err != nil {
			panic(fmt.Errorf("Running rule %v found bug: %w. Please report this bug.", r.name, err))
		}
		m.logger.Printf("Running applicable rule %v: %s, result: %v", r.name, r.description, res.resultDescription)
		for _, sc := range res.stateChanges {
			m.logger.Printf("State change: %s", sc.description)
			sc.fn(&m.st)
		}
		if res.action != nil {
			return res.action
		}
	}

	// Running all rules MUST always result in an action being chosen
	panic("no action chosen after running all rules; bot is defective! please report this bug!")
}
