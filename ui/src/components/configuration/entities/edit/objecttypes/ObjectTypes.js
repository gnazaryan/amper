import React from 'react';
import HostManager from "../../../../../HostManager";
import {sessionManager} from "../../../../../SessionManager";
import DataStore from "../../../../data/DataStore";
import FormPanel from "../../../../form/FormPanel";
import GridPanel from "../../../../grid/GridPanel";
import Convenience from "../../../../help/Convenience";
import RelationshipPicker from '../../../../relationshippicker/RelationshipPicker'
import './ObjectTypes.css';
import Dialog from '../../../../dialog/Dialog'

export default class ObjectTypes extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            entityId: this.props.entityId,
            objectTypeFields: [],
            loaded: false,
            createDialogOpen: false,
            addFieldDialogOpen: false,
        };
        this.gridPanel = React.createRef();
        this.fieldsGridPanel = React.createRef();
        this.objectTypeFormPanel = React.createRef();
        this.addFieldFormPanel = React.createRef();
    }

    isTab() {
        return true;
    }

    getId() {
        return "objectTypes";
    }

    getLabel() {
        return "Object Types"
    }

    createHandler() {
        this.setState({
            createDialogOpen: true,
        });
    }

    deleteHandler() {
        let selectedItems = this.gridPanel.current.getSelection();
        if (selectedItems.length > 0) {
            fetch(`${HostManager.amperHost()}objectTypes/deleteObjectType`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify({objectTypeId: selectedItems[0].id})
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    this.gridPanel.current.reload();
                }
            })
        }
    }

    getParentToolItems() {
        return [
            {
                label: 'Create',
                handler: this.createHandler.bind(this)
            }, {
                label: 'Delete',
                handler: this.deleteHandler.bind(this)
            }
        ];
    }

    addField() {
        if (this.state.objectTypeId) {
            this.setState({
                addFieldDialogOpen: true,
            });
        }
    }

    removeField() {
        let selectedItems = this.fieldsGridPanel.current.getSelection();
        if (selectedItems.length > 0 && selectedItems[0].objectTypeId === this.state.objectTypeId &&
            !Convenience.isSystemField(selectedItems[0].key)) {
            fetch(`${HostManager.amperHost()}objectTypes/removeObjectTypeField`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(selectedItems[0])
            })
                .then(res => res.json())
                .then((result) => {
                    if (result.success) {
                        this.gridPanel.current.reload(()=>{
                            this.onSelectionChange();
                        });
                    }
                })
        }
    }

    getChildToolItems() {
        return [
            {
                label: 'Add',
                handler: this.addField.bind(this)
            }, {
                label: 'Remove',
                handler: this.removeField.bind(this)
            }
        ];
    }

    getGridDataModel() {
        return [{
            label: 'Label',
            key: 'label'
        },{
            label: 'Api Name',
            key: 'apiName'
        },{
            label: 'Extends',
            key: 'extendsToLabel'
        }];
    }

    getFieldsGridDataModel() {
        return [{
            label: 'Label',
            key: 'label'
        },{
            label: 'Api Name',
            key: 'apiName'
        },{
            label: 'Field Type',
            key: 'type'
        }];
    }

    getDataStore() {
        return new DataStore({
            url: `${HostManager.amperHost()}objectTypes/getObjectTypes`,
            requestMethod: "POST",
            parameters: {
                objectId: this.props.entity.id,
            },
        });
    }

    getFieldsDataStore() {
        return new DataStore({
            data: this.state.objectTypeFields,
            cached: true,
        });
    }

    onSelectionChange() {
        let selectedItems = this.gridPanel.current.getSelection();
        if (selectedItems.length > 0) {
            this.setState({
                objectTypeFields: selectedItems[0].objectTypeFields,
                objectTypeId: selectedItems[0].id,
            })
        } else {
            this.setState({
                objectTypeFields: null,
                objectTypeId: null,
            })
        }
    }

    closeCreateDialog() {
        this.setState({
            createDialogOpen: false,
        });
    }

    handleCreateDialogOk() {
        this.objectTypeFormPanel.current.submit();
    }

    getControlItems() {
        return [
            {
                handler: this.handleCreateDialogOk.bind(this),
                label: 'Ok',
            }, {
                handler: this.closeCreateDialog.bind(this),
                label: 'Cancel',
            }
        ];
    }

    closeAddFieldDialog() {
        this.setState({
            addFieldDialogOpen: false,
        });
    }

    handleAddFieldDialog() {
        this.addFieldFormPanel.current.submit();
    }

    getAddFieldControlItems() {
        return [
            {
                handler: this.handleAddFieldDialog.bind(this),
                label: 'Add',
            }, {
                handler: this.closeAddFieldDialog.bind(this),
                label: 'Cancel',
            }
        ];
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
        if (result.length > 0) {
            result = result + '_amp';
        }
        this.objectTypeFormPanel.current.updateFieldValue('apiName', result)
    }

    getObjectTypeFormConfig() {
        return [[{
                type: "id",
                label: "Object Id",
                name: "objectId",
                value: this.props.entity.id,
            }],[{
                type: 'text',
                label: 'Label',
                name: 'label',
                required: true,
                inputWidth: 350,
                onChange: (event, value) => {
                    this.onLabelChange(event, value);
                },
            }], [{
                type: 'text',
                label: 'Api Name',
                name: 'apiName',
                required: true,
                disabled: true,
                inputWidth: 350,
            }], [{
                type: 'select',
                label: 'Extends',
                name: 'extendsTo',
                keyField: 'id',
                required: true,
                inputWidth: 350,
                value: this.gridPanel.current.getData()[0].id,
                options: this.gridPanel.current.getData(),
                valueFormatter: (value) => {
                    return parseInt(value);
                }
            }]];
    }

    createObjectTypeSuccess() {
        this.gridPanel.current.reload();
        this.closeCreateDialog();
    }

    addFieldToObjectTypeSuccess() {
        this.gridPanel.current.reload(()=>{
            this.onSelectionChange();
        });
        this.closeAddFieldDialog();
    }

    getObjectTypeDialog() {
        if (this.state.createDialogOpen) {
            return <Dialog title={'Create Object Type'} open={this.state.createDialogOpen}
                           controls={this.getControlItems()} closeHandler={this.closeCreateDialog.bind(this)}>
                <FormPanel ref={this.objectTypeFormPanel}
                           parent={this}
                           url={`${HostManager.amperHost()}objectTypes/createObjectType`}
                           onSuccess={this.createObjectTypeSuccess.bind(this)}
                           items={this.getObjectTypeFormConfig()}>
                </FormPanel>
            </Dialog>;
        }
    }

    getAddFieldToObjectTypeFormConfig() {
        return [[{
            type: "id",
            label: "Object Type Id",
            name: "objectTypeId",
            value: this.state.objectTypeId,
        }],[{
            type: 'finder',
            label: 'Fields',
            name: 'fieldId',
            parameters: {
                objectId: this.props.entity.id,
            },
            hideItems: this.fieldsGridPanel.current.getData(),
            hideProperty: 'fieldId',
            inputWidth: 350,
            src: `${HostManager.amperHost()}entities/getFields`,
        }]];
    }

    getAddFieldDialog() {
        if (this.state.addFieldDialogOpen) {
            return <Dialog title={'Add field to the object type'} open={this.state.addFieldDialogOpen}
                           controls={this.getAddFieldControlItems()} closeHandler={this.closeAddFieldDialog.bind(this)}
                           width={500} height={150}>
                <FormPanel ref={this.addFieldFormPanel}
                           parent={this}
                           url={`${HostManager.amperHost()}fields/addObjectTypeField`}
                           onSuccess={this.addFieldToObjectTypeSuccess.bind(this)}
                           items={this.getAddFieldToObjectTypeFormConfig()}>
                </FormPanel>
            </Dialog>;
        }
    }

    render() {
        return (
            <div className={'objectTypes'}>
                <RelationshipPicker parentToolItems={this.getParentToolItems()}
                                    childToolItems={this.getChildToolItems()}>
                    <GridPanel
                        ref={this.gridPanel}
                        selectable={true}
                        selectMode={'row'}
                        onSelectionChange={this.onSelectionChange.bind(this)}
                        dataModel={this.getGridDataModel()}
                        dataStore={this.getDataStore()}
                    />
                    <GridPanel
                        ref={this.fieldsGridPanel}
                        selectable={true}
                        selectMode={'row'}
                        dataModel={this.getFieldsGridDataModel()}
                        dataStore={this.getFieldsDataStore()}
                    />
                </RelationshipPicker>
                {this.getObjectTypeDialog()}
                {this.getAddFieldDialog()}
            </div>
        );
    }
}
