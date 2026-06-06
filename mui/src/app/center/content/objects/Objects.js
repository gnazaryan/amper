import React, { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import { DataGridPremium } from '../../../components/x-data-grid-premium';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import HostManager from "../../../../HostManager";
import DataStore from "../../../data/DataStore";
import {sessionManager} from "../../../../SessionManager";
import AmperConstatns from "../../../util/AmperConstants";
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import EditIcon from '@mui/icons-material/Edit';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import { useNavigate } from 'react-router-dom'
import { breadcrumbs } from '../../Breadcrambs'
import LinearProgress from '@mui/material/LinearProgress';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import {post} from '../../../data/Submit'
import Typography from '@mui/material/Typography';
import {makeApieName} from '../../../amper/Instruments'


export default function Objects({toast}) {

    const initialState = {
      loading: true,
      createObjectDialogOpen: false,
      data: [],
      craeteObjectForm: {
        title: undefined,
        titlePlural: undefined,
        apiName: '',
      },
      createObjectFormError: undefined,
      objectPaging: {
        page: 0,
        pageSize: 50,
    }
  };
    const [state, setState] = useState(initialState);
    const [selectedRowData, setSelectedRowData] = React.useState([]);
    const navigate = useNavigate();

    useEffect(() => {
      if (state.loading) {
        getDataStore().load((result)=> {          
            setState({
              ...state,
              loading: false,
              data: result.data || [],
            })
        });
      }
    });

    const columns = [
        { 
            field: 'id', 
            headerName: 'ID',
            hide: true,
        },
        {
          field: 'title',
          headerName: 'Label',
          flex: 1,
        },
        {
          field: 'titlePlural',
          headerName: 'Label Plural',
          flex: 1,
        },
        {
          field: 'apiName',
          headerName: 'Api name',
          flex: 1,
        },
      ];

      const create = () => {
        setState({
          ...state,
          createObjectDialogOpen: true,
        });
      };

      const remove = () => {
        if (selectedRowData.length > 0) {
          post(`${HostManager.amperHost()}entities/deleteEntity`, {
            entityId: selectedRowData[0].id,
          }, () => {
            setState(initialState);
          }, () => {
            setState({
              initialState
            });
          });
        }
      };

      const manage = () => {
        if (selectedRowData.length > 0) {
          setTimeout(()=> {
            navigate(breadcrumbs.configuration.objects.edit.key + '?objectId=' + selectedRowData[0].id);
          }, 1);
        }
      };

      const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}entities/getEntities`,
            requestMethod: "POST",
            parameters: {
                sessionId: sessionManager.getSessionId(),
                start: 0,
                limit: AmperConstatns.INTEGER.MAX_VALUE
            }
        });
    }

    const handleCreaeteObjectClose = () => {
      setState({
        ...state,
        createObjectDialogOpen: false,
      });
    };

    const handleObjectLabelChange = (event) => {
      const {
        target: { value, name },
      } = event;
      setState({
        ...state,
        craeteObjectForm : {
          ...state.craeteObjectForm,
          title: value,
          apiName: makeApieName(value),
        }
      });
    };

    const handleObjectLabelPluralChange = (event) => {
      const {
        target: { value, name },
      } = event;
      setState({
        ...state,
        craeteObjectForm : {
          ...state.craeteObjectForm,
          titlePlural: value,
        }
      });
    };

    const handleCreaeteObjectSubmit = () => {
      if (state.craeteObjectForm.title && state.craeteObjectForm.titlePlural && state.craeteObjectForm.apiName) {
        post(`${HostManager.amperHost()}entities/create`, state.craeteObjectForm, (result) => {
          setState(initialState);
        }, (result) => {
          setState({
            ...state,
            createObjectFormError: result.error,
          });
        });
      }
    };

    const setUserPaginationModel = (pagingModel) => {
      setState({
          ...state,
          objectPaging: pagingModel
      });
  };

    const getError = (error) => {
      if (error) {
        return <Typography sx={{m: 1}} color="error" variant="caption" display="block">
          {error}
        </Typography>;
      }
  };

    const getCreateObjectTypeDialog = () => {
      return (<Dialog open={state.createObjectDialogOpen} onClose={handleCreaeteObjectClose}>
        <DialogTitle>Create Object Type</DialogTitle>
        <DialogContent>
          <DialogContentText>
            To create an object, specify the lable, which would auto generate the api name and then fill in the plural form of the label.
          </DialogContentText>
          <TextField
            sx={{mt: 3}}
            autoFocus
            name="label"
            onChange={handleObjectLabelChange}
            value={state.craeteObjectForm.title}
            error={!state.craeteObjectForm.title}
            label="Label"
            fullWidth
            required
            variant="filled"
            color="primary"
            size="large"
          />
           <TextField
            sx={{mt: 3}}
            autoFocus
            name="labelPlural"
            onChange={handleObjectLabelPluralChange}
            value={state.craeteObjectForm.titlePlural}
            error={!state.craeteObjectForm.titlePlural}
            label="Label plural"
            fullWidth
            required
            variant="filled"
            color="primary"
            size="large"
          />
          <TextField
            sx={{mt: 3}}
            autoFocus
            name="apiName"
            value={state.craeteObjectForm.apiName}
            error={!state.craeteObjectForm.apiName}
            required
            disabled
            label="Api name"
            fullWidth
            variant="filled"
            color="primary"
            size="large"
          />
          {getError(state.createObjectFormError)}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCreaeteObjectClose}>Cancel</Button>
          <Button onClick={handleCreaeteObjectSubmit} disabled={!state.craeteObjectForm.title || !state.craeteObjectForm.titlePlural || !state.craeteObjectForm.apiName}>Ok</Button>
        </DialogActions>
      </Dialog>);
    };

  return (
        <Box sx={{ height: 'calc(100% - 35px)', width: 'calc(100% - 35px)' }}>
            {getCreateObjectTypeDialog()}
            <Stack direction="row" spacing={1} sx={{ mb: 1 }}>
                <Button size="small" onClick={create} startIcon={<AddCircleOutlineIcon/>}>
                    Create
                </Button>
                <Button size="small" onClick={manage} startIcon={<EditIcon/>} disabled={selectedRowData.length == 0}>
                    Manage
                </Button>
                <Button size="small" onClick={remove} startIcon={<RemoveCircleOutlineIcon/>} disabled={selectedRowData.length == 0}>
                    Remove
                </Button>
            </Stack>
            <DataGridPremium
                rows={state.data}
                loading={state.loading}
                slots={{
                    loadingOverlay: LinearProgress,
                }}
                columns={columns}
                pagination
                pageSizeOptions={[50, 100, 500, 1000, 5000, 50000]}
                paginationModel={state.objectPaging}
                onPaginationModelChange={setUserPaginationModel}
                onRowSelectionModelChange={(ids) => {
                    const selectedIDs = new Set(ids);
                    const rowData = state.data.filter((row) =>
                      selectedIDs.has(row.id)
                    )
                    setSelectedRowData(rowData);
                  }}
            />
        </Box>
    );
}
