<template>
  <v-container class="fill-height" fluid min-height="480">

    <v-row class="gameplay-area fill-height">
      <v-col cols="1" class="panel my-pieces bordered" ref="pieceDeck">
        <piece
          class="mx-auto my-8"
          v-if="myPieces.length !== 0"
          v-for="(p,idx) in myPieces"
          :key="idx"
          :color="allPlayers[playerID]?.color || '#ffffff'"
          :blocks="p.body"
          :square-size="16"
          @click.stop="handlePieceClick($event, p, idx)"
          @contextmenu.prevent
          :disabled="Boolean(allPlayers[playerID]?.status & (1<<2))"
          />
      </v-col>

      <v-col cols="9" class="panel board-view bordered">
        <div class="hint-btn">
          <v-btn
            :color="allPlayers[playerID]?.color || '#ffffff'"
            @click.stop="getHint"
            :disabled="hints <= 0 || hintRequested"
          >
            Hint
          </v-btn>
        </div>
        <div class="board fill-height mx-auto" ref="boardRef">
          <div v-for="row, i in board" :key="i" class="board-row">
            <board-square
              v-for="pid, j in row"
              :key="j"
              :owner="pid"
              :color="allPlayers[pid&0xffff]?.color"
              @mouseover="hoverX = i; hoverY = j; calculateOverlap(i, j)"
              @mouseout="hoverX = null; hoverY = null; clearHighlight()"
            />
          </div>
        </div>
      </v-col>

      <v-col cols="2" class="panel players bordered">
        <player-card
          v-for="player, idx in playerOrder"
          :key="idx"
          :player="player"
          :myTurn="whoseTurn === player.pid"
        >
          <template v-slot:timer>
            <timer
              :hide="((player.status&(1<<3)) === 0) && player.time === 0"
              :time="player.time"
              :active="whoseTurn === player.pid && (gameStatus & 0b111) == 0b011"
            />
          </template>
        </player-card>
      </v-col>
    </v-row>

    <piece
      v-if="selectedPiece"
      class="selected-piece"
      ref="selectedPieceRef"
      :color="allPlayers[playerID]?.color || '#ffffff'"
      :blocks="selectedPiece?.block.body || []"
      :square-size="squareSize"
      :style="{
        display: selectedPiece === null ? 'none' : 'inline-block',
        top: `${cursorY + offsetY - (squareSize / 2)}px`,
        left: `${cursorX + offsetX - (squareSize / 2)}px`,
      }"
      @click.stop="handlePieceClick($event, selectedPiece.block, selectedPiece.index)"
      @contextmenu.prevent
      @change="nextTick(() => snapPieceToCursor())"
    />

    <chat
      class="chat-drawer bordered"
      :messages="liveChat"
      :pid="playerID"
      :players="allPlayers"
      @send="msg => ws.send(msg)"
    />

  </v-container>

</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, computed, watch } from 'vue'
import BoardSquare from './BoardSquare.vue';
import Chat from './Chat.vue'
import Piece from './Piece.vue'
import PlayerCard from './PlayerCard.vue'
import Timer from './Timer.vue'
import { useRouter } from 'vue-router';
import { useStore } from '@/stores/store';
import { MessageType } from '@/api';


const store = useStore();
const router = useRouter()
const squareSize = ref(0);

const board = ref([]);
const boardSize = ref(0);
const boardRef = ref(null)
const myPieces = ref([]);
const pieceDeck = ref(null);

const allPlayers = ref([]);
const playerOrder = computed(() => {
  let sorted = []
  for (const pid in allPlayers.value) {
    let order = pid - whoseTurn.value
    if (order < 0) {
      order += Object.keys(allPlayers.value).length
    }
    allPlayers.value[pid].order = order
    sorted.push(allPlayers.value[pid])
  }
  sorted.sort((p1,p2) => p1.order - p2.order)
  
  return sorted
});

const playerID = ref(null);
const whoseTurn = ref(null);
const gameStatus = ref(0);

const selectedPiece = ref(null);
const selectedPieceRef = ref(null);
const selectedPieceOverlap = ref([]);

const hints = ref(0);
const hintRequested = ref(false);
const hintCoords = ref([]);

const hoverX = ref(null);
const hoverY = ref(null);

const cursorX = ref(null);
const offsetX = ref(null);
const cursorY = ref(null);
const offsetY = ref(null);

const ws = ref(null);

const liveChat = ref([]);

// const boardHTML = ref([]);
  // nextTick(() => boardHTML.value = Array.from(boardRef.value?.children || []).map(htmlCol => Array.from(htmlCol?.children)));
const boardHTML = computed(() => Array.from(boardRef.value?.children || []).map(htmlCol => Array.from(htmlCol?.children)));

const IsOccupied = ([x,y]) => !Boolean(board.value[x][y] & (1<<29));
const IsValid = ([x,y]) => !Boolean(board.value[x][y] & (1<<31));
const SquarePID = ([x,y]) => (board.value[x][y] & 0xffff);
const IsMyOrigin = coords => IsValid(coords) && Boolean(board.value[coords[0]][coords[1]] & (1<<30)) && SquarePID(coords) === playerID.value;
const IsOtherOrigin = coords => IsValid(coords) && Boolean(board.value[coords[0]][coords[1]] & (1<<30)) && SquarePID(coords) !== playerID.value;
const OccupiedByMe = coords => IsValid(coords) && IsOccupied(coords) && SquarePID(coords) === playerID.value;

function calculateOverlap(i, j) {
  if (selectedPiece.value) {

    const nullCoords = [null,null];

    const GetNeighborCoords = (dir, [x, y]) => {
      if (x === null || y === null) {
        return nullCoords;
      }
      switch (dir) {
        case "up": return x > 0 ? [x-1,y] : nullCoords;
        case "down": return x < boardSize.value-1 ? [x+1,y] : nullCoords;
        case "left": return y > 0 ? [x,y-1] : nullCoords;
        case "right": return y < boardSize.value-1 ? [x,y+1] : nullCoords;
        default: return nullCoords;
      }
    }

    const GetOverlapCoords = () => {
      let validOverlap = [];
      for (const c of selectedPieceRef.value.blocksInternal) {
        const overlapX = i + c.x - selectedPiece.value.origin[0];
        const overlapY = j + c.y - selectedPiece.value.origin[1];

        if (
          overlapX >= 0 && overlapX < boardSize.value &&
          overlapY >= 0 && overlapY < boardSize.value
        ) {
          validOverlap.push([overlapX, overlapY]);
        } else {
          return null;
        }
      }
      return validOverlap;
    }

    const HasSideNeighbor = coords => {
      if (coords == nullCoords) {
        return false;
      }
      const leftCoords = GetNeighborCoords("left", coords);
      const rightCoords = GetNeighborCoords("right", coords);
      const downCoords = GetNeighborCoords("down", coords);
      const upCoords = GetNeighborCoords("up", coords);

      return (
        (leftCoords !== nullCoords && OccupiedByMe(leftCoords)) ||
        (rightCoords !== nullCoords && OccupiedByMe(rightCoords)) ||
        (downCoords !== nullCoords && OccupiedByMe(downCoords)) ||
        (upCoords !== nullCoords && OccupiedByMe(upCoords))
      );
    }


    const HasCornerNeighbor = coords => {
      if (coords == nullCoords) {
        return false;
      }
      const leftUpCoords = GetNeighborCoords("left", GetNeighborCoords("up", coords));
      const leftDownCoords = GetNeighborCoords("left", GetNeighborCoords("down", coords));
      const rightUpCoords = GetNeighborCoords("right", GetNeighborCoords("up", coords));
      const rightDownCoords = GetNeighborCoords("right", GetNeighborCoords("down", coords));

      return (
        (leftUpCoords !== nullCoords && OccupiedByMe(leftUpCoords)) ||
        (leftDownCoords !== nullCoords && OccupiedByMe(leftDownCoords)) ||
        (rightUpCoords !== nullCoords && OccupiedByMe(rightUpCoords)) ||
        (rightDownCoords !== nullCoords && OccupiedByMe(rightDownCoords))
      );
    }
    
    const overlapCoords = GetOverlapCoords();

    if (
      overlapCoords !== null &&
      overlapCoords.every(coords => !IsOccupied(coords)) &&
      overlapCoords.every(coords => IsValid(coords)) &&
      overlapCoords.every(coords => !HasSideNeighbor(coords)) &&
      overlapCoords.every(coords => !IsOtherOrigin(coords)) &&
      overlapCoords.some(coords => HasCornerNeighbor(coords) || IsMyOrigin(coords))
    ) {
      // This is a valid placement

      for (const [x,y] of overlapCoords) {
        const sq = boardHTML.value[x][y];
        sq.classList.add("highlighted");
      }

      selectedPieceOverlap.value = overlapCoords;
      return;
    }
  }
  selectedPieceOverlap.value = [];
}

function onResize() {
  const sq = document.querySelector(".board-square")
  squareSize.value = Math.round(sq.getBoundingClientRect().width);
};

onBeforeUnmount(() => {
  store.disconnectSocket();

  window.removeEventListener('resize', onResize);
  document.onmousemove = null;
  document.onkeydown = null;

  router.beforeEach((to, from, next) => {
    store.setGameActive(false);
    next();
  });
});

onMounted(() => {

  store.setGameActive(true);

  nextTick(() => window.addEventListener('resize', onResize));

  document.onmousemove = (event) => {
    cursorX.value = event.pageX;
    cursorY.value = event.pageY;
  };

  document.onkeydown = (event) => {
    if (event.key === "Escape") {
      if (selectedPiece.value) {
        dropPiece(true);
      }
    }
  };

  document.onclick = (event) => {
    if (selectedPiece.value) {
      if (pieceDeck.value.$el.contains(event.target)) {
        dropPiece(true);
      } else {
        selectedPieceRef.value.handlePieceClick(event);
      }
    }
  };

  document.oncontextmenu = (event) => {
    if (selectedPiece.value) {
      selectedPieceRef.value.handlePieceClick(event);
      return false;
    }
    return true;
  };
  
  // open a websocket for game updates
  const new_ws = store.connectSocket();
  new_ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    // console.log(msg);

    switch (msg.type) {

      case MessageType.ChatMesssage:
        liveChat.value.unshift({
          origin: msg.data.origin,
          msg: msg.data.message,
        });
        break;
      
      case MessageType.BoardState:
        boardSize.value = msg.data.length;
        board.value = msg.data;
        nextTick(() => onResize());
        break;
      
      case MessageType.PrivateGameState:
        playerID.value = msg.data.pid; 
        myPieces.value = msg.data.pieces.sort((p1,p2) => p1.hash - p2.hash);
        hints.value = msg.data.hints;
        break;
      
      case MessageType.GameStatus:
        whoseTurn.value = msg.data.turn;
        gameStatus.value = msg.data.status;
        break;
      
      case MessageType.PlayerUpdate:
        allPlayers.value = msg.data.reduce((acc, p) => {
          acc[p.pid] = {
            pid: p.pid,
            name: p.name,
            color: `#${p.color.toString(16).padStart(6, '0')}`,
            status: p.status,
            time: p.timeMs,
          };
          return acc;
        }, {});
        break;

      default:
        console.log("unknown message type ", msg);
        break;
    }
  }

  new_ws.onerror = (e) => {
    store.revokeToken();
    router.push({ path: "/join" });
  }
  
  ws.value = new_ws;
});

function snapPieceToCursor() {
  // Find left most block
  let x = Infinity;
  let y = Infinity;
  const blocks = selectedPieceRef.value.$el.children;
  for (const block of blocks) {
    const rect = block.getBoundingClientRect();
    if (rect.top < x) {
      x = rect.top;
    }
  }

  let snappedToBlock;

  for (const block of blocks) {
    const rect = block.getBoundingClientRect();
    if (rect.top === x && rect.left <= y) {
      y = rect.left;
      snappedToBlock = block;
    }
  }
  const rect = selectedPieceRef.value.$el.getBoundingClientRect();
  offsetY.value = rect.top - x;
  offsetX.value = rect.left - y;

  // Find the coordinate of where we snapped to
  selectedPiece.value.origin = [
    Math.floor(parseInt(snappedToBlock.style.top) / squareSize.value),
    Math.floor(parseInt(snappedToBlock.style.left) / squareSize.value),
  ];

  clearHighlight();
  if (hoverX.value !== null && hoverY.value !== null) {
    calculateOverlap(hoverX.value, hoverY.value);
  }
};

async function pickupPiece(evt, piece, idx) {
  const block = evt.target.parentElement;
  if (!block.classList.contains("piece")) {
    return;
  }

  selectedPiece.value = {
    block: piece,
    index: idx,
    elem: block,
  };

  block.classList.add("hidden");

  await nextTick(() => snapPieceToCursor());
};

function dropPiece(discard) {
  if (selectedPiece.value !== null) {
    const block = selectedPiece.value.elem;
    if (discard) {
      block.classList.remove("hidden")
    } else {
      block.classList.add("removed")
    }
    selectedPiece.value = null;
    
  }
};


function placePiece() {
  if (isMyTurn()) {
    if (selectedPieceOverlap.value.length > 0) {
      issueBoardUpdate(selectedPieceOverlap.value)
        .then(() => {
          clearHint();
          dropPiece(false);
        })
        .catch(() => {});
    }
  }
};

function handlePieceClick(evt, piece, idx) {
  if (selectedPiece.value === null) {
    pickupPiece(evt, piece, idx);
  } else {
    placePiece();
  }
};

function clearHighlight() {
  for (const [x,y] of selectedPieceOverlap.value) {
    const sq = boardHTML.value[x][y];
    sq.classList.remove("highlighted");
  }
};

function clearHint() {
  hintRequested.value = false;
  for (const [x,y] of hintCoords.value) {
    const sq = boardHTML.value[x][y];
    sq.classList.remove("hinted");
  }
};

function isMyTurn() {
  return whoseTurn.value === playerID.value;
};

function issueBoardUpdate(piece) {
  let placement = [];
  for (const [x,y] of piece) {
    placement.push({x:x, y:y});
  }
  return store.placePiece(placement);
};

function getHint() {
  store.requestHint()
    .then((hint) => {
      hintRequested.value = true;
      hints.value -= 1;
      hintCoords.value = [[hint.data.x, hint.data.y]];      
      for (const [x,y] of hintCoords.value) {
        const sq = boardHTML.value[x][y];
        sq.classList.add("hinted");
      }
    })
    .catch(() => {});
};

watch(selectedPiece, (newPiece) => {
  if (newPiece === null) {
    clearHighlight()
  }
});

</script>

<style scoped>

.gameplay-area {
  background-color: rgba(255,255,255,0.9);
}

.bordered {
  border: 1px solid gray;
  border-radius: 4px;
}

.chat-drawer {
  position: absolute;
  bottom: 0;
  right: 12px;
  width: 344px;
}

.panel {
  height: 100%;
  padding: 0.5em;
}

.my-pieces {
  overflow-y: auto;
  overflow-x: hidden;
}

.board-view {
  position: relative;
  overflow-x: auto;
}

.board {
  aspect-ratio: 1/1;
}

.players {
  overflow-y: auto;
}

.board-row {
  display: flex;
}

.hint-btn {
  padding: 0 1em;
  position: absolute;
  right: 0;
}

.highlighted {
  border-color: yellow !important;
}

.hinted {
  border-color: greenyellow;
}

.selected-piece {
  position: absolute;
  pointer-events: none;
}

.hidden {
  visibility: hidden;
  opacity: 0;
}

.removed {
  height: 0 !important;
  padding: 0;
  margin: 0;
}


.game-state {
  border-radius: 4px;
  padding-left: 0.5em;
}

</style>
