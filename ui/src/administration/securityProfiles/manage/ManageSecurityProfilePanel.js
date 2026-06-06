import React from 'react'
import './ManageSecurityProfilePanel.css'
import Button from "../../../components/button/Button";
import TabPanel from "../../../components/tab/TabPanel";
import EventRegistry from "../../../components/event/EventRegistery";
import HostManager from "../../../HostManager";
import {sessionManager} from "../../../SessionManager";
import ManageSecurityProfileGeneralPanel from "./general/ManageSecurityProfileGeneralPanel"
import ManageSecurityProfileObjectsPanel from "./objects/ManageSecurityProfileObjectsPanel"

export default function ManageSecurityProfilePanel(props) {

    const onUserSaveSuccess = () => {
        EventRegistry.fire("viewChange", null, ["securityProfiles"]);
    };

    const submit = () => {
        fetch(`${HostManager.amperHost()}securityProfiles/create`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
            body: JSON.stringify()
        })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {

                }
            })
    };

    const cancel = () =>  {

    };

    const onChange = (id, values) => {debugger;
        switch (id) {
            case 'general':
                const general = {
                    name: values.name,
                    description: values.description,
                    general: JSON.stringify({
                        manageObjects: values.manageObjects === true,
                        manageUsers: values.manageUsers === true,
                        manageProfiles: values.manageProfiles === true,
                    })
                };
                break;
        }
    };

    return (<div className={'createSecurityProfilePanel'}>
        <TabPanel title="Security Profiles > Add Security Profile">
            <ManageSecurityProfileGeneralPanel onChange={onChange} securityProfile={props.securityProfile}></ManageSecurityProfileGeneralPanel>
            <ManageSecurityProfileObjectsPanel onChange={onChange} securityProfile={props.securityProfile}></ManageSecurityProfileObjectsPanel>
            <Button onClick={submit} label={'Save'}></Button>
            <Button onClick={cancel} label={'Cancel'}></Button>
        </TabPanel>
    </div>);
}