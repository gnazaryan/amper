import HostManager from "../../HostManager";
import {sessionManager} from "../../SessionManager";
import AmperConstants from "../util/AmperConstants";

class Convenience {

    static isSystemField(fieldKey) {
        if (AmperConstants.SYSTEM_FIELDS.ID === fieldKey ||
            AmperConstants.SYSTEM_FIELDS.NAME === fieldKey ||
            AmperConstants.SYSTEM_FIELDS.STATUS === fieldKey) {
            return true;
        }
        return false;
    }

    static hasValue(value) {
        return value !== '' && value != undefined;
    }

    static EMAIL_VAILDATOR_REGEX = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    static isValid(value, validator) {
        if (validator) {
            switch(validator) {
                case 'email':
                    return Convenience.EMAIL_VAILDATOR_REGEX.test(value);
            }
        }
        return true
    }

    static getUrlParameterValueFromQuery(url, parameterName) {
        const urlParts = url.split('?');
        if (urlParts.length == 2) {
            const urlParams = urlParts[1].split('&');
            if (urlParams.length > 0) {
                for (let i = 0; i < urlParams.length; i++) {
                    const urlParam = urlParams[i].split('=');
                    if (urlParam.length == 2) {
                        if (urlParam[0] === parameterName) {
                            return urlParam[1];
                        }
                    }
                }
            }
        }
        return null;
    }

    static getUrlParameterValue(parameterName) {
        return new URLSearchParams(window.location.search).get(parameterName);
    }

    static makeUrl(baseUrl, parameters) {
    let query = '';
    for (var key in parameters) {
      if (query !== '') {
        query = query + '&';
      }
      query = query + key + '=' +  parameters[key]
    }
    return baseUrl + '?' + query;
    }

    static containsNullOrEmpty(object, values) {
        let result = true;
        for (let i = 0; i < values.length; i++) {
          let value = values[i];
          if (!object || !object[value] || object[value] === '') {
            result = false;
            break;
          }
        }
        return result;
    }

    static REMOTE_VALID_TIMEOUT = {

    };
    static isRemoteValid(remoteValidation, value, callback, beforeCallback) {
        if (Convenience.REMOTE_VALID_TIMEOUT[remoteValidation]) {
            clearTimeout(Convenience.REMOTE_VALID_TIMEOUT[remoteValidation]);
        }
        Convenience.REMOTE_VALID_TIMEOUT[remoteValidation] =  setTimeout(() => {
            beforeCallback();
            fetch(`${HostManager.amperHost()}${remoteValidation}?value=` + value, {
                method: 'get',
                headers: {'Content-Type':'application/json', sessionId: sessionManager.getSessionId()},
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    callback(result.valid);
                } else {
                    callback(false);
                }
            });
        }, 1000);
    }

    static getRandomInt(min, max) {
        min = Math.ceil(min);
        max = Math.floor(max);
        return Math.floor(Math.random() * (max - min) + min); //The maximum is exclusive and the minimum is inclusive
    }
}

export default Convenience;
