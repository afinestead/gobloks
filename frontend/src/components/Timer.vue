<template>
  <div class="timer px-4 text-center">
    <span v-if="timeH > 0">{{ timeH.toFixed(0) }}:{{ timeM }}:{{ timeS.toFixed(0) }}</span>
    <span v-else-if="timeM > 0">{{ timeM }}:{{ (timeS % 60).toString(10).padStart(2, '0') }}</span>
    <span v-else-if="timeS >= 20">0:{{ timeS.toFixed(0) }}</span>
    <span v-else>{{ (timeMs / 1000).toFixed(1) }}</span>
  </div>
</template>


<script setup>

import { computed, ref, onMounted, watch } from 'vue';

const props = defineProps({
  time: Number,
  active: Boolean
});

const timeInternal = ref(props.time);
const timeMs = computed(() => timeInternal.value);
const timeS  = computed(() => Math.floor(timeMs.value / 1000));
const timeM  = computed(() => Math.floor(timeS.value / 60));
const timeH  = computed(() => Math.floor(timeM.value / 60));

const timer = ref(null);

onMounted(() => {
  startTime();
});

function startTime() {
  clearInterval(timer.value);
  console.log("start time", props.active, timeInternal.value);
  if (props.active && timeInternal.value > 0) {
    timer.value = setInterval(() => {
      // hold value at 0.1s- 0.0 will be sent by the server
      // this prevents a potentially false 0 from being displayed
      timeInternal.value = Math.max(100, timeInternal.value - 100);
      if (timeInternal.value <= 100) {
        clearInterval(timer.value);
      }
    }, 100);
  }
}

watch(() => props.time, () => {
  clearInterval(timer.value);
  timeInternal.value = props.time;
  startTime(); // for 1 player games
});

watch(() => props.active, () => {
  if (!props.active) {
      clearInterval(timer.value);
  } else {
    startTime();
  }
});


</script>

<style scoped>
.timer {
  background-color: black;
  color: white;
  width: 5em;
}

</style>