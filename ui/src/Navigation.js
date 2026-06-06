import EventRegistry from "./components/event/EventRegistery";

class Navigation {

    static HISTORY = [];

    static push(view, args) {
        Navigation.HISTORY.push({
            view: view,
            args: args,
        })
    }

    static back() {
        if (Navigation.HISTORY.length > 0) {
            Navigation.HISTORY.pop();
            const lastItem = Navigation.HISTORY[Navigation.HISTORY.length - 1];
            EventRegistry.fire("menuItemChange", this, [lastItem.view, lastItem.args]);
        }
    }
}

export default Navigation;