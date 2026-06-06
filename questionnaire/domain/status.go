package domain

type Status string

const (
	StatusUploaded   Status = "uploaded"
	StatusQueued     Status = "queued"
	StatusParsing    Status = "parsing"
	StatusParsed     Status = "parsed"
	StatusFailed     Status = "failed"
	StatusAssigned   Status = "assigned"
)

func (s Status) Valid() bool {
	switch s {
	case StatusUploaded, StatusQueued, StatusParsing, StatusParsed, StatusFailed, StatusAssigned:
		return true
	}
	return false
}

func (s Status) CanTransitionTo(next Status) bool {
	transitions := map[Status][]Status{
		StatusUploaded: {StatusQueued, StatusFailed},
		StatusQueued:   {StatusParsing, StatusFailed},
		StatusParsing:  {StatusParsed, StatusFailed},
		StatusParsed:   {StatusAssigned},
		StatusFailed:   {StatusQueued},
		StatusAssigned: {StatusParsed},
	}
	allowed, ok := transitions[s]
	if !ok {
		return false
	}
	for _, st := range allowed {
		if st == next {
			return true
		}
	}
	return false
}
