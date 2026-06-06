import React, { useEffect, } from 'react';
import useState from 'react-usestateref'
import Box from '@mui/material/Box';
import { DataGridPremium,GridToolbarContainer, GridToolbarQuickFilter } from '../../../../components/x-data-grid-premium';
import LinearProgress from '@mui/material/LinearProgress';
import Button from '@mui/material/Button';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import DialogActions from '@mui/material/DialogActions';
import TextField from '@mui/material/TextField';
import Convenience from '../../../../help/Convenience';
import HostManager from '../../../../../HostManager';
import { sessionManager } from '../../../../../SessionManager';
import DataStore from '../../../../data/DataStore';
import AutocompleteRemote from '../../../../components/fields/AutocompleteRemote';
import { AppContext } from '../../../../../App';
import UserSelect from '../../../../components/user/UserSelect';
import Avatar from '@mui/material/Avatar';

export default function ChatConfiguration() {
    const app = React.useContext(AppContext);
    const [state, setState, stateRef] = useState(() => {
        return {
            chatChannelGrouploading: true,
            chatChannelloading: false,
            chatChannelUsersloading: false,
            groupData: [],
            channelData: [],
            channelUserData: [],
            groupSelectedRowData: [],
            channelSelectedRowData: [],
            channelUserSelectedRowData: [],
            submitting: false,
            newGroupDialogOpen: false,
            newChannelDialogOpen: false,
            newChannelUserDialogOpen: false,
            newGroupForm: {},
            newChannelForm: {},
            newChannelUserForm: {},
            userSearch: '',
            userPaging: {
                page: 0,
                pageSize: 50,
            }
        };
    });

    useEffect(() => {
        if (state.chatChannelGrouploading) {
            getChatChannelGroupDataStore().load((result)=> {       
                setState({
                  ...stateRef.current,
                  chatChannelGrouploading: false,
                  groupData: result.data || [],
                })
            });
        }
    }, [state.chatChannelGrouploading]);

    useEffect(() => {
        if (state.chatChannelloading && state.groupSelectedRowData.length > 0) {
            getChatChannelDataStore().load((result)=> {       
                setState({
                  ...stateRef.current,
                  chatChannelloading: false,
                  channelData: result.data || [],
                })
            });
        }
    }, [state.chatChannelloading]);

    useEffect(() => {
        if (state.chatChannelUsersloading && state.channelSelectedRowData.length > 0) {
            getChatChannelUsersDataStore().load((result)=> {
                setState({
                  ...stateRef.current,
                  chatChannelUsersloading: false,
                  channelUserData: result.data || [],
                  channelUsersRowsCount: result.totalCount,
                });
            });
        }
    }, [state.chatChannelUsersloading]);

    const getChatChannelGroupDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}chat/fetchChatChannelGroups`,
            requestMethod: "POST",
        });
    };

    const getChatChannelDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}chat/fetchChatChannels`,
            requestMethod: "POST",
            parameters: {
                groupId: state.groupSelectedRowData.length > 0 ? parseInt(state.groupSelectedRowData[0].id) : null,
            }
        });
    };

    const getChatChannelUsersDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}chat/fetchChatChannelUsers`,
            requestMethod: "POST",
            parameters: {
                search: state.userSearch,
                channelId: state.channelSelectedRowData.length > 0 ? parseInt(state.channelSelectedRowData[0].id) : null,
                start: state.userPaging.page * state.userPaging.pageSize,
                limit: state.userPaging.pageSize,
            }
        });
    };

    const newChannelGroup = () => {
        setState({
            ...state,
            newGroupDialogOpen: true,
        });
    };

    const removeChannelGroup = () => {
        if (state.groupSelectedRowData.length > 0) {
            fetch(`${HostManager.amperHost()}chat/removeChatChannelGroup`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify({
                    groupId: parseInt(state.groupSelectedRowData[0].id),
                })
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    setState({
                        ...state,
                        chatChannelGrouploading: true,
                    })
                } else {
                    app.toast('error', "not able to remove the channel group, you must first remove all channels from the group")
                }
            });
        }
    };

    const removeChannel = () => {
        if (state.channelSelectedRowData.length > 0) {
            fetch(`${HostManager.amperHost()}chat/removeChatChannel`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify({
                    channelId: parseInt(state.channelSelectedRowData[0].id),
                })
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    setState({
                        ...state,
                        chatChannelloading: true,
                    })
                } else {
                    app.toast('error', "not able to remove the channel")
                }
            });
        }
    };

    const newChannel = () => {
        setState({
            ...state,
            newChannelDialogOpen: true,
        });
    };

    const CustomToolbarGroup = () => {
        return (
            <Box sx={{
                display: 'flex',
                flexDirection: 'row',
                bgcolor: 'background.paper',
                }}>
                <Box sx={{ flexGrow: 1, mt: 1 }}>
                    <Button size="small" sx={{pl: 1}} onClick={newChannelGroup} startIcon={<AddCircleOutlineIcon/>}>
                        New Channel Group
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={removeChannelGroup} startIcon={<RemoveCircleOutlineIcon/>}>
                        Remove
                    </Button>
                </Box>
        </Box>
        );
      };

      const CustomToolbarChannel = () => {
        return (
            <Box sx={{
                display: 'flex',
                flexDirection: 'row',
                bgcolor: 'background.paper',
                }}>
                <Box sx={{ flexGrow: 1, mt: 1 }}>
                    <Button size="small" sx={{pl: 1}} onClick={newChannel} startIcon={<AddCircleOutlineIcon/>}>
                        New Channel
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={removeChannel} startIcon={<RemoveCircleOutlineIcon/>}>
                        Remove
                    </Button>
                </Box>
        </Box>
        );
      };

      const newChannelUser = () => {
        setState({
            ...state,
            newChannelUserDialogOpen: true,
        });
      };

      const removeChannelUser = () => {
        if (state.channelSelectedRowData.length > 0 && state.channelUserSelectedRowData.length > 0) {
            fetch(`${HostManager.amperHost()}chat/removeChatChannelUser`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify({
                    channelId: parseInt(state.channelSelectedRowData[0].id),
                    userIds: state.channelUserSelectedRowData.map(user => user.id)
                })
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    setState({
                        ...state,
                        chatChannelUsersloading: true,
                    })
                } else {
                    app.toast('error', "not able to remove the channel user")
                }
            });
        }
      };

      const CustomToolbarChannelUsers = () => {
        return (
            <Box sx={{
                display: 'flex',
                flexDirection: 'row',
                bgcolor: 'background.paper',
                }}>
                <Box sx={{ flexGrow: 1, mt: 1 }}>
                    <Button size="small" sx={{pl: 1}} onClick={newChannelUser} disabled={state.channelSelectedRowData.length < 1} startIcon={<AddCircleOutlineIcon/>}>
                        New User
                    </Button>
                    <Button size="small" sx={{pl: 1}} onClick={removeChannelUser} disabled={state.channelUserSelectedRowData.length < 1 || state.channelSelectedRowData.length < 1} startIcon={<RemoveCircleOutlineIcon/>}>
                        Remove
                    </Button>
                </Box>
                <Box sx={{ flexGrow: 0 }}>
                    <GridToolbarContainer>
                        <GridToolbarQuickFilter sx={{pl: 1, width: 400,}} debounceMs={1000}/>
                    </GridToolbarContainer>
                </Box>
        </Box>
        );
      };

      const newGroupColumns = [
        { field: 'id',
          headerName: 'ID',
          width: 90,
          sortable: false,
          hide: true,
        },{
          field: 'name',
          headerName: 'Name',
          width: 200,
        }
      ];

      
      const newChannelColumns = [
        { field: 'id',
          headerName: 'ID',
          width: 90,
          sortable: false,
          hide: true,
        },{
          field: 'name',
          headerName: 'Name',
          width: 200,
        },{
            field: 'groupName',
            headerName: 'Group Name',
            width: 200,
        }, {
            field: 'groupId',
            headerName: 'Group Id',
            width: 80,
            hide: true,
        },{
            field: 'amperName',
            headerName: 'Amper Name',
            width: 200,
        }, {
            field: 'amperId',
            headerName: 'Amper Id',
            width: 80,
            hide: true,
        }
      ];

      const getImageSource = (base64Image) => {
        if (Convenience.hasValue(base64Image)) {
          return 'data:image/png;base64,' + base64Image;
        }
        return '/static/images/avatar/2.jpg';
      };

      const newChannelUserColumns = [
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
          }},{
          field: 'username',
          headerName: 'Username',
          width: 200,
        },{
            field: 'firstName',
            headerName: 'First Name',
            width: 200,
        }, {
            field: 'lastName',
            headerName: 'Last Name',
            width: 200,
        }
      ];

    const closeDialogReset = () => {
        setState({
            ...state,
            newGroupDialogOpen: false,
            errorMessage: '',
        });
    };

    const onNewGroupChange = (event, arg1) => {
        const {
            target: { value, name },
          } = event;
          setState({
            ...state,
            newGroupForm: {
                ...state.newGroupForm,
                [name]: value,
            },
        });
    };

    const createNewGroup = () => {
        if (Convenience.hasValue(state.newGroupForm.name)) {
            setState({
                ...state,
                submitting: true,
            });
            fetch(`${HostManager.amperHost()}chat/createChatGroup`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(state.newGroupForm)
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    setState({
                        ...state,
                        newGroupForm: {},
                        submitting: false,
                        newGroupDialogOpen: false,
                        chatChannelGrouploading: true,
                    })
                } else {
                    setState({
                        ...state,
                        errorMessage: result.error,
                    });                
                }
            });
        }
    };

    const getNewIChatChannelGroupDialog = () => {
        return <Dialog onClose={closeDialogReset} open={state.newGroupDialogOpen}>
            <DialogTitle>New Channel Group</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    The dialog can be used to create new channel groups, it can be used later to group the channels in the chat
                </DialogContentText>
                <Box sx={{ width: '100%', visibility: state.submitting ? 'visible' : 'hidden' }}>
                    <LinearProgress />
                </Box>
                <Box sx={{ height: '100%', width: '100%' }}>
                    <TextField disabled={state.submitting} name="name" label="Name" value={state.newGroupForm.name} onChange={onNewGroupChange} required error={!Convenience.hasValue(state.newGroupForm.name)} variant="filled" fullWidth sx={{ mt: 1}}/>
                </Box>
                <Box style={{color: 'red', visibility: state.errorMessage ? 'visible' : 'hidden', marginTop: '4px'}}>{state.errorMessage}</Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={closeDialogReset}>Cancel</Button>
                <Button onClick={createNewGroup} autoFocus disabled={!Convenience.hasValue(state.newGroupForm.name)}>Create</Button>
            </DialogActions>
        </Dialog>;
    };

    const closeNewChannelDialogReset = () => {
        setState({
            ...state,
            newChannelDialogOpen: false,
            newChannelForm: {},
            amperId: null,
            groupId: null,
            errorMessage: '',
        });
    };

    const onNewChannelChange = (event, arg1) => {
        const {
            target: { value, name },
          } = event;
          setState({
            ...state,
            newChannelForm: {
                ...state.newChannelForm,
                [name]: value,
            },
        });
    };

    const onAutocompleteChange = (name, event, value) => {
        setState({
            ...state,
            newChannelForm: {
                ...state.newChannelForm,
                [name]: value ? parseInt(value.id) : null
            },
            [name]: value
        });
    };

    const createNewChannel = () => {
        if (Convenience.hasValue(state.newChannelForm.name) && Convenience.hasValue(state.newChannelForm.amperId) && Convenience.hasValue(state.newChannelForm.groupId)) {
            setState({
                ...state,
                submitting: true,
            });
            fetch(`${HostManager.amperHost()}chat/createChatChannel`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(state.newChannelForm)
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    setState({
                        ...state,
                        newChannelForm: {},
                        amperId: null,
                        groupId: null,
                        submitting: false,
                        newChannelDialogOpen: false,
                        chatChannelloading: state.groupSelectedRowData.length > 0,
                    })
                } else {
                    setState({
                        ...state,
                        errorMessage: result.error,
                    });                
                }
            });
        }
    };

    const getNewChatChannelDialog = () => {
        return <Dialog onClose={closeNewChannelDialogReset} open={state.newChannelDialogOpen}>
            <DialogTitle>New Channel</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    The dialog can be used to create new channels, it can be used later to associate a user to channel
                </DialogContentText>
                <Box sx={{ width: '100%', visibility: state.submitting ? 'visible' : 'hidden' }}>
                    <LinearProgress />
                </Box>
                <Box sx={{ height: '100%', width: '100%' }}>
                    <TextField disabled={state.submitting} name="name" label="Name" value={state.newChannelForm.name} onChange={onNewChannelChange} required error={!Convenience.hasValue(state.newChannelForm.name)} variant="filled" fullWidth sx={{ mt: 1}}/>
                    <AutocompleteRemote name="amperId" value={state.amperId} required={true} error={state.amperId == null} label="Amper instance" onChange={onAutocompleteChange} url={`${HostManager.amperHost()}amper/fetch`} parameters={{start: 0, limit: Number.MAX_SAFE_INTEGER, type: 'amperInstance'}} keyIdentifier="id" labelIdentifier="name"/>
                    <AutocompleteRemote name="groupId" value={state.groupId} required={true} error={state.groupId == null} label="Channel Group" onChange={onAutocompleteChange} url={`${HostManager.amperHost()}chat/fetchChatChannelGroups`} parameters={{start: 0, limit: Number.MAX_SAFE_INTEGER}} keyIdentifier="id" labelIdentifier="name"/>
                </Box>
                <Box style={{color: 'red', visibility: state.errorMessage ? 'visible' : 'hidden', marginTop: '4px'}}>{state.errorMessage}</Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={closeNewChannelDialogReset}>Cancel</Button>
                <Button onClick={createNewChannel} autoFocus disabled={!Convenience.hasValue(state.newChannelForm.name) || !Convenience.hasValue(state.newChannelForm.amperId) || !Convenience.hasValue(state.newChannelForm.groupId)}>Create</Button>
            </DialogActions>
        </Dialog>;
    };

    const closeNewChannelUserDialogReset = () => {
        setState({
            ...state,
            newChannelUserDialogOpen: false,
            newChannelUserForm: {},
            errorMessage: '',
        });
    };

    const addNewChannelUser = () => {
        if (state.newChannelUserForm.userIds.length > 0 && state.channelSelectedRowData.length > 0) {
            state.newChannelUserForm.channelId = parseInt(state.channelSelectedRowData[0].id);
            setState({
                ...state,
                submitting: true,
            });
            fetch(`${HostManager.amperHost()}chat/addUsersToChannel`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json', sessionId: sessionManager.getSessionId()},
                body: JSON.stringify(state.newChannelUserForm)
            })
            .then(res => res.json())
            .then((result) => {
                if (result.success) {
                    setState({
                        ...state,
                        newChannelUserForm: {},
                        submitting: false,
                        newChannelUserDialogOpen: false,
                        chatChannelUsersloading: state.channelSelectedRowData.length > 0,
                    });
                } else {
                    setState({
                        ...state,
                        errorMessage: result.error,
                    });                
                }
            });
        }
    };

    const onUserelectionChange = (users) => {
        const userIds = users.map(user => user.id);
        setState({
            ...state,
            newChannelUserForm: {
                userIds: userIds,
            },
        })
    };

    const getNewChatChannelUserDialog = () => {
        return <Dialog onClose={closeNewChannelUserDialogReset} open={state.newChannelUserDialogOpen}>
            <DialogTitle>New Channel</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    The dialog can be used to add new user to the selected channel, the selected users will be subscribed to the given channel
                </DialogContentText>
                <Box sx={{ width: '100%', visibility: state.submitting ? 'visible' : 'hidden' }}>
                    <LinearProgress />
                </Box>
                <Box sx={{ height: '100%', width: '100%' }}>
                    <UserSelect onSelectionChange={onUserelectionChange} includeSelf={true}></UserSelect>
                </Box>
                <Box style={{color: 'red', visibility: state.errorMessage ? 'visible' : 'hidden', marginTop: '4px'}}>{state.errorMessage}</Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={closeNewChannelUserDialogReset}>Cancel</Button>
                <Button onClick={addNewChannelUser} autoFocus disabled={!(state.newChannelUserForm.userIds != null && state.newChannelUserForm.userIds.length > 0)}>Add</Button>
            </DialogActions>
        </Dialog>;
    };

    const onPageSizeChange = (pageSize) => {debugger
        setState({
            ...state,
            chatChannelUsersloading: true,
            userPaging: {
                ...state.userPaging,
                limit: pageSize,
            }
        });
    };

    const onPageChange = (page) => {
        setState({
            ...state,
            chatChannelUsersloading: true,
            userPaging: {
                ...state.userPaging,
                start: page * state.limit,
            } 
        });
    };

    const setUserPaginationModel = (pagingModel) => {
        setState({
            ...state,
            chatChannelUsersloading: state.channelSelectedRowData.length > 0,
            userPaging: pagingModel
        });
    };

    const onFilterChange = (filterModel) => {
        setState({
            ...state,
            userSearch: filterModel.quickFilterValues,
            chatChannelUsersloading: true,
        });
      };

  return (
        <Box sx={{ height: '100%', width: 'calc(100% - 25px)'}}>
            {getNewIChatChannelGroupDialog()}
            {getNewChatChannelDialog()}
            {getNewChatChannelUserDialog()}
            <Box style={{height: '100%', width: '100%', display: 'flex', flexDirection: 'row', justifyContent: 'center'}}>
                <Box sx={{ height: '100%', width: '300px', mr: '5px'}}>
                    <DataGridPremium
                        initialState={{
                            columns: {
                            columnVisibilityModel: {
                                // Hide column id, the other columns will remain visible
                                id: false,
                            },
                            },
                        }}
                        disableColumnFilter
                        loading={state.chatChannelGrouploading}
                        slots={{
                            loadingOverlay: LinearProgress,
                            toolbar: CustomToolbarGroup,
                        }}
                        slotProps={{
                            toolbar: {
                                showQuickFilter: true,
                                quickFilterProps: { debounceMs: 500 },
                            },
                        }}
                        onRowSelectionModelChange={(ids) => {
                            const selectedIDs = new Set(ids);
                            const rowData = state.groupData.filter((row) =>
                                selectedIDs.has(row.id)
                            )

                            setState({
                                ...state,
                                groupSelectedRowData: rowData,
                                channelSelectedRowData: [],
                                chatChannelloading: true,
                                channelUserData: [],
                            });
                        }}
                        filterMode="client"
                        sortingMode="client"
                        rows={state.groupData}
                        columns={newGroupColumns}
                        pageSize={Number.MAX_SAFE_INTEGER}/>
                </Box>
                <Box sx={{ height: '100%', width: '400px', mr: '5px'}}>
                    <DataGridPremium
                        initialState={{
                            columns: {
                            columnVisibilityModel: {
                                // Hide columns id, groupId and traderName, the other columns will remain visible
                                id: false,
                                groupId: false,
                                amperId: false,
                            },
                            },
                        }}
                        disableColumnFilter
                        loading={state.chatChannelloading}
                        slots={{
                            loadingOverlay: LinearProgress,
                            toolbar: CustomToolbarChannel,
                        }}
                        slotProps={{
                            toolbar: {
                                showQuickFilter: true,
                                quickFilterProps: { debounceMs: 500 },
                            },
                        }}
                        onRowSelectionModelChange={(ids) => {
                            const selectedIDs = new Set(ids);
                            const rowData = state.channelData.filter((row) =>
                                selectedIDs.has(row.id)
                            )

                            setState({
                                ...state,
                                channelSelectedRowData: rowData,
                                chatChannelUsersloading: true,
                            });
                        }}
                        filterMode="client"
                        sortingMode="client"
                        rows={state.channelData}
                        columns={newChannelColumns}
                        pageSize={Number.MAX_SAFE_INTEGER}/>
                </Box>
                <Box sx={{ height: '100%', width: 'calc(100% - 710px)'}}>
                    <DataGridPremium
                        initialState={{
                            columns: {
                            columnVisibilityModel: {
                                // Hide column id, the other columns will remain visible
                                id: false,
                            },
                            },
                        }}
                        disableColumnFilter
                        loading={state.chatChannelUsersloading}
                        slots={{
                            loadingOverlay: LinearProgress,
                            toolbar: CustomToolbarChannelUsers,
                        }}
                        slotProps={{
                            toolbar: {
                                showQuickFilter: true,
                                quickFilterProps: { debounceMs: 1000 },
                            },
                        }}
                        onRowSelectionModelChange={(ids) => {
                            const selectedIDs = new Set(ids);
                            const rowData = state.channelUserData.filter((row) =>
                                selectedIDs.has(row.id)
                            )

                            setState({
                                ...state,
                                channelUserSelectedRowData: rowData,
                            });
                        }}
                        paginationMode="server"
                        filterMode="server"
                        sortingMode="server"
                        rows={state.channelUserData}
                        columns={newChannelUserColumns}
                        pagination
                        pageSizeOptions={[50, 100, 500, 1000, 5000, 50000]}
                        rowCount={state.channelUsersRowsCount}
                        onPaginationModelChange={setUserPaginationModel}
                        paginationModel={state.userPaging}
                        onFilterModelChange={onFilterChange}
                        checkboxSelection
                        disableColumnSorting/>
                </Box>
            </Box>
        </Box>
    );
}
