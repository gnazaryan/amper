import {sessionManager} from "../../SessionManager";

class DataStore {

    constructor(configuration) {
        this.configuration = configuration;
    }

    validate() {
        if (!this.configuration || !this.configuration.url) {
            throw new Error('invalid configuration supplied for Data Store');
        }
    }

    load(callBack) {
        if (this.configuration.data) {
            callBack({
                data: this.configuration.data,
                success: true,
            })
        } else {
            fetch(this.configuration.url, {
                method: (this.configuration.requestMethod ? this.configuration.requestMethod : 'get'),
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Access-Control-Allow-Origin':'*',
                    sessionId: sessionManager.getSessionId()
                },
                body: JSON.stringify(this.configuration.parameters),
                mode:'cors'
            })
            .then(res => res.json())
            .then((result) => {
                if (result.authenticated < 0) {
                    sessionManager.invalidateSession();
                    window.location.reload();
                }
                callBack(result);
            }).catch(err => {
                console.debug("Error in fetch", err);
            });
        }
    }
}

export default DataStore;