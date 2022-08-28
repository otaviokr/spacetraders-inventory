package web

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "spacetraderinventory"
)

var (
	Credits = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "credits",
			Help:      "How much credits user has",
		},
		[]string{"username"})

	GameStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "game_status",
			Help:      "Indicates if the game is up, running and available",
		},
		[]string{"username"})

	ShipCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "shipcount",
			Help:      "Total of ships user has",
		},
		[]string{"username"})

	StructureCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "structurecount",
			Help:      "Total of structure user has",
		},
		[]string{"username"})

	ShipLoad = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "shipload",
			Help:      "Unused space in ship cargo",
		},
		[]string{"username", "id", "class", "manufacturer", "type", "maxcargo", "plating", "speed", "weapons"})

	UserRank = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "userrank",
			Help:      "User rank in leaderboard",
		},
		[]string{"username"})
)
