import React from 'react';
import './Menu.css';
import '../Main.css';
import Button from "../button/Button";
import EventRegistry from '../event/EventRegistery.js';
import {sessionManager} from "../../SessionManager";
import Icon from '../icon/Icon'

export default class Menu extends React.Component {

  constructor(props) {
      super(props);
      this.state = {active: true, menuOpen: false};
      this.profileMenu = React.createRef();
      this.profilePicture = React.createRef();
      this.profileDropDown = React.createRef();
      this.notificationWindow = React.createRef();
      this.fullNameContainer = React.createRef();
      this.emailConteiner = React.createRef();
      this.logoutConteiner = React.createRef();
      this.logOutButton = React.createRef();
      this.menuITems = [this.profileMenu, this.profilePicture,
          this.profileDropDown, this.notificationWindow,
          this.fullNameContainer, this.emailConteiner,
            this.logoutConteiner, this.logOutButton];
  }


  onMenuItemActivate(name, context) {
    this.setState({activeName: name});
    EventRegistry.fire("menuItemChange", [name]);
  }

    startAnimation(callback) {
        requestAnimationFrame(() => {
            requestAnimationFrame(() => {
                callback();
            });
        });
    }

  menuActuatorClickHandler(event) {
      this.startAnimation(() => {
          this.setState({
              animate: true,
              active: !this.state.active,
          })
      });
      EventRegistry.fire("menuActuate", []);
  }

  getLogoActuatorContent() {
      if(this.state.active) {
          return '☁ amper';
      } else {
          return '☁';
      }
  }

  openProfileMenu(event) {
    this.setState({
        menuOpen: !this.state.menuOpen,
    })
  }

  closeProfileMenu(event) {
    if (this.state.menuOpen) {
        let inProfileMenuClicked = false;
        for (let i = 0; i < this.menuITems.length; i++) {
            const item = this.menuITems[i];
            if (item.current == event.target || (item.current.image && item.current.image.current == event.target)) {
                inProfileMenuClicked = true;
            }
        }
        if (!inProfileMenuClicked) {
            this.setState({
                menuOpen: false,
            });
        }
    }
  }

    logOut() {
        sessionManager.invalidateSession();
        this.props.parent.handleLogOut();
    }

  render() {
      const logoActuatorClassName = this.state.active ? 'logoActuatorActive' : 'logoActuatorPassive';
      const user = sessionManager.getUser();
      document.body.onclick = this.closeProfileMenu.bind(this);
      return (
      <div className="menu">
          <div className={logoActuatorClassName}>
              {this.getLogoActuatorContent()}
          </div>
          <div className={this.state.active ? 'menuLeftSectionPassive' : 'menuLeftSectionActive'}>
              <div onClick={this.menuActuatorClickHandler.bind(this)} className={'menuActuator'}>&#9776;</div>
          </div>
          <div ref={this.profileMenu} className={'menuUserLoginSection noselect'} onClick={this.openProfileMenu.bind(this)}>
              <Icon ref={this.profilePicture} className={'profilePictureIcon'} pointer={true} width={'32px'} height={'32px'} src={'/images/user 32.png'}></Icon>
              <Icon ref={this.profileDropDown} className={'loginDropdownIcon'} pointer={true} width={'12px'} height={'12px'} src={'/images/white-down-arrow 16.png'}></Icon>
          </div>
          <div ref={this.notificationWindow} className={'loginProfileMenu'} style={{display : (this.state.menuOpen ? 'block' : 'none')}}>
              <div ref={this.fullNameContainer} className={'userFullName'}>
                  {user.firstName + ' ' + user.lastName}
              </div>
              <div ref={this.emailConteiner} className={'userFullName'}>
                  {user.email}
              </div>
              <div ref={this.logoutConteiner} className={'logOutButtonContainer'}>
                <Button ref={this.logOutButton} onClick={this.logOut.bind(this)} label="Log out"/>
              </div>
          </div>
      </div>
    );
  }
}
