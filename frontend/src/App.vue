<template>
  <v-app fluid class="bkgnd fill-height">
    <v-toolbar app color="primary" dark>
      <v-toolbar-title>GoBloks</v-toolbar-title>
    </v-toolbar>  
   
    <router-view></router-view>
  </v-app>
</template>

<script setup>
import { onBeforeMount } from "vue";
import { useRouter } from 'vue-router';

const router = useRouter();

onBeforeMount(() => {
  router.push({ path: "/join" });
});

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
