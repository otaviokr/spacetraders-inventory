package component

import (
	"bytes"
	"context"
	"fmt"

	"github.com/otaviokr/spacetraders-inventory/web"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
)

// User contains the essential information to authenticate in the game, but also to map the response from user details.
type User struct {
	token    string
	tracer   trace.Tracer
	webProxy web.Proxy
	Details  UserDetails `yaml:"user"`
	Error    Error       `yaml:"error"`
}

// UserDetails is the response from the User Detail API.
type UserDetails struct {
	Credits        int64  `yaml:"credits"`
	JoinedAt       string `yaml:"joinedAt"`
	ShipCount      int64  `yaml:"shipCount"`
	StructureCount int64  `yaml:"structureCount"`
	Username       string `yaml:"username"`
}

// NewUser creates a new instance of component.User.
func NewUser(ctx context.Context, tracer trace.Tracer, token string) (*User, error) {
	return NewUserCustomProxy(ctx, tracer, web.NewWebProxy(token), token)
}

// NewUserCustomProxy creates a new instance of component.User, using a provided custom web.WebProxy.
func NewUserCustomProxy(ctx context.Context, tracer trace.Tracer, proxy web.Proxy, token string) (*User, error) {
	userCtx, span := tracer.Start(ctx, "User login")
	defer span.End()
	user := User{
		token:    token,
		tracer:   tracer,
		webProxy: proxy}
	if err := user.GetDetails(userCtx); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}
	return &user, nil
}

// GetDetails will get the user details from the game.
func (u *User) GetDetails(ctx context.Context) error {
	_, span := u.tracer.Start(
		ctx,
		"Get User Details",
		trace.WithAttributes(
			attribute.Key("user.username").String(u.Details.Username)))
	defer span.End()

	data, err := u.webProxy.GetUserDetails()
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("data").String(string(data)))
		span.SetStatus(codes.Error, err.Error())
		return err
	} else {
		u.Error.Code = -1
		u.Error.Message = ""
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&u); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if len(u.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", u.Error.Code, u.Error.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// GetShips will get information about all ships the user owns.
func (u *User) GetShips(ctx context.Context) (*ShipList, error) {
	_, span := u.tracer.Start(
		ctx,
		"Get Ships Owned by User",
		trace.WithAttributes(
			attribute.Key("user.username").String(u.Details.Username)))
	defer span.End()

	data, err := u.webProxy.GetShipList()
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("data").String(string(data)))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	} else {
		u.Error.Code = -1
		u.Error.Message = ""
	}

	var shipList ShipList
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&shipList); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(u.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", u.Error.Code, u.Error.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &shipList, nil
}

// GetLeaderboard will get the leaderboard and the user rank.
func (u *User) GetLeaderboard(ctx context.Context) (*Leaderboard, error) {
	_, span := u.tracer.Start(
		ctx,
		"Get Leaderboard",
		trace.WithAttributes(
			attribute.Key("user.username").String(u.Details.Username)))
	defer span.End()

	data, err := u.webProxy.GetLeaderboard()
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("data").String(string(data)))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	} else {
		u.Error.Code = -1
		u.Error.Message = ""
	}

	var leaderboard Leaderboard
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&leaderboard); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(u.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", u.Error.Code, u.Error.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &leaderboard, nil
}

// GetStatus will get the game status (i.e., if it is running and available).
func (u *User) GetStatus(ctx context.Context) (int, error) {
	_, span := u.tracer.Start(
		ctx,
		"Get Game Status",
		trace.WithAttributes(
			attribute.Key("user.username").String(u.Details.Username)))
	defer span.End()

	data, err := u.webProxy.GetGameStatus()
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("data").String(string(data)))
		span.SetStatus(codes.Error, err.Error())
		return -1, err
	} else {
		u.Error.Code = -1
		u.Error.Message = ""
	}

	var status map[string]string
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&status); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return -1, err
	}

	if len(u.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", u.Error.Code, u.Error.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return -1, err
	}

	switch status["status"] {
	case "spacetraders is currently online and available to play":
		return 1, nil
	default:
		return -1, nil
	}
}
