package session

import (
	"log"
	"mafia-grpc/internal/event"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/role"
	"mafia-grpc/internal/state"
	"mafia-grpc/internal/util"
)

type Sender struct {
	gsession *GameSession
}

func (s Sender) HandleJE(joinEvent *event.JoinEvent, ctx event.EventHandleContext) {
	player_uuid := util.BytesToUUID(joinEvent.Player.Uuid)
	for uuid, psession := range s.gsession.players {
		hide := uuid != player_uuid
		copyJoinEvent := *joinEvent
		copyJoinEvent.Player = player.CopyPlayer(copyJoinEvent.Player, hide, hide)
		psession.EventChan <- event.JoinEventToGeneralEvent(&copyJoinEvent)
	}
}

func (s Sender) HandleLE(leftEvent *event.LeftEvent, ctx event.EventHandleContext) {
	player_uuid := util.BytesToUUID(leftEvent.Player.Uuid)
	for uuid, psession := range s.gsession.players {
		hide := uuid != player_uuid
		copyLeftEvent := *leftEvent
		copyLeftEvent.Player = player.CopyPlayer(copyLeftEvent.Player, hide, hide)
		psession.EventChan <- event.LeftEventToGeneralEvent(&copyLeftEvent)
	}
}

func (s Sender) HandleSDE(startDayEvent *event.StartDayEvent, ctx event.EventHandleContext) {
	for _, psession := range s.gsession.players {
		psession.EventChan <- event.StartDayEventToGeneralEvent(startDayEvent)
	}
}

func (s Sender) HandleSNE(startNightEvent *event.StartNightEvent, ctx event.EventHandleContext) {
	for _, psession := range s.gsession.players {
		psession.EventChan <- event.StartNightEventToGeneralEvent(startNightEvent)
	}
}

func (s Sender) HandleVE(voteEvent *event.VoteEvent, ctx event.EventHandleContext) {
	if ctx.State == nil {
		log.Fatal("missing game state")
	}
	player_uuid := util.BytesToUUID(voteEvent.Requester.Uuid)
	night := (ctx.State.PartOfTheDay == state.PartOfTheDay_NIGHT)
	for uuid, psession := range s.gsession.players {
		hideUuid := uuid != player_uuid
		copyVoteEvent := *voteEvent
		copyVoteEvent.Requester = player.CopyPlayer(copyVoteEvent.Requester, hideUuid, true)
		if copyVoteEvent.Vote != nil {
			copyVoteEvent.Vote = player.CopyPlayer(copyVoteEvent.Vote, hideUuid, true)
		}

		if !night || (psession.Player.Role == role.Role_MAFIA) {
			psession.EventChan <- event.VoteEventToGeneralEvent(&copyVoteEvent)
		}
	}
}

func (s Sender) HandleRVE(repeatVotingEvent *event.RepeatVotingEvent, ctx event.EventHandleContext) {
	if ctx.State == nil {
		log.Fatal("missing game state")
	}
	night := (ctx.State.PartOfTheDay == state.PartOfTheDay_NIGHT)
	for _, psession := range s.gsession.players {
		if !night || (psession.Player.Role == role.Role_MAFIA) {
			psession.EventChan <- event.RepeatVotingEventToGeneralEvent(repeatVotingEvent)
		}
	}
}

func (s Sender) HandleRE(revealedEvent *event.RevealedEvent, ctx event.EventHandleContext) {
	player_uuid := util.BytesToUUID(revealedEvent.Player.Uuid)
	for uuid, psession := range s.gsession.players {
		hideUuid := uuid != player_uuid
		copyRevealedEvent := *revealedEvent
		copyRevealedEvent.Player = player.CopyPlayer(copyRevealedEvent.Player, hideUuid, false)
		psession.EventChan <- event.RevealedEventToGeneralEvent(&copyRevealedEvent)
	}
}

func (s Sender) HandleGEE(gameEndEvent *event.GameEndEvent, ctx event.EventHandleContext) {
	for _, psession := range s.gsession.players {
		psession.EventChan <- event.GameEndEventToGeneralEvent(gameEndEvent)
	}
}

func (s Sender) Preprocess(e *event.Event, ctx event.EventHandleContext) {}
