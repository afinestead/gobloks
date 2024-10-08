package game

import (
	"errors"
	"fmt"
	"gobloks/internal/sockets"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

/* TODO: Namespace these flags somehow */
// Game status flags (bitset)
const (
	FULL        types.Flags = (1 << 0)
	IN_PROGRESS types.Flags = (1 << 1)
	COMPLETE    types.Flags = (1 << 2)
)

type GameState struct {
	board  *Board
	turn   types.PlayerID
	status types.Flags
}

func (gs *GameState) Copy() *GameState {
	return &GameState{
		board:  gs.board.Copy(),
		turn:   gs.turn,
		status: gs.status,
	}

}

type Game struct {
	gid            types.GameID
	lock           *sync.Mutex
	config         types.GameConfig
	startingPieces PieceSet
	socketManager  *sockets.SocketManager
	lastActive     time.Time
	evalEngine     *EvalEngine
	state          *GameState
	players        map[types.PlayerID]*Player
}

func InitGame(gid types.GameID, config types.GameConfig) *Game {

	pieces, setPixels, err := GeneratePieceSet(config.BlockDegree) // TODO: cache
	if err != nil {
		return nil
	}

	// playerStates := make(map[types.PlayerID]*PlayerState, config.Players)
	// pids := make([]types.PlayerID, config.Players)

	// for ii := 1; ii <= len(pids); ii++ {
	// 	pid := types.PlayerID(ii)
	// 	pids[ii-1] = pid
	// 	playerStates[pid] = nil
	// }
	players := make(map[types.PlayerID]*Player, config.Players)
	pids := make([]types.PlayerID, config.Players)

	for ii := 1; ii <= len(pids); ii++ {
		pid := types.PlayerID(ii)
		pids[ii-1] = pid
		players[pid] = nil
	}

	board, err := NewBoard(pids, setPixels, config.Density)
	if err != nil {
		return nil
	}

	engine := InitEvalEngine(1)
	go engine.Start()

	go func() {
		for eval := range engine.chResult {
			fmt.Println("Eval result", eval)
		}
	}()

	return &Game{
		gid:            gid,
		lock:           &sync.Mutex{},
		config:         config,
		startingPieces: pieces,
		socketManager:  sockets.InitSocketManager(len(pids)),
		lastActive:     time.Now(),
		evalEngine:     engine,
		state: &GameState{
			board,
			pids[0],
			0,
		},
		players: players,
	}
}

func (g *Game) GetPlayers() map[types.PlayerID]*Player {
	return g.players
}

func (g *Game) IsStale() bool {
	var cleanupAfter time.Duration
	if g.config.TimeControl == 0 {
		// Cleanup untimed games after a week of inactivity
		cleanupAfter = 7 * 24 * time.Hour
	} else {
		cleanupAfter = time.Duration(5*g.config.Players*g.config.TimeControl*1000) * time.Millisecond
	}
	fmt.Println("Cleanup after", cleanupAfter, "currently", time.Since(g.lastActive))
	return time.Since(g.lastActive) > cleanupAfter
}

func (g *Game) nextTurn() {
	for _, player := range g.players {
		if player != nil && !player.state.status.Has(DISABLED) && player.possiblePlacements.Next == nil {
			player.state.status.Set(DISABLED)
		}
	}

	var nextUp types.PlayerID = PID_NONE
	for i := 0; i < len(g.players); i++ {
		maybeNext := types.PlayerID((int(g.state.turn)+i)%len(g.players)) + 1
		if g.players[maybeNext] != nil && !g.players[maybeNext].state.status.Has(DISABLED) {
			nextUp = maybeNext
			fmt.Println("Next up:", nextUp)
			break
		}
	}

	g.state.turn = nextUp
}

func (g *Game) updateValidPlacements(player *Player, placement utilities.Set[types.Point]) {
	var wg sync.WaitGroup

	updatePlayerPlacements := func(p *Player) {
		defer wg.Done()
		prev := p.possiblePlacements
		for plc := p.possiblePlacements.Next; plc != nil; plc = plc.Next {
			removed := false
			for pt := range placement {
				if plc.Value.Has(pt) {
					prev.Next = plc.Next // remove this from the list- placement is no longer valid
					removed = true
					break
				}
			}
			if !removed {
				prev = plc // else prev stays the same
			}
		}
	}

	// TODO: use goroutines to speed this up
	for _, p := range g.players {
		if p == nil {
			continue
		}
		wg.Add(1)
		go updatePlayerPlacements(p)
	}
	// NOTE: could technically continue, we only need to wait for the player who just placed
	//       this wouldn't be too much of an optimization though
	wg.Wait()

	asPiece := PieceFromPoints(placement)
	prev := player.possiblePlacements
	for plc := player.possiblePlacements.Next; plc != nil; plc = plc.Next {

		if PieceFromPoints(plc.Value).IsSame(asPiece) {
			prev.Next = plc.Next // remove this from the list- placement is no longer valid
			continue
		}
		removed := false
		for pt := range placement {
			if plc.Value.Has(pt) ||
				plc.Value.Has(pt.GetAdjacent(types.UP)) ||
				plc.Value.Has(pt.GetAdjacent(types.DOWN)) ||
				plc.Value.Has(pt.GetAdjacent(types.RIGHT)) ||
				plc.Value.Has(pt.GetAdjacent(types.LEFT)) {
				prev.Next = plc.Next // remove this from the list- placement is no longer valid
				removed = true
				break
			}
		}
		if !removed {
			prev = plc
		}
	}

	newCorners := g.state.board.findCorners(
		placement.ToSlice(),
		types.Owner(player.state.pid),
	)

	last := prev
	for _, pt := range newCorners {
		last.Next = g.state.board.getPlacementsAtPoint(
			pt,
			types.Owner(player.state.pid),
			player.state.pieces.Copy(),
		).Next
		for ; last.Next != nil; last = last.Next {
		}
	}
}

func (g *Game) updateGameState(player *Player) bool {
	defer g.sendPlayerList() // broadcast updated player list
	defer g.sendGameStatus() // broadcast updated game status

	if g.state.turn == player.state.pid {
		g.nextTurn() // advance turn if necessary
		if g.state.turn == PID_NONE && (g.config.TurnBased || g.config.Players == 1) {
			g.endGame()
			return true // game over
		} else if g.config.TimeControl > 0 {
			nextPlayer, _ := g.getPlayer(g.state.turn)
			nextPlayer.playerTimer.Start()
		}
	}
	return false // game not over
}

func (g *Game) endGame() {
	winners := g.determineWinners()
	winString := "Game over! "
	if len(winners) > 1 {
		for i := 0; i < len(winners)-1; i++ {
			winners[i].state.status.Set(DRAWN)

			winString += winners[i].name
			if len(winners) > 2 {
				winString += ", "
			} else {
				winString += " "
			}
		}
		winString += "and " + winners[len(winners)-1].name + " tied!"
	} else {
		winners[0].state.status.Set(WINNER)
		winString += winners[0].name + " wins!"
	}
	g.sendGameMessage(winString)
	g.state.status.Set(COMPLETE)
	g.evalEngine.Stop()

	// Stop all player timers
	for _, player := range g.players {
		if player != nil {
			if player.connectionTimer != nil {
				player.connectionTimer.Pause()
			}
			if g.config.TimeControl > 0 {
				player.playerTimer.Pause()
			}
		}
	}
}

func (g *Game) receiveMessages(player *Player) {
	for {
		var inMsg types.SocketData
		err := g.socketManager.Recv(player.socket, &inMsg)
		if err != nil {
			fmt.Println(err)
			break
		}
		if inMsg.Type == sockets.CHAT_MESSAGE {
			g.socketManager.Broadcast(&inMsg)
		}
	}

	// handle socket disconnection type events
	g.lock.Lock()
	defer g.lock.Unlock()
	g.socketManager.Disconnect(player.socket)
	player.socket = nil
	player.state.status.Clear(CONNECTED)
	fmt.Println("Disconnected", player.state.pid)

	g.sendPlayerList()
	g.sendGameStatus()

	player.connectionTimer = utilities.InitTimer(15000, 0, func(...any) {
		g.lock.Lock()
		defer g.lock.Unlock()
		player.state.status.Set(DISABLED) // Remove player from active set
		player.playerTimer.Pause()        // stop timer if applicable
		g.updateGameState(player)
		g.sendGameMessage(fmt.Sprintf("%s has left the game", player.name))
	})
	player.connectionTimer.Start()
}

func (g *Game) sendGameMessage(msg string) {
	g.socketManager.Broadcast(&types.SocketData{
		Type: sockets.CHAT_MESSAGE,
		Data: &types.ChatMessage{
			Origin:  types.RESERVED,
			Message: msg,
		},
	})
}

func (g *Game) sendPlayerList() {
	players := make([]types.PlayerConfig, 0, len(g.players))
	for pid, player := range g.players {
		if player != nil {
			players = append(players, types.PlayerConfig{
				PID:    pid,
				Name:   player.name,
				Color:  player.color,
				Status: player.state.status,
				Time:   player.playerTimer.TimeLeftMs(),
			})
		}
	}

	g.socketManager.Broadcast(&types.SocketData{
		Type: sockets.PLAYER_UPDATE,
		Data: players,
	})
}

func (g *Game) sendGameStatus() {
	g.socketManager.Broadcast(&types.SocketData{
		Type: sockets.GAME_STATUS,
		Data: &types.PublicGameState{Turn: g.state.turn, Status: g.state.status},
	})
}

func (g *Game) getPlayer(pid types.PlayerID) (*Player, error) {
	player, ok := g.players[pid]
	if !ok {
		return nil, errors.New("invalid player id")
	}
	return player, nil
}

func (g *Game) AddPlayer(name string, color uint) (types.PlayerID, error) {
	/* Assign the new player a PID, if there is one available */
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.state.status.Has(FULL) {
		return 0, errors.New("game full")
	}

	var ii int
	for ii = 1; ii <= len(g.players); ii++ {
		pid := types.PlayerID(ii)
		if g.players[pid] == nil {
			g.players[pid] = &Player{
				name:  name,
				color: color,
				state: &PlayerState{
					pid:    pid,
					status: JOINED,
					pieces: g.startingPieces.Copy(),
				},
				socket: nil,
				playerTimer: utilities.InitTimer(
					g.config.TimeControl*1000,
					g.config.TimeBonus*1000,
					g.handleTimeout,
					pid,
				),
				connectionTimer: nil,
				possiblePlacements: g.state.board.getPlacementsAtPoint(
					g.state.board.getOrigin(pid),
					types.Owner(pid),
					g.startingPieces.Copy(),
				),
				hints: g.config.Hints,
			}
			break
		}
	}

	fmt.Println("Added player ", ii)
	g.lastActive = time.Now()

	if ii == len(g.players) {
		fmt.Println("Game is full")
		g.state.status.Set(FULL)
	}

	return types.PlayerID(ii), nil
}

func (g *Game) ConnectSocket(socket *websocket.Conn, pid types.PlayerID) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	player, err := g.getPlayer(pid)
	if err != nil {
		return err
	}

	if !player.state.status.Has(JOINED) {
		return errors.New("invalid player status")
	}

	player.socket = g.socketManager.Connect(socket)
	player.state.status.Set(CONNECTED)

	// begin receiving messages on this socket
	go g.receiveMessages(player)

	if player.connectionTimer != nil {
		player.connectionTimer.Pause()
	}

	fmt.Println("Connected player ", pid)

	var playerPieces []types.PublicPiece

	for piece := range player.state.pieces {
		playerPieces = append(playerPieces, types.PublicPiece{
			Hash: piece.Hash(),
			Body: piece.ToPoints().ToSlice(),
		})
	}

	// Send all players the current player list and status to sync up
	g.sendPlayerList()
	g.sendGameStatus()
	// Send the player their PID, pieces, and current board on connection
	g.socketManager.Send(
		player.socket,
		&types.SocketData{
			Type: sockets.PRIVATE_GAME_STATE,
			Data: &types.PrivateGameState{PID: pid, Pieces: playerPieces, Hints: player.hints},
		},
	)
	g.socketManager.Send(player.socket, &types.SocketData{Type: sockets.BOARD_STATE, Data: g.state.board.GetRaw()})
	g.sendGameMessage(fmt.Sprintf("%s has joined the game", player.name))

	return nil
}

func (g *Game) PlacePiece(pid types.PlayerID, placement types.Placement) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	player, err := g.getPlayer(pid)
	if err != nil {
		return err
	}

	_, err = g.playerActionValid(player)
	if err != nil {
		return err
	}

	// ptSet := utilities.NewSet(placement)
	// piece := PieceFromPoints(ptSet)
	// if !player.state.pieces.Has(piece) {
	// 	return errors.New("player does not have this piece")
	// }

	internalPlace := utilities.NewSet(placement)
	hasPlacement := false
	for p := player.possiblePlacements.Next; p != nil; p = p.Next {
		if p.Value.Is(internalPlace) {
			hasPlacement = true
			break
		}
	}
	if !hasPlacement {
		return errors.New("invalid placement")
	}

	_, err = g.state.board.Place(internalPlace, pid)
	if err != nil {
		return err
	}

	g.lastActive = time.Now()

	g.state.status.Set(IN_PROGRESS)

	g.socketManager.Broadcast(&types.SocketData{
		Type: sockets.BOARD_UPDATE,
		Data: &types.BoardUpdate{
			Owner:     types.Owner(pid),
			Placement: placement,
		},
	})

	player.playerTimer.Pause() // successfully placed piece, pause timer

	player.state.pieces.Remove(PieceFromPoints(internalPlace))

	g.updateValidPlacements(player, internalPlace)

	g.updateGameState(player)
	// if !gameOver {
	// 	playerStates := make(map[types.PlayerID]*PlayerState, len(g.players))
	// 	for pid, player := range g.players {
	// 		if player != nil {
	// 			playerStates[pid] = player.state
	// 		}
	// 	}
	// 	g.evalEngine.Evaluate(&EvalState{g.state, playerStates})
	// }

	return nil
}

func (g *Game) GetHint(pid types.PlayerID) (types.Point, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	player, err := g.getPlayer(pid)
	if err != nil {
		return types.Point{}, err
	}

	_, err = g.playerActionValid(player)
	if err != nil {
		return types.Point{}, err
	}

	if player.hints == 0 {
		return types.Point{}, errors.New("no more hints")
	}

	player.hints -= 1

	// TODO: don't return the same hint twice in a row
	if player.possiblePlacements.Next != nil {
		for pt := range player.possiblePlacements.Next.Value {
			if g.state.board.HasCorner(pt, types.Owner(pid)) {
				return pt, nil
			}
		}
	}
	return types.Point{}, errors.New("no corners")
}

func (g *Game) playerActionValid(player *Player) (bool, error) {
	if g.config.TurnBased && g.state.turn != player.state.pid {
		return false, errors.New("not your turn")
	}

	if player.state.status.Has(TIMED_OUT | DISABLED) {
		return false, errors.New("player inactive")
	}

	if !g.state.status.Has(FULL) {
		return false, errors.New("waiting for all players")
	}

	return true, nil
}

func (g *Game) determineWinners() []*Player {
	winners := make([]*Player, 0, len(g.players))
	scores := make(map[*Player]int, len(g.players))
	minScore := 0xffffffff

	for _, player := range g.players {
		if player != nil {
			score := 0
			for piece := range player.state.pieces {
				score += int(piece.Size())
			}
			scores[player] = score
			if score <= minScore {
				minScore = score
			}
		}
	}

	for player, score := range scores {
		fmt.Println(player.name, score)
		if score == minScore {
			winners = append(winners, player)
		}
	}
	return winners
}

func (g *Game) handleTimeout(args ...any) {
	g.lock.Lock()
	defer g.lock.Unlock()

	pid := args[0].(types.PlayerID)
	player, _ := g.getPlayer(pid)
	player.state.status.Set(TIMED_OUT | DISABLED)
	if player.connectionTimer != nil {
		player.connectionTimer.Pause()
	}
	g.updateGameState(player)
}
