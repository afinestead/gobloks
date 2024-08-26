<template>
  <div class="bkgnd">
    <router-view></router-view>
  </div>
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
  padding: 2em;
  overflow: auto;
}

.bkgnd {
  height: 100%;
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
