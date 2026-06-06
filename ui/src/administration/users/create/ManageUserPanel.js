import React from 'react'
import './ManageUserPanel.css'
import HostManager from "../../../HostManager";
import EventRegistry from "../../../components/event/EventRegistery";
import FormPanel from "../../../components/form/FormPanel";
import Convenience from "../../../components/help/Convenience";
import Loading from "../../../components/loading/Loading";

export default function ManageUserPanel(props) {
    const userFormPanel = React.createRef();

    const onUserSaveSuccess = () => {
        EventRegistry.fire("viewChange", null, ["users"]);
    };

    const getProfilePhotoUrl = (user) => {
        if (user.id) {
            Convenience.makeUrl('users/download', {
                fileName: user.photo,
                userId: user.id,
            })
        }
    };

    const getFormConfig = (user) => {
        return [
            [{
                type: "id",
                label: "Id",
                name: "id",
                value: user.id,
            }],
            [{
                type: 'iamge',
                label: 'Profile picture',
                name: 'photo',
                value: user.photo,
                rowSpan: 1,
                src: getProfilePhotoUrl(user),
            }], [{
                type: 'text',
                label: 'First Name',
                name: 'firstName',
                value: user.firstName,
                required: true,
            },{
                type: 'text',
                label: 'Last Name',
                name: 'lastName',
                value: user.lastName,
                required: true,
            }], [{
                type: 'text',
                label: 'Middle Name',
                name: 'middleName',
                value: user.middleName,
                required: false,
            }], [{
                type: 'text',
                label: 'Username',
                name: 'username',
                value: user.username,
                remoteValidation: 'users/isValidUserName',
                required: true,
                editable: false,
            }], [{
                type: 'text',
                label: 'Email',
                name: 'email',
                validator: 'email',
                value: user.email,
                required: true,
            }], [{
                type: 'finder',
                label: 'Security Profile',
                name: 'profileId',
                value: {
                    id: user.profileId,
                    title: user.profileName,
                },
                required: true,
                src: `${HostManager.amperHost()}securityProfiles/getSecurityProfiles`,
            }]
        ];
    };

    return (<div className="createUserPanel">
        <FormPanel ref={userFormPanel}
                   title="Add User"
                   url={`${HostManager.amperHost()}users/create`}
                   editUrl={`${HostManager.amperHost()}users/edit`}
                   items={getFormConfig(props.user || {})}
                   onSuccess={onUserSaveSuccess}
                   mode={props.mode}
                   submitLabel={"Save"}>
        </FormPanel>
    </div>);
}