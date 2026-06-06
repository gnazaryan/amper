import { breadcrumbs } from "./app/center/Breadcrambs";
import { parseBoolean } from "./app/util/BooleanUtil";
import { gettStoreValue } from "./app/amper/Instruments";

class SessionManager {

    constructor() {
        this.init();
    }

    init() {
        this.initUsers();
        this.initSettings();
    } 

    initUsers() {
        let cookies = document.cookie.split(';');
        this.user = {};
        if (cookies.length > 0) {
            for (let i = 0; i < cookies.length; i++) {
                let cookie = cookies[i].split('=');
                if (cookie.length = 2) {
                    const key = cookie[0].trim();
                    if (key.startsWith('user_')) {
                        try {
                            this.user[key.substring(5, key.length)] = JSON.parse(decodeURIComponent(cookie[1]));
                        } catch(e) {
                        }
                    }
                }
            }
        }
        breadcrumbs.profile.label = this.user.firstName + ' ' + this.user.lastName;
        //breadcrumbs.profile.path = '/' + this.user.username;
        //breadcrumbs.alternativePaths[this.user.username] = breadcrumbs.profile.key;
        //breadcrumbs.profile.settings.path = breadcrumbs.profile.path + '/' + breadcrumbs.profile.settings.key
    }

    initSettings() {
        let cookies = document.cookie.split(';');
        this.settings = {};
        if (cookies.length > 0) {
            for (let i = 0; i < cookies.length; i++) {
                let cookie = cookies[i].split('=');
                if (cookie.length = 2) {
                    const key = cookie[0].trim();
                    if (key.startsWith('setting_')) {
                        try {
                            this.settings[key.substring(8, key.length)] = JSON.parse(decodeURIComponent(cookie[1]));
                        } catch(e) {
                            
                        }
                    }
                }
            }
        }
    }
    getExpireyDate() {
        var date = new Date();
        date.setTime(date.getTime() + (3*60*60*1000));
        return date;
    }
    setUser(user) {
        var date = this.getExpireyDate();
        for (let key in user) {
            this.setCookie("user_" + key, user[key])
        }
        this.initUsers();
    }
    setSettings(settings) {
        for (let key in settings) {
            this.setSetting(key, settings[key]);
        }
        this.initSettings();
    }
    setSetting(key, value) {
        this.setCookie("setting_" + key, value);
        this.settings[key] = value;
    }
    setCookie(key, value) {
        var date = this.getExpireyDate();
        const expired =  `${key}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`;
        document.cookie = expired;
        const cookie = `${key}=${encodeURIComponent(JSON.stringify(value))}; expires=${date.toGMTString()}; path=/`;
        document.cookie = cookie;
    }
    invalidateSession() {
        const cookies = document.cookie.split(";");
        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i];
            const eqPos = cookie.indexOf("=");
            const name = eqPos > -1 ? cookie.substring(0, eqPos) : cookie;
            document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
        }
    }

    getSessionId() {
        return this.user.sessionId;
    }

    getUser() {
        const photo = gettStoreValue('user_photo');
        return {
            ...this.user,
            photo,
        };
    }

    getSettings() {
        return this.settings;
    }

    isExpanded() {
        return parseBoolean(this.settings['expanded']);
    }
}

export let sessionManager = new SessionManager();