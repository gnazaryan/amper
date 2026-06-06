import React, { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import { sessionManager } from '../../../SessionManager';
import Fab from '@mui/material/Fab';
import Convenience from '../../help/Convenience';
import { DataGridPremium, GridToolbarContainer, useGridApiContext } from '../x-data-grid-premium';
import Grid from '@mui/material/Grid2';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import { uniqueId } from '../../amper/Instruments';
import { post } from '../../data/Submit'
import { AppContext } from '../../../App';
import HostManager from '../../../HostManager';
import DialogTitle from '@mui/material/DialogTitle';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import CircularProgress from '@mui/material/CircularProgress';
import makeStyles from '@mui/styles/makeStyles';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';

const SyncAllSwitchStyle = makeStyles((theme) => ({  
    switch_track: {
        backgroundColor: "#4CAF50",
    },
    switch_base: {
        color: "#4CAF50",
        "&.Mui-checked": {
            color: "#2196f3"
        },
        "&.Mui-checked + .MuiSwitch-track": {
            backgroundColor: "#2196f3",
        }
    },
    switch_primary: {
        "&.Mui-checked": {
            color: "#2196f3",
        },
        "&.Mui-checked + .MuiSwitch-track": {
            backgroundColor: "#2196f3",
        },
    },
}));

function PasswordEditInputCell(props) {
    const { id, value, field } = props;
    const apiRef = useGridApiContext();

    const handleChange = (event, newValue) => {
      apiRef.current.setEditCellValue({ id, field, value: event.target.value });
      props.row.password = event.target.value;
    };
  
    const handleRef = (element) => {
      if (element) {
        const input = element.querySelector(`input[value="${value}"]`);
  
        input?.focus();
      }
    };
  
    return (
      <Box sx={{ display: 'flex', alignItems: 'center'}}>
        <TextField
          ref={handleRef}
          name="password"
          type="password"
          value={value}
          onChange={handleChange}
        />
      </Box>
    );
  }

export default function Settings(props) {
    const app = React.useContext(AppContext);

    const [state, setState] = useState({
        data: props.data ? props.data : {},
        selectedRowData: [],
        addEmail: false,
        addEmailStep: 0,
        addEmailCredentials: {},
        addEmailMailboxes: [],
    });

    const configureHandler = (result) => {
        if (!result.success) {
            app.toast('warning', result.error)
            setState({
                ...state,
                addEmail: false,
            });
        } else {
            setState({
                ...state,
                addEmailStep: 2,
                addEmailMailboxes: result.data
            });
        }
    };
    
    React.useEffect(() => {
        if (state.addEmailStep == 1) {
            post(`${HostManager.myHost()}email/configure`, {
              email: state.addEmailCredentials.email,
              password: state.addEmailCredentials.password
            }, configureHandler, configureHandler);
        }
      }, [state.addEmailStep]);

    useEffect(() => {
        setState({
            ...state,
            data: props.data,
        });
    }, [props.data]);

    const save = () => {
        post(`${HostManager.amperHost()}profile/saveConfiguration`, {
            name: 'settings',
            value: {
                email: getEmailSettings(),
            },
        }, (result) => {
            if (result.success) {
                setState({
                    ...state,
                });
                app.toast('info', "settings saved successfully")
            } else {
                app.toast('warning', result.error)
            }
        }, (result) => {
            if (result) {
                app.toast('warning', result.error)
            }
        });
    };

    const getEmailSettingsColumns = () => {
        return [
            {
                field: 'id',
                headerName: 'Id',
                width: 20,
                editable: false,
                hide: true,
            },
            {
                field: 'label',
                headerName: 'Label',
                width: 150,
                editable: true,
            },
            {
              field: 'email',
              headerName: 'Email',
              width: 150,
              editable: true,
            },
            {
              field: 'password',
              headerName: 'Password',
              renderCell: () => "********",
              renderEditCell: renderEditInputCell,
              width: 150,
              editable: true,
            }];
    };

    const renderEditInputCell = (params) => {
        return <PasswordEditInputCell {...params}/>;
      };

    const addEmail = () => {
        /*setState({
            ...state,
            data: [
                ...(getEmailSettings()),
                ...[{id: uniqueId(), label: 'Email label', email: '', password: ''}]
            ]
        });*/
        setState({
            ...state,
            addEmail: true,
        });
    };

    const removeEmail = () => {
        const emails = getEmailSettings();

        for (let l = 0; l < state.selectedRowData.length; l++) {
            const index = emails.findIndex(x => x.id === state.selectedRowData[l].id);
            emails.splice(index, 1);
        }
        setState({
            ...state,
            data: {
                settings: {
                    ...state.data.settings,
                    email: emails,
                }
                
            },
        });
    };

    const getEmailSettings = () => {
        return state.data && state.data.settings && state.data.settings.email ? JSON.parse(JSON.stringify(state.data.settings.email)) : [];
    };

    const getToolbar = () => {
        return (
          <GridToolbarContainer>
            <Button variant="text" endIcon={<AddCircleOutlineIcon />} onClick={addEmail}>
                Add
            </Button>
            <Button variant="text" endIcon={<RemoveCircleOutlineIcon />} onClick={removeEmail}>
                Remove
            </Button>
          </GridToolbarContainer>
        );
      };

      const handleAddEmailClose = () => {
        setState({
            ...state,
            addEmail: false,
            addEmailStep: 0,
        });
      };

      const handleAddEmail = () => {
        const emails = JSON.parse(JSON.stringify(state.data.settings.email));
        for (let i = 0; i < emails.length; i++) {
            if (emails[i].email === state.addEmailCredentials.email) {
                emails.splice(i, 1);
            }
        }
        const email = {
            id: uniqueId(),
            label: state.addEmailCredentials.label,
            email: state.addEmailCredentials.email,
            password: state.addEmailCredentials.password,
            mailboxes: state.addEmailMailboxes,
        };
        emails.push(email);
        setState({
            ...state,
            data: {
                ...state.data,
                settings: {
                    ...state.data.settings,
                    email: emails,
                }
            },
            addEmail: false,
            addEmailStep: 0,
        });
      };

      const handleAddEmailNext = () => {
        setState({
            ...state,
            addEmailStep: 1,
        });
      };

      const onAddEmailChange = (event) => {
        const {
          target: { value, name },
        } = event;
        setState({
          ...state,
          addEmailCredentials: {
            ...state.addEmailCredentials,
            [name]: value
          },
        });
      };

      const handleEmailSyncNumberChange = (event) => {
        const addEmailMailboxes = state.addEmailMailboxes;
        for (let i = 0; i < addEmailMailboxes.length; i++) {
            if (addEmailMailboxes[i].label == event.target.name) {
                addEmailMailboxes[i].syncNumber = parseInt(event.target.value);
            }
        }

        setState({
          ...state,
          addEmailMailboxes,
        });
      };

      const handleEmailSyncChange = (event) => {
        const addEmailMailboxes = state.addEmailMailboxes;
        for (let i = 0; i < addEmailMailboxes.length; i++) {
            if (addEmailMailboxes[i].label == event.target.name) {
                addEmailMailboxes[i].all = !addEmailMailboxes[i].all;
            }
        }

        setState({
          ...state,
          addEmailMailboxes,
        });
      };

      const syncAllSwitchStyle = SyncAllSwitchStyle();

      const getAddEmailContent = () => {
        let result = null;
        if (state.addEmailStep == 0) {
            result = [<DialogContentText>
                Enter username and password for your email provider using below inputs and configure mailboxes.
            </DialogContentText>,
            <Box component="form"
                sx={{
                    '& .MuiTextField-root': { m: 1, width: '25ch' },
                }}
                noValidate
                autoComplete="off">
                    <div>
                        <TextField
                            error={!Convenience.hasValue(state.addEmailCredentials.label)}
                            name="label"
                            onChange={onAddEmailChange}
                            label="Label"/>
                    </div>
                    <div>
                        <TextField
                            error={!Convenience.hasValue(state.addEmailCredentials.email)}
                            name="email"
                            onChange={onAddEmailChange}
                            label="Email"/>
                        <TextField
                            error={!Convenience.hasValue(state.addEmailCredentials.password)}
                            name="password"
                            onChange={onAddEmailChange}
                            label="Password"
                            type="password"
                            helperText="Passphrase to access your email"
                        />
                    </div>
            </Box>]
        } else if (state.addEmailStep == 1) {
            result =  [<div style={{textAlign: 'center'}}><CircularProgress/></div>];
        } else if (state.addEmailStep == 2) {
            const mailboxes = [];
            for (let i = 0; i < state.addEmailMailboxes.length; i++) {
                const label = state.addEmailMailboxes[i].label;
                const all = state.addEmailMailboxes[i].all;
                const syncNumber = state.addEmailMailboxes[i].syncNumber;
                mailboxes.push(
                    <Stack direction="row" spacing={1} alignItems="center">
                        <Typography style={{width: '200px'}}>{label}</Typography>
                        <Typography>All</Typography>
                        <Switch color="primary" classes={{
                                track: syncAllSwitchStyle.switch_track,
                                switchBase: syncAllSwitchStyle.switch_base,
                                colorPrimary: syncAllSwitchStyle.switch_primary,
                            }} checked={!all} name={label} onChange={handleEmailSyncChange}/>
                        <Typography>Last</Typography>
                        <TextField
                            name={label}
                            value={syncNumber}
                            type="number"
                            variant="standard"
                            size='small'
                            onChange={handleEmailSyncNumberChange}/>
                    </Stack>);
            }
            result = [<DialogContentText>
                Select mailboxes of your email account to sync in your Amper e-mail section, choose all or last N emails to download
            </DialogContentText>,
            <Box component="form"
                sx={{
                    '& .MuiTextField-root': { m: 1, width: '25ch' },
                }}
                noValidate
                autoComplete="off">
                {mailboxes}
            </Box>]
        }
        return <DialogContent>
                {result}
            </DialogContent>;
      };

      const getAddEmailDialog = () => {//password: qxfqilaoxdxvnsku
        return <Dialog onClose={handleAddEmailClose} open={state.addEmail} maxWidth="sm" maxHeight="sm" fullWidth fullHeight>
            <DialogTitle>Add email</DialogTitle>
            {getAddEmailContent()}
            <DialogActions>
                <Button onClick={handleAddEmailClose}>Cancel</Button>
                {state.addEmailStep == 2 ? <Button onClick={handleAddEmail}>Add</Button> : <Button onClick={handleAddEmailNext} disabled={!Convenience.hasValue(state.addEmailCredentials.email) && !Convenience.hasValue(state.addEmailCredentials.password)}>Next</Button>}
            </DialogActions>
        </Dialog>;
      };

    return <Box width="100%" sx={{position: 'relative'}}>
        <Box width="100%" height="auto">
            <Grid container spacing={2} width="100%" height="100%" sx={{mt: 3}}>
                <Grid size={6} height="300px">
                    <Typography variant="subtitle1" gutterBottom sx={{ml: 1}}>
                        Use the below list to enter the email addresses and credentials for configuring your mailbox.
                    </Typography>
                    <DataGridPremium
                        width="100%" height="100%"
                        slots={{
                            toolbar: getToolbar,
                        }}
                        onCellEditCommit={(cellData) => {
                            const { id, field, value } = cellData;
                            const emails = getEmailSettings();
                            for (let i = 0; i < emails.length; i++) {
                                if (emails[i].id == id) {
                                    emails[i][field] = value;
                                }
                            }
                            setState({
                                ...state,
                                data: {
                                    ...state.data,
                                    settings: {
                                        ...state.data.settings,
                                        email: emails,
                                    }
                                }
                            });
                        }}
                        onRowSelectionModelChange={(ids) => {
                            const selectedIDs = new Set(ids);
                            const rowData = getEmailSettings().filter((row) =>
                                selectedIDs.has(row.id)
                            )
        
                            setState({
                                ...state,
                                selectedRowData: rowData,
                            });
                        }}
                        rows={getEmailSettings()}
                        columns={getEmailSettingsColumns()}
                        initialState={{
                            pagination: {
                                paginationModel: {
                                    pageSize: 5,
                                },
                            },
                        }}
                        pageSizeOptions={[5]}
                        checkboxSelection
                        disableRowSelectionOnClick
                    />
                </Grid>
                <Grid size={6}>

                </Grid>
                <Grid size={4}>

                </Grid>
                <Grid size={8}>

                </Grid>
            </Grid>
            {getAddEmailDialog()}
    </Box>
        <Fab sx={{ position: 'absolute', bottom: 0, right: 0, }} onClick={save} color="primary">
            Save
        </Fab>
    </Box>
}