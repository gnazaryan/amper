import React from 'react';
import Convenience from "../help/Convenience";
import './CheckBox.css';
import '../Main.css';

export default class CheckBox extends React.Component {

	constructor(props) {
    super(props);
		this.state = {
		    value: (props.value ? props.value : ''),
            valid: true,
        };
		this.onChange = this.onChange.bind(this);
	}

	onChange(event) {
		 this.setState({value: event.target.value, valid: true});
		 if (this.props.parent) {
			 if (!this.props.parent.inputValues) {
				 this.props.parent.inputValues = {};
			 }
			 this.props.parent.inputValues[this.props.id] = event.target.checked;
		 }
        if(this.props.onChange) {
            this.props.onChange(event, event.target.checked);
        }
	}

    updateValue(value) {
        this.setState({value: value, valid: Convenience.hasValue(value)});
    }

    validate() {
        return true;
    }

	render() {
		const id = this.props.id;
		const name = this.props.name;
		const label = this.props.label;
		const type = this.props.type || "checkbox";
		const disabled = this.props.disabled || false;
		let classNames = this.state.valid ? 'checkBoxField' : 'checkBoxField invalid';
        if (this.props.required) {
            classNames += ' fieldRequired';
        }
		const style = {
		    //width: this.props.inputWidth ? (this.props.inputWidth + 'px') : '300px',
        };
		return (
			<div className="checkBoxFieldMainContainer">
				<label className="checkBoxFieldLabel" htmlFor={id}>{label}</label>
                <br/>
				<input id={id} name={this.props.name} className={classNames} style={style} value={this.state.value}
				 onChange={this.onChange} type={type} disabled={disabled}></input>
			</div>
		);
	}
}
