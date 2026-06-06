import React from 'react';
import './ConfigurationWestPanel.css';
import '../Main.css'
import WestPanelItem from './WestPanelItem.js'
import EventRegistry from '../event/EventRegistery.js';

export default class ConfigurationWestPanel extends React.Component {

    constructor(props) {
        super(props);
        this.state = {activeName: "entities"};
    }

    onMenuItemActivate(name, context) {
        this.setState({activeName: name});
        EventRegistry.fire("viewChange", context, [name])
    }

    render() {
        return (
            <div className="westPanel">
                <span className="section noselect">Data structure</span>
                <WestPanelItem parent={this} name="entities"
                               active={this.state.activeName === "entities"} label="Objects"></WestPanelItem>
                <span className="section noselect">View setup</span>
                <WestPanelItem parent={this} name="views"
                               active={this.state.activeName === "views"} label="Views"></WestPanelItem>
                <span className="section noselect">User management</span>
                <WestPanelItem parent={this} name="users"
                               active={this.state.activeName === "users"} label="Users"></WestPanelItem>
                <WestPanelItem parent={this} name="userPermission"
                               active={this.state.activeName === "userPermission"} label="User permission"></WestPanelItem>
                <span className="section noselect">Finite state machine</span>
                <WestPanelItem parent={this} name="entityStateCycles"
                               active={this.state.activeName === "entityStateCycles"} label="Entity state cycles"></WestPanelItem>
                <WestPanelItem parent={this} name="states"
                               active={this.state.activeName === "states"} label="States"></WestPanelItem>
                <span className="section noselect">Actions</span>
                <WestPanelItem parent={this} name="entityAction"
                               active={this.state.activeName === "entityAction"} label="Entity actions"></WestPanelItem>
                <WestPanelItem parent={this} name="stateActions"
                               active={this.state.activeName === "stateActions"} label="State actions"></WestPanelItem>
            </div>
        );
    }
}
