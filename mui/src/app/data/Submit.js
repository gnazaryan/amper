import {sessionManager} from "../../SessionManager";

export const post = (url, parameters, success, failure) => {
    fetch(url, {
        method: 'POST',
        headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
        body: JSON.stringify(parameters)
    })
    .then(res => res.json())
    .then((result) => {
        if (result && result.success) {
            if (success) {
                success(result);
            }
        } else {
            if (result.authenticated < 0) {
                sessionManager.invalidateSession();
                window.location.reload();
            } else if (failure) {
                failure(result);
            }
        }
    })
    .catch((error) => {
        if (failure) {
            failure({error: error});
        }
    });
};

export const postBlob = (url, blob, success, failure) => {
    fetch(url, {
        method: 'POST',
        headers: {sessionId: sessionManager.getSessionId()},
        body: blob
    })
    .then(res => res.json())
    .then((result) => {
        if (result && result.success) {
            if (success) {
                success(result);
            }
        } else {
            if (failure) {
                failure(result);
            }
        }
    })
    .catch((error) => {
        if (failure) {
            failure({error: error});
        }
    });
};

export const download = (url) => {
    var action = document.createElement('A');
    action.href = url;
    action.download = url.substr(url.lastIndexOf('/') + 1);
    document.body.appendChild(action);
    action.click();
    document.body.removeChild(action);
};