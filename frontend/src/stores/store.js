import { ApiClient, DefaultApi } from '../api/index'
import { defineStore } from 'pinia'

function getAPIlocation() {
    return `${window.location.hostname}:8888`;
}

export const useStore = defineStore("store", {
    state: () => {
        return {
            apiLocation: getAPIlocation(),
            api: new DefaultApi(new ApiClient(`http://${getAPIlocation()}`)),
            token: sessionStorage.getItem("accessToken"),
            inGame: false,
            ws: null,
        }
    },
    actions: {
        async createGame(gameConfig) {
            const r = await this.api.createGame(gameConfig);
            return r.data;
        },
        async joinGame(gameId, name, color) {
            const r = await this.api.joinGame(gameId, name, color);
            this.token = r.response.headers['access-token'];
            sessionStorage.setItem("accessToken", this.token);
        },
        async placePiece(placement) {
            return this.api.place(this.token, placement);
        },
        async requestHint() {
            return this.api.hint(this.token);
        },
        setGameActive(active) { this.inGame = active; },
        revokeToken() {
            this.disconnectSocket();
            this.inGame = false;
            this.token = null;
            sessionStorage.removeItem("accessToken");
        },
        connectSocket() {
            // NOTE: Smuggling access token to server via websocket protocol header
            //       https://stackoverflow.com/questions/4361173/http-headers-in-websockets-client-api
            this.ws = new WebSocket(`ws://${this.apiLocation}/ws?access_token=${this.token}`);
            return this.ws;
        },
        disconnectSocket() {
            if (this.ws) {
                this.ws.close();
                this.ws = null;
            }
        }
    },
});