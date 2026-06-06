import React from 'react';
import './ui.css';
import ManageSecurityProfilePanel from "../../administration/securityProfiles/manage/ManageSecurityProfilePanel";
import SecurityProfilesCentralPanel from "../../administration/securityProfiles/SecurityProfilesCentralPanel";
import ManageUserPanel from "../../administration/users/create/ManageUserPanel";
import Menu from '../menu/Menu.js';
import WestPanel from "../west/WestPanel";
import ConfigEntitiesCentralPanel from '../configuration/entities/ConfigEntitiesCentralPanel';
import CreateEntityPanel from  '../configuration/entities/create/CreateEntityPanel'
import UsersCentralPanel from '../../administration/users/UsersCentralPanel'
import EditEntityPanel from "../configuration/entities/edit/EditEntityPanel";
import EventRegistry from '../event/EventRegistery.js';
import Navigation from '../../Navigation';

export default class UI extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
        view: 'securityProfiles',
        westActive: true,
    };
    EventRegistry.on('menuItemChange', this.updateCenterView, this);
    EventRegistry.on('viewChange', this.updateCenterView, this);
    EventRegistry.on('menuActuate', this.actuateWestMenu, this);
  }

  updateCenterView(item, options) {
    this.setState({
        view: item,
        options
    })
  }

  updateState(sessionId) {
    this.setState({sessionId: sessionId});
  }

    getCentralPanel() {
        Navigation.push(this.state.view, this.state.options);
        switch(this.state.view) {
            case 'entities':
                return <ConfigEntitiesCentralPanel sessionId={this.props.sessionId}></ConfigEntitiesCentralPanel>;
            case 'createEntity':
                return <CreateEntityPanel sessionId={this.props.sessionId}></CreateEntityPanel>;
            case 'modifyEntity':
                return <EditEntityPanel sessionId={this.props.sessionId} options={this.state.options}></EditEntityPanel>;
            case 'createUser':
                return <ManageUserPanel sessionId={this.props.sessionId}></ManageUserPanel>;
            case 'viewUser':
                return <ManageUserPanel sessionId={this.props.sessionId} user={this.state.options.user} mode={'view'}></ManageUserPanel>;
            case 'editUser':
                return <ManageUserPanel sessionId={this.props.sessionId} user={this.state.options.user} mode={'edit'}></ManageUserPanel>;
            case 'users':
                return <UsersCentralPanel sessionId={this.props.sessionId} options={this.state.options}></UsersCentralPanel>;
            case 'securityProfiles':
                return <SecurityProfilesCentralPanel sessionId={this.props.sessionId} options={this.state.options}></SecurityProfilesCentralPanel>;
            case 'manageSecurityProfile':
                return <ManageSecurityProfilePanel sessionId={this.props.sessionId}></ManageSecurityProfilePanel>;
        }
    }

  actuateWestMenu() {
      this.setState(
          {
              westActive: !this.state.westActive,
          }
      );
  }

  getWestPanel() {
      const westClass = this.state.westActive ? 'westActive' : 'westPassive';
      return <div className={westClass}>
          <WestPanel sessionId={this.props.sessionId} active={this.state.westActive}></WestPanel>
      </div>;
  }

  handleLogOut() {
    this.props.parent.handleLogOut()
  }

  render() {
    return (
      <div className="ui">
        <div className="north">
          <Menu parent={this}></Menu>
        </div>
        <div className="middle">
            {this.getWestPanel()}
            <div className="center">
                {this.getCentralPanel()}
            </div>
            {/*<div className="east">
              East Panel
            </div>*/}
        </div>
      </div>
    );
  }
}
