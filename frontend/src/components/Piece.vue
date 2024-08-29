<template>
  <div class="piece" :style="pieceStyle">
    <div
      v-for="(block, idx) in blocksInternal"
      :key="idx"
      class="block"
      :style="[
        blockStyle,
        {
          backgroundColor: color,
          top: `${block.x*squareSize}px`,
          left: `${block.y*squareSize}px`,
        }
      ]"
      @click.stop="handlePieceClick"
      @mousedown.right="handlePieceClick"
      @contextmenu.prevent
    />
  </div>
</template>

<script setup>

import { computed, ref, onMounted, watch } from 'vue';

const props = defineProps({
  squareSize: {
    type: Number,
    default: 0,
  },
  blocks: {
    type: Array,
    default: () => [],
  },
  color: {
    type: String,
    default: ""
  },
});

const emit = defineEmits(["click", "change"]);

const blocksInternal = ref([]);
const clicks = ref(0);
const clickTimer = ref(null);
const dblClickDelay = 200;

const pieceLen = computed(() => {
  let maxX = -Infinity;
  let maxY = -Infinity;
    
  for (const b of blocksInternal.value) {    
    maxX = Math.max(maxX, b.x);
    maxY = Math.max(maxY, b.y);
  }
  return [maxX + 1, maxY + 1];
});

const pieceStyle = computed(() => ({
  height: `${pieceLen.value[0]*props.squareSize}px`,
  width: `${pieceLen.value[1]*props.squareSize}px`,
}));

const blockStyle = computed(() => ({
  height: `${props.squareSize-1}px`,
  width: `${props.squareSize-1}px`,
}));

defineExpose({
  handlePieceClick,
  blocksInternal,
});

onMounted(() => {
  blocksInternal.value = props.blocks;
});

watch(() => props.blocks, (newBlocks) => blocksInternal.value = newBlocks);

function handlePieceClick(evt) {
  if (evt.button === 2) {
    flipPiece("x");
  } else {
    clicks.value++;
    if (clicks.value === 1) {
      clickTimer.value = setTimeout(() => {
        emit("click", evt);
        clicks.value = 0;
      }, dblClickDelay);
    } else {
      // Double click
      rotatePiece(90);
      clearTimeout(clickTimer.value);
      emit("change");
      clicks.value = 0;
    }
  }
}

function translatePiece(blocks) {
  let minX = Infinity;
  let minY = Infinity;
  for (const b of blocks) {
    minX = Math.min(minX, b.x);
    minY = Math.min(minY, b.y);
  }
  return blocks.map(b => ({x: b.x-minX, y: b.y-minY}));
}

function flipPiece(ax) {
  const p = 
    ax === "x" ? blocksInternal.value.map(b => ({x: b.x, y: -b.y})) :
                 blocksInternal.value.map(b => ({x: -b.x, y: b.y}));
  blocksInternal.value = translatePiece(p);
  emit("change");
}


function rotatePiece(deg) {
  const rad = deg * Math.PI / 180;
  const cos = Math.cos(rad);
  const sin = Math.sin(rad);
  const p = blocksInternal.value.map(b => ({x: Math.round(b.x*cos - b.y*sin), y: Math.round(b.x*sin + b.y*cos)}));
  blocksInternal.value = translatePiece(p);
  emit("change");
}
</script>

<style scoped>
.piece {
  position: relative;
}

.block {
  position: absolute;
  border: 1px solid black;
  margin: 1px;
  box-sizing: border-box;
}
</style>