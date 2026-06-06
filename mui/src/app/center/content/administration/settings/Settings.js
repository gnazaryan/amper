import React, { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import IconButton from '@mui/material/IconButton';
import Button from '@mui/material/Button';
import FolderIcon from '@mui/icons-material/Folder';
import InputAdornment from '@mui/material/InputAdornment';
import KeyIcon from '@mui/icons-material/Key';
import Fab from '@mui/material/Fab';
import SaveIcon from '@mui/icons-material/Save';
import { post } from '../../../../data/Submit';
import HostManager from '../../../../../HostManager';
import { AppContext } from '../../../../../App';
import DataStore from '../../../../data/DataStore';
import LinearProgress from '@mui/material/LinearProgress';
import { sessionManager } from '../../../../../SessionManager';
import Tab from '@mui/material/Tab';
import TabContext from '@mui/lab/TabContext';
import TabList from '@mui/lab/TabList';
import TabPanel from '@mui/lab/TabPanel';
import { DataGridPremium, GridToolbarContainer } from '../../../../components/x-data-grid-premium';
import AddIcon from '@mui/icons-material/Add';
import RemoveIcon from '@mui/icons-material/Remove';
import { clone } from '../../../../amper/Instruments';
import { uniqueId } from '../../../../amper/Instruments';

function ImapToolbar(props) {
    const { addNewRow, removeRow } = props;
  
    const addClick = () => {
        addNewRow({id: uniqueId(),  name: '', domain: '', serverName: '', port: 25 });
    };

    const renoveClick = () => {
        removeRow();
    };
  
    return (
      <GridToolbarContainer>
        <Button color="primary" startIcon={<AddIcon />} onClick={addClick}>
          Add
        </Button>
        <Button color="primary" startIcon={<RemoveIcon />} onClick={renoveClick}>
          Remove
        </Button>
      </GridToolbarContainer>
    );
  }

  function SmtpToolbar(props) {
    const { addNewSmtpRow, removeSmtpRow } = props;
  
    const addClick = () => {
        addNewSmtpRow({id: uniqueId(),  name: '', domain: '', serverName: '', port: 25 });
    };

    const renoveClick = () => {
        removeSmtpRow();
    };
  
    return (
      <GridToolbarContainer>
        <Button color="primary" startIcon={<AddIcon />} onClick={addClick}>
          Add
        </Button>
        <Button color="primary" startIcon={<RemoveIcon />} onClick={renoveClick}>
          Remove
        </Button>
      </GridToolbarContainer>
    );
  }
  
export default function Settings() {
    const app = React.useContext(AppContext);
    const [state, setState] = useState(() => {
        return {
            loading: true,
            settings: {
                rootDirectory: '',
                adobeLicenseKey: '',
                imap: {
                    domains: []
                },
                smtp: {
                    domains: []
                }
            },
            tab: '1',
        };
    });

    useEffect(() => {
        if (state.loading) {
            getDiscoverDataStore().load((result) => {
                setState({
                    ...state,
                    settings: result.data,
                    loading: false,
                });
            });
        }
    }, [state.loading])
    const onChange = (event) => {
        const {name, value} = event.target;
        setState({
            ...state,
            settings: {
                ...state.settings,
                [name]: value,
            }
        });
    };

    const save = () => {
        post(`${HostManager.amperHost()}settings/save`, {
            settings: state.settings,
        }, (result) => {
            sessionManager.setSetting('adobeLicenseKey', state.settings.adobeLicenseKey)
            app.toast('info', `The settings were successfully updated.`);
        }, (result) => {
            app.toast('info', `The settings were not updated successfully with error: ${result.error}`);
        });
    };

    const getDiscoverDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}settings/fetch`,
            requestMethod: "POST",
        });
    };

    const getProgress = () => {
        return <LinearProgress sx={{mb: 2, mr: 4, visibility: state.loading ? 'visible' : 'hidden'}}/>
    };

    const openAdobePage = () => {
        window.open('https://developer.adobe.com/console');
    }

    const handleChange = (event, newValue) => {
        setState({
            ...state,
            tab: newValue,
        })
    };
    
    const columns = [
        { field: 'name', headerName: 'Name', width: 180, editable: true },
        { field: 'domain', headerName: 'Domain', width: 180, editable: true },
        { field: 'serverName', headerName: 'Server name', width: 180, editable: true },
        {
          field: 'port',
          headerName: 'Port',
          type: 'number',
          editable: true,
          align: 'left',
          headerAlign: 'left',
        },
        {
          field: 'auth',
          headerName: 'Authentication method',
          width: 220,
          editable: true,
          type: 'singleSelect',
          valueOptions: ['Plain auth', 'Login auth'],
        },
      ];

      const addNewRow = (row) => {
        let domains = clone(state.settings.imap.domains);
        domains.push(row);
        setState({
            ...state,
            settings: {
                ...state.settings,
                imap: {
                    ...state.settings.imap,
                    domains: domains,
                }
            }
            });
      };

      const addNewSmtpRow = (row) => {
        let domains = clone(state.settings.smtp.domains);
        domains.push(row);
        setState({
            ...state,
            settings: {
                ...state.settings,
                smtp: {
                    ...state.settings.smtp,
                    domains: domains,
                }
            }
            });
      };

      const removeRow = () => {
        const resultDomains = [];
        const selection = new Set(state.selectedIds);
        const domains = state.settings.imap.domains;
        for (let i = 0; i < domains.length; i++) {
            const domain = domains[i];
            if (!selection.has(domain.id)) {
                resultDomains.push(domain);
            }
        }
        setState({
            ...state,
            settings: {
                ...state.settings,
                imap: {
                    ...state.settings.imap,
                    domains: resultDomains,
                }
            }
            });
      };

      const removeSmtpRow = ()=> {
        const resultDomains = [];
        const selection = new Set(state.selectedSmtpIds);
        const domains = state.settings.smtp.domains;
        for (let i = 0; i < domains.length; i++) {
            const domain = domains[i];
            if (!selection.has(domain.id)) {
                resultDomains.push(domain);
            }
        }
        setState({
            ...state,
            settings: {
                ...state.settings,
                smtp: {
                    ...state.settings.smtp,
                    domains: resultDomains,
                }
            }
            });
      };

  return (<Box sx={{ width: '100%', typography: 'body1' }}>
    <TabContext value={state.tab}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <TabList onChange={handleChange} aria-label="lab API tabs example">
            <Tab label="General" value="1" />
            <Tab label="Imap" value="2" />
            <Tab label="Smtp" value="3" />
        </TabList>
        </Box>
        <TabPanel value="1">
            <Box component="form"
                sx={{
                    mt: 0,
                    width: '100%'
                }}
                noValidate
                autoComplete="off">
                    {getProgress()}
                <div>
                    <TextField
                        required
                        name="rootDirectory"
                        label="Root directory"
                        size="large"
                        sx={{mb: 3, width: 400}}
                        onChange={onChange}
                        value={state.settings.rootDirectory}
                        InputProps={{
                            endAdornment: <InputAdornment position="end">
                                <IconButton type="button" edge="end" variant="contained" component="label">
                                    <FolderIcon color="primary"/>
                                </IconButton>
                            </InputAdornment>,}}>
                        </TextField>
                </div>
                <div>
                    <TextField
                            required
                            name="adobeLicenseKey"
                            label="Adobe DC key"
                            size="large"
                            sx={{mb: 3, width: 400}}
                            onChange={onChange}
                            value={state.settings.adobeLicenseKey}
                            InputProps={{
                                endAdornment: <InputAdornment position="end">
                                    <IconButton type="button" edge="end" variant="contained" component="label" onClick={openAdobePage}>
                                        <KeyIcon color="primary"/>
                                    </IconButton>
                                </InputAdornment>,}}>
                    </TextField>
                </div>
            </Box>
        </TabPanel>
        <TabPanel value="2">
            <div style={{ height: 300, width: '100%' }}>
                <DataGridPremium 
                    onCellEditCommit={(cellData) => {
                        const { id, field, value } = cellData;
                        const domains = state.settings.imap.domains;
                        
                        for (let i = 0; i < domains.length; i++) {
                            if (domains[i].id == id) {
                                domains[i][field] = value;
                            }
                        }
                        setState({
                            ...state,
                            settings: {
                                ...state.settings,
                                imap: {
                                    ...state.settings.imap,
                                    domains: domains,
                                }
                            }
                        });
                    }}
                    onRowSelectionModelChange={(ids) => {
                        setState({
                            ...state,
                            selectedIds: ids,
                        });
                    }}
                    rows={state.settings.imap.domains} 
                    columns={columns}
                    slots={{
                        toolbar: ImapToolbar,
                    }}
                    slotProps={{
                        toolbar: {
                            addNewRow,
                            removeRow
                        }
                    }}/>
            </div>
        </TabPanel>
        <TabPanel value="3">
            <div style={{ height: 300, width: '100%' }}>
                <DataGridPremium 
                    onCellEditCommit={(cellData) => {
                        const { id, field, value } = cellData;
                        const domains = state.settings.smtp.domains;
                        
                        for (let i = 0; i < domains.length; i++) {
                            if (domains[i].id == id) {
                                domains[i][field] = value;
                            }
                        }
                        setState({
                            ...state,
                            settings: {
                                ...state.settings,
                                smtp: {
                                    ...state.settings.smtp,
                                    domains: domains,
                                }
                            }
                        });
                    }}
                    onRowSelectionModelChange={(ids) => {
                        setState({
                            ...state,
                            selectedSmtpIds: ids,
                        });
                    }}
                    rows={state.settings.smtp.domains} 
                    columns={columns}
                    slots={{
                        toolbar: SmtpToolbar,
                    }}
                    slotProps={{
                        toolbar: {
                            addNewSmtpRow,
                            removeSmtpRow
                        }
                    }}/>
            </div>
        </TabPanel>
    </TabContext>
        <Fab color="primary" aria-label="add" sx={{
            position: 'absolute',
            top: 100,
            right: 26,}} onClick={save}>
            <SaveIcon/>
        </Fab>
    </Box>
    );
}
