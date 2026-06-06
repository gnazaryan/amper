import React from 'react';
import './ManageSecurityProfileGeneralPanel.css';
import Button from "../../../../components/button/Button";
import HostManager from "../../../../HostManager";
import FormPanel from "../../../../components/form/FormPanel";

export default class ManageSecurityProfileGeneralPanel extends React.Component {

    constructor(props) {
        super(props);
        this.state = {

        };
        this.securityProfileFormPanel = React.createRef();
    }

    isTab() {
        return true;
    }

    getId() {
        return "general";
    }

    getLabel() {
        return "General"
    }

    getFormConfig(profile) {
        return [
            [{
                type: "id",
                label: "Id",
                name: "id",
            }],[{
                type: 'text',
                label: 'Name',
                name: 'profileName',
                value: profile.name,
                required: true,
            }], [{
                type: 'text',
                label: 'Description',
                name: 'description',
                value: profile.description,
                required: true,
            }], [{
                type: 'boolean',
                label: 'Manage Objects',
                name: 'manageObjects',
                required: true,
            }, {
                type: 'boolean',
                label: 'Manage Users',
                name: 'manageUsers',
                required: true,
            }, {
                type: 'boolean',
                label: 'Manage Profiles',
                name: 'manageProfiles',
                required: true,
            }]
        ];
    };

    onChange(values) {
        this.props.onChange(this.getId(), values);
    }

    getGeneralContent() {
        return <FormPanel ref={this.securityProfileFormPanel}
                          url={`${HostManager.amperHost()}securityProfiles/create`}
                          editUrl={`${HostManager.amperHost()}securityProfiles/edit`}
                          items={this.getFormConfig(this.props.securityProfile || {})}
                          onSuccess={this.onGeneralSaveSuccess}
                          onChange={this.onChange.bind(this)}
                          mode={this.props.mode}>
        </FormPanel>;
    };

    render() {
        return (
            <div className="general">
                {this.getGeneralContent()}
            </div>
        )
    }
}
