import React, { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import { DataGridPremium } from '../../../../components/x-data-grid-premium';
import LinearProgress from '@mui/material/LinearProgress';
import HostManager from "../../../../../HostManager";
import DataStore from "../../../../data/DataStore";
import Convenience from "../../../../help/Convenience";
import { useLocation } from 'react-router-dom'
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
import {post} from '../../../../data/Submit'
import Autocomplete from '@mui/material/Autocomplete';
import Typography from '@mui/material/Typography';
import {makeApieName} from '../../../../amper/Instruments'

export default function ObjectTypes() {
    const {search} = useLocation();
    const objectId = Convenience.getUrlParameterValueFromQuery(search, 'objectId')

    const [selectedRowData, setSelectedRowData] = React.useState([]);
    const [selectedOtfRowData, setSelectedOtfRowData] = React.useState([]);
    const initialState = {
        loading: true,
        fieldsLoading: false,
        data: [],
        fields: [],
        createObjectTypeDialogOpen: false,
        createObjectTypeFieldDialogOpen: false,
        craeteObjectTypeForm: {
          label: '',
          apiName: '',
          extendsTo: undefined,
          objectId: parseInt(objectId),
        },
        addFieldForm: {
          field: undefined,
        },
        craeteObjectTypeFormError: '',
        addFieldToObjectTypeFormError: '',
        objectTypePaging: {
            page: 0,
            pageSize: 50,
        },
        objectTypeFieldsPaging: {
            page: 0,
            pageSize: 50,
        }
    };
    const [state, setState] = useState(initialState);
    
    const objectTypeFields = selectedRowData.length > 0 ? selectedRowData[0].objectTypeFields : [];
    const objectTypeFieldsIds = new Set(objectTypeFields.map(function(item) {
      return item.fieldId;
    }));
    const objectTypeFieldsFiltered = state.fields.filter((objectTypeField) => !objectTypeFieldsIds.has(objectTypeField.id));
    useEffect(() => {
        if (state.loading) {
          getDataStore().load((result)=> {          
              setState({
                ...state,
                loading: false,
                data: result.data,
              })
          })
        }
        if (state.fieldsLoading) {
          getFieldsDataStore().load((result)=> {          
            setState({
              ...state,
              fieldsLoading: false,
              fields: result.data,
            })
        })
        }
      });

    const objectTypeColumns = [
        { 
            field: 'id', 
            headerName: 'ID',
            hide: true,
        },
        {
          field: 'label',
          headerName: 'Label',
          flex: 1,
        },
        {
          field: 'apiName',
          headerName: 'Api name',
          flex: 1,
        },
        {
          field: 'extendsToLabel',
          headerName: 'Extends',
          flex: 1,
        },
      ];

      const objectTypeFieldsColumn = [{
            field: 'label',
            headerName: 'Label',
            flex: 1,
        },{
          field: 'objectTypeLabel',
          headerName: 'Object type',
          flex: 1,
        },{
            field: 'apiName',
            headerName: 'Api Name',
            flex: 1,
        },{
            field: 'type',
            headerName: 'Field Type',
            flex: 1,
        }];

    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}objectTypes/getObjectTypes`,
            requestMethod: "POST",
            parameters: {
                objectId: parseInt(objectId),
            },
        });
    }

    const getFieldsDataStore = () => {
      return new DataStore({
          url: `${HostManager.amperHost()}entities/getFields`,
          requestMethod: "POST",
          parameters: {
              objectId: parseInt(objectId),
          },
      });
    }

    const create = () => {
      if (!state.loading) {
        setState({
          ...state,
          createObjectTypeDialogOpen: true,
        })
      }
    };

    const handleCreaeteObjectTypeClose = () => {
      setState({
        ...state,
        createObjectTypeDialogOpen: false,
      })
    };

    const handleCreaeteObjectTypeSubmit = () => {
      if (state.craeteObjectTypeForm.extendsTo && state.craeteObjectTypeForm.apiName && state.craeteObjectTypeForm.label) {
        post(`${HostManager.amperHost()}objectTypes/createObjectType`, state.craeteObjectTypeForm, (result) => {
          setState({
            ...state,
            createObjectTypeDialogOpen: false,
            loading: true,
          })
        }, (result) => {
          setState({
            ...state,
            craeteObjectTypeFormError: result.error,
          });
        })
      }
    };

    const remove = () => {
      if (selectedRowData.length > 0) {
        post(`${HostManager.amperHost()}objectTypes/deleteObjectType`, {objectTypeId: selectedRowData[0].id}, () => {
          setState(initialState)
        })
      }
    };
    
    const add = () => {
      if (selectedRowData.length > 0) {
        setState({
          ...state,
          createObjectTypeFieldDialogOpen: true,
          fieldsLoading: true,
        });
      }
    };

    const removeOtf =() => {
      if (selectedOtfRowData.length > 0 && selectedRowData.length > 0) {
        post(`${HostManager.amperHost()}fields/deleteObjectTypeField`, {
          fieldId: selectedOtfRowData[0].fieldId,
          objectTypeId: selectedRowData[0].id
        }, () => {
          setState(initialState)
        }, () => {
          setState(initialState)
        })
      }
    };

    const handleObjectTypeLabelChange = (event) => {
      const {
        target: { value, name },
      } = event;
      setState({
        ...state,
        craeteObjectTypeForm : {
          ...state.craeteObjectTypeForm,
          label: value,
          apiName: makeApieName(value),
        }
      });
    };

    const handleAddObjectTypeFieldClose = () => {
      setState({
        ...state,
        createObjectTypeFieldDialogOpen: false,
        addFieldForm: {
          field: undefined,
        },
        addFieldToObjectTypeFormError: undefined,
      });
    };

    const handleAddFieldToObjectTypeChange = (event, newValue) => {
      setState({
        ...state,
        addFieldForm: {
          field: newValue
        }
      });
    };

    const handleObjectTypeExtendsToChange = (event, newValue) => {
      setState({
        ...state,
        craeteObjectTypeForm: {
          ...state.craeteObjectTypeForm,
          extendsTo: newValue.id
        },
        extendsToRecord: newValue,
      });
    };

    const handleAddFieldToObjectTypeSubmit = () => {
      if (state.addFieldForm.field && selectedRowData.length > 0) {
        post(`${HostManager.amperHost()}fields/addObjectTypeField`, {
          objectTypeId: selectedRowData[0].id,
          fieldId: state.addFieldForm.field.id,
        }, (result) => {
          setState(initialState);
        }, (result) => {
          setState({
            ...state,
            addFieldToObjectTypeFormError: result.error,
          });
        });
      }
    };

    const setObjectTypeFieldsPaginationModel = (pagingModel) => {
        setState({
            ...state,
            objectTypeFieldsPaging: pagingModel
        });
    };

    const setObjectTypePaginationModel = (pagingModel) => {
      setState({
          ...state,
          objectTypePaging: pagingModel
      });
  };

    const getError = (error) => {
        if (error) {
          return <Typography sx={{m: 1}} color="error" variant="caption" display="block">
            {error}
          </Typography>;
        }
    };

    const getAddObjectTypeFieldDialog = () => {
      return (
        <Dialog open={state.createObjectTypeFieldDialogOpen} onClose={handleAddObjectTypeFieldClose}>
          <DialogTitle>Add Field</DialogTitle>
          <DialogContent>
            <DialogContentText>
            To add a field to the object type, select the field in the drop down box
            </DialogContentText>
            <Autocomplete
              sx={{mt: 3}}
              onChange={handleAddFieldToObjectTypeChange}
              loading = {state.fieldsLoading}
              options={objectTypeFieldsFiltered}
              renderInput={(params) => <TextField
                variant="standard"
                {...params}
                error={!state.addFieldForm.field}
                name="addField"
                required
                label="Select field" />}
            />
            {getError(state.addFieldToObjectTypeFormError)}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleAddObjectTypeFieldClose}>Cancel</Button>
            <Button onClick={handleAddFieldToObjectTypeSubmit} disabled={!state.addFieldForm.field}>Ok</Button>
          </DialogActions>
        </Dialog>
      );
    };

    const getCreateObjectTypeDialog = () => {
      return (
        <Dialog open={state.createObjectTypeDialogOpen} onClose={handleCreaeteObjectTypeClose}>
          <DialogTitle>Create Object Type</DialogTitle>
          <DialogContent>
            <DialogContentText>
              To create an object type, specify the lable, api name and then select what other object type is it going to inherit from.
            </DialogContentText>
            <TextField
              sx={{mt: 3}}
              autoFocus
              name="label"
              onChange={handleObjectTypeLabelChange}
              value={state.craeteObjectTypeForm.label}
              error={!state.craeteObjectTypeForm.label}
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
              name="apiName"
              value={state.craeteObjectTypeForm.apiName}
              error={!state.craeteObjectTypeForm.apiName}
              required
              disabled
              label="Api name"
              fullWidth
              variant="filled"
              color="primary"
              size="large"
            />
            <Autocomplete
              sx={{mt: 3}}
              onChange={handleObjectTypeExtendsToChange}
              options={state.data}
              renderInput={(params) => <TextField
                variant="standard"
                {...params}
                error={!state.craeteObjectTypeForm.extendsTo}
                name="extendsTo"
                required
                label="Extends to" />}
            />
            {getError(state.craeteObjectTypeFormError)}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCreaeteObjectTypeClose}>Cancel</Button>
            <Button onClick={handleCreaeteObjectTypeSubmit} disabled={!state.craeteObjectTypeForm.extendsTo || !state.craeteObjectTypeForm.apiName || !state.craeteObjectTypeForm.label}>Ok</Button>
          </DialogActions>
        </Dialog>
      );

    };

    return (
        <Box sx={{ height: 'calc(100% - 40px)', width: '100%', flexGrow: 1, display: 'flex',
        flexDirection: 'row',}}>
            {getCreateObjectTypeDialog()}
            {getAddObjectTypeFieldDialog()}
            <Box sx={{mr: 3, ml: -3, height: '100%', width: '50%'}}>
                <Stack direction="row" spacing={1} sx={{ mb: 1 }}>
                    <Button size="small" onClick={create} startIcon={<AddCircleOutlineIcon/>}>
                        Create
                    </Button>
                    <Button size="small" disabled={selectedRowData.length == 0} onClick={remove} startIcon={<RemoveCircleOutlineIcon/>}>
                        Remove
                    </Button>
                </Stack>
                <DataGridPremium
                    rows={state.data}
                    slots={{
                        loadingOverlay: LinearProgress,
                    }}
                    loading={state.loading}
                    columns={objectTypeColumns}
                    pageSize={50}
                    pageSizeOptions={[50, 100, 500]}
                    paginationModel={state.objectTypePaging}
                    pagination
                    onPaginationModelChange={setObjectTypePaginationModel}
                    sx={{color: 'secondary.gridText'}}
                    onRowSelectionModelChange={(ids) => {
                        const selectedIDs = new Set(ids);
                        const rowData = state.data.filter((row) =>
                          selectedIDs.has(row.id)
                        )
                        setSelectedRowData(rowData);
                      }}
                      localeText={{
                        footerRowSelected: (count) => `${count} row selected`,
                        MuiTablePagination: {
                          labelDisplayedRows: ({ from, to, count }) =>
                            `${from} - ${to} of more than ${count}`,
                        },
                      }}
                />
            </Box>
            <Box sx={{mr: 3, height: '100%', width: '50%' }}>
                <Stack direction="row" spacing={1} sx={{ mb: 1 }}>
                    <Button size="small" disabled={!(selectedRowData.length > 0)} onClick={add} startIcon={<AddCircleOutlineIcon/>}>
                        ADD
                    </Button>
                    <Button size="small" disabled={!(selectedOtfRowData.length > 0 && selectedRowData.length > 0 && selectedOtfRowData[0].objectTypeId === selectedRowData[0].id)} onClick={removeOtf} startIcon={<RemoveCircleOutlineIcon/>}>
                        Remove
                    </Button>
                </Stack>
                <DataGridPremium
                    slots={{
                        loadingOverlay: LinearProgress,
                    }}
                    rows={objectTypeFields}
                    columns={objectTypeFieldsColumn}
                    pageSize={50}
                    sx={{color: 'secondary.gridText' }}
                    pageSizeOptions={[50, 100, 500]}
                    paginationModel={state.objectTypeFieldsPaging}
                    pagination
                    onPaginationModelChange={setObjectTypeFieldsPaginationModel}
                    onRowSelectionModelChange={(ids) => {
                      const selectedIDs = new Set(ids);
                      const rowData = objectTypeFields.filter((row) =>
                        selectedIDs.has(row.id)
                      );
                      setSelectedOtfRowData(rowData);
                    }}
                />
            </Box>
        </Box>
        
    );
}