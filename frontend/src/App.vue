<template>
  <v-app fluid class="bkgnd fill-height">
    <v-toolbar app color="primary" dark>
      <v-toolbar-title>GoBloks</v-toolbar-title>
      <v-spacer></v-spacer>
      
      <v-btn text v-if="store.inGame" @click="exitGame">Exit Game<v-icon class="mx-2">mdi-export</v-icon></v-btn>
    </v-toolbar>  
   
    <router-view></router-view>
  </v-app>
</template>

<script setup>
import { onBeforeMount } from "vue";
import { useRouter } from 'vue-router';
import { useStore } from '@/stores/store';

const router = useRouter();
const store = useStore()

onBeforeMount(() => {
  router.push({ path: "/join" });
});

function exitGame() {
  store.revokeToken();
  router.push({ path: "/join" });
}

</script>

<style>

#app {
  height: 100vh;
}

.bkgnd {
  background-color: #ffffff50;
}

#app:before {
  content: "";
  position: fixed;
  width: 300%;
  height: 300%;
  top: -100%;
  left: -100%;
  z-index: -1;
  background: url("./assets/blokus.svg") repeat;
  background-size: 40%;
  transform: rotate(-30deg);
}
</style>
