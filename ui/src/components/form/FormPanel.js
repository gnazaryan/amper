import React from 'react';
import './FormPanel.css';
import '../Main.css';
import Navigation from "../../Navigation";
import Button from '../button/Button';
import CheckBox from "../field/CheckBox";
import ImageField from "../field/ImageField";
import NumberField from "../field/NumberField";
import TextField from '../field/TextField';
import Finder from '../field/Finder';
import Select from '../field/Select';
import ID from '../field/ID';
import {sessionManager} from "../../SessionManager";
import Loading from "../loading/Loading";

export default class FormPanel extends React.Component {

    constructor(props) {
        super(props);
        this.myRef = [];
        this.values = {};
        this.state = {
            valid: true,
            loading: false,
            mode: props.mode || 'create',
        };
        this.formPanelError = React.createRef();
    }

    getTitle() {
        if(this.props.title) {
            return (
                <div className="formTitle">
                    {this.props.title}
                </div>
            );
        }
    }

    updateFieldValue(name, value) {
        let counter = 0;
        for (let i = 0; i < this.props.items.length; i++) {
            if (Array.isArray(this.props.items[i])) {
                for (let l = 0; l < this.props.items[i].length; l++) {
                    if(this.props.items[i][l].name === name) {
                        this.myRef[counter].updateValue(value);
                        this.values[name] = value;
                        break;
                    }
                }
            } else {
                if(this.props.items[i].name === name) {
                    this.myRef[i].updateValue(value);
                    this.values[name] = value;
                    break;
                }
            }
            counter++;
        }
    }

    getInputField(index, item) {
        let result = null;
        if (item.dependsOnName && this.values[item.dependsOnName] !== item.dependsOnValue) {
            return;
        }
        let changeHandler = this.handleChange;
        if (item.forceUpdate) {
          changeHandler = this.handleChangeForceUpdate;
        }
        switch (item.type) {
            case 'finder':
                result = <Finder
                    id={item.id || item.name}
                    name={item.name}
                    src={item.src}
                    ref={(ref) => { this.myRef[index] = ref; return true; }}
                    label={item.label}
                    onChange={changeHandler.bind(this)}
                    disabled={item.disabled === true || this.state.mode === 'view'}
                    required={item.required === true}
                    inputWidth={item.inputWidth}
                    parameters={item.parameters}
                    hideItems={item.hideItems}
                    hideProperty={item.hideProperty}
                    value={item.value}
                    mode={this.state.mode}>
                </Finder>;
                if (!this.values[item.name] && item.value) {
                    this.values[item.name] = item.value.id;
                }
                break;
            case 'iamge':
                result = <ImageField
                    id={item.id || item.name}
                    name={item.name}
                    label={item.label}
                    value={item.value}
                    width={item.width}
                    height={item.height}
                    src={item.src}
                    disabled={item.disabled === true || this.state.mode === 'view'}
                    onChange={changeHandler.bind(this)}
                    mode={this.state.mode}
                />;
                if (!this.values[item.name] && item.value) {
                    this.values[item.name] = item.value;
                }
                break;
            case 'text':
                result = <TextField
                    id={item.id || item.name}
                    inputWidth={item.inputWidth}
                    name={item.name}
                    ref={(ref) => { this.myRef[index] = ref; return true; }}
                    label={item.label}
                    onChange={changeHandler.bind(this)}
                    disabled={item.disabled === true || this.state.mode === 'view' || (this.state.mode === 'edit' && item.editable === false)}
                    required={item.required === true}
                    validator={item.validator}
                    remoteValidation={item.remoteValidation}
                    value={item.value}
                    mode={this.state.mode}>
                         </TextField>;
                if (!this.values[item.name] && item.value) {
                    this.values[item.name] = item.value;
                }
                break;
            case 'number':
                result = <NumberField
                    id={item.id || item.name}
                    name={item.name}
                    ref={(ref) => { this.myRef[index] = ref; return true; }}
                    label={item.label}
                    onChange={changeHandler.bind(this)}
                    disabled={item.disabled === true || this.state.mode === 'view'}
                    required={item.required === true}
                    inputWidth={item.inputWidth}
                    value={item.value}
                    mode={this.state.mode}/>
                if (!this.values[item.name] && item.value) {
                    this.values[item.name] = item.value;
                }
                break;
            case 'boolean':
                result = <CheckBox
                    id={item.id || item.name}
                    name={item.name}
                    ref={(ref) => { this.myRef[index] = ref; return true; }}
                    label={item.label}
                    onChange={changeHandler.bind(this)}
                    disabled={item.disabled === true || this.state.mode === 'view'}
                    required={item.required === true}
                    inputWidth={item.inputWidth}
                    value={item.value}>
                </CheckBox>;
                if (!this.values[item.name] && item.value) {
                    this.values[item.name] = item.value;
                }
                break;
            case 'select':
                result = <Select
                    id={item.id || item.name}
                    name={item.name}
                    ref={(ref) => { this.myRef[index] = ref; return true; }}
                    label={item.label}
                    keyField={item.keyField}
                    onChange={changeHandler.bind(this)}
                    disabled={item.disabled === true || this.state.mode === 'view'}
                    required={item.required === true}
                    options={item.options}
                    inputWidth={item.inputWidth}
                    value={item.value}
                    valueFormatter={item.valueFormatter}>
                </Select>;
                if (!this.values[item.name] && item.value) {
                    this.values[item.name] = item.value;
                }
                break;
            case 'id':
                result = <ID
                    id={item.id || item.name}
                    name={item.name}
                    ref={(ref) => { this.myRef[index] = ref; return true; }}
                    label={item.label}
                    disabled={item.disabled === true || this.state.mode === 'view'}
                    required={item.required === true}
                    visible={false}
                    inputWidth={item.inputWidth}
                    value={item.value}>
                </ID>;
                if (!this.values[item.name] && item.value) {
                    this.values[item.name] = item.value;
                }
                break;
        }
        return result;
    }

    getFormContent() {
        let result = [];
        if (this.props.items) {
            let f = 0;
            for (let l = 0; l < this.props.items.length; l++) {
                let itemSet = this.props.items[l];
                const row = [];
                for (let i = 0; i < itemSet.length; i++) {
                    const inputField = this.getInputField(f, itemSet[i]);
                    f++;
                    if (inputField) {
                        row.push(
                            <td className="formPanelColumn" rowspan={itemSet[i].rowSpan || 1}>
                                {inputField}
                            </td>
                        );
                    }
                }
                result.push(<tr>
                        {row}
                    </tr>
                );
            }
        }
        result.push(<tr>
            <td>
                <div className="formPanelError" ref={this.formPanelError}></div>
            </td>
        </tr>);
        return (
            <table className="tablePanel">{result}</table>
        );
    }

    handleKeyDown(event) {
        for (let i = 0; i < this.props.items.length; i++) {
            const item = this.props.items[i];
            if (item.name === event.target.name && item.onKeyDown) {
                return item.onKeyDown(event);
            }
        }
    }

    handleChangeForceUpdate(event, value) {
      this.handleChange(event, value, true);
    }

    handleChange(event, value, forceUpdate) {
        const name = event.target.getAttribute('name');
        this.values[name] = value;
        for (let i = 0; i < this.props.items.length; i++) {
            const item = this.props.items[i];
            if (item && Array.isArray(item)) {
                for (let l = 0; l < item.length; l++) {
                    if (item[l].name === name && item[l].onChange) {
                        item[l].onChange(event, value);
                    }
                }
            } else {
                if (item.name === name && item.onChange) {
                    item.onChange(event, value);
                }
            }
        }
        if (this.props.onChange) {
            this.props.onChange(this.values);
        }
        if (forceUpdate) {
          this.forceUpdate();
        }
    }

    valid() {
        let result = true;
        for (let i = 0; i < this.props.items.length; i++) {
           const input = this.myRef[i];
           if (input && !input.validate()) {
               result = false;
           }
        }
        return result;
    }

    edit() {
        this.setState({
            mode: 'edit',
        });
    }

    getValues() {
        return this.values;
    }

    submit() {
        if(this.valid()) {
            this.setState({
                loading: true,
            });

            fetch(this.state.mode === 'edit' ? this.props.editUrl : this.props.url, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(this.values)
            })
                .then(res => res.json())
                .then((result) => {
                    if (result.success) {
                        if(this.props.onSuccess) {
                            this.setState({
                                loading: false,
                            });
                            this.props.onSuccess();
                        }
                    } else {
                        this.formPanelError.current.innerHTML = result.error;
                    }
                })
        } else {
            this.formPanelError.current.innerHTML = "Please fill in all required fields";
        }
    }

    getSubmitButton() {
        if (this.state.mode === 'view') {
            return <Button onClick={this.edit.bind(this)} label={'Edit'}></Button>;
        } if (this.props.submitLabel) {
            return <Button onClick={this.submit.bind(this)} label={this.props.submitLabel}></Button>;
        }
    }

    cancel() {
        Navigation.back();
    }

    getCancelButton() {
        if (this.state.mode != 'view' && this.props.submitLabel) {
            return <Button onClick={this.cancel.bind(this)} label={'Cancel'}></Button>;
        }
    }

    render() {
        return (
            <div className="formPanel">
                {this.getTitle()}
                <form className="formPanel">
                    {this.getFormContent()}
                </form>
                {this.getSubmitButton()}
                {this.getCancelButton()}
                {this.state.loading ? <Loading label={'Registering...'}/> : ''}
            </div>
        );
    }
}
