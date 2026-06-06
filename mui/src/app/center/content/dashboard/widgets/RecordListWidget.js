import React, { useState, useEffect, useRef } from 'react';
import Widget from './Widget';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Autocomplete from '@mui/material/Autocomplete';
import DataStore from "../../../../data/DataStore";
import AmperConstatns from "../../../../util/AmperConstants";
import HostManager from "../../../../../HostManager";
import RecordList from '../../../../components/records/recordlist/RecordList';

export default function RecordListWidget(props) {
  const {dashboardId, interactions, widget, onUpdateWidget, onWidgetStateChange, toast} = props;
  const configuration = JSON.parse(widget.configuration);
    const initialState = {
        widget: widget,
        configureOpen: false,
        objectsLoading: true,
        interactionsLoading: (configuration.object != null && configuration.object.apiName != null),
        objectTypesLoading: false,
        form: {
            label: widget.label,
            description: widget.description,
            object: configuration.object,
            interactions: configuration.interactions || [],
            objectType: configuration.objectType,
        },
        objects: [],
        interactions: [],
        objectTypes: [],
        craeteDashboardWidgetFormError: undefined,
        metadata: null,
        configuration,
    };
    const [state, setState] = useState(initialState);
    const myRef = useRef();

    const onInteraction = (metadata, records) => {
      if (state.metadata != null) {
        let fieldApiName = null;
        for (let i = 0; i < state.metadata.Fields.length; i++) {
          if (state.metadata.Fields[i].type === 'REFERENCE' &&
          state.metadata.Fields[i].objectReference === metadata.Object.id) {
            fieldApiName = state.metadata.Fields[i].apiName;
          }
        }
        if (fieldApiName != null) {
          const filter = {
            field: fieldApiName,
            columnField: fieldApiName,
            id: Math.random() * 100000000000000,
            operatorValue: "hasAnyOf",
            value: records.map(item => item[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER]),
            originalValue: records,
          };
          /*let filterModel = configuration.state.filterModel;
          if (filterModel) {
            filterModel = {
              items: [],
            };
          } else if (filterModel.items == null) {
            filterModel.items = [];
          }
          filterModel.items.push(filter);*/
          myRef.current.addFilter(filter);
        }
      }
    };
    
    if (state.configuration && state.configuration.interactions) {
      interactions.removeInteractions(widget.id)
      for (let i = 0; i < state.configuration.interactions.length; i++) {
        interactions.registerInteraction(state.configuration.interactions[i].id, widget.id, onInteraction);
      }
    }
    
    useEffect(() => {
      if (state.configureOpen && state.objectsLoading) {
        getDataStore().load((result)=> {          
            setState({
              ...state,
              objectsLoading: false,
              objects: result.data || [],
            });
        });
      }
    }, [state.configureOpen, state.objectsLoading]);
    
    useEffect(() => {
      if (state.configureOpen && state.interactionsLoading) {
        getInteractionsDataStore().load((result)=> {
          const interactions = [];
          if (result.success === true && result.data.widgets.length > 0) {
            for (let i = 0; i < result.data.widgets.length; i++) {
              interactions.push({
                label: result.data.widgets[i].label,
                id: result.data.widgets[i].id,
              });
            }
          }
          
          setState({
            ...state,
            interactionsLoading: false,
            objectsLoading: false,
            interactions: interactions,
          });
        });
      }
    }, [state.configureOpen, state.interactionsLoading]);

    const getInteractionsDataStore = () => {
      return new DataStore({
          url: `${HostManager.amperHost()}widgets/interactions`,
          requestMethod: "POST",
          parameters: {
              dashboardId: parseInt(dashboardId),
              widgetId: widget.id,
              objectApiName: state.form.object.apiName,
          }
      });
  };

    const handleInputChange = (event) => {
        const {
          target: { value, name },
        } = event;
        const form = state.form;
        form[name] = value;
        setState({
          ...state,
          form,
        });
      };

    const handleClose = () => {
      setState({
        ...state,
        configureOpen: false,
        interactionsLoading: (configuration.object != null && configuration.object.apiName != null),
    });
    };

    const handleSubmit = () => {
      if (state.form.label && state.form.description && state.form.object && state.form.object.id) {
        const newConfiguration = {
          ...JSON.parse(state.widget.configuration),
          object: {
            id: state.form.object.id,
            apiName: state.form.object.apiName,
            label: state.form.object.label,
          },
          interactions: state.form.interactions,
          configured: true,
        };
        const updatedWidget = {
          ...state.widget,
          label: state.form.label,
          description: state.form.description,
          configuration: JSON.stringify(newConfiguration),
        };
        
        onUpdateWidget(updatedWidget, (success) => {
          setState({
            ...state,
            widget: updatedWidget,
            configuration: newConfiguration,
            configureOpen: false,
          });
          if (myRef.current && myRef.current.reset) {
            myRef.current.reset();
          }
        });
      }
    };

    const handleConfigure = () => {
        setState({
            ...state,
            configureOpen: true,
            objectsLoading: true,
        });
    };

    const handleObjectChange = (event, newValue) => {
      let object = undefined;
      if (newValue) {
        object = {
            id: newValue.id,
            apiName: newValue.apiName,
            label: newValue.title || newValue.label,
          };
      }
        setState({
          ...state,
          form: {
            ...state.form,
            object: object,
            interactions: [],
          },
          interactionsLoading: true,
        });
    };

    const handleInteractionChange = (event, newValue) => {
      setState({
        ...state,
        form: {
          ...state.form,
          interactions: newValue
        },
      });
    };

    const getDataStore = () => {
      return new DataStore({
          url: `${HostManager.amperHost()}entities/getEntities`,
          requestMethod: "POST",
          parameters: {
              start: 0,
              limit: AmperConstatns.INTEGER.MAX_VALUE
          }
      });
    };

  const updateWidgetConfigStateSave = (recordListState, delay, toast) => {
    onWidgetStateChange(state.widget, recordListState, delay, toast)
  };

    const getCreateWidgetDialog = () => {
        return (
          <Dialog open={state.configureOpen} onClose={handleClose}>
            <DialogTitle>Update {state.form.label} Widget</DialogTitle>
            <DialogContent>
              <DialogContentText>
              To update the widget, specify the object, interactions, lable and description that would be most descriptive for the users.
              </DialogContentText>
              <TextField
                sx={{mt: 3}}
                autoFocus
                name="label"
                onChange={handleInputChange}
                value={state.form.label}
                error={!state.form.label}
                label="Label"
                fullWidth
                required
                inputProps={{ maxLength: 16 }}
                variant="filled"
                color="primary"
                size="large"
              />
              <TextField
                sx={{mt: 3}}
                autoFocus
                name="description"
                onChange={handleInputChange}
                value={state.form.description}
                error={!state.form.description}
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
              <Autocomplete
                sx={{mt: 3}}
                onChange={handleObjectChange}
                options={state.objects}
                value={state.form.object}
                getOptionLabel={item => item.title || item.label}
                renderInput={(params) => <TextField
                  variant="standard"
                  {...params}
                  error={!state.form.object}
                  name="object"
                  required
                  label="Select object" />}
              />
              <Autocomplete
                sx={{mt: 3}}
                onChange={handleInteractionChange}
                options={state.interactions}
                isOptionEqualToValue={(option, value) => option.id === value.id}
                value={state.form.interactions}
                multiple={true}
                getOptionLabel={item => item.title || item.label}
                renderInput={(params) => <TextField
                  variant="standard"
                  {...params}
                  name="object"
                  label="Select widget interactions" />}
              />
            </DialogContent>
            <DialogActions>
              <Button onClick={handleClose}>Cancel</Button>
              <Button onClick={handleSubmit} disabled={!state.form.label || !state.form.description}>Ok</Button>
            </DialogActions>
          </Dialog>
        );
    };
        
    const onRecordsLoad = (metadata, records) => {
      setState({
        ...state,
        metadata,
      });
    };
  return (
        <Widget {...props} onConfigure={handleConfigure} configuration={state.configuration}>
            {getCreateWidgetDialog()}
            <RecordList id={state.widget.id} ref={myRef} toast={toast} onSelect={interactions.run} onLoad={onRecordsLoad} state={state.configuration.state} object={state.configuration.object} objectType={state.configuration.objectType} onStateChange={updateWidgetConfigStateSave}/>
        </Widget>
    );
}
