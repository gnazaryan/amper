import React, { useState, useEffect, useMemo } from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import LinearProgress from '@mui/material/LinearProgress';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import { useLocation } from 'react-router-dom'
import Typography from '@mui/material/Typography';
import {post} from '../../../data/Submit'
import HostManager from "../../../../HostManager";
import MenuItem from '@mui/material/MenuItem';
import DataStore from "../../../data/DataStore";
import { registerResize, debounce } from '../../../amper/Instruments';
import RecordListWidget from './widgets/RecordListWidget';
import RecordDetailWidget from './widgets/RecordDetailWidget';
import Grid from '@mui/material/Grid2';
import { AppContext } from '../../../../App';

export default function Dashboard({expanded, toast}) {
    const {pathname} = useLocation();
    const pathParts = pathname.split('/');
    const dashboardId = pathParts[3];
    const app = React.useContext(AppContext);

    const handleUpdateWidget = (getWidgets, toastResult) => {
        const widgets = getWidgets();
        if (widgets.length > 0 && widgets[0].id) {
            post(`${HostManager.amperHost()}widgets/update`, {
                Id: parseInt(dashboardId),
                widgets: widgets,
            }, (result) => {
                if (widgets.length > 1) {
                    if (toastResult === true) {
                        app.toast('info', `The dashboard widgets "${widgets.map((list) => list.label).join(", ")}" configuration were successfully updated.`);
                    }
                } else {
                    if (toastResult === true) {
                        app.toast('info', `The dashboard widget "${widgets[0].label}" configuration successfully updated.`);
                    }
                }
            }, (result) => {
                app.toast('error', `The dashboard "${widgets.map((list) => list.label).join(", ")}" configuration updated was not successfull, because ${result.error}`)
                setState(initialState);
            });
        }
    };
    const initialState = {
        loading: true,
        createWidgetDialogOpen: false,
        dashboardViewHeight: window.innerHeight - 170,
        craeteDashboardWidgetForm: {
            dashboardId: parseInt(dashboardId),
        },
        craeteDashboardWidgetFormError: undefined,
        dashboard: undefined,
        saveWidget: debounce(handleUpdateWidget, 3000),
    };
    const [state, setState] = useState(initialState);

    const handleResize = (height, width, expanded) => {
        setState({
          ...state,
          dashboardViewHeight: height - 170,
        })
      };
    useEffect(() => {
        if (state.loading) {
          getDataStore().load((result)=> {
              setState({
                ...state,
                loading: false,
                dashboard: result.data,
              })
          });
        }
        registerResize(handleResize, expanded)
      });
    
    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}widgets/fetch`,
            requestMethod: "POST",
            parameters: {
                dashboardId: parseInt(dashboardId),
            }
        });
    };

    const addWidget = () => {
        setState({
            ...state,
            createWidgetDialogOpen: true,
        });
    };

    const removeDashboard = () => {

    };

    const getProgressBar = () => {
        if (state.loading) {
            return <LinearProgress sx={{mb: 1, mr: 6}}/>;
        }
    };

    const handleCreateWidgetClose = () => {
        setState({
            ...state,
            createWidgetDialogOpen: false,
        });
    };

    const getWidgetDefouldConfiguration = (type) => {
        const result = {
            type: type,
        }
        if (type == 'recordList') {
            result.state = {
                columnsOrder: ['name_sys', 'objectType_sys', 'status_sys'],
                hiddenColumns: ['id', 'identifier_sys'],
            };
        } else if (type == 'recordDetail') {

        }
        return result;
    };

    const handleCreateDashboardWidgetSubmit = () => {
        if (state.craeteDashboardWidgetForm.label && state.craeteDashboardWidgetForm.description) {
          post(`${HostManager.amperHost()}widgets/add`, {
            ...state.craeteDashboardWidgetForm,
            configuration: JSON.stringify(getWidgetDefouldConfiguration(state.craeteDashboardWidgetForm.type))
          }, (result) => {
            setState(initialState);
            app.toast('info', `The dashboard widget "${state.craeteDashboardWidgetForm.label}" was successfully added.`)
          }, (result) => {
            app.toast('error', `The dashboard "${state.craeteDashboardWidgetForm.label}" add was not successfull, please contact the support`)
            setState({
              ...state,
              craeteDashboardWidgetFormError: result.error,
            });
          })
        }
      };

    const getError = (error) => {
        if (error) {
        return <Typography sx={{m: 1}} color="error" variant="caption" display="block">
            {error}
        </Typography>;
        }
    };

    const handleDashboardWidgetInputChange = (event) => {
        const {
          target: { value, name },
        } = event;
        const craeteDashboardWidgetForm = state.craeteDashboardWidgetForm;
        craeteDashboardWidgetForm[name] = value;
        setState({
          ...state,
          craeteDashboardWidgetForm,
        });
      };

    const getCreateWidgetDialog = () => {
        return (<Dialog open={state.createWidgetDialogOpen} onClose={handleCreateWidgetClose}>
          <DialogTitle>Add Widget</DialogTitle>
          <DialogContent>
            <DialogContentText>
              To add a brand new widget, specify the lable, widget type and a description that would be most descriptive for the users.
            </DialogContentText>
            <TextField
              sx={{mt: 3}}
              autoFocus
              name="label"
              onChange={handleDashboardWidgetInputChange}
              value={state.craeteDashboardWidgetForm.label || ''}
              error={!state.craeteDashboardWidgetForm.label}
              label="Label"
              fullWidth
              required
              inputProps={{ maxLength: 16 }}
              variant="filled"
              color="primary"
              size="large"
            />
            <TextField
                select
                sx={{mt: 3}}
                name="type"
                size="large"
                autoFocus
                required
                variant="filled"
                color="primary"
                fullWidth
                value={state.craeteDashboardWidgetForm.type || ''}
                error={!state.craeteDashboardWidgetForm.type}
                label="Type"
                onChange={handleDashboardWidgetInputChange}
                >
                    <MenuItem key="recordList" value="recordList">Record list</MenuItem>
                    <MenuItem key="recordDetail" value="recordDetail">Record detail</MenuItem>
            </TextField>
            <TextField
              sx={{mt: 3}}
              autoFocus
              name="description"
              onChange={handleDashboardWidgetInputChange}
              value={state.craeteDashboardWidgetForm.description || ''}
              error={!state.craeteDashboardWidgetForm.description}
              label="Description"
              fullWidth
              required
              inputProps={{ maxLength: 150 }}
              multiline
              rows={3}
              variant="filled"
              color="primary"
              size="large"
            />
            {getError(state.craeteDashboardWidgetFormError)}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCreateWidgetClose}>Cancel</Button>
            <Button onClick={handleCreateDashboardWidgetSubmit} disabled={!state.craeteDashboardWidgetForm.label || !state.craeteDashboardWidgetForm.description}>Ok</Button>
          </DialogActions>
        </Dialog>);
    };

    const onWidgetStateChange = (originWidget, updatedState, delay, toast) => {
        if (state.dashboard) {
            let updatedWidget = null;
            for (let i = 0; i < state.dashboard.widgets.length; i++) {
                if (state.dashboard.widgets[i].id === originWidget.id) {
                    const configuration = JSON.parse(state.dashboard.widgets[i].configuration);
                    configuration.state = updatedState;
                    state.dashboard.widgets[i].configuration = JSON.stringify(configuration);
                    updatedWidget = state.dashboard.widgets[i]
                    break;
                }
            }
            if (updatedWidget) {
                if (delay === false) {
                    handleUpdateWidget(()=>{
                        return [updatedWidget];
                    }, toast)
                } else {
                    state.saveWidget(updatedWidget, toast);
                }
            }
        }
    };

    const onExpand = (originWidget) => {
        if (state.dashboard) {
            const widgets = state.dashboard.widgets;
            let updatedWidget;
            for (let i = 0; i < widgets.length; i++) {
                updatedWidget = widgets[i];
                if (updatedWidget.id === originWidget.id) {
                    const configuration = JSON.parse(updatedWidget.configuration);
                    if (configuration.widthSX && configuration.widthSX > 47) {
                        return;
                    } else if (configuration.widthSX) {
                        configuration.widthSX++;
                    } else {
                        configuration.widthSX = 25;
                    }
                    updatedWidget.configuration = JSON.stringify(configuration);
                    break;
                }
            }
            if (updatedWidget) {
                state.saveWidget(updatedWidget);
            }
            setState({
                ...state,
                dashboard: state.dashboard,
            });
        }
    };

    const onShrink = (originWidget) => {
        if (state.dashboard) {
            const widgets = state.dashboard.widgets;
            let updatedWidget;
            for (let i = 0; i < widgets.length; i++) {
                updatedWidget = widgets[i];
                if (updatedWidget.id === originWidget.id) {
                    const configuration = JSON.parse(updatedWidget.configuration);
                    if (configuration.widthSX && configuration.widthSX < 13) {
                        return;
                    } else if (configuration.widthSX) {
                        configuration.widthSX--;
                    } else {
                        configuration.widthSX = 23;
                    }
                    updatedWidget.configuration = JSON.stringify(configuration);
                    break;
                }
            }
            if (updatedWidget) {
                state.saveWidget(updatedWidget);
            }
            setState({
                ...state,
                dashboard: state.dashboard,
            });
        }
    };

    const onHighten = (originWidget) => {
        if (state.dashboard) {
            const widgets = state.dashboard.widgets;
            let updatedWidget;
            for (let i = 0; i < widgets.length; i++) {
                updatedWidget = widgets[i];
                if (updatedWidget.id === originWidget.id) {
                    const configuration = JSON.parse(updatedWidget.configuration);
                    if (configuration.height && configuration.height > 1500) {
                        return;
                    } else if (configuration.height) {
                        configuration.height+=25;
                    } else {
                        configuration.height = 525;
                    }
                    updatedWidget.configuration = JSON.stringify(configuration);
                    break;
                }
            }
            if (updatedWidget) {
                state.saveWidget(updatedWidget);
            }
            setState({
                ...state,
                dashboard: state.dashboard,
            });
        }
    };

    const onLower = (originWidget) => {
        if (state.dashboard) {
            const widgets = state.dashboard.widgets;
            let updatedWidget;
            for (let i = 0; i < widgets.length; i++) {
                updatedWidget = widgets[i];
                if (updatedWidget.id === originWidget.id) {
                    const configuration = JSON.parse(updatedWidget.configuration);
                    if (configuration.height && configuration.height < 201) {
                        return;
                    } else if (configuration.height) {
                        configuration.height-=25;
                    } else {
                        configuration.height = 525;
                    }
                    updatedWidget.configuration = JSON.stringify(configuration);
                    break;
                }
            }
            if (updatedWidget) {
                state.saveWidget(updatedWidget);
            }
            setState({
                ...state,
                dashboard: state.dashboard,
            });
        }
    };

    const onRemove = (widget) => {
        if (dashboardId && widget.id) {
            post(`${HostManager.amperHost()}widgets/remove`, {
              dashboardId: parseInt(dashboardId),
              widgetId: widget.id,
            }, (result) => {
                app.toast('info', `The widget "${widget.label}" was successfully removed.`)
                setState(initialState);
            }, (result) => {
                app.toast('error', `The widget "${widget.label}" remove was not successfull`)
                setState(initialState);
            })
          }
    };

    const onUpdateWidget = (widget,  callback) => {
        const oldConfiguration = JSON.parse(widget.configuration);
        if (oldConfiguration.type) {
            const newConfiguration = getWidgetDefouldConfiguration(oldConfiguration.type);
            oldConfiguration.state = newConfiguration.state
            oldConfiguration.configured = true;
            widget.configuration = JSON.stringify(oldConfiguration)       
        }
        post(`${HostManager.amperHost()}widgets/update`, {
            Id: parseInt(dashboardId),
            widgets: [widget],
        }, (result) => {
            callback(true)
            app.toast('info', `The dashboard widget "${widget.label}" successfully updated.`)
        }, (result) => {
            callback(true)
            app.toast('error', `The dashboard "${widget.label}" updated was not successfull, because ${result.error}`)
            setState(initialState);
        });
    };

    const hooks = {
        onUpdateWidget,
        onRemove,
        onHighten,
        onLower,
        onShrink,
        onExpand
    };

    const interactions = useMemo(() => {
        return {
            listeners: {},
            removeInteractions: (listenerId) => {
                for (const [index, value] of Object.entries(interactions.listeners)) {
                    if (value[listenerId] != null) {
                        delete value[listenerId];
                    }
                }
            },
            registerInteraction: (providerId, listenerId, callback) => {
                if (interactions.listeners[providerId] == null) {
                    interactions.listeners[providerId] = {};
                }
                interactions.listeners[providerId] = {
                    ...interactions.listeners[providerId],
                    [listenerId]: callback,
                };
            },
            run: (id, arg, arg1) => {
                const widgets = interactions.listeners[id]
                if (widgets != null) {
                    for (const [index, callback] of Object.entries(widgets)) {
                        if (callback) {
                            callback(arg, arg1);
                        }
                    }    
                }
            }
        }
    });

    const getWidgets = () => {
        const result = [];
        if (state.dashboard) {
            const widgets = state.dashboard.widgets;
            if (widgets) {
                for (let i = 0; i < widgets.length; i++) {
                    const widget = widgets[i];
                    const configuration = JSON.parse(widget.configuration);
                    const height = configuration.height ? configuration.height : 500;
                    const widthSX = configuration.widthSX ? configuration.widthSX : 24;
                    let widgetItem;
                    switch(configuration.type) {
                        case 'recordList': 
                            widgetItem = <RecordListWidget key={widget.id} interactions={interactions} toast={toast} {...hooks} widget={widget} dashboardId={dashboardId} height={height} onWidgetStateChange={onWidgetStateChange}></RecordListWidget>;
                            break;
                        case 'recordDetail':
                            widgetItem = <RecordDetailWidget key={widget.id} interactions={interactions} toast={toast} {...hooks} widget={widget} dashboardId={dashboardId} widgets={widgets} height={height}></RecordDetailWidget>;
                            break;
                    }
                    result.push(<Grid key={i} size={widthSX} sx={{p: 1}}>{widgetItem}</Grid>);
                }
            }
        }
        return result;
    };
    
    return (<Box sx={{ height: '100%', width: 'calc(100% - 25px)' }}>
        {getCreateWidgetDialog()}
        <Stack direction="row" spacing={1} sx={{ mb: 0 }}>
            <Button size="medium" onClick={addWidget} startIcon={<AddCircleOutlineIcon/>}>
                New widget
            </Button>
        </Stack>
        {getProgressBar()}
        <Box sx={{ overflowX: 'hidden', overflowY: 'auto' }} width="100%" height={state.dashboardViewHeight}>
            <Grid container columns={48}>
                {getWidgets()}
            </Grid>
        </Box>
    </Box>);
}