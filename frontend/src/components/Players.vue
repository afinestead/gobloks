<template>
<div class="players-container">
    <v-card
        class="player-card"
        v-for="player in players"  
        :key="player.player_id"
        :width="200"
    >
        <span class="player-name text-h6" :style="{color: colors[player.player_id]}">
            <v-icon v-show="player.player_id === playerID" size="x-small" color="yellow">mdi-star</v-icon>
            {{ player.name }}
            <span v-show="player.player_id === playerID">
                <v-btn 
                    @click="colorPickerActive = !colorPickerActive"
                    icon="mdi-format-color-fill"
                    variant="plain"
                    size="small"
                />
                <v-dialog width="500" v-model="colorPickerActive">
                    <template v-slot:default>
                    <v-color-picker
                        :modes="['rgb']"
                        dot-size="10"
                        hide-inputs
                        v-model="colors[player.player_id]"
                    />
                    </template>
                </v-dialog>
            </span>
        </span>
    </v-card>
</div>

</template>
  
<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
    playerID: Number,
    players: Array,
    colors: Object,
});

const emit = defineEmits(["update:color"])

// const playerColors = computed(() => props.players?.reduce((acc, p) => ({...acc, [p.player_id]: p.color ? `#${p.color.toString(16).padStart(6, '0')}` : "#ffffff"}), {}));
// const playerNames = computed(() => props.players?.reduce((acc, p) => ({...acc, [p.player_id]: p.name}), {}));

const colorPickerActive = ref(false);

watch(colorPickerActive, (isActive) => {
  if (!isActive) {
    const color_update = props.colors[props.playerID];
    emit("update:color", color_update)
    // console.log("push player update");
    // ws.value.send(`{"update": {"color": ${color_as_int}}}`);
  }
});

</script>
  
<style scoped>
.player-card {
    margin: 1em 0;
}

.player-name {

}
</style>