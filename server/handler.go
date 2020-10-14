package server

type Handler func(conn *ContextHandler, msg interface{})
