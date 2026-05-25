package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	activeRoomsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "video_active_rooms",
		Help: "Number of active video call rooms.",
	})
	activeParticipantsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "video_active_participants",
		Help: "Total connected participants across all rooms.",
	})
	CallsTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "video_calls_total",
		Help: "Total video calls started.",
	})
)
