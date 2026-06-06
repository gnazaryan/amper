import React from 'react';
import './Icon.css';

export default class Icon extends React.Component {

	constructor(props) {
        super(props);
        this.tooltip = React.createRef();
        this.image = React.createRef();
	}

	render() {
		const text = this.props.text;
		const src = this.props.src;
		const width = this.props.width || "12px";
		const height = this.props.height || "12px";
		let classNames = '';
		if (this.props.pointer) {
		    classNames = 'iconPointer';
        }
		if (this.props.className) {
		    classNames+=(' ' + this.props.className);
        }
		return (
            <img ref={this.image} onClick={this.props.handler} className={classNames} width={width} heigth={height} src={src}/>
		);
	}
}
