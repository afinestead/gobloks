<template>
  <!-- TODO:
  If a piece is rotated/flipped while hovering, the highlighting doesn't update
  -->
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
          />
      </v-col>

      <v-col cols="9" class="panel board-view bordered">
        <div class="board fill-height mx-auto" ref="boardRef">
          <div v-for="row, i in board" :key="i" class="board-row">
            <board-square
              v-for="pid, j in row"
              :key="j"
              :owner="pid"
              :players="allPlayers"
              @mouseover="calculateOverlap(i, j)"
            />
          </div>
        </div>
      </v-col>

      <v-col cols="2" class="panel players bordered">
        <player-card
          v-for="player, idx in allPlayers"
          :key="idx"
          :player="player"
          :myTurn="whoseTurn === player.pid"
        >
          <template v-slot:timer>
            <timer :time="player.time"/>
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

    <v-expansion-panels v-model="chatPanel" class="chat-drawer bordered">
      <v-expansion-panel expand-icon="mdi-chevron-up" collapse-icon="mdi-chevron-down">
        <v-expansion-panel-title static color="primary">
          <v-badge v-if="unreadMessages" color="error" :content="unreadMessages > 9 ? `9+` : unreadMessages">
            <v-icon>mdi-email-outline</v-icon>
          </v-badge>
          <v-icon v-else>mdi-email-outline</v-icon>
          <span class="chat-title">Chat</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text class="chat-box-expansion">
          <chat
            :messages="liveChat"
            :pid="playerID"
            :players="allPlayers"
            @send="msg => ws.send(msg)"
          />
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>

  </v-container>

 

</template>

<script setup>
import { nextTick, onMounted, ref, computed, watch } from 'vue'
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
const playerID = ref(null);
const whoseTurn = ref(null);

const selectedPiece = ref(null);
const selectedPieceRef = ref(null);
const selectedPieceOverlap = ref(null);
const cursorX = ref(null);
const offsetX = ref(null);
const cursorY = ref(null);
const offsetY = ref(null);

const ws = ref(null);

const chatPanel = ref(0);
const chatOpen = ref(true);
const unreadMessages = ref(0);
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

    // TODO: be smarter about recomputing css
    clearHighlight();

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
  selectedPieceOverlap.value = null;
}

function onResize() {
const sq = document.querySelector(".board-square")
squareSize.value = Math.round(sq.getBoundingClientRect().width);
};

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
    }
    return false;
  };
  
  // open a websocket for game updates
  const new_ws = store.connectSocket();
  new_ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    console.log(msg);

    switch (msg.type) {

      case MessageType.ChatMesssage:
        if (!chatOpen.value) {
          unreadMessages.value += 1;
        }
        liveChat.value.unshift({
          origin: msg.data.origin,
          msg: msg.data.message,
        });
        break;
      
      case MessageType.PublicGameState:
        boardSize.value = msg.data.board.length;
        board.value = msg.data.board;
        whoseTurn.value = msg.data.turn;
        nextTick(() => onResize());
        break;
      
      case MessageType.PrivateGameState:
        playerID.value = msg.data.pid; 
        myPieces.value = msg.data.pieces.sort((p1,p2) => p1.hash - p2.hash);
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
    if (selectedPieceOverlap.value !== null) {
      issueBoardUpdate(selectedPieceOverlap.value)
        .then(() => dropPiece(false))
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

function clearHighlight(x,y) {
  document.querySelectorAll(".board-square").forEach(sq => sq.classList.remove("highlighted"));
};

function isMyTurn() {
  // TODO
  return true;
};

function issueBoardUpdate(piece) {
  let placement = [];
  for (const [x,y] of piece) {
    placement.push({x:x, y:y});
  }
  return store.placePiece(placement);
};

watch(selectedPiece, (newPiece) => {
  if (newPiece === null) {
    clearHighlight()
  }
});

watch(chatPanel, (opened) => {
  chatOpen.value = opened === 0;
  if (chatOpen.value) {
    unreadMessages.value = 0;
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

.chat-box-expansion > * {
  padding: 0;
  height: 300px;
}

.chat-title {
  margin-left: 2em;
}

.panel {
  height: 100%;
}

.my-pieces {
  padding: 0.5em;
  overflow-y: auto;
  overflow-x: hidden;
}

.board-view {
  overflow-x: auto;
}

.board {
  aspect-ratio: 1/1;
  padding: 1em;
}

.players {
  padding: 0.5em;
  overflow-y: auto;
}

.board-row {
  display: flex;
}

.highlighted {
  border-color: yellow;
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
