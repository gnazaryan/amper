import React, { useState, useRef, useImperativeHandle, useEffect } from 'react';
import {
  GridToolbarContainer,
  GridToolbarFilterButton,
  GridToolbarColumnsButton,
  GridToolbarExport,
  GridToolbarDensitySelector
} from '@mui/x-data-grid';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import RecordCreate from '../create/RecordCreate';
import { post } from '../../../data/Submit';
import HostManager from '../../../../HostManager';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import AdjustIcon from '@mui/icons-material/Adjust';
import { AppContext } from '../../../../App';

export const RecordToolbar = ({ setFilterButtonEl, metadata, toast, refresh, onRemove, onUpdate, parentRef }) => {
  const app = React.useContext(AppContext);

    const initialState = () => {
      return {
        createDialogOpen: false,
        recordCreateSubmitted: false,
        recordCreateValid: false,
        updateEnabled: false,
      };
    };

    const [state, setState] = useState(initialState);

    const createRef = useRef(null);

    useEffect(() => {
      if (state.recordCreateSubmitted) {
        if (refresh) {
          refresh();
        }
      }
    }, [state.recordCreateSubmitted]);

    const openCreateDialog = () => {
        setState({
          ...state,
          createDialogOpen: true,
        })
    };

    const closeCreateDialog = () => {
        setState(initialState());
    };

    const submitCreateRecord = () => {
      if (createRef.current.onSubmit) {
        createRef.current.onSubmit();
      }
      post(`${HostManager.amperHost()}records/add`, {
        apiName: metadata.Object.apiName,
        payload: JSON.stringify(state.payload)
      }, (result) => {
        if (createRef.current.onSubmitComplete) {
          createRef.current.onSubmitComplete();
        }
        setState({
          ...state,
          createDialogOpen: false,
          recordCreateSubmitted: true,
        });
        app.toast('info', `The object "${metadata.Object.apiName}" record was successfully created.`)
      }, (result) => {
        if (createRef.current.onSubmitComplete) {
          createRef.current.onSubmitComplete();
        }
        setState({
          ...state,
          createDialogOpen: false,
          recordCreateSubmitted: false,
        });
        app.toast('error', `The object "${metadata.Object.apiName}" record creation was not successfull, please contact the support`)
      });
    };

    const onRecordCreateChange = (payload, valid) => {
      setState({
        ...state,
        payload,
        payloadValid: valid,
      });
    };

    const removeRecords = () => {
      if (onRemove) {
        onRemove();
      }
    };

    useImperativeHandle(parentRef, () => ({
      enableUpdate() {
        setState({
          ...state,
          updateEnabled: true
        });
      },
      disableUpdate() {
        setState({
          ...state,
          updateEnabled: false
        });
      }
  }));

    const getCreateWidgetDialog = () => {
      return (<Dialog open={state.createDialogOpen} onClose={closeCreateDialog} fullWidth={true} maxWidth={'sm'}>
        <DialogTitle>Create Record</DialogTitle>
        <DialogContent>
          <DialogContentText>
            To create a brand new record, input the the values for the displayed fields below.
          </DialogContentText>
          <RecordCreate ref={createRef} metadata={metadata} onChange={onRecordCreateChange}></RecordCreate>
        </DialogContent>
        <DialogActions>
          <Button onClick={closeCreateDialog}>Cancel</Button>
          <Button onClick={submitCreateRecord} disabled={!state.payloadValid}>Ok</Button>
        </DialogActions>
      </Dialog>);
  };
  
    return <Box sx={{
      display: 'flex',
      flexDirection: 'row',
      bgcolor: 'background.paper',
      }}>
        <Box sx={{ flexGrow: 1 }}>
          <GridToolbarContainer>
            <GridToolbarFilterButton ref={setFilterButtonEl} label="sdsd"/>
            <GridToolbarColumnsButton label="sdsd"/>
            <GridToolbarExport label="sdsd"/>
            <GridToolbarDensitySelector label="sdsd"/>
          </GridToolbarContainer>
        </Box>
        <Box sx={{ flexGrow: 0 }}>
          {getCreateWidgetDialog()}
            <Button
              onClick={onUpdate}
              disabled={!state.updateEnabled}
              key={'amperWidgetUpdateButton'}
              size="small"
              color="primary"
              aria-label="UPDATE"
              startIcon={<AdjustIcon color={state.updateEnabled ? 'primary' : 'inactive'} fontSize="small" />} sx={{m: '4px'}}>
              UPDATE              
            </Button>
          <Button
              onClick={removeRecords}
              key={'amperWidgetRemoveButton'}
              size="small"
              color="primary"
              aria-label="REMOVE"
              startIcon={<RemoveCircleOutlineIcon color='primary' fontSize="small" />} sx={{m: '4px'}}>
              REMOVE              
            </Button>
          <Button
              onClick={openCreateDialog}
              key={'amperWidgetCreateButton'}
              size="small"
              color="primary"
              aria-label="CREATE"
              startIcon={<AddCircleOutlineIcon color='primary' fontSize="small" />} sx={{m: '4px'}}>
              CREATE              
            </Button>
        </Box>
    </Box>
};
  