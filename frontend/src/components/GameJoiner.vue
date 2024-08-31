<template>
  <v-card class="join-card" width="344" elevation="24">
    <v-card-text>
      <v-text-field
        v-model="gameId"
        label="Game ID"
        placeholder="GAME"
        @input="gameId = gameId.toUpperCase(); errorMessage = ''"
        :error="errorMessage.length !== 0"
        :error-messages="errorMessage"
        :loading="joining"
      />
      <v-text-field
        v-model="playerName"
        label="Name yourself"
        @keydown.enter="tryJoin"
      >
        <template v-slot:append-inner>
          <v-menu
            v-model="colorPickerActive"
            :close-on-content-click="false"
          >
            <template v-slot:activator="{ props }">
              <v-btn 
                icon="mdi-format-color-fill"
                variant="plain"
                v-bind="props"
                :color="playerColor"
              />
            </template>
            <v-color-picker
              :modes="['rgb']"
              hide-inputs
              hide-canvas
              hide-sliders
              show-swatches
              v-model="playerColor"
            />
          </v-menu>
        </template>
      </v-text-field>
    </v-card-text>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        @click.stop="tryJoin"
        :color="playerColor"
        size="large"
        variant="tonal"
        :disabled="!playerName || !gameId"
      >
        Join
      </v-btn>
      <v-spacer></v-spacer>
    </v-card-actions>
  </v-card>
  <!-- <v-dialog width="500" v-model="colorPickerActive">
    <template v-slot:default>
      
    </template>
  </v-dialog> -->
  
</template>

<script setup>
import { onBeforeMount, onMounted, ref } from "vue";
import { useRouter, useRoute } from 'vue-router';
import { useStore } from "@/stores/store"


const router = useRouter();
const route = useRoute();
const store = useStore();

function randomColor() {
const inRange = () => Math.floor(Math.random() * 0xff).toString(16).padStart(2, "0");
return `#${inRange()}${inRange()}${inRange()}`;
}

const gameId = ref(route.query.game || "");
const playerName = ref("");
const playerColor = ref(randomColor());

const colorPickerActive = ref(false);
const joining = ref(false);
const errorMessage = ref("");

onBeforeMount(() => {
if (store.token) {
  router.push({ path: "/play" });
}
});

function tryJoin() {
  joining.value = true;
  const colorInt = parseInt(playerColor.value.substring(1), 16);
  store.joinGame(gameId.value, playerName.value, colorInt)
    .then(() => router.push({ path: "/play" }))
    .catch(e => {
      if (e.status === 409) {
        errorMessage.value = "This game is full";
      } else {
        errorMessage.value = "Unable to join game";
      }
    })
    .finally(() => joining.value = false)
}

</script>

<style scoped>

.join-card {
margin: auto;
}
</style>