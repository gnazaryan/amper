import Convenience from "./components/help/Convenience";

class HostManager {
    static AMPER_HOSTS = ["http://localhost:7777/", "http://localhost:7777/", "http://localhost:7777/"];

    static amperHost() {
        return HostManager.AMPER_HOSTS[Convenience.getRandomInt(0, HostManager.AMPER_HOSTS.length)];
    }
}

export default HostManager;