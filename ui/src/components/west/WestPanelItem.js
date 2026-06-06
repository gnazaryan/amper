import React from 'react';
import './WestPanelItem.css';
import '../Main.css';

export default class WestPanelItem extends React.Component {

    constructor(props) {
        super(props);
        this.onClick = this.onClick.bind(this);
        this.state = {active: this.props.active};
    }

    onClick(e) {
        if (this.props.parent) {
            this.props.parent.onMenuItemActivate(this.props.name)
        }
        this.setState({
            active: !this.state.active,
        })
    }

    getArrowIcon() {
        let icon = "/images/right-arrow 16.png";
        if (this.state.active) {
            icon = "/images/down-arrow 16.png";
        }
        return <span className="westPanelItemArrow" width="8px" heigth="8px"><img width="12px" heigth="12px" src={icon}/></span>
    }

    render() {
        const active = this.props.active ? "westPanelItemActive noselect pointer" : "westPanelItemInactive noselect pointer";
        return (
            <div>
            <div className={active}
                 href='#' onClick={this.onClick.bind(this)}>
                <img className={'westPanelItemIcon'} width="20px" height="20px" src={this.props.iconSrc}/>
                <span className={'westPanelItemLabel'}>
                    {this.props.label}
                </span>
                {this.props.hidden ? "" : this.getArrowIcon()}
            </div>
            <div>
            </div>
                {(!this.props.hidden && this.state.active) ? this.props.children : ""}
            </div>
        );
    }
}
