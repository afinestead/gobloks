<template>
  <v-card
    v-if="player.name"
    :title="player.name"
    class="mx-auto my-5"
    :color="`${player.color}`"
    :disabled="connected"
  >
    <template v-slot:prepend>
      <v-icon :class="connected ? 'connected' : 'disconnected'" size="x-small">mdi-checkbox-blank-circle</v-icon>
    </template>
  </v-card>
  
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  player: Object,
  myTurn: Boolean,
});

const connected = computed(() => props.player.status);

</script>

<style scoped>
.connected {
  color: green;
}

.disconnected {
  visibility: hidden;
}
</style>