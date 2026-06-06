import React from 'react';
import DataStore from "../data/DataStore";
import Convenience from "../help/Convenience";
import './Field.css';
import './Finder.css';
import '../Main.css';

export default class Finder extends React.Component {

	constructor(props) {
    super(props);
		this.state = {
		    value: (props.value ? props.value : ''),
            hideItems: props.hideItems,
            hideProperty: props.hideProperty,
            valid: true,
            finderFocusIn: false,
        };
        this.finderInputReference = React.createRef();
		this.onChange = this.onChange.bind(this);
	}

	onChange(event) {

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

    onFocusIn() {
        this.setState({
            finderFocusIn: true,
        });
        if (!this.state.data) {
            this.getDataStore().load((result) => {
                if (result.success) {
                    this.setState({
                        data: result.data,
                    });
                }
            });
        }
    }

    onFocusOut(event) {
	    setTimeout(()=> {
            this.setState({
                finderFocusIn: false,
            });
        }, 200)
    }

    getDataStore() {
        return new DataStore({
            url: this.props.src,
            requestMethod: "POST",
            parameters: {
                ...this.props.parameters,
                'sessionId': this.props.sessionId,
                start: 0,
                limit: 50,
            },
        });
    }

    onItemClick(event) {
        const id = parseInt(event.target.id);
        const apiName = event.target.getAttribute("apiName");
        const title = event.target.innerHTML;
        this.setState({
            value: {
                id: id,
                apiName: apiName,
                title: title,
            },
            finderFocusIn: false,
        });
        if (this.props.parent) {
            if (!this.props.parent.inputValues) {
                this.props.parent.inputValues = {};
            }
            this.props.parent.inputValues[this.props.id] = id;
        }
        if(this.props.onChange) {
            this.props.onChange(event, id);
        }
    }

    isHidden(item) {
	    if (this.state.hideItems) {
            for (let i = 0; i < this.state.hideItems.length; i++) {
                if (this.state.hideItems[i][this.state.hideProperty] === item.id) {
                    return true;
                }
            }
        }
	    return false;
    }

    getFinderItems() {
	    let result = [];
	    if (this.state.data) {
	        for (let i = 0; i < this.state.data.length; i++) {
	            let item = this.state.data[i];
	            if (!this.isHidden(item)) {
                    result.push(<div name={this.props.name} onClick={this.onItemClick.bind(this)} id={item.id}
                                     apiName={item.apiName}
                                     className={'finderWindowItem'}>{item.title || item.label || item.name}</div>);
                }
            }
        }
	    return result;
    }

	render() {
		const id = this.props.id;
		const name = this.props.name;
		const label = this.props.label;
		const type = this.props.type || "text";
		const disabled = this.props.disabled || false;
		let classNames = this.state.valid ? 'field ' : 'field invalid ';
        if (this.props.required) {
            classNames += ' fieldRequired';
        }
		const classNames1 = this.state.finderFocusIn ? 'finderVisible finderWindow ' : 'finderInvisible finderWindow ';
        const style = {
            width: this.props.inputWidth ? (this.props.inputWidth + 'px') : '300px',
        };
        const styleWindow = {
            width: this.props.inputWidth ? ((this.props.inputWidth) + 'px') : '300px',
        };
		return (
			<div className="fieldMainContainer">
				<label className="fieldLabel" htmlFor={id}>{label}</label>
                <br/>
                <input id={id} autocomplete="off" ref={this.finderInputReference} name={this.props.name} className={classNames}
                       value={this.state.value.title} id={this.state.value.id} apiName={this.state.value.apiName}
                       onChange={this.onChange} onBlur={this.onFocusOut.bind(this)} onFocus={this.onFocusIn.bind(this)}
                       type={type} disabled={disabled} style={style}>
                </input>
                <div className={classNames1} style={styleWindow}>
                    {this.getFinderItems()}
                </div>
			</div>
		);
	}
}
