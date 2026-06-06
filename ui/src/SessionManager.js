class SessionManager {

    constructor() {
        this.initUsers();
    }

    initUsers() {
        let cookies = document.cookie.split(';');
        this.user = {};
        if (cookies.length > 0) {
            for (let i = 0; i < cookies.length; i++) {
                let cookie = cookies[i].split('=');
                if (cookie.length = 2) {
                    this.user[cookie[0].trim()] = cookie[1];
                }
            }
        }
    }
    setUser(user) {
        var date = new Date();
        date.setTime(date.getTime() + (3*60*60*1000));
        for (let key in user) {
            document.cookie = key + "=" + user[key] + "; expires=" + date.toGMTString();
        }
        this.initUsers();
    }

    invalidateSession() {
        for (let key in this.user) {
            document.cookie = key + "=" + this.user[key] + "; expires=" + new Date().toGMTString();
        }
    }

    getSessionId() {
        return this.user.sessionId;
    }

    getUser() {
        return this.user;
    }
}

export let sessionManager = new SessionManager();