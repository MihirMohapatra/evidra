package workflow

import "fmt"

type TransitionDef struct {
	From string
	To   string
	Name string
}

type StateMachine struct {
	transitions map[string]map[string]string
}

func New(defs []TransitionDef) *StateMachine {
	sm := &StateMachine{transitions: make(map[string]map[string]string)}
	for _, d := range defs {
		if sm.transitions[d.From] == nil {
			sm.transitions[d.From] = make(map[string]string)
		}
		sm.transitions[d.From][d.To] = d.Name
	}
	return sm
}

func (sm *StateMachine) CanTransition(from, to string) bool {
	_, ok := sm.transitions[from][to]
	return ok
}

func (sm *StateMachine) Transition(from, to string) error {
	if _, ok := sm.transitions[from][to]; !ok {
		return fmt.Errorf("invalid transition: %s -> %s", from, to)
	}
	return nil
}

func (sm *StateMachine) ActionName(from, to string) string {
	return sm.transitions[from][to]
}

func (sm *StateMachine) Sources(state string) []string {
	var result []string
	for from := range sm.transitions {
		if _, ok := sm.transitions[from][state]; ok {
			result = append(result, from)
		}
	}
	if result == nil {
		return []string{}
	}
	return result
}

func (sm *StateMachine) Targets(state string) []string {
	toMap, ok := sm.transitions[state]
	if !ok {
		return []string{}
	}
	var result []string
	for to := range toMap {
		result = append(result, to)
	}
	return result
}
