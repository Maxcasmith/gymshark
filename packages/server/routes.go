package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	createordercommand "gymshark/server/commands/CreateOrderCommand"
	replacepackscommand "gymshark/server/commands/ReplacePacksCommand"
	"gymshark/server/entities"
	orderaggregator "gymshark/server/helpers/OrderAggregator"
	getorderquery "gymshark/server/queries/GetOrderQuery"
	getpackquery "gymshark/server/queries/GetPackQuery"
	orderrepo "gymshark/server/repositories/Order"
	packrepo "gymshark/server/repositories/Pack"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response interface {
	[]entities.Pack | entities.Order | ErrorResponse
}

func routes(dbDir string) *http.ServeMux {
	packs := packrepo.New(dbDir)
	orders := orderrepo.New(dbDir)

	getPacks := &getpackquery.GetPackQueryHandler{Packs: packs}
	getOrder := &getorderquery.GetOrderQueryHandler{Orders: orders}
	createOrder := &createordercommand.CreateOrderCommandHandler{Packs: packs, Orders: orders}
	replacePacks := &replacepackscommand.ReplacePacksCommandHandler{Packs: packs}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/packs", func(writer http.ResponseWriter, request *http.Request) {
		packList, err := getPacks.Handle(getpackquery.GetPackQuery{})
		respond(writer, http.StatusOK, packList, err)
	})

	mux.HandleFunc("PUT /api/packs", func(writer http.ResponseWriter, request *http.Request) {
		var command replacepackscommand.ReplacePacksCommand
		if err := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<10)).Decode(&command); err != nil {
			writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
			return
		}
		packList, err := replacePacks.Handle(command)
		respond(writer, http.StatusOK, packList, err)
	})

	mux.HandleFunc("POST /api/packs/reset", func(writer http.ResponseWriter, request *http.Request) {
		packList, err := replacePacks.Handle(replacepackscommand.ReplacePacksCommand{Sizes: replacepackscommand.DefaultPackSizes})
		respond(writer, http.StatusOK, packList, err)
	})

	mux.HandleFunc("POST /api/orders", func(writer http.ResponseWriter, request *http.Request) {
		var command createordercommand.CreateOrderCommand
		if err := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 1<<10)).Decode(&command); err != nil {
			writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
			return
		}
		order, err := createOrder.Handle(command)
		respond(writer, http.StatusCreated, order, err)
	})

	mux.HandleFunc("GET /api/orders/{id}", func(writer http.ResponseWriter, request *http.Request) {
		orderID, err := strconv.Atoi(request.PathValue("id"))
		if err != nil {
			writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "id must be a number"})
			return
		}
		order, err := getOrder.Handle(getorderquery.GetOrderQuery{OrderID: orderID})
		respond(writer, http.StatusOK, order, err)
	})

	return mux
}

func respond[TBody Response](writer http.ResponseWriter, status int, body TBody, err error) {
	switch {
	case errors.Is(err, orderaggregator.ErrInvalidOrder), errors.Is(err, replacepackscommand.ErrInvalidPacks):
		writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, getorderquery.ErrNotFound):
		writeJSON(writer, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case err != nil:
		log.Print(err)
		writeJSON(writer, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	default:
		writeJSON(writer, status, body)
	}
}

func writeJSON[TBody Response](writer http.ResponseWriter, status int, body TBody) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(body); err != nil {
		log.Print(err)
	}
}
