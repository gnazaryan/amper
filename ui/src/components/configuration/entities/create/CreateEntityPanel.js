import React from 'react';
import './CreateEntityPanel.css';
import HostManager from "../../../../HostManager";
import EventRegistry from "../../../event/EventRegistery";
import FormPanel from "../../../form/FormPanel";

export default class CreateEntityPanel extends React.Component {

    constructor(props) {
        super(props);
        this.entityFormPanel = React.createRef();
        this.state = {};
    }

    onTitleChange(event, value) {
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
        this.entityFormPanel.current.updateFieldValue('apiName', result)
    }

    getFormConfig() {
        return [
            [{
                type: 'text',
                label: 'Title',
                name: 'title',
                required: true,
                onChange: (event, value) => {
                    this.onTitleChange(event, value);
                },
            }],[{
                type: 'text',
                label: 'Title plural',
                name: 'titlePlural',
                required: true,
            }],[{
                type: 'text',
                label: 'Api name',
                name: 'apiName',
                required: true,
                disabled: true,
                onKeyDown: (event) => {
                    /*var key = event.keyCode;
                    if(!((key >= 65 && key <= 90) || key == 8)) {
                        event.preventDefault();
                    }*/
                },
            }]
        ];
    }

    onEntitySaveSuccess() {
        EventRegistry.fire("viewChange", this, ["entities"])
    }

    render() {
        return (
            <div className="createEntity">
                <FormPanel ref={this.entityFormPanel}
                           parent={this}
                           title="Create entity"
                           url={`${HostManager.amperHost()}entities/create`}
                           items={this.getFormConfig()}
                           onSuccess={this.onEntitySaveSuccess}
                           submitLabel={"Create"}>
                </FormPanel>
            </div>
        );
    }
}
