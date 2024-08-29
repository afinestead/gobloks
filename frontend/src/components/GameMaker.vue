<template>
  <v-card
    class="mx-auto"
    width="344"
  >
    <v-card-text>
      <v-row>
        <v-col cols="4" class="my-auto">
          <v-label>Players</v-label>
        </v-col>
        <v-col>
          <v-text-field
            width="80"
            v-model="nPlayers"
            :rules="[rules.required]"
            type="number"
            min="1"
            max="32"
            step="1"
            outlined
            density="compact"
            hide-details
            variant="solo-filled"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="4" class="my-auto">
          <v-label>Block Degree</v-label>
        </v-col>
        <v-col cols="4">
          <v-text-field
            width="80"
            v-model="blockDeg"
            type="number"
            min="1"
            max="8"
            step="1"
            outlined
            density="compact"
            hide-details
            variant="solo-filled"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="4" class="my-auto">
          <v-label>Difficulty</v-label>
        </v-col>
        <v-col>
          <v-slider
            v-model="density"
            :max="1.0"
            :min="0.7"
            class="align-center"
            hide-details
            density="compact"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="4" class="my-auto">
          <v-label>Take turns</v-label>
        </v-col>
        <v-col>
          <v-switch
            v-model="turns"
            class="align-center"
            hide-details
            density="compact"
            :disabled="nPlayers < 2"
          />
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        @click.stop="tryCreate"
        size="large"
        variant="tonal"
        :disabled="!nPlayers || !blockDeg"
      >
        Create
      </v-btn>
      <v-spacer></v-spacer>
    </v-card-actions>
  </v-card>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router';
import { useStore } from "@/stores/store"


const router = useRouter();
const store = useStore();

const nPlayers = ref(4);
const blockDeg = ref(5);
const density = ref(0.85);
const turns = ref(true)

const rules = ref({
  required: (v) => !!v || "Required",
})

// watch(nPlayers, (v) => nPlayers.value = Math.min(Math.max(Math.round(v), 1), 32));
// watch(blockDeg, (v) => blockDeg.value = Math.min(Math.max(Math.round(v), 1), 8));


function tryCreate() {
  store.createGame({
    players: parseInt(nPlayers.value),
    degree: parseInt(blockDeg.value),
    density: density.value,
    turns: turns.value,
  }).then((gid) => {
    router.push({ path: '/join', query: { game: gid } });
  }).catch((e) => {
    console.error(e);
    // errorMessages.value = "Unable to create game";
  });

  console.log("Create game with", nPlayers.value, "players, block size", blockDeg.value, "and density", density.value);
}

</script>

<style scoped>

</style>