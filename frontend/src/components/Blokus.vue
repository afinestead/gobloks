<template>
  <!-- TODO:
  If a piece is rotated/flipped while hovering, the highlighting doesn't update
  -->
  <div class="game-container">      
    <div class="my-pieces">
      <piece
        class="unplaced-piece"
        v-if="myPieces.length !== 0"
        v-for="(p,idx) in myPieces"
        :key="idx"
        :color="playerColors[playerID]"
        :blocks="p"
        :square-size="16"
        @click.stop="handlePieceClick($event, p, idx)"
        @contextmenu.prevent
        />
    </div>

    <div v-if="board" ref="boardRef" class="board">
      <div v-for="row, i in board" :key="i" class="board-row">
        <board-square
          v-for="pid, j in row"
          :key="j"
          :owner="pid"
          :colors="playerColors"
          @mouseover="calculateOverlap(i, j)"
        />
      </div>
    </div>

    <div class="interact-area">
      <div class="chat-box">
        <div :class="['game-state', gameMsgStyle]">
          <span v-if="gameStatus === 'waiting'">Waiting for all players to join</span>
          <span v-else-if="gameStatus === 'active'">
            <span class="font-weight-black" :style="{color: playerColors[whoseTurn]}">{{ playerNames[whoseTurn] }}</span>'s turn
          </span>
          <span v-else-if="gameStatus === 'done'">Player 1 wins!</span>
        </div>
        <div class="live-chat">
          <div v-for="msg in liveChat" class="chat-msg">
            <span
              v-if="!(msg.origin & (1<<31))"
              class="font-weight-black"
              :style="{color: playerColors[msg.origin]}"
            >
                {{ playerNames[msg.origin] }}: 
            </span>
            <span :class="(msg.origin & (1<<31)) ? gameMsgStyle : ''">{{ msg.msg }}</span>
          </div>
        </div>
        <div class="my-chat">
          <v-text-field
            v-model="myChat"
            placeholder="Say something..."
            hide-details
            variant="outlined"
            @keydown.enter="sendMessage"
          >
            <template v-slot:append-inner>
              <v-btn 
                size="small"
                variant="plain"
                @click="sendMessage"
                :disabled="!myChat"
                :color="playerColors[playerID]"
                :ripple="false"
              >
                <v-icon size="x-large">mdi-send</v-icon>
              </v-btn>
            </template>
          </v-text-field>
        </div>
      </div>
    </div>

    <piece
      v-if="selectedPiece"
      class="selected-piece"
      ref="selectedPieceRef"
      :color="playerColors[playerID]"
      :blocks="selectedPiece?.block || []"
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
  </div>

 

</template>

<script setup>
import { nextTick, onMounted, ref, reactive, computed, watch } from 'vue'
import BoardSquare from './BoardSquare.vue';
import Piece from './Piece.vue'
import PlayerCard from './PlayerCard.vue'
import Players from './Players.vue'
import { useRouter } from 'vue-router';
import { useStore } from '@/stores/store';
import { MessageType } from '@/api';
import Message from '@/api/model/Message';


const store = useStore();
const router = useRouter()
const squareSize = ref(0);

const gameStatus = ref(null);

const board = ref([]);
const boardSize = ref(0);
const boardRef = ref(null)
const myPieces = ref([]);
const allPlayers = ref([]);
const playerID = ref(null);
const playerColors = computed(() => allPlayers.value?.reduce((acc, p) => ({...acc, [p.pid]: p.color ? `#${p.color.toString(16).padStart(6, '0')}` : "#ffffff"}), {}));
const playerNames = computed(() => allPlayers.value?.reduce((acc, p) => ({...acc, [p.pid]: p.name}), {}));
const whoseTurn = ref(null);

const selectedPiece = ref(null);
const selectedPieceRef = ref(null);
const selectedPieceOverlap = ref(null);
const cursorX = ref(null);
const offsetX = ref(null);
const cursorY = ref(null);
const offsetY = ref(null);

const ws = ref(null);

const myChat = ref("");
const liveChat = ref([]);
const gameMsgStyle = ref([
"font-italic",
"font-weight-thin",
]);

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
      selectedPieceRef.value.handlePieceClick(event);
    }
  };

  document.oncontextmenu = (event) => {
    if (selectedPiece.value) {
      selectedPieceRef.value.handlePieceClick(event);
    }
    return false;
  };
  
  // open a websocket for game updates
  const new_ws = store.socket;
  new_ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    console.log(msg);

    switch (msg.type) {

      case MessageType.ChatMesssage:
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
        myPieces.value = msg.data.pieces;
        break;
      
      case MessageType.PlayerUpdate:
        allPlayers.value = msg.data.players.sort((p1,p2) => p1.pid - p2.pid);  
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
  // debugger; // eslint-disable-line
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

  console.log(placement);
  
  return store.placePiece(placement);
};

function sendMessage() {
  if (myChat.value.length) {
    ws.value.send(
      JSON.stringify({
        type: MessageType.ChatMesssage,
        data: {origin: playerID.value, message: myChat.value},
      })
    );
    myChat.value = "";
  }
}

watch(selectedPiece, (newPiece) => {
  if (newPiece === null) {
    clearHighlight()
  }
});

</script>

<style scoped>

.game-container {
  height: 100%;
  min-height: 480px;
  display: flex;
  justify-content: center;
}

.interact-area {
  width: 420px;
  height: 100%;
}

.my-pieces {
  border: 1px solid gray;
  border-radius: 4px;
  padding: 0.5em;
  height: 100%;
  width: 10%;
  overflow-y: auto;
  overflow-x: hidden;
  background-color: rgba(255,255,255,0.9);
}

.board {
  padding: 1em;
  border: 1px solid gray;
  border-radius: 4px;
  background-color: rgba(255,255,255,0.9);
  height: 100%;
  aspect-ratio: 1/1;
}

.board-row {
  display: flex;
}

.players {
  min-width: 12em;
  margin-right: 1em;
}

.highlighted {
  border-color: yellow;
}

.unplaced-piece {
  margin: 1em auto;
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

.chat-box {
  background-color: rgba(255,255,255,0.9);
  height: 100%;
  display: flex;
  flex-direction: column;
  font-size: x-large;
  min-width: 260px;
}

.game-state {
  border-radius: 4px;
  padding-left: 0.5em;
}

.live-chat {
  display: flex;
  flex: 1;
  flex-direction: column-reverse;
  padding: 0.5em;
  border: 1px solid grey;
  border-radius: 4px;
  overflow-y: auto;
}

</style>
