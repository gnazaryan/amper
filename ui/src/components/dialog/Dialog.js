import React, { useRef, useEffect } from 'react'

import './Dialog.css'
import Button from "../button/Button";
import Icon from "../icon/Icon";

export default function Dialog(props) {
    /*const initialState = {
        open: props.open,
    };
    const DIALOG_ACTIONS = {
        OPEN: 'OPEN',
        CLOSE: 'CLOSE',
    };
    const reducer = (state, action) => {
        switch (action.type) {
            case DIALOG_ACTIONS.OPEN:
                return {
                    ...state,
                    open: true,
                };
            case DIALOG_ACTIONS.CLOSE:
                return {
                    ...state,
                    open: false,
                };
        }
    };
    const [state, dispatch] = useReducer(reducer, initialState);*/
    const DEFAULT_WIDTH = 500;
    const DEFAULT_HEIGHT = 300;
    const amperDialogRef = useRef(null);
    const width = props.width ? props.width : DEFAULT_WIDTH;
    const height = props.height ? props.height : DEFAULT_HEIGHT;

    const style = {
        top: ((document.body.clientHeight / 2 - height / 2) + 'px'),
        right: ((document.body.clientWidth / 2 - width / 2) + 'px'),
        width: (width + 'px'),
        height: (height + 'px'),
        display: props.open ? 'inline-block' : 'none',
    };
    const styleContent = {
        width: (width + 'px'),
        height: ((height - 70) + 'px'),
    };
    const closeDialog = () => {
        if (props.closeHandler) {
            props.closeHandler();
        }
    };

    const getDialogBottomItems = () => {
        const result = [];
        if (props.controls) {
            for (let i = 0; i < props.controls.length; i++) {
                const control = props.controls[i];
                result.push(<Button onClick={control.handler} label={control.label}></Button>)
            }
        }
        return result;
    };
    return (<div className={'amperDialog'} ref={amperDialogRef}
                style={style}>
        <div className="pointer amperDialogTitle noselect">
            <div className={'amperDialogTitleContainer'}>{props.title}</div>
            <div>
                <Icon className={'amperDialogCloseContainer'} handler={closeDialog}
                      pointer={true} width={'16px'}
                      height={'16px'} src={'/images/close.png'}/>
            </div>
        </div>
        <div className={"amperDialogContent"} style={styleContent}>
            {props.children}
        </div>
        <div className={"amperDialogBottomBar"}>
            {getDialogBottomItems()}
        </div>
    </div>);
}