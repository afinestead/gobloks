package types

// Some type definitions
// TODO: Put these in a real spot?
type Direction int
type Axis int
type Owner uint32
type PlayerID uint16
type GameID string

type GameConfig struct {
	Players     uint    `json:"players"`
	BlockDegree uint8   `json:"degree"`
	Density     float64 `json:"density"`
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// from enum import auto, Enum, IntEnum
// from pydantic import BaseModel, constr
// from typing import List, Optional, Set

// class Message(BaseModel):
//     message: str = ""

// class GameID(BaseModel):
//     game_id: constr(
//         min_length=4,
//         max_length=4,
//         to_upper=True,
//         strip_whitespace=True
//     )

// class AccessToken(BaseModel):
//     access_token: str

// class Coordinate(BaseModel):
//     def __hash__(self):
//         return hash((self.x, self.y))

//     def __lt__(self, other):
//         return (self.x < other.x) or (self.y < other.y)

//     def __gt__(self, other):
//         return (self.x > other.x) or (self.y > other.y)

//     def __eq__(self, other):
//         return self.x == other.x and self.y == other.y

//     x: int
//     y: int

// class Piece(BaseModel):
//     shape: Set[Coordinate]

// class EndGameOn(IntEnum):
//     FIRST_PLAYER_OUT = 0
//     LAST_PLAYER_OUT = 1

// class GameStatus(str, Enum):
//     WAITING = "waiting"
//     ACTIVE = "active"
//     DONE = "done"

// class GameState(BaseModel):
//     status: str # TODO: Why won't GameStatus enum work?
//     turn: Optional[int]

// class GameConfig(BaseModel):
//     players: int
//     board_size: int
//     block_size: int
//     # turn_based: bool
//     # turn_timeout: Optional[float]
//     # end_condition: int

// class PlayerProfile(BaseModel):
//     player_id: int
//     color: int
//     name: str
//     pieces: Optional[List[Piece]] = None
