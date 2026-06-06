import React from 'react';
import './MenuItem.css';

export default class MenuItem extends React.Component {

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
    const active = this.props.active ? "active " : "";
    const style = {
        fontSize: this.props.fontSize ? this.props.fontSize : '14px',
        width: this.props.width ? this.props.width : '100px',
        color: this.props.width ? this.props.color : 'white',
    };

    return (
      <div className={active} style={style} href='#' onClick={this.onClick}>
        {this.props.label}
      </div>
    );
  }
}
