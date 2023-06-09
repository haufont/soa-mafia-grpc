package session

import (
	"fmt"
	"log"
	"mafia-grpc/internal/event"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/util"
	"strings"
	"sync"

	"github.com/google/uuid"
)

const (
	EventChanSize = 64
)

type PlayerSession struct {
	Player      player.Player
	GameSession *GameSession
	EventChan   chan<- *event.Event
}

type PlayerSessionManager struct {
	m                sync.RWMutex
	players          map[uuid.UUID]*PlayerSession
	playersName      map[string]int
	playersNameCount map[string]int
}

func NewPlayerSessionManager() PlayerSessionManager {
	return PlayerSessionManager{
		players:          make(map[uuid.UUID]*PlayerSession),
		playersName:      make(map[string]int),
		playersNameCount: make(map[string]int),
	}
}

func (manager *PlayerSessionManager) getUniqueName(username string) string {
	manager.playersNameCount[username] += 1
	manager.playersName[username] += 1
	if manager.playersNameCount[username] > 1 {
		username += fmt.Sprintf("#%d", manager.playersName[username])
	}
	return username
}

func (manager *PlayerSessionManager) NewPlayerSession(username string) (*PlayerSession, <-chan *event.Event) {
	manager.m.Lock()
	defer manager.m.Unlock()
	username = manager.getUniqueName(username)
	uuid, suiid := util.NewUUID()
	eventChan := make(chan *event.Event, EventChanSize)
	session := &PlayerSession{
		Player: player.Player{
			Uuid: suiid,
			Name: username,
		},
		GameSession: nil,
		EventChan:   eventChan,
	}
	manager.players[uuid] = session
	return session, eventChan
}

func (manager *PlayerSessionManager) GetPlayerSession(id uuid.UUID) (*PlayerSession, error) {
	manager.m.RLock()
	defer manager.m.RUnlock()
	session, ok := manager.players[id]
	if !ok {
		return nil, fmt.Errorf("user not found (uuid \"%s\")", id)
	}
	return session, nil
}

func (manager *PlayerSessionManager) RemovePlayerSession(psession *PlayerSession) {
	manager.m.Lock()
	defer manager.m.Unlock()
	suuid := psession.Player.Uuid
	uuid := util.BytesToUUID(suuid)
	if _, ok := manager.players[uuid]; !ok {
		log.Printf("PlayerSessionManager: player has already left or has not connected (uuid: \"%s\")", uuid)
		return
	}
	delete(manager.players, uuid)
	baseName := strings.Split(psession.Player.Name, "#")[0]
	manager.playersNameCount[baseName] -= 1
	if manager.playersNameCount[baseName] == 0 {
		delete(manager.playersName, baseName)
		delete(manager.playersNameCount, baseName)
	}
}
