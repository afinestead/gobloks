package database

import (
	"context"
	"fmt"
	"gobloks/internal/types"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const PostgresInitPath = "internal/database/init.sql"

type DatabaseManager struct {
	pool *pgxpool.Pool
}

func ConnectDB() (*DatabaseManager, error) {
	const DB_URL = "postgres://gobloks_user:password@localhost:5432/gobloks_db"
	pool, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	sqlInit, err := os.ReadFile(PostgresInitPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read init file: %v\n", err)
		return nil, err
	}

	_, err = pool.Exec(context.Background(), string(sqlInit))

	return &DatabaseManager{pool}, err
}

func (db *DatabaseManager) Close() {
	db.pool.Close()
}

func (db *DatabaseManager) AddGame(config *types.GameConfig) (types.GameID, error) {
	var gid types.GameID
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire connection: %v\n", err)
		return 0, err
	}

	defer conn.Release()

	err = conn.QueryRow(
		context.Background(),
		"INSERT INTO gobloks.game (created,last_active,game_status,player_count,public)"+
			"VALUES (NOW(),NOW(),0,0,$1)"+
			"RETURNING id",
		config.Public,
	).Scan(&gid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Add game failed: %v\n", err)
		return 0, err
	}

	fmt.Println("Created game #", gid)

	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO gobloks.game_config (id,players,block_degree,density,turns,time_seconds,time_bonus,hints)"+
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		gid, config.Players, config.BlockDegree, config.Density, config.TurnBased, config.TimeControl, config.TimeBonus, config.Hints,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert GameConfig failed: %v\n", err)
	}
	return gid, err
}

func (db *DatabaseManager) UpdateGamePlayers(gid types.GameID, players uint) error {
	_, err := db.pool.Exec(
		context.Background(),
		"UPDATE gobloks.game SET last_active=NOW(), players=$1 WHERE id=$2",
		players, gid,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Game player update failed: %v\n", err)
	}
	return err
}

func (db *DatabaseManager) UpdateGameStatus(gid types.GameID, status types.Flags) error {
	_, err := db.pool.Exec(
		context.Background(),
		"UPDATE gobloks.game SET last_active=NOW(), game_status=$1 WHERE id=$2",
		status, gid,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Game status update failed: %v\n", err)
	}
	return err
}

type GameWithConfig struct {
	ID           types.GameID `json:"gid"`
	Created      time.Time    `json:"created"`
	Last_Active  time.Time    `json:"active"`
	Player_Count uint         `json:"players"`
	Players      uint         `json:"maxPlayers"`
	Block_Degree uint8        `json:"degree"`
	Density      float64      `json:"density"`
	Turns        bool         `json:"turns"`
	Time_Seconds uint         `json:"timeSeconds"`
	Time_Bonus   uint         `json:"timeBonus"`
	Hints        uint         `json:"hints"`
}

func (db *DatabaseManager) GetWaitingGames(limit, offset uint64) ([]GameWithConfig, error) {
	rows, err := db.pool.Query(
		context.Background(),
		`SELECT g.id, g.created, g.last_active, g.player_count, gc.players, gc.block_degree, gc.density, gc.turns, gc.time_seconds, gc.time_bonus, gc.hints
         FROM gobloks.game g
         JOIN gobloks.game_config gc ON g.id = gc.id
         WHERE g.game_status = 0 AND g.public = true
         ORDER BY g.created DESC
         LIMIT $1 OFFSET $2`,
		limit, offset,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Get games failed: %v\n", err)
		return []GameWithConfig{}, err
	}

	games, err := pgx.CollectRows(rows, pgx.RowToStructByName[GameWithConfig])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect games failed: %v\n", err)
		return []GameWithConfig{}, err
	}
	return games, nil
}
