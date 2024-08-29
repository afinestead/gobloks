<template>
  <div class="timer px-4">
    <span v-if="timeM > 0">{{ timeM }}:{{ (timeS % 60).toFixed(0) }}</span>
    <span v-else-if="timeS > 20">0:{{ timeS.toFixed(0) }}</span>
    <span v-else>{{ timeS.toFixed(2) }}</span>
  </div>
</template>


<script setup>

import { computed, ref, watch } from 'vue';

const props = defineProps({
  time: Number,
  active: Boolean
});

const timeInternal = ref(props.time);
const timeS = computed(() => timeInternal.value / 1000);
const timeM = computed(() => Math.floor(timeS.value / 60));

const timer = ref(null);

watch(() => props.time, () => {
  console.log("Time changed to", props.time);
  
  timeInternal.value = props.time;
});

// watch(() => props.active, () => {
//   console.log("Active changed to", props.active);
  
//   if (props.active) {
//     t = setInterval(() => {
//       timeInternal.value -= 10;
//       if (timeInternal.value <= 0) {
//         clearInterval(t);
//       }
//     }, 10);
//   }
// });


</script>

<style scoped>
.timer {
  background-color: black;
  color: white;
}

</style>