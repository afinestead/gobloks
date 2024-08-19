
import ApiClient from "./ApiClient";
import AccessToken from "./model/AccessToken";
import PlayerProfile from "./model/PlayerProfile";

export default class DefaultApi {
 
  constructor(apiClient) {
    this.apiClient = apiClient || ApiClient.instance;
  }

  getCurrentPlayer(accessToken) {
      let postBody = null;

      let pathParams = {};
      let queryParams = {};
      let headerParams = {'Access-Token': accessToken};
      let formParams = {};

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];

      return this.apiClient.callApi(
        '/player', 'GET',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, PlayerProfile, null
      );
  }

  joinGame(gameId, playerName, playerColor) {
    const postBody = {"name": playerName, "color": playerColor};
    
    const pathParams = {};  
    const queryParams = {"game": gameId};
    const headerParams = {};
    const formParams = {};

    let authNames = [];
    let contentTypes = [];
    let accepts = ['application/json'];
    
    return this.apiClient.callApi(
      '/join', 'POST',
      pathParams, queryParams, headerParams, formParams, postBody,
      authNames, contentTypes, accepts, AccessToken, null
    );
  }
};