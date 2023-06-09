package event

import (
	"log"
	player "mafia-grpc/internal/player"
	"mafia-grpc/internal/state"
)

type EventHandleContext struct {
	Player *player.Player
	State  *state.State
}

type EventHandler interface {
	Preprocess(*Event, EventHandleContext)
	HandleJE(*JoinEvent, EventHandleContext)
	HandleLE(*LeftEvent, EventHandleContext)
	HandleSDE(*StartDayEvent, EventHandleContext)
	HandleSNE(*StartNightEvent, EventHandleContext)
	HandleVE(*VoteEvent, EventHandleContext)
	HandleRVE(*RepeatVotingEvent, EventHandleContext)
	HandleRE(*RevealedEvent, EventHandleContext)
	HandleGEE(*GameEndEvent, EventHandleContext)
}

func Handle(handler EventHandler, e *Event, ctx EventHandleContext) {
	handler.Preprocess(e, ctx)
	joinEvent := e.GetJoinEvent()
	if joinEvent != nil {
		handler.HandleJE(joinEvent, ctx)
		return
	}
	leftEvent := e.GetLeftEvent()
	if leftEvent != nil {
		handler.HandleLE(leftEvent, ctx)
		return
	}
	startDayEvent := e.GetStartDayEvent()
	if startDayEvent != nil {
		handler.HandleSDE(startDayEvent, ctx)
		return
	}
	startNightEvent := e.GetStartNightEvent()
	if startNightEvent != nil {
		handler.HandleSNE(startNightEvent, ctx)
		return
	}
	voteEvent := e.GetVoteEvent()
	if voteEvent != nil {
		handler.HandleVE(voteEvent, ctx)
		return
	}
	repeatVotingEvent := e.GetRepeatVotingEvent()
	if repeatVotingEvent != nil {
		handler.HandleRVE(repeatVotingEvent, ctx)
		return
	}
	revealedEvent := e.GetRevealedEvent()
	if revealedEvent != nil {
		handler.HandleRE(revealedEvent, ctx)
		return
	}
	gameEndEvent := e.GetGameEndEvent()
	if gameEndEvent != nil {
		handler.HandleGEE(gameEndEvent, ctx)
		return
	}
	log.Fatalln("unknown event")
}
