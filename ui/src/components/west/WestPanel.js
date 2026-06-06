import React from 'react';
import './WestPanel.css';
import '../Main.css'
import WestPanelItem from './WestPanelItem.js'
import WestPanelChildItem from './WestPanelChildItem'
import EventRegistry from '../event/EventRegistery.js';

export default class WestPanel extends React.Component {

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
                <WestPanelItem parent={this} name="dashboard"
                               active={false}
                               hidden={!this.props.active}
                               iconSrc="/images/dashboard 32.png"
                               label={this.props.active ? 'Dashboards' : ''}>
                    <WestPanelChildItem label="Overview"
                                        name="overview"
                                        active={'overview' === this.state.activeName}
                                        iconSrc="/images/speedometer.png"
                                        leaf={true}>

                    </WestPanelChildItem>
                </WestPanelItem>
                <WestPanelItem parent={this} name="configuration"
                               active={false}
                               hidden={!this.props.active}
                               iconSrc="/images/settings 32.png"
                               label={this.props.active ? 'Configuration' : ''}>
                    <WestPanelChildItem label="Objects"
                                        name="entities"
                                        active={'entities' === this.state.activeName}
                                        parent={this}
                                        iconSrc="/images/obj.png"
                                        leaf={true}>
                    </WestPanelChildItem>
                </WestPanelItem>
                <WestPanelItem parent={this} name="administration"
                               active={false}
                               hidden={!this.props.active}
                               iconSrc="/images/settings.png"
                               label={this.props.active ? 'Administration' : ''}>
                    <WestPanelChildItem label="Users"
                                        name="users"
                                        active={'users' === this.state.activeName}
                                        parent={this}
                                        iconSrc="/images/man.png"
                                        leaf={true}>
                    </WestPanelChildItem>
                    <WestPanelChildItem label="Security Profiles"
                                        active={'securityProfiles' === this.state.activeName}
                                        name="securityProfiles"
                                        parent={this}
                                        iconSrc="/images/shield.png"
                                        leaf={true}>
                    </WestPanelChildItem>
                </WestPanelItem>
            </div>
        );
    }
}
