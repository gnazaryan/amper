import React, { useEffect, useState, } from 'react';
import LinearProgress from '@mui/material/LinearProgress';
import { DataGridPremium, GridToolbarContainer,
    GridToolbarColumnsButton,
    GridToolbarExport,
    GridToolbarDensitySelector, GridToolbarQuickFilter, } from '../../x-data-grid-premium';
import Button from '@mui/material/Button';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import AdjustIcon from '@mui/icons-material/Adjust';
import Box from '@mui/material/Box';
import DataStore from '../../../data/DataStore';
import HostManager from '../../../../HostManager';
import { AppContext } from '../../../../App';
import { useNavigate, useLocation } from 'react-router-dom'
import { breadcrumbs } from '../../../center/Breadcrambs';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import DialogActions from '@mui/material/DialogActions';
import TextField from '@mui/material/TextField';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Convenience from '../../../help/Convenience';
import { sessionManager } from '../../../../SessionManager';

const NodeList = ({onSelect}) => {
    const app = React.useContext(AppContext);
    const [state, setState] = useState(() => {
        return {
            loading: true,
            data: [],
            search: [],
            newInstanceForm: {
                valid: false,
                type: 'amperDatastoreInstance',
            },
            selectedRowData: [],
            submitting: false,
            removeDialogOpen: false,
        };
    });

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
    }, [state.loading]);

    const {pathname} = useLocation();
    const navigate = useNavigate();
    const newInstancePath = breadcrumbs.administration.nodes.new.path;
    const updateInstancePath = breadcrumbs.administration.nodes.update.path;
    const updateInstanceMode = pathname.startsWith(updateInstancePath);
    const newInstanceMode = pathname.startsWith(newInstancePath);
    const newInstanceDialogOpen = updateInstanceMode || newInstanceMode

    useEffect(() => {
        if (updateInstanceMode) {
            setState({
                ...state,
                newInstanceForm: {
                    ...state.newInstanceForm,
                    ...state.selectedRowData[0],
                    valid: true,
                }
            });
        }
    }, [updateInstanceMode]);

    useEffect(() => {
        if (newInstanceMode) {
            setState({
                ...state,
                newInstanceForm: {
                    valid: false,
                    type: 'amperDatastoreInstance',
                }
            });
        }
    }, [newInstanceMode]);

    if (app) {
        app.registerRefresh("nodes", () => {
            refresh();
        });    
    }

    const refresh = () => {
        setState({
            ...state,
            loading: true,
        });
    };

    const handleNewUserDialogClose = () => {
        navigate(breadcrumbs.administration.nodes.path);
    };

    const closeDialogReset = () => {
        handleNewUserDialogClose();
        setState({
            ...state,
            loading: true,
            errorMessage: '',
        });
    };

    const closeRemoveDialogReset = () => {
        setState({
            ...state,
            loading: true,
            errorMessage: '',
            removeDialogOpen: false,
        });
    };

    const onChange = (event, arg1) => {
        const {
            target: { value, name },
          } = event;
          setState({
            ...state,
            newInstanceForm: {
                ...state.newInstanceForm,
                [name]: value
            },
        });
    };

    const onChangeLimit = (event, arg1) => {
        const {
            target: { value, name },
          } = event;
          setState({
            ...state,
            newInstanceForm: {
                ...state.newInstanceForm,
                [name]: parseInt(value),
            },
        });
    };
    const formValid = () => {
        const form = state.newInstanceForm;
        return Convenience.hasValue(form.name) && Convenience.hasValue(form.type)
            && Convenience.hasValue(form.address) && Convenience.hasValue(form.port)
            && Convenience.hasValue(form.limit) && Convenience.hasValue(form.directory);
    };
    
    const addNewInstance = () => {
        if (formValid()) {
            setState({
                ...state,
                submitting: true,
            });
            fetch(`${HostManager.amperHost()}amper/create`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(state.newInstanceForm)
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    closeDialogReset()
                } else {
                    setState({
                        ...state,
                        errorMessage: result.error,
                    });                
                }
            });
        }
    };

    const updatInstance = () => {
        if (formValid()) {
            setState({
                ...state,
                submitting: true,
            });
            fetch(`${HostManager.amperHost()}amper/edit`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(state.newInstanceForm)
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    closeDialogReset()
                } else {
                    setState({
                        ...state,
                        errorMessage: result.error,
                    });                
                }
            });
        }
    }

    const testInstanceAddress = () => {
        if (formValid()) {
            setState({
                ...state,
                submitting: true,
            });
            fetch(`${HostManager.amperHost()}amper/fetchInstanceInfo`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(state.newInstanceForm)
            })
            .then(res => res.json())
            .then((result) => {
                setState({
                    ...state,
                    submitting: false,
                    newInstanceForm: {
                        ...state.newInstanceForm,
                        identifier: result.data?.identifier,
                        valid: result.success,
                    },
                    errorMessage: result.error,
                });
            });
        }
    };
    
    const confirmRemoveInstance = () => {
        fetch(`${HostManager.amperHost()}amper/remove`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
            body: JSON.stringify(state.selectedRowData[0])
        })
        .then(res => res.json())
        .then((result) => {
            if (result.success) {
                closeRemoveDialogReset()
            } else {
                setState({
                    ...state,
                    errorMessage: result.error,
                });                
            }
        });
    };

    const getRemoveInstanceDialog = () => {
        return <Dialog onClose={closeDialogReset} open={state.removeDialogOpen}>
        <DialogTitle>Remove amper instance</DialogTitle>
        <DialogContent>
            <DialogContentText>
                Are you sure you want to remove the selected '{state.selectedRowData[0]?.name}' amper instance ?
            </DialogContentText>
            <Box sx={{ width: '100%', visibility: state.submitting ? 'visible' : 'hidden' }}>
                <LinearProgress />
            </Box>
            <Box style={{color: 'red', visibility: state.errorMessage ? 'visible' : 'hidden', marginTop: '4px'}}>{state.errorMessage}</Box>
        </DialogContent>
        <DialogActions>
            <Button onClick={closeRemoveDialogReset}>Cancel</Button>
            <Button onClick={confirmRemoveInstance} autoFocus >Remove</Button>
        </DialogActions>
    </Dialog>;
    };
    
    const getNewInstanceDialog = () => {
        return <Dialog onClose={closeDialogReset} open={newInstanceDialogOpen}>
            <DialogTitle>New amper instance</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    The dialog can be used to add a new amper or amper datastore instance, fill in all required fields to create a new amper instance.
                </DialogContentText>
                <DialogContentText>
                    Once address and port are entered, press Test button, to verify the amper instance connection.
                </DialogContentText>
                <Box sx={{ width: '100%', visibility: state.submitting ? 'visible' : 'hidden' }}>
                    <LinearProgress />
                </Box>
                <Box sx={{ height: '100%', width: '100%' }}>
                    <TextField disabled={state.submitting} name="name" label="Name" value={state.newInstanceForm.name} onChange={onChange} required error={!Convenience.hasValue(state.newInstanceForm.name)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <Select
                        name="type"
                        disabled={updateInstanceMode || state.submitting}
                        value={state.newInstanceForm.type}
                        label="Type"
                        onChange={onChange}
                        required error={!Convenience.hasValue(state.newInstanceForm.type)} variant="filled" fullWidth sx={{ mt: 1}}
                        >
                        <MenuItem value="amperInstance">Amper instance</MenuItem>
                        <MenuItem value="amperDatastoreInstance">Amper datastore instance</MenuItem>
                    </Select>
                    <TextField disabled={state.submitting} name="address" label="Address" value={state.newInstanceForm.address} onChange={onChange} required error={!Convenience.hasValue(state.newInstanceForm.address)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <TextField disabled={state.submitting} name="port" label="Port" value={state.newInstanceForm.port} onChange={onChange} required error={!Convenience.hasValue(state.newInstanceForm.port)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <TextField disabled={state.submitting} name="limit" type="number" label="Limit (in bytes)" value={state.newInstanceForm.limit} onChange={onChangeLimit} required error={!Convenience.hasValue(state.newInstanceForm.limit)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <TextField disabled={state.submitting} name="directory" label="Directory" value={state.newInstanceForm.directory} onChange={onChange} required error={!Convenience.hasValue(state.newInstanceForm.directory)} variant="filled" fullWidth sx={{ mt: 1}}/>
                </Box>
                <Box style={{color: 'red', visibility: state.errorMessage ? 'visible' : 'hidden', marginTop: '4px'}}>{state.errorMessage}</Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={closeDialogReset}>Cancel</Button>
                <Button onClick={testInstanceAddress} disabled={!formValid()|| state.submitting}>Test</Button>
                <Button onClick={updateInstanceMode ? updatInstance : addNewInstance} autoFocus disabled={!formValid() || !state.newInstanceForm.valid || state.submitting || state.newInstanceForm.identifier == null}>{updateInstanceMode ? 'Update' : 'Add'}</Button>
            </DialogActions>
        </Dialog>;
    };

    const newInstance = () => {
        navigate(breadcrumbs.administration.nodes.new.key);
    };

    const updateInstance = () => {
        navigate(breadcrumbs.administration.nodes.update.key);
    };

    const removeInstance = () => {
        setState({
            ...state,
            removeDialogOpen: true,
        });
    };

    
    const onFilterChange = (filterModel) => {
        setState({
            ...state,
            search: filterModel.quickFilterValues,
            loading: true,
        });
      };

      const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}amper/fetch`,
            requestMethod: "POST",
        });
    };

      const columns = [
        { field: 'id',
          headerName: 'ID',
          width: 90,
          sortable: false,
          hide: true,
        },
        {
          field: 'identifier',
          headerName: 'Identifier',
          width: 200,
        },{
          field: 'name',
          headerName: 'Name',
          width: 200,
        },{
          field: 'type',
          headerName: 'Type',
          width: 200,
        },{
          field: 'address',
          headerName: 'Address',
          width: 150,
        },{
          field: 'port',
          headerName: 'Port',
          width: 200,
        },{
          field: 'state',
          headerName: 'State',
          width: 150,
        },{
            field: 'usage',
            headerName: 'Usage',
            width: 150,
        },{
            field: 'limit',
            headerName: 'Limit',
            width: 150,
        },{
            field: 'directory',
            headerName: 'Directory',
            width: 250,
        }
      ];

    const CustomToolbar = () => {
        return (
            <Box sx={{
                display: 'flex',
                flexDirection: 'row',
                bgcolor: 'background.paper',
                }}>
                <Box sx={{ flexGrow: 1, mt: 1 }}>
                    <Button size="small" sx={{pl: 1}} onClick={newInstance} startIcon={<AddCircleOutlineIcon/>}>
                        New Instance
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={updateInstance} startIcon={<AdjustIcon/>} disabled={state.selectedRowData.length !== 1}>
                        Update
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={removeInstance} startIcon={<RemoveCircleOutlineIcon/>} disabled={state.selectedRowData.length !== 1}>
                        Remove
                    </Button>
                </Box>
                <Box sx={{ flexGrow: 0 }}>
                    <GridToolbarContainer>
                        <GridToolbarColumnsButton />
                        <GridToolbarDensitySelector />
                        <GridToolbarExport />
                    </GridToolbarContainer>
                </Box>
        </Box>
        );
      };

    return <Box sx={{ height: '100%', width: '100%' }}>
        {getNewInstanceDialog()}
        {getRemoveInstanceDialog()}
        <DataGridPremium
        disableColumnFilter
        loading={state.loading}
        slots={{
            loadingOverlay: LinearProgress,
            toolbar: CustomToolbar,
        }}
        slotProps={{
            toolbar: {
                showQuickFilter: false,
                quickFilterProps: { debounceMs: 500 },
            },
        }}
        onRowSelectionModelChange={(ids) => {
            const selectedIDs = new Set(ids);
            const rowData = state.data.filter((row) =>
                selectedIDs.has(row.id)
            )

            setState({
                ...state,
                selectedRowData: rowData,
            });
            setTimeout(() => {
                if (onSelect) {
                    onSelect(rowData);
                }
            }, 1);
        }}
        onFilterModelChange={onFilterChange}
        filterMode="client"
        sortingMode="client"
        rows={state.data}
        columns={columns}
        pageSize={Number.MAX_SAFE_INTEGER}
        checkboxSelection
        />
    </Box>
};
export default NodeList;