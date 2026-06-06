import React from 'react';
import './Button.css';
import '../Main.css'

export default class Button extends React.Component {

	constructor(props) {
      super(props);
			this.onClick = this.onClick.bind(this);
	}

	onClick(e) {
		if (this.props.onClick) {
			this.props.onClick(e);
		} else {
			alert("the click is not handled for " + this.props.label);
		}
	}

	render() {
		const id = this.props.id;
		const label = this.props.label;
		return (
			<button id={id} className="button noselect pointer" onClick={this.onClick} type="button">{label}</button>
		);
	}
}
