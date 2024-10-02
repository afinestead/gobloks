<template>
  <v-card class="mx-auto my-5">
    <v-card-text>
      <v-row>
        <v-col cols="4" class="my-auto">
          <v-label>Players</v-label>
        </v-col>
        <v-col cols="5">
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
        <v-col cols="5">
          <v-text-field
            width="80"
            v-model="blockDeg"
            type="number"
            min="3"
            max="8"
            step="1"
            outlined
            density="compact"
            hide-details
            variant="solo-filled"
          />
        </v-col>
      </v-row>
      <!-- <v-row>
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
      </v-row> -->
      <!-- <v-row>
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
      </v-row> -->
      <v-row>
        <v-col cols="4" class="my-auto">
          <v-label>Time Control</v-label>
        </v-col>
        <v-col cols="4">
          <v-select
            v-model="timeControl"
            :items="[
              {
                title: '1 min',
                value: '60b0',
              },
              {
                title: '1 | 1',
                value: '60b1',
              },
              {
                title: '2 | 1',
                value: '120b1',
              },
              {
                title: '5 min',
                value: '300b0',
              },
              {
                title: '5 | 5',
                value: '300b5',
              },
              {
                title: '10 min',
                value: '600b0',
              },
              {
                title: '10 | 10',
                value: '600b10',
              },
              {
                title: '30 min',
                value: '1800b0',
              },
              {
                title: '1 day',
                value: '86400b0',
              },
              {
                title: '3 days',
                value: '259200b0',
              },
              {
                title: '7 days',
                value: '604800b0',
              },
            ]"
            hide-details
            dense
            outlined
            variant="solo-filled"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="4" class="my-auto">
          <v-label>Hints</v-label>
        </v-col>
        <v-col cols="5">
          <v-text-field
            width="80"
            v-model="hints"
            :rules="[rules.required]"
            type="number"
            min="0"
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
          <v-label>Private</v-label>
        </v-col>
        <v-col>
          <v-switch
            v-model="privateGame"
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
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router';
import { useStore } from "@/stores/store"


const router = useRouter();
const store = useStore();

const nPlayers = ref(4);
const blockDeg = ref(5);
const density = ref(0.85);
const turns = ref(true)
const timeControl = ref("600b0");
const hints = ref(3);
const privateGame = ref(false);

const rules = ref({
  required: (v) => !!v || "Required",
})

watch(nPlayers, (v) => {
  if (v < 2) {
    turns.value = false;
    privateGame.value = true;
  }
});

function tryCreate() {
  const [time, bonus] = timeControl.value.split('b');

  store.createGame({
    players: parseInt(nPlayers.value),
    degree: parseInt(blockDeg.value),
    density: density.value,
    turns: turns.value,
    timeSeconds: parseInt(time),
    timeBonus: parseInt(bonus),
    hints: parseInt(hints.value),
    public: !privateGame.value,
  }).then((gid) => {
    router.push({ path: '/join', query: { game: gid } });
  }).catch((e) => {
    console.error(e);
    // errorMessages.value = "Unable to create game";
  });
}

</script>

<style scoped>

</style>