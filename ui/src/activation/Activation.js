import React from 'react';
import Convenience from "../components/help/Convenience";
import Loading from "../components/loading/Loading";
import Welcome from '../components/welcome/Welcome.js'
import TextField from '../components/field/TextField.js'
import Button from '../components/button/Button.js'
import '../components/Main.css';
import './Activation.css';
import {sessionManager} from "../SessionManager";

export default class Activation extends React.Component {

    constructor(props) {
        super(props);
        this.passwordField = React.createRef();
        this.passwordConfirmationField = React.createRef();
        this.state = {
            valid: true,
            loading: false,
            name: Convenience.getUrlParameterValue('name'),
        };
        this.inputValues = {};
        sessionManager.invalidateSession();
        this.props.parent.handleLogOut();
    }

    updateState(user) {
        this.props.parent.updateState(user);
    }

    activate() {
        if (!this.passwordField.current.validate() || !this.passwordConfirmationField.current.validate()
            || this.inputValues.password !== this.inputValues.confirmPassword) {
            let errorMessage;
            if (this.inputValues.password !== this.inputValues.confirmPassword ||
            !this.inputValues.password || !this.inputValues.confirmPassword) {
                errorMessage = 'Both password and confirmation password are required and must match';
            }
            this.setState({
                valid: false,
                errorMessage: errorMessage,
            });
            return;
        }
        this.setState({
            loading: true,
        });
        fetch("http://localhost:8080/users/activate", {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                activationCode: this.props.activationCode,
                ...this.inputValues
            })
        })
        .then(res => res.json())
        .then((result) => {
            if (result.success) {
                window.location.replace(window.location.origin);
            }
            this.setState({
                loading: false,
            });
        });
    }

    componentDidMount() {

    }

    render() {
        return (
            <div className="activateMainContainer">
                <div className="activateInnerContainer">
                    <Welcome name={this.state.name}/>
                    <TextField ref={this.passwordField} id="password" parent={this} label="Password" type="password" required={true} marginBottom={10}/>
                    <TextField ref={this.passwordConfirmationField} id="confirmPassword" parent={this} label="Confirm password" required={true} type="password" marginBottom={10}/>
                    <div className={'activationErrorMessage'}>{this.state.errorMessage}</div>
                    <Button onClick={this.activate.bind(this)} label="Activate"/>
                    {this.state.loading ? <Loading label={'Activating...'}/> : ''}
                </div>
            </div>
        );
    }
}
