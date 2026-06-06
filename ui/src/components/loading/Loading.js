import React, {useRef, useEffect, useReducer} from 'react'
import './Loading.css'
import {sessionManager} from "../../SessionManager";

export default function Loading(props) {

    const loadingRef = useRef(null);
    const upperCloud = useRef(null);
    const lowerCloud = useRef(null);
    const parentMask = useRef(null);

    const sun = useRef(null);
    const initialState = {
        x: 0,
        y: 0,
    };
    const LOADING_ACTIONS = {
        MOVE_CLOUDS: 'MOVE_CLOUDS',
    };
    const reducer = (state, action) => {
        switch (action.type) {
            case LOADING_ACTIONS.MOVE_CLOUDS:
                return {
                    ...state,
                    x: action.x,
                    y: action.y,
                };
        }
    };
    const [state, dispatch] = useReducer(reducer, initialState);

    useEffect(() => {
        const parent = loadingRef.current.parentElement;
        const parentRect = parent.getBoundingClientRect();

        loadingRef.current.style.left = (parentRect.x + (parentRect.width / 2) - 50) + 'px';
        loadingRef.current.style.top = (parentRect.y + (parentRect.height / 2) - 25) + 'px';

        animate(0, 1);
        /*dispatch({
            type: LOADING_ACTIONS.MOVE_CLOUDS,
            x: state.x + 1,
        });*/

        parentMask.current.style.top = parentRect.y + 'px';
        parentMask.current.style.left = parentRect.x + 'px';
        parentMask.current.style.width = parentRect.width + 'px';
        parentMask.current.style.height = parentRect.height + 'px';
    });

    const animate = (degree, increment) => {
        setTimeout(() => {
            if (upperCloud.current && lowerCloud.current) {
                let upperOffsetLeft = upperCloud.current.offsetLeft;
                let lowerOffsetLeft = lowerCloud.current.offsetLeft;
                if (upperOffsetLeft > 100) {
                    upperOffsetLeft = -25;
                }
                if (lowerOffsetLeft > 100) {
                    lowerOffsetLeft = -25
                }
                upperCloud.current.style.left = upperOffsetLeft + 1 + 'px';
                lowerCloud.current.style.left = lowerOffsetLeft + 1 + 'px';

                if (degree > 30) {
                    increment = -1;
                } else if (degree < -30) {
                    increment = 1;
                }
                sun.current.style.webkitTransform = 'rotate(' + degree + 'deg)';
                sun.current.style.mozTransform = 'rotate(' + degree + 'deg)';
                sun.current.style.msTransform = 'rotate(' + degree + 'deg)';
                sun.current.style.oTransform = 'rotate(' + degree + 'deg)';
                sun.current.style.transform = 'rotate(' + degree + 'deg)';
                degree += increment;

                /*var fontSize = parseFloat(window.getComputedStyle(sun.current).getPropertyValue('font-size'));
                fontSize+=increment;
                if (fontSize > 35) {
                    increment = -1;
                } else if (fontSize < 34) {
                    increment = 1;
                }
                sun.current.style.fontSize = fontSize + 'px';*/
                animate(degree, increment);
            }
        }, 500);
    };

    return ([
        <div ref={parentMask} className={'parentMask'}></div>,
        <div ref={loadingRef} className={'loadingMainContainer'}>
            <div ref={sun} className={'sunLoading'}>☀</div>
            <div ref={upperCloud} className={'upperCloudLoading'}>☁</div>
            <div ref={lowerCloud} className={'lowerCloudLoading'}>☁</div>
            <div className={'pleaseWaitLoading'}>{props.label || 'Please wait...'}</div>
            {/*<div ref={hause} className={'houseLoading'}>⛪</div>
                <div className={'treeLoading'}>🌲</div>
                <div className={'tree1Loading'}>🌲</div>
                <div className={'tree2Loading'}>🌲</div>
                <div className={'cactusLoading'}>🌳</div>*/}
        </div>]);
}