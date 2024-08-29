
import ApiClient from "./ApiClient";
import AccessToken from "./model/AccessToken";
import PlayerProfile from "./model/PlayerProfile";

export default class DefaultApi {
 
  constructor(apiClient) {
    this.apiClient = apiClient || ApiClient.instance;
  }

  createGame(config) {
    const postBody = config;
    
    const pathParams = {};  
    const queryParams = {};
    const headerParams = {};
    const formParams = {};

    const authNames = [];
    const contentTypes = [];
    const accepts = ['application/json'];

    return this.apiClient.callApi(
      '/create', 'POST',
      pathParams, queryParams, headerParams, formParams, postBody,
      authNames, contentTypes, accepts, String, null
    );
  }
  
  joinGame(gameId, playerName, playerColor) {
    const postBody = {"name": playerName, "color": playerColor};
    
    const pathParams = {};  
    const queryParams = {"game": gameId};
    const headerParams = {};
    const formParams = {};

    const authNames = [];
    const contentTypes = [];
    const accepts = ['application/json'];
    
    return this.apiClient.callApi(
      '/join', 'POST',
      pathParams, queryParams, headerParams, formParams, postBody,
      authNames, contentTypes, accepts, AccessToken, null
    );
  }

  place(accessToken, placement) {
    const postBody = {"coords": placement};

    const pathParams = {};
    const queryParams = {};
    const headerParams = {'Access-Token': accessToken};
    const formParams = {};

    const authNames = [];
    const contentTypes = ['application/json'];
    const accepts = ['application/json'];

    return this.apiClient.callApi(
      '/place', 'PUT',
      pathParams, queryParams, headerParams, formParams, postBody,
      authNames, contentTypes, accepts, Object, null
    );
  }
};