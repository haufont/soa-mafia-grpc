package server

import (
	"context"
	"errors"
	"mafia-grpc/api"
	"mafia-grpc/internal/event"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/session"
	"mafia-grpc/internal/util"
)

type ServerConfig struct {
	GameSessionConfig session.GameSessionConfig
}

type Server struct {
	api.UnimplementedMafiaServiceServer

	config               ServerConfig
	playerSessionManager session.PlayerSessionManager
	gameSessionManager   session.GameSessionManager
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		config:               config,
		playerSessionManager: session.NewPlayerSessionManager(),
		gameSessionManager:   session.NewGameSessionManager(config.GameSessionConfig),
	}
}

func (s *Server) addPlayer(username string) (psession *session.PlayerSession, eventChan <-chan *event.Event, ctx context.Context) {
	psession, eventChan = s.playerSessionManager.NewPlayerSession(username)
	s.gameSessionManager.AddPlayer(psession)
	ctx = psession.GameSession.Ctx
	return
}

func (s *Server) removePlayer(psession *session.PlayerSession) {
	s.playerSessionManager.RemovePlayerSession(psession)
	gsession := psession.GameSession
	gsession.RemovePlayer(psession)
}

func (s *Server) Join(req *api.ReqJoin, outs api.MafiaService_JoinServer) (err error) {
	username := req.Username
	err = player.ValidatePlayerName(username)
	if err != nil {
		return err
	}
	psession, eventChan, ctx := s.addPlayer(username)
	defer s.removePlayer(psession)
	s.ServePlayer(ctx, eventChan, outs)
	return nil
}

func (s *Server) ServePlayer(ctx context.Context, eventChan <-chan *event.Event, outs api.MafiaService_JoinServer) {
	for {
		var event *event.Event
		select {
		case <-ctx.Done():
			return
		case <-outs.Context().Done():
			return
		case event = <-eventChan:
		}
		if event == nil {
			return
		}
		err := outs.Send(event)
		if err != nil {
			return
		}
	}
}

func (s *Server) GetState(ctx context.Context, req *api.ReqGetState) (rsp *api.RspGetState, err error) {
	uuid := util.BytesToUUID(req.PlayerUuid)
	psession, err := s.playerSessionManager.GetPlayerSession(uuid)
	if err != nil {
		return
	}
	if psession.GameSession == nil {
		return nil, errors.New("no game session found for this player")
	}
	rsp = &api.RspGetState{
		State: psession.GameSession.GetState(psession),
	}
	return
}

func (s *Server) Kill(ctx context.Context, req *api.ReqKill) (rsp *api.RspKill, err error) {
	uuid := util.BytesToUUID(req.RequesterUuid)
	psession, err := s.playerSessionManager.GetPlayerSession(uuid)
	if err != nil {
		return
	}
	if psession.GameSession == nil {
		return nil, errors.New("no game session found for this player")
	}
	err = psession.GameSession.Kill(psession, req.TargetPlayerName)
	return &api.RspKill{}, err
}

func (s *Server) Check(ctx context.Context, req *api.ReqCheck) (rsp *api.RspCheck, err error) {
	uuid := util.BytesToUUID(req.RequesterUuid)
	psession, err := s.playerSessionManager.GetPlayerSession(uuid)
	if err != nil {
		return
	}
	if psession.GameSession == nil {
		return nil, errors.New("no game session found for this player")
	}
	isMafia, err := psession.GameSession.Check(psession, req.TargetPlayerName)
	return &api.RspCheck{IsMafia: isMafia}, err
}

func (s *Server) Publish(ctx context.Context, req *api.ReqPublish) (rsp *api.RspPublish, err error) {
	uuid := util.BytesToUUID(req.RequesterUuid)
	psession, err := s.playerSessionManager.GetPlayerSession(uuid)
	if err != nil {
		return
	}
	if psession.GameSession == nil {
		return nil, errors.New("no game session found for this player")
	}
	err = psession.GameSession.Publish(psession, req.TargetPlayerName)
	return &api.RspPublish{}, err
}
