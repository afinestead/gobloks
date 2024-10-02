<template>
  <v-container class="fill-height" fluid min-height="480">
    <v-row class="fill-height" justify="center">
      <v-col class="active-games game-info bordered" cols="3">
        <v-card v-for="game in games" :key="game.gid" class="mb-4">
          <v-card-text>
            <v-row>
              <v-col cols="4" class="my-auto">
                <v-icon>mdi-account</v-icon>
                {{ game.players }} / {{ game.maxPlayers }}
              </v-col>
              <v-col cols="4" class="my-auto game-degree">
                <v-img v-if="game.degree === 3" :src="block3" class="inline-svg"/>
                <v-img v-else-if="game.degree === 4" :src="block4" class="inline-svg"/>
                <v-img v-else-if="game.degree === 5" :src="block5" class="inline-svg"/>
                <v-img v-else-if="game.degree === 6" :src="block6" class="inline-svg"/>
                <v-img v-else-if="game.degree === 7" :src="block7" class="inline-svg"/>
                <v-img v-else-if="game.degree === 8" :src="block8" class="inline-svg"/>
              </v-col>
              <v-col cols="4" class="my-auto">
                <v-icon>mdi-timer</v-icon>
                <span v-if="timeD(game.timeSeconds) > 0">{{ timeD(game.timeSeconds).toFixed(0) }} {{ timeD(game.timeSeconds) > 1 ? 'days' : 'day'  }}</span>
                <span v-else-if="timeH(game.timeSeconds) > 0">{{ timeH(game.timeSeconds).toFixed(0) }}:{{ timeM(game.timeSeconds) }}</span>
                <span v-else-if="timeM(game.timeSeconds) > 0">{{ timeM(game.timeSeconds) }}:{{ (timeS(game.timeSeconds) % 60).toString(10).padStart(2, '0') }}</span>
                <span v-else>0:{{ timeS(game.timeSeconds).toFixed(0) }}</span>
                <span>+{{ timeS(game.timeBonus) }}</span>
              </v-col>
            </v-row>
          </v-card-text>
          <v-card-actions>
            <v-spacer/>
            <v-btn color="primary" @click="null">Join</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col class="game-info bordered" cols="4">
        <game-maker/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useStore } from '@/stores/store';
import { useRoute } from 'vue-router';

import GameMaker from '@/components/GameMaker.vue';

import block3 from '@/assets/block3.svg';
import block4 from '@/assets/block4.svg';
import block5 from '@/assets/block5.svg';
import block6 from '@/assets/block6.svg';
import block7 from '@/assets/block7.svg';
import block8 from '@/assets/block8.svg';

const store = useStore();
const router = useRoute();

const games = ref([]);

const timeS  = (v) => Math.floor(v);
const timeM  = (v) => Math.floor(timeS(v) / 60);
const timeH  = (v) => Math.floor(timeM(v) / 60);
const timeD  = (v) => Math.floor(timeH(v) / 24);

onMounted(async() => {
  console.log('Lobby component mounted');
  const ws = store.connectLobbySocket();
  ws.onmessage = (event) => {
    console.log(event.data);
  };

  games.value = await store.listGames(0);
  console.log(games.value);
});

</script>

<style scoped>

.game-info {
  height: 100%;
  background-color: rgba(255,255,255,0.9);
}

.active-games {
  overflow-y: auto;
  
}

.bordered {
  border: 1px solid gray;
  border-radius: 4px;
}

.game-degree {
  display: flex;
  align-items: center;
} 

.inline-svg {
  height: 1.2em;
}

</style>