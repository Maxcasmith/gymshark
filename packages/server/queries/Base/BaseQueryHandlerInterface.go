package base

import "gymshark/server/entities"

type Query interface {
	QueryName() string
}

type Result interface {
	entities.Order | []entities.Pack
}

type QueryHandlerInterface[TQuery Query, TResult Result] interface {
	Handle(query TQuery) (TResult, error)
}
