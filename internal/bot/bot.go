package bot

import (
	"context"
	"log"
	"mafia-grpc/internal/client"
	"mafia-grpc/internal/event"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/role"
	"mafia-grpc/internal/state"
	"mafia-grpc/internal/util"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Bot struct {
	m        sync.RWMutex
	state    *state.State
	maxRetry uint
	lag      time.Duration
	client   *client.Client
	ctx      context.Context
	waiter   util.Waiter
}

func NewBot(lag time.Duration) Bot {
	ctx, waiter := util.NewStopCtx()
	return Bot{
		maxRetry: 5,
		lag:      lag,
		ctx:      ctx,
		waiter:   waiter,
	}
}

func (b *Bot) StartClient(options client.ClientOptions) error {
	options.EventHandlers = append(options.EventHandlers, b)
	client, err := client.NewClient("bot", b.ctx, options)
	if err != nil {
		return err
	}
	b.client = client
	return nil
}

func (b *Bot) Serve() {
	b.waiter.Wait()
}

// Избыточно, но так проще
func (b *Bot) Preprocess(*event.Event, event.EventHandleContext) {
	if b.client == nil {
		return
	}
	newState, err := b.client.GetState()
	if err != nil {
		return
	}
	b.m.Lock()
	defer b.m.Unlock()
	b.state = newState
}

func (b *Bot) getState() *state.State {
	b.m.RLock()
	defer b.m.RUnlock()
	copyState := *b.state
	return &copyState
}

func (b *Bot) randLag() {
	lag := b.lag + (b.lag*time.Duration(rand.Int31n(100)))/100
	time.Sleep(lag)
}

func (b *Bot) randKill(s *state.State) (err error) {
	livePlayersNumber := 0
	for _, player := range s.Players {
		if !player.Dead {
			livePlayersNumber++
		}
	}
	k := rand.Uint32() % (uint32(livePlayersNumber) + 1)
	for _, player := range s.Players {
		if !player.Dead && k == 0 {
			err = b.client.Kill(player.Name)
			if err == nil {
				color.Yellow("[BotAction] Voted for %s", player.Name)
			}
			return
		}
		k--
	}
	err = b.client.Kill("")
	if err == nil {
		color.Yellow("[BotAction] Abstained from voting")
	}
	return
}

func (b *Bot) randCheck(s *state.State) (err error) {
	livePlayersNumber := 0
	for _, player := range s.Players {
		if !player.Dead {
			livePlayersNumber++
		}
	}
	k := rand.Uint32() % (uint32(livePlayersNumber) + 1)
	for _, player := range s.Players {
		if !player.Dead && k == 0 {
			var isMafia bool
			isMafia, err = b.client.Check(player.Name)
			if err == nil {
				color.Yellow("[BotAction] Checked %s, result: %t", player.Name, isMafia)
			}
			return
		}
		k--
	}
	_, err = b.client.Check("")
	if err == nil {
		color.Yellow("[BotAction] Abstained from checking")
	}
	return
}

func (b *Bot) publishAll(s *state.State) {
	for _, player := range s.Players {
		if !player.Dead {
			_ = b.client.Publish(player.Name)
			time.Sleep(time.Microsecond * 100)
		}
	}
}

func (b *Bot) makeADayMove(p *player.Player, s *state.State) error {
	err := b.randKill(s)
	if err != nil {
		return err
	}
	if p.Role == role.Role_COMMISAR {
		b.publishAll(s)
	}
	return nil
}

func (b *Bot) makeANightMove(p *player.Player, s *state.State) (err error) {
	switch p.Role {
	case role.Role_MAFIA:
		err = b.randKill(s)
	case role.Role_COMMISAR:
		err = b.randCheck(s)
	}
	return
}

func (b *Bot) checkAlive(p *player.Player, s *state.State) bool {
	for _, player := range s.Players {
		if string(p.Uuid) == string(player.Uuid) {
			return !player.Dead
		}
	}
	return false
}

func (b *Bot) makeAMove(p *player.Player) {
	b.randLag()
	for i := uint(0); i < b.maxRetry; i++ {
		s := b.getState()
		if !b.checkAlive(p, s) {
			return
		}
		var err error
		switch s.PartOfTheDay {
		case state.PartOfTheDay_DAY:
			err = b.makeADayMove(p, s)
		case state.PartOfTheDay_NIGHT:
			err = b.makeANightMove(p, s)
		case state.PartOfTheDay_UNKNOWN:
			log.Fatalf("Unknown part of the day")
		}
		if err == nil {
			return
		}
		if i+1 == b.maxRetry {
			log.Fatalf("%s", err)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (b *Bot) HandleJE(*event.JoinEvent, event.EventHandleContext) {}
func (b *Bot) HandleLE(*event.LeftEvent, event.EventHandleContext) {}

func (b *Bot) HandleSDE(e *event.StartDayEvent, ctx event.EventHandleContext) {
	b.makeAMove(ctx.Player)
}

func (b *Bot) HandleSNE(e *event.StartNightEvent, ctx event.EventHandleContext) {
	b.makeAMove(ctx.Player)
}

func (b *Bot) HandleVE(*event.VoteEvent, event.EventHandleContext) {}

func (b *Bot) HandleRVE(e *event.RepeatVotingEvent, ctx event.EventHandleContext) {
	b.makeAMove(ctx.Player)
}

func (b *Bot) HandleRE(*event.RevealedEvent, event.EventHandleContext) {}
func (b *Bot) HandleGEE(*event.GameEndEvent, event.EventHandleContext) {}
