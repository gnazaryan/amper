class EventRegistry {

    static on(name, callBack, context) {
        window.eventRegistry[name] = {
            callBack: callBack,
            context: context
        };
    }

    static fire(name, callee, values) {
        if(window.eventRegistry[name]) {
            let contruct = window.eventRegistry[name];
            contruct.callBack.apply(contruct.context, values);
        }
    }
}

export default EventRegistry;
