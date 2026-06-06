import React, { useState, useEffect, useRef } from 'react';
import { breadcrumbs, breadcrumbDashboard } from '../../Breadcrambs'
import TileView from '../TileView';
import { useLocation, useNavigate } from 'react-router-dom'
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import {post} from '../../../data/Submit'
import HostManager from "../../../../HostManager";
import TextField from '@mui/material/TextField';
import DataStore from "../../../data/DataStore";
import LinearProgress from '@mui/material/LinearProgress';

export default function Dashboards({expanded, toast}) {
  const {pathname} = useLocation();
  const navigate = useNavigate();
  const dashboardDialogPath = breadcrumbs.dashboard.add.path;
  const initialState = {
      loading: true,
      craeteDashboardForm: {
          configuration: JSON.stringify({
            icon: '@mui/icons-material/ViewQuilt'
          })
      },
      data: [],
      craeteDashboardFormError: undefined,
      snackBarOpen: false,
  };
  const [state, setState] = useState(initialState);

  useEffect(() => {
    if (state.loading) {
      getDataStore().load((result)=> {
        breadcrumbs.resetDashboard();
        for (let i = 0; i < result.data.length; i++) {
          const dashboard = result.data[i];
          breadcrumbs.addDashboard(dashboard);
        }
          setState({
            ...state,
            loading: false,
            data: result.data,
          })
      });
    }
  });

  const getDataStore = () => {
    return new DataStore({
        url: `${HostManager.amperHost()}dashboards/fetch`,
        requestMethod: "POST",
    });
  };

  const getError = (error) => {
      if (error) {
        return <Typography sx={{m: 1}} color="error" variant="caption" display="block">
          {error}
        </Typography>;
      }
  };

  const handleCreateDashboardClose = () => {
    navigate(breadcrumbs.dashboard.path);
  };

  const handleDashboardInputChange = (event) => {
    const {
      target: { value, name },
    } = event;
    const craeteDashboardForm = state.craeteDashboardForm;
    craeteDashboardForm[name] = value;
    setState({
      ...state,
      craeteDashboardForm,
    });
  };

  const handleCreateDashboardSubmit = () => {
    if (state.craeteDashboardForm.label && state.craeteDashboardForm.description) {
      post(`${HostManager.amperHost()}dashboards/add`, state.craeteDashboardForm, (result) => {
        setState(initialState);
        toast('info', `The dashboard "${state.craeteDashboardForm.label}" was successfully added.`)
        navigate(breadcrumbs.dashboard.path);
      }, (result) => {
        toast('info', `The dashboard "${state.craeteDashboardForm.label}" add was not successfull, please contact the support`)
        setState({
          ...state,
          craeteDashboardFormError: result.error,
        });
      })
    }
  };

  const removeDashboardHandler = (dashboardId, label) => {
    if (dashboardId) {
      post(`${HostManager.amperHost()}dashboards/remove`, {
        dashboardId: parseInt(dashboardId)
      }, (result) => {
        toast('info', `The dashboard "${label}" was successfully removed.`)
        setState(initialState);
      }, (result) => {
        toast('error', `The dashboard "${label}" remove was not successfull`)
        setState(initialState);
      })
    }
  };

  const getDialog = () => {
    return (<Dialog open={pathname.startsWith(dashboardDialogPath)} onClose={handleCreateDashboardClose}>
      <DialogTitle>Add Dashboard</DialogTitle>
      <DialogContent>
        <DialogContentText>
          To create a brand new dashboard, specify the lable and a description that would be most descriptive for the users.
        </DialogContentText>
        <TextField
          sx={{mt: 3}}
          autoFocus
          name="label"
          onChange={handleDashboardInputChange}
          value={state.craeteDashboardForm.label}
          error={!state.craeteDashboardForm.label}
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
          onChange={handleDashboardInputChange}
          value={state.craeteDashboardForm.description}
          error={!state.craeteDashboardForm.description}
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
        {getError(state.craeteDashboardFormError)}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleCreateDashboardClose}>Cancel</Button>
        <Button onClick={handleCreateDashboardSubmit} disabled={!state.craeteDashboardForm.label || !state.craeteDashboardForm.description}>Ok</Button>
      </DialogActions>
    </Dialog>);
  };

  const getProgressBar = () => {
      if (state.loading) {
          return <LinearProgress sx={{mb: 1, mr: 6}}/>;
      }
  };

  return (<Box sx={{ width: '100%', height: '100%'}}>
        {getDialog()}
        {getProgressBar()}
        <TileView expanded={expanded} removeHandler={removeDashboardHandler} items={breadcrumbs.dashboard} order={[breadcrumbs.dashboard.overview.key, breadcrumbs.dashboard.add.key]}></TileView>
      </Box>
    );
}
