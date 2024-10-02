import { createRouter, createWebHashHistory } from 'vue-router';
import GameJoiner from '@/components/GameJoiner.vue';
import Blokus from '@/components/Blokus.vue';
import Lobby from '@/components/Lobby.vue';

const routes = [
    { path: '/', redirect: '/join' },
    { path: '/join', component: GameJoiner },
    { path: '/lobby', component: Lobby },
    { path: '/play', component: Blokus },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes,
});

export default router;