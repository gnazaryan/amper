import React from 'react';
import './EditEntityPanel.css';
import TabPanel from "../../../tab/TabPanel";
import Detail from "./details/Detail";
import Fields from "./filds/Fields";
import ObjectTypes from "./objecttypes/ObjectTypes";

export default class EditEntityPanel extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            entity: this.props.options.entity,
        };
    }

    componentDidMount() {

    }

    getFormConfig() {
        return [
            {
                type: 'text',
                label: 'Title',
                name: 'title'
            },{
                type: 'text',
                label: 'Title plural',
                name: 'titlePlural'
            },{
                type: 'text',
                label: 'Api name',
                name: 'apiName',
                disabled: true,
            }
        ];
    }

    render() {
        return (
            <div className="editEntity">
                <TabPanel title="Objects > Modify object">
                    <Detail entity={this.state.entity}></Detail>
                    <ObjectTypes entity={this.state.entity}></ObjectTypes>
                    <Fields entity={this.state.entity}></Fields>
                </TabPanel>
            </div>
        );
    }
}
