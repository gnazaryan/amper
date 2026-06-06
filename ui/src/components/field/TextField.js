import React from 'react';
import {sessionManager} from "../../SessionManager";
import Convenience from "../help/Convenience";
import './Field.css';
import '../Main.css';
import Loading from "../loading/Loading";

export default class TextField extends React.Component {

	constructor(props) {
        super(props);
		this.state = {
		    value: (props.value ? props.value : ''),
            valid: true,
            remoteValid: true,
        };
		this.onChange = this.onChange.bind(this);
	}

	onChange(event) {
		 this.setState({value: event.target.value, valid: this.isValid(event.target.value)});
		 if (this.props.parent) {
			 if (!this.props.parent.inputValues) {
				 this.props.parent.inputValues = {};
			 }
			 this.props.parent.inputValues[this.props.id] = event.target.value;
		 }
        if(this.props.onChange) {
            this.props.onChange(event, event.target.value);
        }
        if (this.props.remoteValidation) {
            this.remoteValidate(event.target.value);
        }
	}

    updateValue(value) {
        this.setState({value: value, valid: Convenience.hasValue(value)});
    }

    remoteValidate(value) {
        Convenience.isRemoteValid(this.props.remoteValidation, value,(valid) => {
            this.setState({
                valid: valid,
                remoteValid: valid,
            });
        }, () => {
            this.setState({
                valid: null,
            });
        });
    }

    isValid(value) {
	    return (this.props.required !== true || Convenience.hasValue(value)) &&
            Convenience.isValid(value, this.props.validator) &&
            (this.props.remoteValidation == null || this.state.remoteValid);
    }

    validate() {
        const valid = this.isValid(this.state.value);
        this.setState({
            valid: valid,
        });
        return valid;
    }

	render() {
		const id = this.props.id;
		const name = this.props.name;
		const label = this.props.label;
		const type = this.props.type || "text";
		const disabled = this.props.disabled || false;
		let classNames = this.state.valid ? 'field' : 'field invalid';
		if (this.props.required) {
		    classNames += ' fieldRequired';
        }
        const style = {
            width: this.props.inputWidth ? (this.props.inputWidth + 'px') : '300px',
            marginBottom: this.props.marginBottom ? (this.props.marginBottom + 'px') : '0px',
        };
		return (
			<div className="fieldMainContainer">
				<label className="fieldLabel" htmlFor={id}>{label}</label>
                <br/>
				<input id={id} name={this.props.name} className={classNames} value={this.state.value}
				 onChange={this.onChange} type={type} disabled={disabled} style={style}></input>
                {this.state.valid === null ? <Loading/> : ''}
			</div>
		);
	}
}
