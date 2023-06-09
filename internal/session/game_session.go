package session

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mafia-grpc/internal/event"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/role"
	"mafia-grpc/internal/state"
	"mafia-grpc/internal/util"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameSessionConfig struct {
	PlayersInSession uint
	OpenVoting       bool
}

type GameSession struct {
	m               sync.RWMutex
	config          GameSessionConfig
	id              uuid.UUID
	players         map[uuid.UUID]*PlayerSession
	roleAlloctor    role.RoleAllocator
	freeSlots       uint
	sender          *Sender
	state           state.State
	started         bool
	ended           bool
	checked         bool
	numberOfAttemts int32
	Ctx             context.Context
	cancel          context.CancelFunc
}

func NewGameSession(config GameSessionConfig) *GameSession {
	ctx, cancel := context.WithCancel(context.Background())
	gsession := &GameSession{
		id:           uuid.New(),
		config:       config,
		players:      make(map[uuid.UUID]*PlayerSession),
		roleAlloctor: role.NewRoleAllocator(config.PlayersInSession),
		freeSlots:    config.PlayersInSession,
		state: state.State{
			Voices: make(map[string]string),
		},
		started:         false,
		ended:           false,
		checked:         false,
		numberOfAttemts: 1,
		Ctx:             ctx,
		cancel:          cancel,
	}
	gsession.sender = &Sender{gsession: gsession}
	return gsession
}

func (gsession *GameSession) sendEvent(e *event.Event, ctx event.EventHandleContext) {
	event.Handle(gsession.sender, e, ctx)
}

func (gsession *GameSession) blockingSendEvent(e *event.Event, ctx event.EventHandleContext) {
	gsession.m.RLock()
	defer gsession.m.RUnlock()
	gsession.sendEvent(e, ctx)
}

func (gsession *GameSession) sendDayEvent(player *player.Player) {
	gsession.sendEvent(event.StartDayEventToGeneralEvent(
		&event.StartDayEvent{KilledPlayer: player},
	), event.EventHandleContext{})
}

func (gsession *GameSession) sendNightEvent(player *player.Player) {
	gsession.sendEvent(event.StartNightEventToGeneralEvent(
		&event.StartNightEvent{KilledPlayer: player},
	), event.EventHandleContext{})
}

func (gsession *GameSession) startGame() {
	gsession.started = true
	gsession.state.PartOfTheDay = state.PartOfTheDay_DAY
}

func (gsession *GameSession) endGame(team event.GameEndEvent_Team) {
	gsession.sendEvent(event.GameEndEventToGeneralEvent(
		&event.GameEndEvent{Team: team},
	), event.EventHandleContext{})
	gsession.ended = true
	go func() {
		time.Sleep(5 * time.Second)
		gsession.cancel()
	}()
}

func (gsession *GameSession) countingOfVotes() string {
	candidates := make(map[string]int)
	for _, vote := range gsession.state.Voices {
		if vote != "" {
			candidates[vote] += 1
		}
	}

	maxVotes := 0
	mcount := 0
	for _, votes := range candidates {
		if maxVotes == votes {
			mcount += 1
		}
		if maxVotes < votes {
			maxVotes = votes
			mcount = 1
		}
	}

	defer func() {
		for k := range gsession.state.Voices {
			delete(gsession.state.Voices, k)
		}
	}()

	if mcount != 1 {
		return ""
	}
	for candidates, votes := range candidates {
		if maxVotes == votes {
			return candidates
		}
	}

	log.Fatalf("unreachable code")
	return ""
}

func (gsession *GameSession) checkEndGame() bool {
	mafia := 0
	town := 0
	for _, psession := range gsession.players {
		if psession.Player.Dead {
			continue
		}
		if psession.Player.Role == role.Role_MAFIA {
			mafia++
		} else {
			town++
		}
	}
	if mafia >= town {
		gsession.endGame(event.GameEndEvent_Black)
		return true
	} else if mafia == 0 {
		gsession.endGame(event.GameEndEvent_Red)
		return true
	}
	return false
}

func (gsession *GameSession) checkEndStage() {
	gsession.m.Lock()
	defer gsession.m.Unlock()

	if gsession.checkEndGame() {
		return
	}

	day := gsession.state.PartOfTheDay == state.PartOfTheDay_DAY
	night := gsession.state.PartOfTheDay == state.PartOfTheDay_NIGHT
	votes := len(gsession.state.Voices)
	needVotes := 0
	needCheck := false
	for _, psession := range gsession.players {
		if !psession.Player.Dead {
			if day || psession.Player.Role == role.Role_MAFIA {
				needVotes++
			}
			if night && psession.Player.Role == role.Role_COMMISAR {
				needCheck = true
			}
		}
	}
	if votes != needVotes || (needCheck && !gsession.checked) {
		return
	}

	candidate := gsession.countingOfVotes()

	var killedPlayer *player.Player
	if candidate != "" {
		for _, psession := range gsession.players {
			if psession.Player.Name == candidate {
				psession.Player.Dead = true
				killedPlayer = &psession.Player
				break
			}
		}
		if gsession.checkEndGame() {
			return
		}
	} else if gsession.numberOfAttemts < event.MaxVotingNumberAttemts {
		gsession.numberOfAttemts += 1
		go func() {
			stateCopy := gsession.state
			gsession.blockingSendEvent(
				event.RepeatVotingEventToGeneralEvent(
					&event.RepeatVotingEvent{Attempt: gsession.numberOfAttemts},
				),
				event.EventHandleContext{
					State: &stateCopy,
				},
			)
		}()
		return
	}

	gsession.numberOfAttemts = 1

	if day {
		gsession.state.PartOfTheDay = state.PartOfTheDay_NIGHT
		gsession.sendNightEvent(killedPlayer)
	} else {
		gsession.checked = false
		gsession.state.PartOfTheDay = state.PartOfTheDay_DAY
		gsession.sendDayEvent(killedPlayer)
	}
}

func (gsession *GameSession) inProgress() error {
	if !gsession.started {
		return errors.New("the game hasn't started yet")
	}
	if gsession.ended {
		return errors.New("the game is already over")
	}
	return nil
}

func (gsession *GameSession) findPlayer(playerName string) *player.Player {
	for _, psession := range gsession.players {
		if psession.Player.Name == playerName {
			return &psession.Player
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////

func (gsession *GameSession) AddPlayer(psession *PlayerSession) (full bool) {
	gsession.m.Lock()
	defer gsession.m.Unlock()
	full = (gsession.freeSlots == 0)
	if full {
		log.Fatal("game session is already full")
		return
	}
	role, err := gsession.roleAlloctor.Allocate()
	if err != nil {
		log.Fatal(err)
	}
	gsession.freeSlots -= 1
	full = (gsession.freeSlots == 0)
	if full {
		gsession.startGame()
	}
	psession.Player.Role = role
	psession.GameSession = gsession
	buuid := psession.Player.Uuid
	uuid := util.BytesToUUID(buuid)
	gsession.players[uuid] = psession
	playerCopy := psession.Player
	go func() {
		gsession.blockingSendEvent(event.JoinEventToGeneralEvent(&event.JoinEvent{
			Player: &playerCopy,
		}), event.EventHandleContext{})
		if full {
			gsession.sendDayEvent(nil)
		}
	}()
	return
}

func (gsession *GameSession) RemovePlayer(psession *PlayerSession) {
	gsession.m.Lock()
	defer gsession.m.Unlock()
	suuid := psession.Player.Uuid
	uuid := util.BytesToUUID(suuid)
	if _, ok := gsession.players[uuid]; !ok {
		log.Printf("GameSession: player has already left or has not connected (uuid: \"%s\")", uuid)
	}
	delete(gsession.players, uuid)
	delete(gsession.state.Voices, psession.Player.Name)
	if !gsession.started {
		gsession.roleAlloctor.Rollback(psession.Player.Role)
		gsession.freeSlots += 1
	}
	started := gsession.started
	go func() {
		gsession.blockingSendEvent(event.LeftEventToGeneralEvent(&event.LeftEvent{
			Player: &psession.Player,
		}), event.EventHandleContext{})
		if started {
			gsession.checkEndStage()
		}
	}()
}

func (gsession *GameSession) GetState(requester *PlayerSession) *state.State {
	gsession.m.RLock()
	defer gsession.m.RUnlock()
	uuid := util.BytesToUUID(requester.Player.Uuid)
	result := gsession.state
	for puuid, psession := range gsession.players {
		hideRole := !(uuid == puuid || (requester.Player.Role == role.Role_MAFIA && psession.Player.Role == role.Role_MAFIA))
		result.Players = append(
			result.Players,
			player.CopyPlayer(&psession.Player, (uuid != puuid), hideRole),
		)
	}
	return &result
}

func (gsession *GameSession) Kill(requester *PlayerSession, playerName string) error {
	gsession.m.Lock()
	defer gsession.m.Unlock()

	var vote *player.Player

	err := gsession.inProgress()
	if err != nil {
		return err
	}

	if requester.Player.Dead {
		return errors.New("you're dead and can't vote")
	}

	night := gsession.state.PartOfTheDay == state.PartOfTheDay_NIGHT
	if night && requester.Player.Role != role.Role_MAFIA {
		return errors.New("you can't vote now")
	}

	if _, ok := gsession.state.Voices[requester.Player.Name]; ok {
		return errors.New("you have already voted in this stage")
	}

	if playerName != "" {
		vote = gsession.findPlayer(playerName)
		if vote == nil {
			return fmt.Errorf("%s player not found or already out", playerName)
		}
		if vote.Dead {
			return fmt.Errorf("%s player is already dead", playerName)
		}
	}
	gsession.state.Voices[requester.Player.Name] = playerName

	if !gsession.config.OpenVoting {
		go func() {
			gsession.checkEndStage()
		}()
		return nil
	}
	if vote != nil {
		voteCopy := *vote
		vote = &voteCopy
	}
	requesterPlayer := requester.Player
	stateCopy := gsession.state
	go func() {
		gsession.blockingSendEvent(
			event.VoteEventToGeneralEvent(
				&event.VoteEvent{Requester: &requesterPlayer, Vote: vote},
			),
			event.EventHandleContext{
				State: &stateCopy,
			},
		)
		gsession.checkEndStage()
	}()
	return nil
}

func (gsession *GameSession) Check(requester *PlayerSession, playerName string) (bool, error) {
	gsession.m.Lock()
	defer gsession.m.Unlock()

	err := gsession.inProgress()
	if err != nil {
		return false, err
	}

	if requester.Player.Dead {
		return false, errors.New("you're dead and can't check")
	}

	night := gsession.state.PartOfTheDay == state.PartOfTheDay_NIGHT
	if !night || requester.Player.Role != role.Role_COMMISAR {
		return false, errors.New("you can't check now (you're not a commissar or it's not night)")
	}

	if gsession.checked {
		return false, errors.New("you have already checked someone at this stage")
	}

	if playerName == "" {
		gsession.checked = true
		go gsession.checkEndStage()
		return false, nil
	}

	player := gsession.findPlayer(playerName)
	if player == nil {
		return false, fmt.Errorf("%s player not found or already out", playerName)
	}

	isMafia := player.Role == role.Role_MAFIA

	player.Checked = isMafia
	gsession.checked = true
	go gsession.checkEndStage()
	return isMafia, nil
}

func (gsession *GameSession) Publish(requester *PlayerSession, playerName string) error {
	gsession.m.Lock()
	defer gsession.m.Unlock()

	err := gsession.inProgress()
	if err != nil {
		return err
	}

	if requester.Player.Dead {
		return errors.New("you're dead and can't publish")
	}

	day := gsession.state.PartOfTheDay == state.PartOfTheDay_DAY
	if !day || requester.Player.Role != role.Role_COMMISAR {
		return errors.New("you can't publish now (you're not a commissar or it's not day)")
	}

	player := gsession.findPlayer(playerName)
	if player == nil {
		return fmt.Errorf("%s player not found or already out", playerName)
	}

	if !player.Checked {
		return fmt.Errorf("%s player has not been checked or is not a mafia", playerName)
	}

	if player.Revealed {
		return fmt.Errorf("the role of the %s player has already been revealed", playerName)
	}

	player.Revealed = true

	playerCopy := *player
	go func() {
		gsession.blockingSendEvent(event.RevealedEventToGeneralEvent(
			&event.RevealedEvent{Player: &playerCopy},
		), event.EventHandleContext{})
	}()
	return nil
}

////////////////////////////////////////////////////////////////////////////////////

type GameSessionManager struct {
	m                 sync.RWMutex
	lastGameSession   *GameSession
	gameSessionConfig GameSessionConfig
}

func NewGameSessionManager(gameSessionConfig GameSessionConfig) GameSessionManager {
	return GameSessionManager{
		lastGameSession:   nil,
		gameSessionConfig: gameSessionConfig,
	}
}

func (manager *GameSessionManager) AddPlayer(player *PlayerSession) {
	manager.m.Lock()
	defer manager.m.Unlock()
	if manager.lastGameSession == nil {
		manager.lastGameSession = NewGameSession(manager.gameSessionConfig)
	}
	full := manager.lastGameSession.AddPlayer(player)
	player.GameSession = manager.lastGameSession
	if full {
		manager.lastGameSession = NewGameSession(manager.gameSessionConfig)
	}
}
