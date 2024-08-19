import { createRouter, createWebHashHistory } from 'vue-router';
import App from '@/App.vue';
import GameJoiner from '@/components/GameJoiner.vue';
import GameMaker from '@/components/GameMaker.vue';
import Blokus from '@/components/Blokus.vue';

import { useStore } from '@/stores/store';

const routes = [
    { path: '/', component: App },
    { path: '/join', component: GameJoiner },
    { path: '/create', component: GameMaker },
    { path: '/play', component: Blokus },
    // {
    //     path: '/play',
    //     component: Blokus,
    //     beforeEnter: () => {
    //         const store = useStore();

    //         console.log("entering");
    //         return "/join";
    //         return store.dispatch("getGameState")
    //             .then(r => {
    //                 if (r.status === "done") {
    //                     throw new Error("Game complete");
    //                 }
    //                 return true; // Else stay here
    //             })
    //             .catch(() => {
    //                 store.revokeToken();
    //                 return "/join"
    //             });
    //     },
    // },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes,
});

// router.beforeEach(requireGameToken);

export default router;