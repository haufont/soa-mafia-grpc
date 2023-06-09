package cli

import (
	"mafia-grpc/internal/event"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/role"
	"mafia-grpc/internal/util"
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type EventPrinter struct {
}

func (EventPrinter) HandleJE(e *event.JoinEvent, ctx event.EventHandleContext) {
	var playerUuid uuid.UUID
	if ctx.Player != nil {
		playerUuid = util.BytesToUUID(ctx.Player.Uuid)
	}
	color.Cyan("[Event] Join player (%s)", player.PlayerToString(e.GetPlayer(), playerUuid))
}

func (EventPrinter) HandleLE(e *event.LeftEvent, ctx event.EventHandleContext) {
	var playerUuid uuid.UUID
	if ctx.Player != nil {
		playerUuid = util.BytesToUUID(ctx.Player.Uuid)
	}
	color.Cyan("[Event] Left player (%s)\n", player.PlayerToString(e.GetPlayer(), playerUuid))
}

func (EventPrinter) HandleSDE(e *event.StartDayEvent, ctx event.EventHandleContext) {
	if e.KilledPlayer == nil {
		color.Cyan("[Event] The day has started")
	} else {
		color.Cyan("[Event] The day has started. Mafia killed \"%s\" player", e.KilledPlayer.Name)
	}
}

func (EventPrinter) HandleSNE(e *event.StartNightEvent, ctx event.EventHandleContext) {
	if e.KilledPlayer == nil {
		color.Cyan("[Event] The night has started")
	} else {
		color.Cyan("[Event] The night has started. The town executed \"%s\" player", e.KilledPlayer.Name)
	}
}

func (EventPrinter) HandleVE(e *event.VoteEvent, ctx event.EventHandleContext) {
	if e.Requester == nil || e.Requester.Name == "" {
		color.Red("[Event] Invalid vote event")
		return
	}

	if e.Vote == nil || e.Vote.Name == "" {
		color.Cyan("[Event] %s player abstained from voting", e.Requester.Name)
	} else {
		color.Cyan("[Event] %s player voted for %s", e.Requester.Name, e.Vote.Name)
	}
}

func (EventPrinter) HandleRVE(e *event.RepeatVotingEvent, ctx event.EventHandleContext) {
	if e.Attempt == event.MaxVotingNumberAttemts {
		color.Cyan("[Event] Repeat voting because there is no majority. Attempt: %d (final)", e.Attempt)
	} else {
		color.Cyan("[Event] Repeat voting because there is no majority. Attempt: %d", e.Attempt)
	}
}

func (EventPrinter) HandleRE(e *event.RevealedEvent, ctx event.EventHandleContext) {
	if e.Player == nil {
		color.Red("[Event] Invalid revealed event")
		return
	}
	color.Cyan("[Event] %s player revealed. Role: %s", e.Player.Name, role.RoleToString(e.Player.Role))
}

func (EventPrinter) HandleGEE(e *event.GameEndEvent, ctx event.EventHandleContext) {
	if e.Team == event.GameEndEvent_Red {
		color.Cyan("[Event] Town won")
	} else {
		color.Cyan("[Event] Mafia won")
	}
	os.Exit(0)
}

func (EventPrinter) Preprocess(e *event.Event, ctx event.EventHandleContext) {
}

func NewEventPrinter() event.EventHandler {
	return EventPrinter{}
}
