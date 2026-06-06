import Convenience from "./app/help/Convenience";
import { sessionManager } from "./SessionManager";

class HostManager {

    static AMPER_HOSTS_MAP = {
        1: 'dev.amper.cloud:7777'
    };

    /*
     * This method is designed to return a host address by the host given id
     */
    static amperHostById(instanceId) {
        return 'http://' + HostManager.AMPER_HOSTS_MAP[instanceId] + '/';
    }

    /*
     * This function must be used whenever it doesn't matter which host we are 
     * refering too, normally request which are not managing file system based data
     */
    static amperHost() {
        const amperHosts = Object.values(HostManager.AMPER_HOSTS_MAP);
        return 'http://' + amperHosts[Convenience.getRandomInt(0, amperHosts.length)] + '/';
    }

    /*
     * This function must be used whenever we are trying to reach the 
     * user assigned host address, normally requests managing file system based data
     */
    static myHost() {
        const user = sessionManager.getUser();
        return 'http://' + HostManager.AMPER_HOSTS_MAP[user.amperId] + '/';
    }
}

export default HostManager;