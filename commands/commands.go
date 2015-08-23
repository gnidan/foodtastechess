package commands

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/queries"
	"foodtastechess/users"
)

type Commands interface {
	ExecCommand(
		name string, userId users.Id, params map[string]interface{},
	) (bool, string)

	events() events.Events
	queries() queries.ClientQueries
}

type CommandsService struct {
	Queries queries.ClientQueries `inject:"clientQueries"`
	Events  events.Events         `inject:"events"`
}

func New() Commands {
	return new(CommandsService)
}

func (s *CommandsService) ExecCommand(name string, userId users.Id, params map[string]interface{}) (bool, string) {
	var (
		ctx context
		cmd command
		ok  bool
		msg string
	)
	ctx, ok, msg = makeContext(name, userId, params)
	if !ok {
		return false, msg
	}

	cmd, ok = commandMap[ctx.name]
	if !ok {
		return false, "Unknown Command"
	}

	for _, validator := range cmd.validators {
		ok, msg = validator(ctx, s)
		if !ok {
			return false, fmt.Sprintf(
				"Command Invalid: %s",
				msg,
			)
		}
	}

	for _, event := range cmd.gen(ctx, s) {
		err := s.Events.Receive(event)
		if err != nil {
			return false, fmt.Sprintf("Event Generation Error: %v", err)
		}
	}

	return true, ""
}

func makeContext(name string, userId users.Id, params map[string]interface{}) (context, bool, string) {
	ctx := new(context)

	ctx.name = name
	ctx.userId = userId

	if iface, ok := params["game_id"]; ok {
		ctx.gameId, ok = iface.(game.Id)
		if !ok {
			return *ctx, false, "Invalid Game Id"
		}
	}

	if iface, ok := params["move"]; ok {
		ctx.move, ok = iface.(game.AlgebraicMove)
		if !ok {
			return *ctx, false, "Invalid Move"
		}
	}

	if iface, ok := params["color"]; ok {
		ctx.colorChoice, ok = iface.(game.Color)
		if !ok {
			return *ctx, false, fmt.Sprintf(
				"Invalid color. Must be %v or %v",
				game.White, game.Black,
			)
		}
	}

	return *ctx, true, ""
}

func (s *CommandsService) events() events.Events          { return s.Events }
func (s *CommandsService) queries() queries.ClientQueries { return s.Queries }
