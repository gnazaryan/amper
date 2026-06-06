import React from 'react'
import './DisplayMessage.css'

export default function DisplayMessage(props) {

    return (<div className={'displayMessage'}>
        {props.children}
    </div>);
}