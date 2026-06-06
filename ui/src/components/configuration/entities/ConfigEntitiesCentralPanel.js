import React from 'react';
import './ConfigEntitiesCentralPanel.css';
import '../../Main.css';
import HostManager from "../../../HostManager";
import {sessionManager} from "../../../SessionManager";
import DataStore from "../../data/DataStore";
import EventRegistry from "../../event/EventRegistery";
import GridPanel from "../../grid/GridPanel";
import AmperConstatns from "../../util/AmperConstants";

export default class ConfigEntitiesCentralPanel extends React.Component {

    constructor(props) {
        super(props);
        this.gridPanel = React.createRef();
    }

    createHandler() {
        EventRegistry.fire("menuItemChange", this, ['createEntity']);
    }

    editHandler = () => {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length == 1) {
            EventRegistry.fire("menuItemChange", this, ['modifyEntity', {
                entity: selection[0],
            }]);
        }
    }

    deleteHandler = () => {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length == 1) {
            fetch(`${HostManager.amperHost()}entities/deleteEntity`, {
                method: 'POST',
                headers: {'Content-Type':'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify({
                  entityId: selection[0].id,
                }),
            })
                .then(res => res.json())
                .then((result) => {
                    if (result.success) {
                        this.gridPanel.current.reload();
                    }
                })
        }
    }

    getToolItems() {
        return [
            {
                label: 'Create',
                handler: this.createHandler
            }, {
                label: 'Edit',
                handler: this.editHandler
            }, {
                label: 'Delete',
                handler: this.deleteHandler
            }

        ];
    }

    getGridDataModel() {
        return [{
            label: 'Title',
            key: 'title'
        },{
            label: 'Title plural',
            key: 'titlePlural'
        },{
            label: 'API name',
            key: 'apiName'
        }];
    }

    getDataStore() {
        return new DataStore({
            url: `${HostManager.amperHost()}entities/getEntities`,
            requestMethod: "POST",
            parameters: {
                sessionId: sessionManager.getSessionId(),
                start: 0,
                limit: AmperConstatns.INTEGER.MAX_VALUE
            },
            dataModel: this.getGridDataModel()
        });
    }

    render() {
        return (
            <div className="centralPanel">
                <GridPanel
                    ref={this.gridPanel}
                    title="Objects"
                    selectMode={'row'}
                    selectable={true}
                    toolItems={this.getToolItems()}
                    dataModel={this.getGridDataModel()}
                    dataStore={this.getDataStore()}
                ></GridPanel>
            </div>
        );
    }
}
