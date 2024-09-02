<template>
  <v-card
    v-if="player.name"
    :title="player.name"
    :class="['mx-auto', 'mb-5', {highlight: myTurn}]"
    :color="`${player.color}50`"
    :disabled="!isActive"
    :variant="myTurn ? 'elevated' : 'flat'"
  >
    <template v-slot:prepend>
      <v-icon v-if="connected" size="x-small" color="green">mdi-checkbox-blank-circle</v-icon>
      <v-icon v-else="connected" size="x-small">mdi-checkbox-blank-circle-outline</v-icon>
    </template>
    <template v-slot:append>
      <slot name="timer"></slot>
    </template>
  </v-card>
  
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  player: Object,
  myTurn: Boolean,
});

const connected = computed(() => (props.player.status & (1<<1)) !== 0);
const isActive = computed(() => (props.player.status & (1<<2)) === 0);


</script>

<style scoped>

.highlight {
  border: 5px solid #ffff00ff;
}
</style>