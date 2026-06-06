import React from 'react';
import Convenience from "../help/Convenience";
import './Field.css';
import '../Main.css';

export default class NumberField extends React.Component {

	constructor(props) {
    super(props);
		this.state = {
		    value: (props.value ? props.value : ''),
            valid: true,
        };
		this.onChange = this.onChange.bind(this);
	}

	onChange(event) {
		const value = event.target.value ? parseInt(event.target.value) : null;
		 this.setState({value: value, valid: value != null});
		 if (this.props.parent) {
			 if (!this.props.parent.inputValues) {
				 this.props.parent.inputValues = {};
			 }
			 this.props.parent.inputValues[this.props.id] = value;
		 }
        if(this.props.onChange) {
            this.props.onChange(event, value);
        }
	}

    updateValue(value) {
        this.setState({value: value, valid: Convenience.hasValue(value)});
    }

    validate() {
        if(this.props.required === true) {
            const valid =  Convenience.hasValue(this.state.value);
            this.setState({
                valid: valid,
            });
            return valid;
        }
        return true;
    }

    onkeypress(event) {
        var ev = event || window.event;
        if(ev.charCode < 48 || ev.charCode > 57) {
            ev.preventDefault();
            ev.stopPropagation();
        } else if(event.target.value * 10 + ev.charCode - 48 > event.target.max) {
            ev.preventDefault();
            ev.stopPropagation();
        } else {
            return true;
        }
    }

	render() {
		const id = this.props.id;
		const name = this.props.name;
		const label = this.props.label;
		const disabled = this.props.disabled || false;
		let classNames = this.state.valid ? 'field' : 'field invalid';
        if (this.props.required) {
            classNames += ' fieldRequired';
        }
        const min = this.props.min || 0;
		const max = this.props.max || 1000;
        const maxLength = (max + '').length;
        const style = {
            width: this.props.inputWidth ? (this.props.inputWidth + 'px') : '300px',
        };
		return (
			<div className="fieldMainContainer">
				<label className="fieldLabel" htmlFor={id}>{label}</label>
                <br/>
				<input id={id} name={this.props.name} className={classNames} value={this.state.value || ""}
				 onChange={this.onChange} maxLength={maxLength} onKeyPress={this.onkeypress} type={'number'}
                       disabled={disabled} min={min} max={max} style={style}></input>
			</div>
		);
	}
}
