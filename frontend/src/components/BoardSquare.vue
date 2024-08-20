<template>
  <div
    class="board-square"
    :style="{
      backgroundColor: isOccupied ? ownerColor : isOrigin ? `${ownerColor}50`: '#ffffff',
      visibility: isHidden ? 'hidden' : 'unset',
    }"
  />
</template>

<script setup>

import { computed } from 'vue';

const props = defineProps({
  owner: Number,
  colors: Object,
});

const ownerColor = computed(() => props.colors[props.owner&0xffff]);
const isOccupied = computed(() => !(props.owner & (1<<29)));
const isOrigin = computed(() => props.owner & (1<<30));
const isHidden = computed(() => props.owner & (1<<31));

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