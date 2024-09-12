<template>
  <div
    class="board-square"
    :style="{ backgroundColor: sqColor, visibility: isHidden ? 'hidden' : 'unset' }"
  />
</template>

<script setup>

import { computed } from 'vue';

const props = defineProps({
  owner: Number,
  players: Object,
  color: String,
  hint: String,
});

const isOccupied = computed(() => !(props.owner & (1<<29)));
const isOrigin = computed(() => props.owner & (1<<30));
const isHidden = computed(() => props.owner & (1<<31));

const sqColor = computed(() => {
  const ownerColor = props.color || '#ffffff';
  if (isOccupied.value) {
    return `${ownerColor}ff`;
  } else if (isOrigin.value) {    
    return `${ownerColor}50`;
  } else if (props.hint) {
    return `${props.hint}50`;
  } else {
    return '#ffffffff';
  }
});
</script>

<style scoped>

.board-square {
  border-width: 1px;
  border-radius: 10%;
  border-style: solid;
  box-sizing: border-box;
  flex: 1;
  aspect-ratio: 1 / 1;
}
</style>