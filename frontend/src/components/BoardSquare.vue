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
});

const ownerColor = computed(() => props.players[props.owner&0xffff]?.color || '#ffffff');
const isOccupied = computed(() => !(props.owner & (1<<29)));
const isOrigin = computed(() => props.owner & (1<<30));
const isHidden = computed(() => props.owner & (1<<31));

const sqColor = computed(() => {
  if (isOccupied.value) {
    return `${ownerColor.value}ff`;
  } else if (isOrigin.value) {
    return `${ownerColor.value}50`;
  } else {
    return '#ffffffff';
  }
});

</script>

<style scoped>

.board-square {
  border-width: 1px;
  border-style: solid;
  margin: 1px;
  box-sizing: border-box;
  flex: 1;
  aspect-ratio: 1 / 1;
}
</style>