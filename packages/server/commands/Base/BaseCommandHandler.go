package base

import "gymshark/server/entities"

type Command interface {
	CommandName() string
}

type Result interface {
	entities.Order | []entities.Pack
}

type CommandHandlerInterface[TCommand Command, TResult Result] interface {
	Handle(command TCommand) (TResult, error)
}
