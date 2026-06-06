import React from 'react';
import './WestPanelItem.css';
import '../Main.css';

export default class WestPanel extends React.Component {

    constructor(props) {
        super(props);
        this.onClick = this.onClick.bind(this);
    }

    onClick(e) {
        if (this.props.parent) {
            this.props.parent.onMenuItemActivate(this.props.name)
        }
    }

    render() {
        const active = this.props.active ? "westPanelItemActive noselect pointer" : "westPanelItemInactive noselect pointer";
        return (
            <div className={active}  href='#' onClick={this.onClick}>
                {this.props.label}
            </div>
        );
    }
}
