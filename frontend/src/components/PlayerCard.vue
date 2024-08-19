<template>
  <div class="player-card-container">
    <v-card class="player-card" :variant="myTurn ? 'outlined': 'plain'">
      <span class="player-name text-h6" :style="{color: color}">
        <slot name="prepend"/>
        {{ name }}
        <slot name="append"/>
      </span>
    </v-card>
  </div>
  
</template>

<script setup>

defineProps({
  name: String,
  color: String,
  myTurn: Boolean,
});

</script>

<style scoped>
.player-card {
  margin: 1em 0;
}

.player-name {

}
</style>