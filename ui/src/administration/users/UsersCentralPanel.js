import React from 'react';
import './UsersCentralPanel.css';
import '../../components/Main.css';
import HostManager from '../../HostManager';
import {sessionManager} from "../../SessionManager";
import DataStore from "../../components/data/DataStore";
import Dialog from "../../components/dialog/Dialog";
import DisplayMessage from "../../components/display/DisplayMessage";
import EventRegistry from "../../components/event/EventRegistery";
import GridPanel from "../../components/grid/GridPanel";

export default class UsersCentralPanel extends React.Component {

    constructor(props) {
        super(props);
        this.gridPanel = React.createRef();
        this.state = {
            removeUserDialogOpen: false,
        };
    }

    createHandler() {
        EventRegistry.fire("menuItemChange", this, ['createUser']);
    }

    editHandler = () => {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length == 1) {
            EventRegistry.fire("menuItemChange", this, ['editUser', {
                user: selection[0],
            }]);
        }
    }

    deleteHandler = () => {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length == 1) {
            this.setState({
                removeUserDialogOpen: true,
            });
        }
    }

    getToolItems() {
        return [
            {
                label: 'Add',
                handler: this.createHandler
            }, {
                label: 'Edit',
                handler: this.editHandler
            }, {
                label: 'Remove',
                handler: this.deleteHandler
            }

        ];
    }

    getGridDataModel() {
        return [{
            label: 'First Name',
            key: 'firstName'
        },{
            label: 'Last Name',
            key: 'lastName'
        },{
            label: 'Username',
            key: 'username'
        },{
            label: 'Email',
            key: 'email'
        },{
            label: 'Profile',
            key: 'profileName'
        }];
    }

    getDataStore() {
        return new DataStore({
            url: `${HostManager.amperHost()}users/getUsers`,
            requestMethod: "POST",
            parameters: {
                'sessionId': this.props.sessionId,
                start: 0,
                limit: 50,
            },
            dataModel: this.getGridDataModel()
        });
    }

    handleRemoveUserDialog() {
        const selection = this.gridPanel.current.getSelection();
        if (selection.length == 1) {
            fetch(`${HostManager.amperHost()}users/remove?userId=` + selection[0].id, {
                method: 'get',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    this.gridPanel.current.reload();
                }
                this.closeRemoveUserDialog()
            })
        }
    }

    closeRemoveUserDialog() {
        this.setState({
            removeUserDialogOpen: false,
        });
    }

    getRemoveUserControlItems() {
        return [
            {
                handler: this.handleRemoveUserDialog.bind(this),
                label: 'Remove',
            }, {
                handler: this.closeRemoveUserDialog.bind(this),
                label: 'Cancel',
            }
        ];
    }

    getRemoveUserDialog() {
        if (this.state.removeUserDialogOpen) {
            return <Dialog title={'Remove user'} open={this.state.removeUserDialogOpen}
                           controls={this.getRemoveUserControlItems()} closeHandler={this.closeRemoveUserDialog.bind(this)}
                           width={500} height={110}>
                <DisplayMessage>
                    Are you sure you want to proceed with the remove user operation ?
                </DisplayMessage>
            </Dialog>;
        }
    }

    onDoubleClick(item) {
        EventRegistry.fire("menuItemChange", this, ['viewUser', {
            user: item,
        }]);
    }

    render() {
        return (
            <div className="centralPanel">
                <GridPanel
                    ref={this.gridPanel}
                    title="Users"
                    selectMode={'row'}
                    selectable={true}
                    onDoubleClick={this.onDoubleClick.bind(this)}
                    toolItems={this.getToolItems()}
                    dataModel={this.getGridDataModel()}
                    dataStore={this.getDataStore()}
                ></GridPanel>
                {this.getRemoveUserDialog()}
            </div>
        );
    }
}
