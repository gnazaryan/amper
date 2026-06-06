import React from 'react';
import './ManageSecurityProfileObjectsPanel.css';
import HostManager from "../../../../HostManager";
import FormPanel from "../../../../components/form/FormPanel";

export default class ManageSecurityProfileObjectsPanel extends React.Component {

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
        return "objects";
    }

    getLabel() {
        return "Objects"
    }

    getFormConfig() {
        return [
            [{
                type: "id",
                label: "Id",
                name: "id",
            }]
        ];
    };

    onGeneralSaveSuccess() {

    }

    getGeneralContent() {
        return <FormPanel ref={this.securityProfileFormPanel}
                          url={`${HostManager.amperHost()}securityProfiles/create`}
                          editUrl={`${HostManager.amperHost()}securityProfiles/edit`}
                          items={this.getFormConfig(this.props.securityProfile || {})}
                          onSuccess={this.onGeneralSaveSuccess}
                          mode={this.props.mode}>
        </FormPanel>;
    };

    getValues() {
        return[];
    }
    render() {
        return (
            <div className="general">
                {this.getGeneralContent()}
            </div>
        )
    }
}
