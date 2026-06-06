import React from 'react';
import './Login.css';
import HostManager from "../../HostManager";
import Welcome from '../welcome/Welcome.js'
import TextField from '../field/TextField.js'
import Button from '../button/Button.js'
import Convenience from '../help/Convenience.js';

export default class Login extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            errorMessage: '',
        };
    }

  updateState(user) {
      this.props.parent.updateState(user);
  }

  logIn() {
    if (!Convenience.containsNullOrEmpty(this.inputValues, ['username', 'password'])) {
      this.setState({
          errorMessage: 'Both username and password are required fields',
      });
      return;
    }
    fetch(`${HostManager.amperHost()}users/login`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
          ...this.inputValues
      })
    })
    .then(res => res.json())
    .then((result) => {
        if (result) {
            if (result.success) {
              this.updateState(result.user);
            } else {
                this.setState({
                    errorMessage: result.error,
                });
            }
        } else {
            this.setState({
                errorMessage: 'Something went wrong, please contact your service provider for more details',
            });
        }
    })
  }

  render() {
    return (
      	<div className="loginMainContainer">
            <div className="loginInnerContainer">
                <Welcome/>
                <TextField id="username" parent={this} label="Username" marginBottom={10}/>
                <TextField id="password" parent={this} label="Password" type="password" marginBottom={10}/>
                <div className={'activationErrorMessage'}>{this.state.errorMessage}</div>
                <Button onClick={this.logIn.bind(this)} label="Log in"/>
            </div>
  		</div>
    );
  }
}
