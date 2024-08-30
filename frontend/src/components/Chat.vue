<template>
  <v-expansion-panels v-model="chatPanel">
    <v-expansion-panel expand-icon="mdi-chevron-up" collapse-icon="mdi-chevron-down">
      <v-expansion-panel-title static color="primary">
        <v-badge v-if="unreadMessages" color="error" :content="unreadMessages > 9 ? `9+` : unreadMessages">
          <v-icon>mdi-email-outline</v-icon>
        </v-badge>
        <v-icon v-else>mdi-email-outline</v-icon>
        <span class="chat-title">Chat</span>
      </v-expansion-panel-title>
      <v-expansion-panel-text class="chat-box-expansion">
        <div class="chat-box">
          <div class="live-chat">
            <div v-for="msg in messages" class="chat-msg">
              <span
                v-if="!(msg.origin & (1<<31))"
                class="font-weight-black"
                :style="{color: players[msg.origin].color}"
              >
                  {{ players[msg.origin].name }}: 
              </span>
              <span :class="(msg.origin & (1<<31)) ? gameMsgStyle : ''">{{ msg.msg }}</span>
            </div>
          </div>
          <div>
            <v-text-field
              v-model="myChat"
              placeholder="Say something..."
              hide-details
              variant="outlined"
              @keydown.enter="sendMessage"
            >
              <template v-slot:append-inner>
                <v-btn 
                  size="small"
                  variant="plain"
                  @click="sendMessage"
                  :disabled="!myChat"
                  :color="players[pid]?.color"
                  :ripple="false"
                >
                  <v-icon size="x-large">mdi-send</v-icon>
                </v-btn>
              </template>
            </v-text-field>
          </div>
        </div>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script setup>
import { ref, watch } from 'vue'
import { MessageType } from '@/api';

const emit = defineEmits(["send"]);

const props = defineProps({
  pid: Number,
  messages: Array<Object>{},
  players: Object,
});

const gameMsgStyle = ref([
  "font-italic",
  "font-weight-thin",
]);

const myChat = ref("");
const chatPanel = ref(0);
const unreadMessages = ref(0);

function sendMessage() {
  if (myChat.value.length) {
    
    emit(
      "send",
      JSON.stringify({
          type: MessageType.ChatMesssage,
          data: {origin: props.pid, message: myChat.value},
        })
    );
    myChat.value = "";
  }
}

watch(() => props.messages.length, () => {
  if (chatPanel.value === 0) {
    unreadMessages.value = 0;
  } else {
    unreadMessages.value += 1;
  }
});

watch(() => chatPanel.value, () => {
  if (chatPanel.value === 0) {
    unreadMessages.value = 0;
  }
});

</script>

<style scoped>
.chat-box {
  height: 100%;
  display: flex;
  flex-direction: column;
  font-size: medium;
}

.live-chat {
  display: flex;
  flex: 1;
  flex-direction: column-reverse;
  padding: 0.5em;
  overflow-y: auto;
}

.chat-box-expansion > * {
  padding: 0;
  height: 300px;
}

.chat-title {
  margin-left: 2em;
}
</style>