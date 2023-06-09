package event

const (
	MaxVotingNumberAttemts = 1
)

func JoinEventToGeneralEvent(j *JoinEvent) *Event {
	return &Event{Event: &Event_JoinEvent{JoinEvent: j}}
}

func LeftEventToGeneralEvent(l *LeftEvent) *Event {
	return &Event{Event: &Event_LeftEvent{LeftEvent: l}}
}

func StartDayEventToGeneralEvent(s *StartDayEvent) *Event {
	return &Event{Event: &Event_StartDayEvent{StartDayEvent: s}}
}

func StartNightEventToGeneralEvent(s *StartNightEvent) *Event {
	return &Event{Event: &Event_StartNightEvent{StartNightEvent: s}}
}

func VoteEventToGeneralEvent(v *VoteEvent) *Event {
	return &Event{Event: &Event_VoteEvent{VoteEvent: v}}
}

func RepeatVotingEventToGeneralEvent(r *RepeatVotingEvent) *Event {
	return &Event{Event: &Event_RepeatVotingEvent{RepeatVotingEvent: r}}
}

func RevealedEventToGeneralEvent(r *RevealedEvent) *Event {
	return &Event{Event: &Event_RevealedEvent{RevealedEvent: r}}
}

func GameEndEventToGeneralEvent(g *GameEndEvent) *Event {
	return &Event{Event: &Event_GameEndEvent{GameEndEvent: g}}
}
