/**
 * FastAPI
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';


class Coordinate {
    
    constructor(x, y) { 
        
        Coordinate.initialize(this, x, y);
    }

    
    static initialize(obj, x, y) { 
        obj['x'] = x;
        obj['y'] = y;
    }

    
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new Coordinate();

            if (data.hasOwnProperty('x')) {
                obj['x'] = ApiClient.convertToType(data['x'], Object);
            }
            if (data.hasOwnProperty('y')) {
                obj['y'] = ApiClient.convertToType(data['y'], Object);
            }
        }
        return obj;
    }

    
    static validateJSON(data) {
        // check to make sure all required properties are present in the JSON string
        for (const property of Coordinate.RequiredProperties) {
            if (!data[property]) {
                throw new Error("The required field `" + property + "` is not found in the JSON data: " + JSON.stringify(data));
            }
        }

        return true;
    }


}

Coordinate.RequiredProperties = ["x", "y"];


Coordinate.prototype['x'] = undefined;


Coordinate.prototype['y'] = undefined;






export default Coordinate;

