package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func NewHTTPServer(addr string) *http.Server {
	mux := http.NewServeMux()
	s := newHTTPServer()
	mux.HandleFunc("POST /records", s.handleProduce)
	mux.HandleFunc("GET /records/{offset}", s.handleConsume)
	return &http.Server{Addr: addr, Handler: mux}
}

// Private
type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{Log: NewLog()}
}

// Interface
type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

// Handlers
func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.ParseUint(r.PathValue("offset"), 10, 64)
	if err != nil {
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}

	record, err := s.Log.Read(offset)
	if errors.Is(err, ErrOffsetNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
