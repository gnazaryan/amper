import React, { useState, useEffect, useRef } from 'react';
import Widget from './Widget';
import RecordDetails from '../../../../components/records/details/RecordDetails';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Autocomplete from '@mui/material/Autocomplete';
import DataStore from "../../../../data/DataStore";
import HostManager from "../../../../../HostManager";

export default function RecordDetailWidget(props) {
  const {dashboardId, interactions, widget, onUpdateWidget, onWidgetStateChange, toast} = props;
  const configuration = JSON.parse(widget.configuration);
  const initialState = () => {
    return {
      configuration,
      configureOpen: false,
      interactionsLoading: false,
      interactions: [],
      widget,
      form: {
        label: widget.label,
        description: widget.description,
        interactions: configuration.interactions || [],
      },
    };
  };
  const [state, setState] = useState(initialState);
  const myRef = useRef();

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

  const onInteraction = (metadata, records) => {
    if (metadata != null && records.length > 0) {
      myRef.current.setRecord(metadata, records[records.length - 1]);
    }
  };
  
  if (state.configuration && state.configuration.interactions) {
    interactions.removeInteractions(widget.id)
    for (let i = 0; i < state.configuration.interactions.length; i++) {
      interactions.registerInteraction(state.configuration.interactions[i].id, widget.id, onInteraction);
    }
  }

  const getInteractionsDataStore = () => {
    return new DataStore({
        url: `${HostManager.amperHost()}widgets/interactions`,
        requestMethod: "POST",
        parameters: {
            dashboardId: parseInt(dashboardId),
            widgetId: widget.id,
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

  const handleConfigurationClose = () => {
    setState({
      ...state,
      configureOpen: false,
      interactionsLoading: false,
    });
  };

  const handleConfigure = () => {
    setState({
        ...state,
        configureOpen: true,
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

  const handleConfigurationSubmit = () => {
    if (state.form.label && state.form.description ) {
      const newConfiguration = {
        ...JSON.parse(state.widget.configuration),
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

  const getConfigureWidgetDialog = () => {
      return (
        <Dialog open={state.configureOpen} onClose={handleConfigurationClose}>
          <DialogTitle>Update {state.form.label} Widget</DialogTitle>
          <DialogContent>
            <DialogContentText>
              To update the widget, specify the interactions, lable and description that would be most descriptive for the users.
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
            <Button onClick={handleConfigurationClose}>Cancel</Button>
            <Button onClick={handleConfigurationSubmit} disabled={!state.form.label || !state.form.description}>Ok</Button>
          </DialogActions>
        </Dialog>
      );
  };

  return (
    <Widget {...props} onConfigure={handleConfigure} configuration={state.configuration}>
        {getConfigureWidgetDialog()}
        <RecordDetails key={widget.id} ref={myRef}></RecordDetails>
    </Widget>
    );
}
