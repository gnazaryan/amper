import React from 'react';
import './SecurityProfilesCentralPanel.css';
import '../../components/Main.css';
import HostManager from "../../HostManager";
import {sessionManager} from "../../SessionManager";
import DataStore from "../../components/data/DataStore";
import EventRegistry from "../../components/event/EventRegistery";
import GridPanel from "../../components/grid/GridPanel";

export default class SecurityProfilesCentralPanel extends React.Component {

    constructor(props) {
        super(props);
        this.gridPanel = React.createRef();
    }

    createHandler() {
        EventRegistry.fire("menuItemChange", this, ['manageSecurityProfile']);
    }

    editHandler = () => {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length == 1) {
            EventRegistry.fire("menuItemChange", this, ['modifyEntity', {
                entityId: selection[0].id,
            }]);
        }
    }

    deleteHandler = () => {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length == 1) {
            fetch(`${HostManager.amperHost()}entities/deleteEntity?entityId=` + selection[0].id, {
                method: 'get',
                headers: {'Content-Type':'application/json', sessionId: sessionManager.getSessionId()},
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
                label: 'Add',
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
            label: 'Name',
            key: 'name'
        }, {
            label: 'Description',
            key: 'description'
        }];
    }

    getDataStore() {
        return new DataStore({
            url: `${HostManager.amperHost()}securityProfiles/getSecurityProfiles`,
            requestMethod: "POST",
            parameters: {
                'sessionId': this.props.sessionId,
                start: 0,
                limit: 50,
            },
            dataModel: this.getGridDataModel()
        });
    }

    render() {
        return (
            <div className="centralPanel">
                <GridPanel
                    ref={this.gridPanel}
                    title="Security Profiles"
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
