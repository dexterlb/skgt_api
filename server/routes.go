package server

import "github.com/julienschmidt/httprouter"

func (s *Server) transports(params httprouter.Params) (interface{}, error) {
	transports, err := s.backend.Transports()

	return transports, err
}
