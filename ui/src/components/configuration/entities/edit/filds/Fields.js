import React from 'react';
import './Fields.css';
import HostManager from "../../../../../HostManager";
import {sessionManager} from "../../../../../SessionManager";
import DataStore from "../../../../data/DataStore";
import FormPanel from "../../../../form/FormPanel";
import GridPanel from "../../../../grid/GridPanel";

export default class Fields extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            mode: 'viewFields',
            entityId: this.props.entity.id,
        };
        this.entityFormPanel = React.createRef();
        this.gridPanel = React.createRef();
    }

    componentDidMount() {
        this.getDataStore().load((data) => {
            if (data.success) {
                this.setState({
                    fields: data.fields,
                });
            }
        });
    }

    isTab() {
        return true;
    }

    getId() {
        return "fields";
    }

    getLabel() {
        return "Fields"
    }

    getDataStore() {
        return new DataStore({
            url: `${HostManager.amperHost()}entities/getFields`,
            requestMethod: "POST",
            parameters: {
                objectId: this.props.entity.id,
            },
        });
    }

    createHandler() {
        this.setState({
            mode: 'createField',
        });
    }

    deleteHandler() {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length >= 1) {
            const entityId = selection[0].objectId;
            const fieldIds = [];
            for (let i = 0; i < selection.length; i++) {
                fieldIds.push(selection[i].id);
            }
            fetch(`${HostManager.amperHost()}fields/deleteField`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify({
                    entityId: entityId,
                    fieldIds: fieldIds,
                })
            })
                .then(res => res.json())
                .then((result) => {
                    if (result.success) {
                        this.gridPanel.current.reload();
                    }
                })
        }
    }

    getToolItems() {
        return [
            {
                label: 'Create',
                handler: this.createHandler.bind(this)
            } , {
                label: 'Delete',
                handler: this.deleteHandler.bind(this)
            }
        ];
    }

    getGridDataModel() {
        return [{
            label: 'Label',
            key: 'label'
        },{
            label: 'Field key',
            key: 'apiName'
        },{
            label: 'Type',
            key: 'type'
        },{
            label: 'Required',
            key: 'required',
            render: function (rawValue) {
                return rawValue ? "True" : "False";
            }
        },{
            label: 'Status',
            key: 'status',
            render: function (rawValue) {
                return rawValue ? "Active" : "Inactive";
            }
        }];
    }

    getFieldsView() {
        return <GridPanel
            ref={this.gridPanel}
            selectable={true}
            selectMode={'row'}
            toolItems={this.getToolItems()}
            dataModel={this.getGridDataModel()}
            dataStore={this.getDataStore()}
        />;
    }

    onLabelChange(event, value) {
        let apiName = value.toLowerCase();
        let result = '';
        let spaceFound = false;
        for(let i = 0; i < apiName.length; i++) {
            const character = apiName.charAt(i);
            if (character.match(/[a-z]/i)) {
                result = result + (spaceFound ? character.toUpperCase(): character);
                spaceFound = false;
            } else if(character === ' ') {
                spaceFound = true;
            }
        }
        this.entityFormPanel.current.updateFieldValue('apiName', result + "__amp")
    }

    getCreateFieldsFormConfig() {
        return [
            [{
                type: "id",
                label: "Id",
                name: "entityId",
                value: this.state.entityId,
            }], [{
                type: 'text',
                label: 'Label',
                name: 'label',
                required: true,
                onChange: (event, value) => {
                    this.onLabelChange(event, value);
                },
            }], [{
                type: 'text',
                label: 'Field key',
                name: 'apiName',
                required: true,
                disabled: true,
                onKeyDown: (event) => {
                    /*var key = event.keyCode;
                    if(!((key >= 65 && key <= 90) || key == 8)) {
                        event.preventDefault();
                    }*/
                },
            }], [{
                type: 'select',
                label: 'Type',
                name: 'dataType',
                forceUpdate: true,
                required: true,
                options: [{
                    label: 'Number',
                    key: 'NUMBER',
                },{
                    label: 'Text',
                    key: 'TEXT',
                },{
                    label: "Reference",
                    key: "REFERENCE",
                },{
                    label: 'Boolean',
                    key: 'BOOLEAN',
                },{
                    label: "Date",
                    key: "DATE",
                },{
                    label: "Date Time",
                    key: "DATETIME",
                }],
                value: 'NUMBER'
            }],[{
                type: 'number',
                label: 'Text Length',
                name: 'textLength',
                dependsOnName: 'dataType',
                dependsOnValue: 'TEXT',
                required: true,
                value: 256,
            }],[{
                type: 'finder',
                label: 'Reference',
                name: 'objectReference',
                required: true,
                dependsOnName: 'dataType',
                dependsOnValue: 'REFERENCE',
                src: `${HostManager.amperHost()}entities/getEntities`,
            }],[{
                type: 'boolean',
                label: 'Required',
                name: 'required',
                required: true,
            }],[{
                type: 'boolean',
                label: 'Status',
                name: 'status',
                required: true,
            }]
        ];
    }

    onFieldCreated() {
        this.setState({
            mode: 'viewFields',
        });
    }

    getCreateFields() {
        return <FormPanel ref={this.entityFormPanel}
                          parent={this}
                          onSuccess={this.onFieldCreated.bind(this)}
                          url={`${HostManager.amperHost()}fields/createField`}
                          items={this.getCreateFieldsFormConfig()}
                          submitLabel={"Create"}>
        </FormPanel>;
    }

    getFieldsModeView() {
        const viewMode = this.state.mode;
        switch(viewMode) {
            case 'viewFields':
                return this.getFieldsView();
            case 'editField':
                break;
            case 'createField':
                return this.getCreateFields();
                break;
        }
    }

    render() {
        return (
            <div className="objectFieldsContainer">
                {this.getFieldsModeView()}
            </div>
        );
    }
}
