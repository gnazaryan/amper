import React from 'react';
import Convenience from "../help/Convenience";
import './Field.css';
import './Select.css';
import '../Main.css';

export default class TextField extends React.Component {

	constructor(props) {
    super(props);
		this.state = {
		    value: (props.value ? props.value : ''),
			valueFormatter: props.valueFormatter,
            valid: true,
            options: (props.options ? props.options: []),
        };
		this.onChange = this.onChange.bind(this);
	}

	onChange(event) {
		const value = this.state.valueFormatter ? this.state.valueFormatter(event.target.value) : event.target.value;
		 this.setState({
             value: value,
             valid: Convenience.hasValue(value),
		 });

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

    getOptions() {
	    const result = [];
	    const keyField = this.props.keyField ? this.props.keyField : 'key';
	    if (this.state.options) {
	        for (let i = 0; i < this.state.options.length; i++) {
	            let option = this.state.options[i];
                result.push(
                    <option value={option[keyField]}>{option.label}</option>
                )
            }
        }
	    return result;
    }

    validate() {
			return Convenience.hasValue(this.state.value)
    }

	render() {
		const id = this.props.id;
		const name = this.props.name;
		const label = this.props.label;
		const disabled = this.props.disabled || false;
		let classNames = this.state.valid ? 'field fieldSelect' : 'invalid field fieldSelect';
        if (this.props.required) {
            classNames += ' fieldRequired';
        }
        const style = {
            width: this.props.inputWidth ? (this.props.inputWidth + 'px') : '300px',
        };
		return (
			<div className="fieldMainContainer">
				<label className="fieldLabel" htmlFor={id}>{label}</label>
                <br/>
                <select id={id} name={this.props.name} className={classNames}
                        onChange={this.onChange} disabled={disabled} style={style}>
                    {this.getOptions()}
                </select>
			</div>
		);
	}
}
