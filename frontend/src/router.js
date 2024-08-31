import { createRouter, createWebHashHistory } from 'vue-router';
import GameJoiner from '@/components/GameJoiner.vue';
import GameMaker from '@/components/GameMaker.vue';
import Blokus from '@/components/Blokus.vue';

const routes = [
    { path: '/', redirect: '/join' },
    { path: '/join', component: GameJoiner },
    { path: '/create', component: GameMaker },
    { path: '/play', component: Blokus },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes,
});

export default router;