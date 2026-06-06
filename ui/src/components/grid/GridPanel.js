import React from 'react';
import './GridPanel.css';
import '../Main.css';
import Convenience from "../help/Convenience";

export default class GridPanel extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            loaded: false,
            selection: [],
        };
        this.data = [
            {
                entityKey: "customer",
                entityName: "Customer"
            },
            {
                entityKey: "organization",
                entityName: "Organization"
            },
            {
                entityKey: "employee",
                entityName: "Employee"
            }
        ]
    }

    getTitle() {
        if(this.props.title) {
            return (
                <div className="gridTitle">
                    {this.props.title}
                </div>
            );
        }
    }

    getToolPanelItems() {
        let result = [];
        if (this.props.toolItems) {
            for (let i = 0; i < this.props.toolItems.length; i++) {
                let toolItem = this.props.toolItems[i];
                result.push(<div className="pointer noselect gridToolsItemInactive" onClick={toolItem.handler}>{toolItem.label}</div>);
            }
            return (
                <div className="tools">
                    {result}
                </div>
            );
        }
    }

    getGridContent() {
        let header = [];
        if (this.props.dataModel) {
            if(this.props.selectable && (!this.props.selectMode || this.props.selectMode === 'base')) {
                header.push(<th width="20px"><input type="checkbox" id="selectAll" name="selectAll"></input></th>);
            }

            for (let i = 0; i < this.props.dataModel.length; i++) {
                let modelItem = this.props.dataModel[i];
                header.push(<th>{modelItem.label}</th>);
            }
        }
        let content = [];
        if (this.state.loaded && this.data && this.props.dataModel) {
            for (let l = 0; l < this.data.length; l++) {
                let dataItem = this.data[l];
                content.push(this.getGridRow(dataItem, this.props.dataModel))
            }
        }

        return (
            <div className="gridContent">
                <table className="gridTable">
                    <tr>
                        {header}
                    </tr>
                    {content}
                </table>
            </div>
        );
    }

    onSelectionChange(event) {
        const checked = event.target.checked;
        const id = event.target.id;
        const index = this.state.selection.indexOf(id);
        if (checked) {
            if (index < 0) {
                this.state.selection.push(id);
            }
        } else {
            if (index > -1) {
                this.state.selection.splice(index, 1);
            }
        }
    }

    getData() {
        return this.data;
    }

    getSelection() {
        let result = [];
        if (this.data && this.state.selection.length > 0) {
            for (let i = 0; i < this.data.length; i++) {
                let item = this.data[i];
                if (item.id && this.state.selection.includes(item.id + "")) {
                    result.push(item);
                }
            }
        } else {
            result = this.state.selection;
        }
        return result;
    }

    onRowClicked(event) {
        let id = event.currentTarget.id;
        if (this.props.selectable && this.props.selectMode === 'row') {
            let selection = [event.currentTarget.id];
            if (this.state.selection.length > 0 && this.state.selection[0] == id) {
                selection = [];
            }
            this.setState({
                selection: selection,
            });
            if (this.props.onSelectionChange) {
                setTimeout(()=> {
                    this.props.onSelectionChange(selection);
                }, 100);
            }
        }
    }

    getRecord(id) {
        if (this.data) {
            for (let i = 0; i < this.data.length; i++) {
                let item = this.data[i];
                if (item.id && item.id == id) {
                    return item;
                }
            }
        }
        return null;
    }

    onRowDblclick(event) {
        let id = event.currentTarget.id;
        if (Convenience.hasValue(id) && this.props.onDoubleClick) {
            const item = this.getRecord(id);
            this.props.onDoubleClick(item);
        }
    }

    getGridRow(dataItem, dataModel) {
        const row = [];
        if(this.props.selectable && (!this.props.selectMode || this.props.selectMode === 'base')) {
            row.push(<td width="20px"><input type="checkbox" onChange={this.onSelectionChange.bind(this)} id={dataItem.id} name="id"></input></td>);
        }
        for(let i = 0; i < dataModel.length; i++) {
            let renderedValue = dataItem[dataModel[i].key];
            if (dataModel[i].render) {
                renderedValue = dataModel[i].render(renderedValue);
            }
            row.push(<td key={dataModel[i].key}>{renderedValue}</td>)
        }
        let classNames = '';
        if (this.props.selectable && this.props.selectMode === 'row') {
            classNames = "gridRowHover";
            if (this.state.selection.length && this.state.selection[0] == dataItem.id) {
                classNames = classNames + ' gridRowHoverSelect';
            }
        }
        return (
            <tr id={dataItem.id} className={classNames} onClick={this.onRowClicked.bind(this)} onDoubleClick={this.onRowDblclick.bind(this)}>
                {row}
            </tr>
        )
    }

    componentDidMount() {
        if (this.props.dataStore) {
            var me = this;
            this.props.dataStore.load(function(result) {
                if (result.success) {
                    me.data = result.data;
                    me.setState({loaded: true});
                } else {
                    me.setState({loaded: false});
                }
            });
        }
    }

    reload(callback) {
        this.setState({loaded: false});
        this.props.dataStore.load((result) => {
            if (result.success) {
                this.data = result.data;
                this.setState({loaded: true});
                if (callback) {
                    callback(this.data);
                }
            }
        });
    }

    render() {
        if (this.props.dataStore.configuration.cached) {
            this.props.dataStore.load((result) => {
                if (result.success) {
                    this.data = result.data;
                }
            });
        }
        return (
            <div className="gridPanel">
                {this.getTitle()}
                {this.getToolPanelItems()}
                {this.getGridContent()}
            </div>
        );
    }
}
