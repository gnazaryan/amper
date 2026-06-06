import React from 'react'
import './RelationshipPicker.css'

export default function RelationshipPicker(props) {

    const getToolPanelItems = (toolItems) => {
        let result = [];
        if (toolItems) {
            for (let i = 0; i < toolItems.length; i++) {
                result.push(<div className="pointer noselect relationshipPanelToolBarItemInactive" onClick={toolItems[i].handler}>{toolItems[i].label}</div>);
            }
            return (
                <div className="tools">
                    {result}
                </div>
            );
        }
    }

    return (<div className={'relationshipPanel'}>
        <div className={'relationshipParentPanel'}>
            <div className={'relationshipPanelToolBar'}>{getToolPanelItems(props.parentToolItems)}</div>
            <div className={'relationshipPanelContent'}>
                {props.children[0]}
            </div>
        </div>
        <div className={'relationshipChildPanel'}>
            <div className={'relationshipPanelToolBar'}>{getToolPanelItems(props.childToolItems)}</div>
            <div className={'relationshipPanelContent'}>
                {props.children[1]}
            </div>
        </div>
    </div>);
}