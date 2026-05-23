package main

import "database/sql"

type app struct {
	cfg      config
	store    userStore
	producer *eventProducer
	metrics  metrics
}

func newApp(cfg config, db *sql.DB) *app {
	return &app{
		cfg:      cfg,
		store:    userStore{db: db},
		producer: newEventProducer(cfg),
	}
}
