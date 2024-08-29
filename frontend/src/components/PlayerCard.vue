<template>
  <v-card
    v-if="player.name"
    :title="player.name"
    :class="['mx-auto', 'my-5']"
    :color="`${player.color}50`"
    :disabled="!isActive"
    :variant="myTurn ? 'elevated' : 'flat'"
  >
    <template v-slot:prepend>
      <v-icon :class="connected ? 'connected' : 'disconnected'" size="x-small">mdi-checkbox-blank-circle</v-icon>
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

const connected = computed(() => (props.player.status & 0x2) !== 0);
const isActive = computed(() => (props.player.status & 0x4) !== 0);

</script>

<style scoped>
.connected {
  color: green;
}

.disconnected {
  visibility: hidden;
}

.highlight {
  border: 5px solid #ffff00f0;
}
</style>