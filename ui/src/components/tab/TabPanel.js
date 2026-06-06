import React from 'react';
import './TabPanel.css';

export default class TabPanel extends React.Component {

    constructor(props) {
        super(props);
        this.state = {};
    }

    getTitle() {
        if(this.props.title) {
            return (
                <div className="tabPanleTitle">
                    {this.props.title}
                </div>
            );
        }
    }

    tabItemChang(event) {
        const id = event.target.id;
        this.setState(
            {
                tabId: id,
            }
        )
    }

    getTabItems() {
        const result = [];
        const children = this.props.children;
        for(let i = 0; i < children.length; i++) {
            const child = children[i];
            if (child.type.prototype.isTab && child.type.prototype.isTab()) {
                if (i > 0) {
                    result.push(<span className="intermediateItem">|</span>)
                }
                const childId = child.type.prototype.getId();
                const tabId = this.state.tabId;
                let className = "tabItem";
                if ((tabId == null && i == 0) || childId === tabId) {
                    className = "tabItemActive";
                }

                result.push(
                    <span className={className} onClick={this.tabItemChang.bind(this)} id={childId}>
                    {child.type.prototype.getLabel()}
                </span>
                )
            }
        }
        return result;
    }

    getActiveTabContent() {
        const children = this.props.children;
        for(let i = 0; i < children.length; i++) {
            const child = children[i];
            if (child.type.prototype.isTab && child.type.prototype.isTab()) {
                const childId = child.type.prototype.getId();
                const tabId = this.state.tabId;
                if (tabId == childId) {
                    return child;
                }
            }
        }
        return this.props.children[0];
    }

    getControls() {
        const result = [];
        const children = this.props.children;
        for(let i = 0; i < children.length; i++) {
            const child = children[i];
            if (child.type.name == 'Button') {
                result.push(child);
            }
        }
        return result;
    }

    render() {
        return (
            <div className="tabPanel">
                {this.getTitle()}
                <div className="tabPanelContent">
                    <div className="tabs">
                        <div className="tabsCenter">
                            {this.getTabItems()}
                        </div>
                    </div>
                    <div className="tabContent">
                        {this.getActiveTabContent()}
                    </div>
                    {this.getControls()}
                </div>
            </div>
        );
    }
}
