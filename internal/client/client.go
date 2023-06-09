package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mafia-grpc/api"
	"mafia-grpc/internal/event"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/state"

	"google.golang.org/grpc"
)

type Client struct {
	ctx           context.Context
	conn          *grpc.ClientConn
	grpcClient    api.MafiaServiceClient
	events        api.MafiaService_JoinClient
	eventHandlers []event.EventHandler
	Player        *player.Player
}

type ClientOptions struct {
	Addr          string
	GrpcOptions   []grpc.DialOption
	EventHandlers []event.EventHandler
}

func NewClient(username string, ctx context.Context, options ClientOptions) (*Client, error) {
	conn, err := grpc.DialContext(ctx, options.Addr, options.GrpcOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed init connection (addr: %s). Error: %s", options.Addr, err)
	}
	grpcClient := api.NewMafiaServiceClient(conn)

	events, err := grpcClient.Join(ctx, &api.ReqJoin{Username: username})
	if err != nil {
		return nil, fmt.Errorf("failed join. Error: %s", err)
	}

	client := &Client{
		ctx:           ctx,
		conn:          conn,
		grpcClient:    grpcClient,
		events:        events,
		eventHandlers: options.EventHandlers,
	}

	e, err := client.getEventSync()
	if err != nil {
		return nil, err
	}

	err = client.setPlayer(e)
	if err != nil {
		return nil, err
	}

	client.handleEvent(e)

	go client.ServeEvents()

	return client, nil
}

func (c *Client) getEventSync() (*event.Event, error) {
	return c.events.Recv()
}

func (c *Client) handleEvent(e *event.Event) {
	for i := range c.eventHandlers {
		event.Handle(c.eventHandlers[i], e, event.EventHandleContext{
			Player: c.Player,
		})
	}
}

func (c *Client) getAndHandleEventSync() (*event.Event, error) {
	e, err := c.getEventSync()
	if err == nil {
		c.handleEvent(e)
	}
	return e, err
}

func (c *Client) setPlayer(e *event.Event) error {
	joinEvent := e.GetJoinEvent()
	if joinEvent == nil {
		return errors.New("invalid first event type")
	}
	if joinEvent.Player == nil || joinEvent.Player.Uuid == nil {
		return errors.New("invalid first event type")
	}
	c.Player = joinEvent.Player
	return nil
}

func (c *Client) getUUID() []byte {
	return c.Player.Uuid
}

func (c *Client) ServeEvents() {
	for {
		_, err := c.getAndHandleEventSync()
		if err != nil {
			log.Fatalf("\nServer closed\n")
		}
	}
}

func (c *Client) GetState() (s *state.State, err error) {
	rsp, err := c.grpcClient.GetState(c.ctx, &api.ReqGetState{PlayerUuid: c.getUUID()})
	if err != nil {
		return
	}
	return rsp.State, nil
}

func (c *Client) Kill(playerName string) (err error) {
	_, err = c.grpcClient.Kill(c.ctx, &api.ReqKill{
		RequesterUuid:    c.getUUID(),
		TargetPlayerName: playerName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Check(playerName string) (isMafia bool, err error) {
	rsp, err := c.grpcClient.Check(c.ctx, &api.ReqCheck{
		RequesterUuid:    c.getUUID(),
		TargetPlayerName: playerName,
	})
	if err != nil {
		return false, err
	}
	return rsp.IsMafia, nil
}

func (c *Client) Publish(playerName string) (err error) {
	_, err = c.grpcClient.Publish(c.ctx, &api.ReqPublish{
		RequesterUuid:    c.getUUID(),
		TargetPlayerName: playerName,
	})
	if err != nil {
		return err
	}
	return nil
}
