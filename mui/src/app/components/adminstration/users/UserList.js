import React, { useEffect, useState, } from 'react';
import { useNavigate, useLocation } from 'react-router-dom'
import { DataGridPremium, GridToolbarContainer,
    GridToolbarColumnsButton,
    GridToolbarExport,
    GridToolbarDensitySelector, GridToolbarQuickFilter, } from '../../x-data-grid-premium';
import Box from '@mui/material/Box';
import LinearProgress from '@mui/material/LinearProgress';
import DataStore from '../../../data/DataStore';
import HostManager from '../../../../HostManager';
import Button from '@mui/material/Button';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import { breadcrumbs } from '../../../center/Breadcrambs';
import { AppContext } from '../../../../App';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import DialogActions from '@mui/material/DialogActions';
import TextField from '@mui/material/TextField';
import AutocompleteRemote from '../../fields/AutocompleteRemote';
import Convenience from '../../../help/Convenience';
import { post } from '../../../data/Submit';
import { parseBoolean } from '../../../util/BooleanUtil';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import Avatar from '@mui/material/Avatar';
import AdjustIcon from '@mui/icons-material/Adjust';

const UserList = ({onSelect, showUserManagement, checkboxSelection}) => {
    const app = React.useContext(AppContext);
    const [state, setState] = useState(() => {
        return {
            loading: true,
            data: [],
            search: [],
            newUserForm: {},
            selectedRowData: [],
            userPaging: {
                page: 0,
                pageSize: 50,
            }
        };
    });
    const {pathname} = useLocation();
    const navigate = useNavigate();
    const newUserPath = breadcrumbs.administration.users.new.path;
    const newUserDialogOpen = pathname.startsWith(newUserPath)
    useEffect(() => {
        if (state.loading) {
            getDataStore().load((result)=> {       
                setState({
                  ...state,
                  loading: false,
                  data: result.data || [],
                  totalCount: result.totalCount,
                })
            });
        }
    }, [state.loading]);
    if (app) {
        app.registerRefresh("users", () => {
            refresh();
        });    
    }

    const refresh = () => {
        setState({
            ...state,
            loading: true,
        });
    };

    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}users/fetch`,
            requestMethod: "POST",
            parameters: {
                start: state.userPaging.page * state.userPaging.pageSize,
                limit: state.userPaging.pageSize,
                search: state.search,
                sorfField: state.sorfField,
                sortDirection: state.sortDirection,
            }
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
          hide: true,
        },{
          field: 'photo',
          headerName: 'Photo',
          sortable: false,
          width: 55,
          renderCell: (props) => {
            const { hasFocus, value } = props;
            //TODO render image based on value
            return  <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 1 }} src={getImageSource(value)} />;
        },
        },
        {
          field: 'firstName',
          headerName: 'First name',
          width: 200,
        },{
          field: 'lastName',
          headerName: 'Last name',
          width: 200,
        },{
          field: 'middleName',
          headerName: 'Middle name',
          width: 150,
          hide: true,
        },{
          field: 'username',
          headerName: 'Username',
          width: 150,
        },{
          field: 'email',
          headerName: 'Email',
          width: 200,
        },{
          field: 'profileName',
          headerName: 'Profile',
          width: 150,
        },{
          field: 'amperId',
          headerName: 'Amper instance',
          sortable: false,
          width: 150,
        },{
            field: 'active',
            headerName: 'Active',
            renderCell: (props) => {
                const { hasFocus, value } = props;
                return parseBoolean(value) ? <CheckIcon color="primary"/> : <CloseIcon color="primary"/>;
            },
            sortable: false,
            width: 80,
          }
      ];

      const CustomToolbar = () => {
        return (
            <Box sx={{
                display: 'flex',
                flexDirection: 'row',
                bgcolor: 'background.paper',
                }}>
                {showUserManagement !== false && <Box sx={{ flexGrow: 1, mt: 1 }}>
                    <Button size="small" sx={{pl: 1}} onClick={newUser} startIcon={<AddCircleOutlineIcon/>}>
                        New User
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={updateUser} startIcon={<AdjustIcon/>} disabled={state.selectedRowData.length < 1}>
                        Update
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={removeUser} startIcon={<RemoveCircleOutlineIcon/>}>
                        Remove
                    </Button>
                </Box>}
                <Box sx={{ flexGrow: 0 }}>
                    <GridToolbarContainer>
                        <GridToolbarColumnsButton />
                        <GridToolbarDensitySelector />
                        <GridToolbarExport />
                        <GridToolbarQuickFilter sx={{pl: 1, width: 400,}} debounceMs={1000}/>
                    </GridToolbarContainer>
                </Box>
        </Box>
        );
      }

      const onFilterChange = (filterModel) => {
        setState({
            ...state,
            search: filterModel.quickFilterValues,
            loading: true,
        });
      };

      const handleSortModelChange = (sortModel) => {
        let sorfField = null;
        let sortDirection = null;
        if (sortModel.length > 0) {
            sorfField = sortModel[0].field;
            sortDirection = sortModel[0].sort;
        }
        setState({
            ...state,
            loading: true,
            sorfField: sorfField,
            sortDirection: sortDirection,
        });
      };

      const newUser = () => {
        navigate(breadcrumbs.administration.users.new.key);
      };

      const updateUser = () => {
        if (state.selectedRowData.length > 0) {
            navigate(breadcrumbs.administration.users.new.key);
            setState({
                ...state,
                manageUserMode: 'update',
                newUserForm: state.selectedRowData[0],
            });
        }
      };

    const handleNewUserDialogClose = () => {
        navigate(breadcrumbs.administration.users.path);
    };

    const formValid = () => {
        const form = state.newUserForm;
        return Convenience.hasValue(form.firstName) && Convenience.hasValue(form.lastName)
            && Convenience.hasValue(form.username) && Convenience.hasValue(form.email)
            && Convenience.hasValue(form.profileId) && Convenience.hasValue(form.amperId);
    };

    const removeUser = () => {
        if (state.selectedRowData.length > 0) {
            post(`${HostManager.amperHost()}users/remove`, {
                userId: state.selectedRowData[0].id,
            }, (result) => {
                setState({
                    ...state,
                    loading: true,
                });
                handleNewUserDialogClose();
              }, (result) => {
                setState({
                    ...state,
                    loading: true,
                });
              });
        }
    };

    const addNewUser = () => {
        const form = state.newUserForm;
        if (formValid()) {
            setState({
                ...state,
                submitting: true,
            });
            post(`${HostManager.amperHost()}${state.manageUserMode == 'update' ? 'users/edit' : 'users/create'}`, {
                ...form,
              }, (result) => {
                setState({
                    ...state,
                    loading: true,
                    submitting: false,
                });
                handleNewUserDialogClose();
              }, (result) => {
                setState({
                    ...state,
                    submitting: false,
                    errorMessage: result.error,
                });
              });
        }
    };

    const onAutocompleteChange = (name, event, value) => {
        setState({
            ...state,
            newUserForm: {
                ...state.newUserForm,
                [name]: value ? value.id : null
            },
            [name]: value
        });
    };

    const onChange =(event, arg1) => {
        const {
            target: { value, name },
          } = event;
          setState({
            ...state,
            newUserForm: {
                ...state.newUserForm,
                [name]: value
            },
        });
    };

    const closeDialogReset = () => {
        handleNewUserDialogClose();
        setState({
            ...state,
            manageUserMode: null,
            errorMessage: '',
        });
    };

    const getNewUserDialog = () => {
        return <Dialog onClose={closeDialogReset} open={newUserDialogOpen}>
            <DialogTitle>New user</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    The dialog can be used to add a new user, fill in all required fields to create a new amper user, later after activating the account, user will be able to configure the account in a more specific details.
                </DialogContentText>
                <Box sx={{ width: '100%', visibility: state.submitting ? 'visible' : 'hidden' }}>
                    <LinearProgress />
                </Box>
                <Box sx={{ height: '100%', width: '100%' }}>
                    <TextField name="firstName" label="First name" value={state.newUserForm.firstName} onChange={onChange} required error={!Convenience.hasValue(state.newUserForm.firstName)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <TextField name="lastName" label="Last name" value={state.newUserForm.lastName} onChange={onChange} required error={!Convenience.hasValue(state.newUserForm.lastName)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <TextField name="middleName" label="Middle name" value={state.newUserForm.middleName} onChange={onChange} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <TextField name="username" label="Username" value={state.newUserForm.username} onChange={onChange} required error={!Convenience.hasValue(state.newUserForm.username)} disabled={state.manageUserMode == 'update'} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <TextField name="email" label="Email" value={state.newUserForm.email} onChange={onChange} required error={!Convenience.hasValue(state.newUserForm.email)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <AutocompleteRemote name="profileId" value={state.profileId} required={true} error={state.profileId == null} label="Profile" onChange={onAutocompleteChange} url={`${HostManager.amperHost()}profiles/fetch`} parameters={{start: 0, limit: Number.MAX_SAFE_INTEGER}} keyIdentifier="id" labelIdentifier="name"/>
                    <AutocompleteRemote name="amperId" value={state.amperId} required={true} error={state.amperId == null} label="Amper instance" onChange={onAutocompleteChange} url={`${HostManager.amperHost()}amper/fetch`} parameters={{start: 0, limit: Number.MAX_SAFE_INTEGER, type: 'amperInstance'}} keyIdentifier="id" labelIdentifier="name"/>
                </Box>
                <Box style={{color: 'red', visibility: state.errorMessage ? 'visible' : 'hidden', marginTop: '4px'}}>{state.errorMessage}</Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={closeDialogReset}>Cancel</Button>
                <Button onClick={addNewUser} autoFocus disabled={!formValid() || state.submitting}>{state.manageUserMode == 'update' ? 'Update' : 'Add'}</Button>
            </DialogActions>
        </Dialog>;
    };

    const setUserPaginationModel = (pagingModel) => {
        setState({
            ...state,
            loading: true,
            userPaging: pagingModel
        });
    };

    return <Box sx={{ height: '100%', width: '100%' }}>
        {getNewUserDialog()}
        <DataGridPremium
        disableColumnFilter
        loading={state.loading}
        slots={{
            loadingOverlay: LinearProgress,
            toolbar: CustomToolbar,
        }}
        slotProps={{
            toolbar: {
                showQuickFilter: true,
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
                onSelect(rowData);
            }, 1);
        }}
        onFilterModelChange={onFilterChange}
        filterMode="server"
        sortingMode="server"
        paginationMode="server"
        onSortModelChange={handleSortModelChange}
        rows={state.data}
        rowCount={state.totalCount}
        columns={columns}
        pageSizeOptions={[50, 100, 500, 1000, 5000, 50000]}
        paginationModel={state.userPaging}
        onPaginationModelChange={setUserPaginationModel}
        pagination
        checkboxSelection={checkboxSelection === false ? false : true}
        />
    </Box>;
};
export default UserList;