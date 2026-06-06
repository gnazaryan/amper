import React, { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import { DataGridPremium } from '../../../../components/x-data-grid-premium';
import HostManager from "../../../../../HostManager";
import DataStore from "../../../../data/DataStore";
import { useLocation } from 'react-router-dom'
import Convenience from "../../../../help/Convenience";
import LinearProgress from '@mui/material/LinearProgress';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import Typography from '@mui/material/Typography';
import {makeApieName} from '../../../../amper/Instruments'
import MenuItem from '@mui/material/MenuItem';
import Grid from '@mui/material/Grid2';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';
import {post} from '../../../../data/Submit'

export default function Fields({toast}) {

    const {search} = useLocation();
    const objectId = Convenience.getUrlParameterValueFromQuery(search, 'objectId')
    const [state, setState] = useState({
        loading: true,
        objectsLoading: false,
        data: [],
        objectsData: [],
        createFieldDialogOpen: false,
        createFieldError: '',
        createFieldForm: {
            label: '',
            apiName: '',
            dataType: 'NUMBER',
            objectReference: undefined,
            required: true,
            status: true,
            textLength: 256,
            entityId: parseInt(objectId),
        },
        fieldsPaging: {
            page: 0,
            pageSize: 50,
        }
    });
    const [selectedRowData, setSelectedRowData] = React.useState([]);

    useEffect(() => {
        if (state.loading) {
          getDataStore().load((result)=> {          
              setState({
                ...state,
                loading: false,
                data: result.data,
              })
          });
        }
        if (state.objectsLoading) {
            getObjectsDataStore().load((result) => {
                setState({
                    ...state,
                    objectsLoading: false,
                    objectsData: result.data,
                })
            });
        }
      });

      const create = () => {
        setState({
            ...state,
            createFieldDialogOpen: true,
            objectsLoading: true,
          })
      };
  
    const remove = () => {
        if (selectedRowData.length > 0) {
            const entityId = selectedRowData[0].objectId;
            const fieldIds = [];
            for (let i = 0; i < selectedRowData.length; i++) {
                fieldIds.push(selectedRowData[i].id);
            }
            post(`${HostManager.amperHost()}fields/deleteField`, {
                entityId: entityId,
                fieldIds: fieldIds,
            }, (result) => {
                setState({
                    ...state,
                    loading: true,
                })
                }, (result) => {toast('error', result.error)
                setState({
                    ...state,
                    loading: true,
                });
            });
        }
    };

    const fieldColumns = [{
            headerName: 'Label',
            field: 'label',
            flex: 1,
        },{
            headerName: 'Field key',
            field: 'apiName',
            flex: 1,
        },{
            headerName: 'Type',
            field: 'type',
            flex: 1,
        },{
            headerName: 'Required',
            field: 'required',
            flex: 1,
            render: function (rawValue) {
                return rawValue ? "True" : "False";
            }
        },{
            headerName: 'Status',
            field: 'status',
            flex: 1,
            render: function (rawValue) {
                return rawValue ? "Active" : "Inactive";
            }
        }
    ];

    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}entities/getFields`,
            requestMethod: "POST",
            parameters: {
                objectId: parseInt(objectId),
            },
        });
    };

    const handleCreateFieldDialogClose = () => {
        setState({
            ...state,
            createFieldDialogOpen: false,
        });
    };

    const handleObjectReferenceChange = (event, value) => {
        setState({
            ...state,
            createFieldForm: {
              ...state.createFieldForm,
              objectReference: value.id
            },
            objectReferenceRecord: value,
          });
    };

    const handleCreateFieldDialogSubmit = () => {
        if (state.createFieldForm.label && state.createFieldForm.apiName && state.createFieldForm.dataType) {
            post(`${HostManager.amperHost()}fields/createField`, state.createFieldForm, (result) => {
              setState({
                ...state,
                createFieldDialogOpen: false,
                loading: true,
              })
            }, (result) => {
              setState({
                ...state,
                createFieldError: result.error,
              });
            })
          }
    };

    const getObjectsDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}entities/getEntities`,
            requestMethod: "POST",
            parameters: {
                limit: 9007199254740991,
                start: 0,
            },
        });
      }

    const getError = (error) => {
        if (error) {
          return <Typography sx={{m: 1}} color="error" variant="caption" display="block">
            {error}
          </Typography>;
        }
    };

    const handleObjectTypeLabelChange = (event) => {
        const {
            target: { value, name },
          } = event;
          setState({
            ...state,
            createFieldForm : {
              ...state.createFieldForm,
              label: value,
              apiName: makeApieName(value),
            }
          });
    };

    const handleInputChange = (event) => {
        const {
            target: { value, name },
          } = event;
          const createFieldForm = state.createFieldForm;
          createFieldForm[name] = value;
          setState({
            ...state,
            createFieldForm,
          });
    };

    const handleSwitchInputChange = (event) => {
        const {
            target: { checked, name },
          } = event;
          setState({
            ...state,
            createFieldForm: {
                ...state.createFieldForm,
                [name]: checked,
            },
          });
    };

    const handleNumberInputChange = (event) => {
        const {
            target: { value, name },
          } = event;
          const createFieldForm = state.createFieldForm;
          createFieldForm[name] = parseInt(value);
          setState({
            ...state,
            createFieldForm,
          });
    };

    const setFieldsPaginationModel = (pagingModel) => {
        setState({
            ...state,
            fieldsPaging: pagingModel
        });
    };

    const getOptionalInputs = () => {
        if (state.createFieldForm.dataType === 'REFERENCE') {
            return (
                <Autocomplete
                    onChange={handleObjectReferenceChange}
                    loading = {state.objectsLoading}
                    options={state.objectsData}
                    value={state.objectReferenceRecord}
                    getOptionLabel={(option) => option.title}
                    renderInput={(params) => <TextField
                        variant="standard"
                        {...params}
                        error={!state.createFieldForm.objectReference}
                        required
                        label="Reference" />}
                  />
            );
        } else if (state.createFieldForm.dataType === 'TEXT' || state.createFieldForm.dataType === 'NUMBER') {
            return (
                <TextField
                    variant="standard"
                    onChange={handleNumberInputChange}
                    name="textLength"
                    label={state.createFieldForm.dataType === 'TEXT' ? "Text Length" : 'Max Value'}
                    type="number"
                    autoFocus
                    fullWidth
                    value={state.createFieldForm.textLength}
                    error={!state.createFieldForm.textLength} />
            );
        }
    };

    const getCreateFieldDialog = () => {
        return (
            <Dialog maxWidth="sm" fullWidth={true} open={state.createFieldDialogOpen} onClose={handleCreateFieldDialogClose}>
                <DialogTitle>Create field</DialogTitle>
                <DialogContent>
                  <DialogContentText>
                  To add a field to the object type, fill in the label, select the data type in the drop down box and configure the required and active properties
                  </DialogContentText>
                  <Grid container spacing={3} sx={{mt: 3}}>
                    <Grid item size={12}>
                        <TextField
                        autoFocus
                        name="label"
                        onChange={handleObjectTypeLabelChange}
                        value={state.createFieldForm.label}
                        error={!state.createFieldForm.label}
                        label="Label"
                        fullWidth
                        required
                        variant="filled"
                        color="primary"
                        size="large"
                        />
                    </Grid>
                    <Grid item size={12}>
                        <TextField
                        autoFocus
                        name="apiName"
                        value={state.createFieldForm.apiName}
                        error={!state.createFieldForm.apiName}
                        required
                        disabled
                        label="Api name"
                        fullWidth
                        variant="filled"
                        color="primary"
                        size="large"
                    />
                    </Grid>
                    <Grid item size={12}>
                        <TextField
                            variant="standard"
                            name="dataType"
                            fullWidth
                            value={state.createFieldForm.dataType}
                            label="Data type"
                            onChange={handleInputChange}
                            select>
                            <MenuItem value="NUMBER">Number</MenuItem>
                            <MenuItem value="TEXT">Text</MenuItem>
                            <MenuItem value="REFERENCE">Reference</MenuItem>
                            <MenuItem value="BOOLEAN">Boolean</MenuItem>
                            <MenuItem value="DATE">Date</MenuItem>
                            <MenuItem value="DATETIME">Date Time</MenuItem>
                        </TextField>
                    </Grid>
                    <Grid item size={12}>
                    {getOptionalInputs()}
                    </Grid>
                    <Grid item size={12}>
                        <FormControlLabel labelPlacement="end" control={<Switch name="required" checked={state.createFieldForm.required} onChange={handleSwitchInputChange}/>} label="Required"/>
                    </Grid>
                    <Grid item size={12}>
                        <FormControlLabel labelPlacement="end" control={<Switch name="status" checked={state.createFieldForm.status} onChange={handleSwitchInputChange}/>} label="Status"/>
                    </Grid>
                </Grid>
                  {getError(state.createFieldError)}
                </DialogContent>
                <DialogActions>
                  <Button onClick={handleCreateFieldDialogClose}>Cancel</Button>
                  <Button onClick={handleCreateFieldDialogSubmit} disabled={!state.createFieldForm.label || !state.createFieldForm.apiName}>Ok</Button>
                </DialogActions>
              </Dialog>
        );
    };

    return (
        <Box sx={{ height: 'calc(100% - 40px)', width: '100%', flexGrow: 1, ml: -3, display: 'flex',
        flexDirection: 'row',}}>
            {getCreateFieldDialog()}
            <Box sx={{ height: '100%', width: '100%'}}>
                <Stack direction="row" spacing={1} sx={{ mb: 1 }}>
                    <Button size="small" onClick={create} startIcon={<AddCircleOutlineIcon/>}>
                        Create
                    </Button>
                    <Button size="small" onClick={remove} startIcon={<RemoveCircleOutlineIcon/>} disabled={selectedRowData.length == 0}>
                        Remove
                    </Button>
                </Stack>
                <DataGridPremium
                    slots={{
                        loadingOverlay: LinearProgress,
                    }}
                    loading={state.loading}
                    onRowSelectionModelChange={(ids) => {
                        const selectedIDs = new Set(ids);
                        const rowData = state.data.filter((row) =>
                            selectedIDs.has(row.id)
                        )
                        setSelectedRowData(rowData);
                    }}
                    rows={state.data}
                    columns={fieldColumns}
                    pageSizeOptions={[50, 100, 500]}
                    paginationModel={state.fieldsPaging}
                    pagination
                    onPaginationModelChange={setFieldsPaginationModel}
                    sx={{color: 'secondary.gridText', }}
                />
            </Box>
        </Box>
    );
}