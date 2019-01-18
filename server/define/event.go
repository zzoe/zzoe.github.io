package define

type Event int

const (
	EventInit Event = iota
	EventCfgChange
	EventResourceClear
)
