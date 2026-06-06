import React from 'react';
import './Detail.css';
import HostManager from "../../../../../HostManager";
import EventRegistry from "../../../../event/EventRegistery";
import FormPanel from "../../../../form/FormPanel";
import Convenience from "../../../../help/Convenience";
import {sessionManager} from "../../../../../SessionManager";

export default class Detail extends React.Component {

    constructor(props) {
        super(props);
    }

    componentDidMount() {
        /*this.load((result) => {
            const entity = result.entity;
            const error = result.error;
            this.setState(
                {
                    loaded: true,
                    entity: entity,
                    error: error,
                }
            )
        });*/
    }

    load(callBack) {
        const url = Convenience.makeUrl(`${HostManager.amperHost()}entities/getEntity`, {
            entityId: this.state.entityId,
        });
        fetch(url, {
            method: 'get',
            headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
        })
        .then(res => res.json())
        .then((result) => {
            callBack(result);
        });
    }

    isTab() {
        return true;
    }

    getId() {
        return "details";
    }

    getLabel() {
        return "Details"
    }


    getFormConfig() {
        return [
            [{
                type: "id",
                label: "Id",
                name: "id",
                value: this.props.entity.id,
            }],[{
                type: 'text',
                label: 'Title',
                name: 'title',
                required: true,
                value: this.props.entity ? this.props.entity.title : null,
            }],[{
                type: 'text',
                label: 'Title plural',
                name: 'titlePlural',
                required: true,
                value: this.props.entity ? this.props.entity.titlePlural : null,
            }],[{
                type: 'text',
                label: 'Api name',
                name: 'apiName',
                required: true,
                disabled: true,
                value: this.props.entity ? this.props.entity.apiName : null,
            }]
        ];
    }

    onEntitySaveSuccess() {
        EventRegistry.fire("viewChange", this, ["entities"])
    }

    getDetailContent() {
      return <FormPanel ref={this.entityFormPanel}
                        parent={this}
                        url={`${HostManager.amperHost()}entities/edit`}
                        items={this.getFormConfig()}
                        onSuccess={this.onEntitySaveSuccess}
                        submitLabel={"Save"}>
      </FormPanel>
    }

    render() {
        return (
            <div className="detail">
                {this.getDetailContent()}
            </div>
        );
    }
}
