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
            token: localStorage.getItem("accessToken"),
        }
    },
    actions: {
        async getCurrentPlayer() {
            return this.api.getCurrentPlayer(this.token);
        },
        async joinGame(gameId, name, color) {
            console.log(name, color);
            const r = await this.api.joinGame(gameId, name, color);
            this.token = r.response.headers['access-token'];
            localStorage.setItem("accessToken", this.token);
        },
        placePiece(placement) {
            // return state.api.placePiecePlacePut(placement, {tokenHeader: state.token});
        },
        revokeToken() {
            this.token = null;
            localStorage.removeItem("accessToken");
        },
    },
    getters: {
        // NOTE: Smuggling access token to server via websocket protocol header
        //       https://stackoverflow.com/questions/4361173/http-headers-in-websockets-client-api
        socket() { return new WebSocket(`ws://${this.apiLocation}/ws?access_token=${this.token}`) },
    }
});