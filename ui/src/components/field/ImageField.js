import React, {useReducer, useRef, useEffect} from 'react'
import './Field.css'
import './imageField.css'
import HostManager from "../../HostManager";
import {sessionManager} from "../../SessionManager";
import Convenience from "../help/Convenience";
import Loading from "../loading/Loading";

export default function ImageField(props) {

    const imageRef = useRef(null);
    const inputRef = useRef(null);
    const initialState = {
        loading: props.src != null,
        value: null,
        src: props.src,
    };
    const IMAGE_FIELD_ACTIONS = {
        IMAGE_LOADED: 'IMAGE_LOADED',
    };
    const reducer = (state, action) => {
        switch (action.type) {
            case IMAGE_FIELD_ACTIONS.IMAGE_LOADED:
                return {
                    ...state,
                    loading: action.loading,
                    value: action.value,
                    src: null,
                };
        }
    };
    const [state, dispatch] = useReducer(reducer, initialState);
    const style = {
        width: (props.width || 200) + 'px',
        height: (props.height || 200) + 'px',
    };
    const browsStyle = {
        width: (props.width || 200) + 'px',
    };

    const loadImage = (url, fileName) => {
        fetch(`${HostManager.amperHost()}${url}`, {
            method: 'get',
            headers: {sessionId: sessionManager.getSessionId()},
        })
        .then((result) => {
            result.blob().then(imageData => {
                setTimeout(() => {
                    if (imageData.size > 0) {
                        imageRef.current.src = URL.createObjectURL(imageData);
                    }
                    if (fileName) {
                        dispatch({
                            type: IMAGE_FIELD_ACTIONS.IMAGE_LOADED,
                            loading: false,
                            value: fileName,
                        });
                        if (props.onChange && inputRef.current) {
                            props.onChange({
                                target: inputRef.current
                            }, fileName);
                        }
                    }
                }, 500);
            });
        });
    };

    useEffect(() => {
        if (state.src) {
            const fileName = Convenience.getUrlParameterValueFromQuery(state.src, 'fileName');
            if (fileName) {
                loadImage(state.src, fileName);
            }
        }
    });

    const onChange = (event) => {
        if (event.target.files && event.target.files.length != 0) {
            dispatch({
                type: IMAGE_FIELD_ACTIONS.IMAGE_LOADED,
                loading: true,
                value: null,
            });
            const formData = new FormData();
            formData.append('file', event.target.files[0]);
            fetch(`${HostManager.amperHost()}files/upload`, {
                method: 'POST',
                headers: {sessionId: sessionManager.getSessionId()},
                body: formData,
            }).then(res => res.json())
              .then((response) => {
                  if (response.success) {
                      loadImage(Convenience.makeUrl('files/download', {
                          fileName: response.fileName,
                      }), response.fileName);
                  }
            })
        }
    };

    const getInputField = () => {
        if (!props.disabled) {
            return <input ref={inputRef} onChange={onChange} id={props.id} disabled={disabled}
                          name={props.name} className={'field imageField'} style={browsStyle} type="file" accept="image/*"/>;
        }
    };

    const disabled = props.disabled || false;
    return (<div className={'imageFieldMainContainer'}>
        <div className={"profileImage"} style={style}>
            <img ref={imageRef} className={'profileImage'} style={style} src={'/images/user 512.png'}/>
            {state.loading ? <Loading label={'Uploading...'}/> : ''}
        </div>
        <div>
            {getInputField()}
        </div>
    </div>);
}