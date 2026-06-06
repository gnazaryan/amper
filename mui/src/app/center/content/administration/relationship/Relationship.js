import React, { useEffect, useState, } from 'react';
import Box from '@mui/material/Box';
import UserList from '../../../../components/adminstration/users/UserList';
import { DataGridPremium, GridToolbarContainer,
    GridToolbarColumnsButton,
    GridToolbarExport,
    GridToolbarDensitySelector, GridToolbarQuickFilter, } from '../../../../components/x-data-grid-premium';
import LinearProgress from '@mui/material/LinearProgress';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import Button from '@mui/material/Button';
import HostManager from '../../../../../HostManager';
import DataStore from '../../../../data/DataStore';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import DialogActions from '@mui/material/DialogActions';
import UserSelect from '../../../../components/user/UserSelect';
import { sessionManager } from '../../../../../SessionManager';
import Avatar from '@mui/material/Avatar';
import Convenience from '../../../../help/Convenience';
import { AppContext } from '../../../../../App';

export default function Relationship() {
    const app = React.useContext(AppContext);
    const [state, setState] = useState(() => {
        return {
            loading: false,
            data: [],
            selectedRowData: [],
            search: [],
            employeeId: null,
            managerId: null,
            userDialogOpen: false,
            submitting: false,
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

    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}users/getUserRelationships`,
            requestMethod: "POST",
            parameters: {
                employeeId: state.employeeId,
            }
        });
    };

    const onUserSelect = (user) => {
        setState({
            ...state,
            loading: true,
            employeeId: user[0].id,
        });
    };

    const newRelationship = () => {
        setState({
            ...state,
            userDialogOpen: true,
        });
    };

    const removeRelationship = () => {
        if (state.selectedRowData.length > 0) {
            fetch(`${HostManager.amperHost()}users/deleteUserRelationship`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify({
                    employeeId: state.selectedRowData[state.selectedRowData.length - 1].employeeId,
                    managerId: state.selectedRowData[state.selectedRowData.length - 1].managerId,
                })
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    setState({
                        ...state,
                        loading: true,
                        selectedRowData: [],
                    });
                } else {
                    app.toast('info', `not able to remove the manager, please try again or contuct the support`);              
                }
            });
        }
    };
    
    const onFilterChange = (filterModel) => {
        setState({
            ...state,
            search: filterModel.quickFilterValues,
        });
    };

    const getImageSource = (base64Image) => {
        if (Convenience.hasValue(base64Image)) {
          return 'data:image/png;base64,' + base64Image;
        }
        return '/static/images/avatar/2.jpg';
      };

    const columns = [
        { field: 'id',
          headerName: 'ID',
          width: 90,
          sortable: false,
          hiden: true,
        },{
            field: 'managerPhoto',
            headerName: 'Manager Photo',
            sortable: false,
            width: 55,
            renderCell: (props) => {
                const { hasFocus, value } = props;
                return  <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 1 }} src={getImageSource(value)} />;
            },
        },{
          field: 'managerId',
          headerName: 'Manager Id',
          width: 200,
        },{
          field: 'managerFirstName',
          headerName: 'Manager First Name',
          width: 200,
        },{
          field: 'managerLastName',
          headerName: 'Manager Last Name',
          width: 200,
        },{
            field: 'employeePhoto',
            headerName: 'Employee Photo',
            width: 60,
        },{
          field: 'employeeId',
          headerName: 'Employee Id',
          width: 200,
        },{
          field: 'employeeFirstName',
          headerName: 'Employee First Name',
          width: 200,
        },{
          field: 'employeeLastName',
          headerName: 'Employee Last Name',
          width: 200,
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
                    <Button size="small" sx={{pl: 1}} onClick={newRelationship} startIcon={<AddCircleOutlineIcon/>} disabled={state.employeeId == null}>
                        New Manager
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={removeRelationship} startIcon={<RemoveCircleOutlineIcon/>} disabled={state.selectedRowData.length !== 1}>
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

      const closeUserDialog = () => {
        setState({
            ...state,
            userDialogOpen: false,
            managerId: null,
        });
      };

      const onUserelectionChange = (users) => {
        let managerId = null;
        if (users.length > 0) {
            managerId = users[0].id;
        }
        setState({
            ...state,
            managerId: managerId,
        });
      };

      const createNewRelationship = () => {
        setState({
            ...state,
            submitting: true,
        });
        fetch(`${HostManager.amperHost()}users/createUserRelationship`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
            body: JSON.stringify({
                employeeId: state.employeeId,
                managerId: state.managerId,
            })
        })
        .then(res => res.json())
        .then((result) => {
            if (result.success) {
                setState({
                    ...state,
                    submitting: false,
                    loading: true,
                    managerId: null,
                    userDialogOpen: false,
                });
            } else {
                setState({
                    ...state,
                    errorMessage: result.error,
                });                
            }
        });
      };

      const getNewRelationshipUserDialog = () => {
        return <Dialog onClose={closeUserDialog} open={state.userDialogOpen}>
            <DialogTitle>Select manager</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    The dialog can be used to select a manager for the choosen employee, use the combo box to slect a user as a manager
                </DialogContentText>
                <Box sx={{ width: '100%', visibility: state.submitting ? 'visible' : 'hidden' }}>
                    <LinearProgress />
                </Box>
                <Box sx={{ height: '100%', width: '100%' }}>
                    <UserSelect onSelectionChange={onUserelectionChange} includeSelf={true} singleSelect={true}></UserSelect>
                </Box>
                <Box style={{color: 'red', visibility: state.errorMessage ? 'visible' : 'hidden', marginTop: '4px'}}>{state.errorMessage}</Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={closeUserDialog}>Cancel</Button>
                <Button onClick={createNewRelationship} autoFocus disabled={state.managerId == null || state.submitting}>Add manager</Button>
            </DialogActions>
        </Dialog>;
    };

  return (
        <Box sx={{ height: '100%', width: 'calc(100% - 25px)', display: 'flex', flexDirection: 'row'}}>
            <Box sx={{display: 'flex', height: '100%', width: '50%'}}>
                <UserList showUserManagement={false} onSelect={onUserSelect} checkboxSelection={false}></UserList>
            </Box>
            <Box sx={{display: 'flex', height: '100%', width: '50%', ml: 1}}>
            {getNewRelationshipUserDialog()}
            <DataGridPremium
                disableColumnFilter
                loading={state.loading}
                initialState={{
                    columns: {
                      columnVisibilityModel: {
                        id: false,
                        managerId: false,
                        employeeId: false,
                      },
                    },
                  }}
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
                }}
                onFilterModelChange={onFilterChange}
                filterMode="client"
                sortingMode="client"
                rows={state.data}
                columns={columns}
                pageSize={Number.MAX_SAFE_INTEGER}
                checkboxSelection={false}
                />
            </Box>
        </Box>
    );
}
