import React from 'react';
import Convenience from "../help/Convenience";
import './Field.css';
import '../Main.css';

export default class ID extends React.Component {

	constructor(props) {
        super(props);
		this.state = {
		    value: (props.value ? props.value : ''),
            valid: true,
        };
	}

    validate() {
	    return true;
    }

	render() {
		const id = this.props.id;
		const name = this.props.name;
		const label = this.props.label;
		const type = this.props.type || "id";
		const disabled = this.props.disabled || true;
		const visible = this.props.visible || false;
		let classNames = this.state.valid ? 'field' : 'field invalid';
        if (this.props.required) {
            classNames += ' fieldRequired';
        }
        const style = {
		    display: visible ? 'block' : 'none'
        };
        const inputStyle = {
            width: this.props.inputWidth ? (this.props.inputWidth + 'px') : '300px',
        };
		return (
			<div className="fieldMainContainer" style={style}>
				<label className="fieldLabel" htmlFor={id}>{label}</label>
                <br/>
				<input id={id} name={this.props.name} className={classNames} value={this.state.value}
				 onChange={this.onChange} type={type} disabled={disabled} style={inputStyle}></input>
			</div>
		);
	}
}
