import React, { useState, useImperativeHandle, forwardRef, useEffect } from 'react';
import Grid from '@mui/material/Grid2';
import Box from '@mui/material/Box';
import Skeleton from '@mui/material/Skeleton';
import Stack from '@mui/material/Stack';
import AmperConstatns from '../../../util/AmperConstants';
import Convenience from '../../../help/Convenience';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import Link from '@mui/material/Link';
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import DataStore from "../../../data/DataStore";
import HostManager from "../../../../HostManager";
import { ReferenceField } from '../fields/ReferenceField';
import { parseBoolean } from '../../../util/BooleanUtil';
import AdjustIcon from '@mui/icons-material/Adjust';
import Tooltip from '@mui/material/Tooltip';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';
import dayjs from 'dayjs';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { DesktopDateTimePicker } from '@mui/x-date-pickers/DesktopDateTimePicker';
import { DesktopDatePicker } from '@mui/x-date-pickers/DesktopDatePicker';
import RemoveRedEyeIcon from '@mui/icons-material/RemoveRedEye';
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import ToggleButton from '@mui/material/ToggleButton';
import { post } from '../../../data/Submit';
import { AppContext } from '../../../../App';

function RecordDetails({metadata, record, onChange}, ref) {
    const app = React.useContext(AppContext);
    const initialState = () => {
        return {
            objectsLoading: true,
            objects: [],
            object: null,
            loading: false,
            noRecord: true,
            mode: 'view',
            reloadRecords: false,
            updatePayload: false,
        };
    }
    const [state, setState] = useState(initialState);

    useEffect(() => {
        if (state.objectsLoading) {
          getDataStore().load((result)=> {          
              setState({
                ...state,
                objectsLoading: false,
                objects: result.data || [],
              });
          });
        }
      }, [state.objectsLoading]);

      useEffect(() => {
        if (state.updatePayload) {
            post(`${HostManager.amperHost()}records/update`, {
                payload: JSON.stringify(state.payload)
              }, (result) => {
                app.toast('info', 'The record was successfully updated.');
                setState({
                    ...state,
                    mode: 'view',
                    record: state.payload,
                    loading: false,
                    updatePayload: false,
                });
              }, (result) => {
                app.toast('warning', 'The record was not successfully updated.');
                setState({
                    ...state,
                    mode: 'view',
                    loading: false,
                    updatePayload: false,
                });
              });
        }
      }, [state.updatePayload]);

    const getDataStore = () => {
    return new DataStore({
        url: `${HostManager.amperHost()}entities/getEntities`,
        requestMethod: "POST",
        parameters: {
            start: 0,
            limit: AmperConstatns.INTEGER.MAX_VALUE
        }
    });
    };
    useImperativeHandle(ref, () => ({
        reset() {
            setState(initialState());
        },
        setRecord(metadata, record) {
            setState({
                ...state,
                loading: false,
                metadata,
                objectType: getObjectType(metadata, record),
                record,
                object: {
                    id: metadata.Object.id,
                    apiName: metadata.Object.apiName,
                    label: metadata.Object.title,
                },
                payload: JSON.parse(JSON.stringify(record)),
                noRecord: false,
                reloadRecords: true
            });
        }
    }));

    const getObjectType = (metadata, record) => {
        let objectType = null;
        if (metadata && record) {
            const objectTypeApiName = record[AmperConstatns.SYSTEM_FIELDS.OBJECT_TYPE]
            if (Convenience.hasValue(objectTypeApiName)) {
                for (let i = 0; i < metadata.ObjectTypes.length; i++) {
                    if (metadata.ObjectTypes[i].apiName === objectTypeApiName) {
                        objectType = metadata.ObjectTypes[i];
                    }
                }
            }
        }
        return objectType;
    };
    
    const valid = () => {
        const fieldsMap = {};
        for (let i = 0; i < state.metadata.Fields.length; i++) {
            let field = state.metadata.Fields[i];
            fieldsMap[field.id] = field;
        }
        for (let l = 0; l < state.objectType.objectTypeFields.length; l++) {
            const field = fieldsMap[state.objectType.objectTypeFields[l].fieldId];
            if (parseBoolean(field.required) && (state.payload[field.apiName] == null || state.payload[field.apiName] === '')) {
                return false;
            }
        }
        return true;
    };

    const getViewFields = () => {
        const result = [];
        if (!state.loading && state.metadata && state.record) {
            const objectTypeApiName = state.record[AmperConstatns.SYSTEM_FIELDS.OBJECT_TYPE]
            if (Convenience.hasValue(objectTypeApiName)) {
                if (state.objectType != null) {
                    let item = null;
                    const fieldsMap = {};
                    for (let i = 0; i < state.metadata.Fields.length; i++) {
                        let field = state.metadata.Fields[i];
                        fieldsMap[field.id] = field;
                    }
                    for (let l = 0; l < state.objectType.objectTypeFields.length; l++) {
                        const field = fieldsMap[state.objectType.objectTypeFields[l].fieldId]
                        if (field != null) {
                            switch(field.type) {
                                case 'REFERENCE':
                                    result.push(<Grid key={l} size={12}>
                                        <Typography variant="h6">
                                            <Stack direction={'row'}>
                                                <Typography sx={{fontWeight: 'bold'}} color="primary.label" variant="h6"> {field.label}:</Typography>
                                                <Typography sx={{ml: 1}} variant="h6" ><Link href="" underline="hover">{state.record[field.apiName + '_name_sys']}</Link></Typography>
                                            </Stack>
                                        </Typography>
                                    </Grid>);
                                    break;
                                case 'BOOLEAN':
                                    result.push(<Grid key={l} size={12}>
                                        <Stack direction={'row'}>
                                            <Typography sx={{fontWeight: 'bold'}} color="primary.label" variant="h6"> {field.label}:</Typography>
                                            <Typography sx={{ml: 1}} variant="h6" >{parseBoolean(state.record[field.apiName]) ? 'Active': 'Inactive'}</Typography>
                                        </Stack>
                                    </Grid>);
                                    break;
                                default:
                                    result.push(<Grid key={l} size={12}>
                                        <Stack direction={'row'}>
                                            <Typography sx={{fontWeight: 'bold'}} color="primary.label" variant="h6"> {field.label}:</Typography>
                                            <Typography sx={{ml: 1}} variant="h6" >{state.record[field.apiName]}</Typography>
                                        </Stack>
                                    </Grid>);
                                    break;
                            }
                        }
                    }
                }
            }
        }
        return result;
    }

    const onInputChange = (event) => {
        const newPayload = {
            ...state.payload,
            [event.target.name]: event.target.value,
        };
        setState({
            ...state,
            payload: newPayload
        });
        if (onChange) {
            onChange(newPayload)
        }
    };
    
    const onChangeSwitch =  (event) => {
        const newPayload = {
            ...state.payload,
            [event.target.name]: event.target.checked,
        }
        setState({
            ...state,
            payload: newPayload,
        });
        if (onChange) {
            onChange(newPayload)
        }
    };

    const onReferenceChange = (name, value) => {
        var newPayload = {
            ...state.payload,
        };
        if (value != null) {
            newPayload[name] = value.identifier_sys;
            newPayload[name + '_name_sys'] = value.name_sys;
        } else {
            delete newPayload[name];
            delete newPayload[name + '_name_sys'];
        };
        setState({
            ...state,
            payload: newPayload,
        });
        if (onChange) {
            onChange(newPayload)
        }
    };

    const NON_MODIFIABLE_FIELDS = ['id', 'identifier_sys', 'objectType_sys']
    const getEditFields = () => {
        const result = [];
        if (!state.loading && state.metadata && state.record) {
            const objectTypeApiName = state.record[AmperConstatns.SYSTEM_FIELDS.OBJECT_TYPE]
            if (Convenience.hasValue(objectTypeApiName)) {
                let objectType = null;
                for (let i = 0; i < state.metadata.ObjectTypes.length; i++) {
                    if (state.metadata.ObjectTypes[i].apiName === objectTypeApiName) {
                        objectType = state.metadata.ObjectTypes[i];
                    }
                }
                if (objectType != null) {
                    let item = null;
                    const fieldsMap = {};
                    for (let i = 0; i < state.metadata.Fields.length; i++) {
                        let field = state.metadata.Fields[i];
                        fieldsMap[field.id] = field;
                    }
                    for (let l = 0; l < objectType.objectTypeFields.length; l++) {
                        const field = fieldsMap[objectType.objectTypeFields[l].fieldId];
                        if (field != null) {
                            const switchValue = NON_MODIFIABLE_FIELDS.includes(field.apiName) ? field.apiName : field.type;
                            switch(switchValue) {
                                case 'id':
                                case 'identifier_sys':
                                case 'objectType_sys':
                                    result.push(<Grid key={l} size={12}>
                                        <Stack direction={'row'}>
                                            <Typography sx={{fontWeight: 'bold'}} color="primary.label" variant="h6"> {field.label}:</Typography>
                                            <Typography sx={{ml: 1}} variant="h6" >{state.record[field.apiName]}</Typography>
                                        </Stack>
                                    </Grid>);
                                    break;
                                case 'TEXT':
                                    result.push(<Grid key={l} size={12}>
                                        <TextField key={field.apiName}
                                            name={field.apiName} 
                                            inputProps={{ maxLength: field.textLength }}
                                            multiline={field.textLength > 256}
                                            maxRows={4}
                                            required={field.required == 1 ? true : false} 
                                            label={field.label} 
                                            onChange={onInputChange}
                                            value={state.payload[field.apiName] || ''}
                                            error={field.required == 1 ? state.payload[field.apiName] == null : false}
                                            variant="standard" 
                                            fullWidth={true}/>
                                        </Grid>);
                                    break;
                                case 'BOOLEAN':
                                    result.push(<Grid key={l} size={12}>
                                        <FormControlLabel key={field.apiName}
                                            name={field.apiName}
                                            value="start"
                                            sx={{ml: 0,}}
                                            control={<Switch color="primary" onChange={onChangeSwitch} defaultChecked={parseBoolean(state.payload[field.apiName])}/>}
                                            label={field.label}
                                            required={field.required == 1 ? true : false}
                                            labelPlacement="start"/>
                                    </Grid>);
                                    break;
                                case 'REFERENCE':
                                    let value = null;
                                    if (state.payload[field.apiName]) {
                                        value = {
                                            id: state.payload[field.apiName],
                                            name_sys: state.payload[field.apiName + '_name_sys']
                                        };
                                    }
                                    result.push(<Grid key={l} size={12}>
                                        <ReferenceField key={field.apiName}
                                            name={field.apiName}
                                            objectId={field.objectReference}
                                            multiple={false}
                                            required={parseBoolean(field.required)} 
                                            label={field.label}
                                            record={value}
                                            onChange={onReferenceChange} 
                                            variant="standard"
                                            fullWidth={true}/>
                                    </Grid>);
                                    break;
                                case 'NUMBER':
                                    result.push(<Grid key={l} size={12}>
                                        <TextField key={field.apiName}
                                            name={field.apiName}
                                            onChange={onInputChange}
                                            inputProps={{ max: field.textLength }}
                                            value={state.payload[field.apiName] || null}
                                            type={'number'}
                                            required={field.required == 1 ? true : false}
                                            error={parseBoolean(field.required) ? state.payload[field.apiName] == null : false}
                                            label={field.label} 
                                            variant="standard" 
                                            fullWidth={true}/>
                                    </Grid>);
                                    break;
                                case 'DATE':
                                    result.push(<Grid key={l} size={12}>
                                        <LocalizationProvider dateAdapter={AdapterDayjs}>
                                        <DesktopDatePicker key={field.apiName} name={field.apiName}
                                            label={field.label}
                                            value={dayjs(state.payload[field.apiName])}
                                            minDate={dayjs('0-01-01')}
                                            onChange={(value, test) => {
                                                const newPayload = {
                                                    ...state.payload,
                                                    [field.apiName]: value.format('YYYY-MM-DD')
                                                };
                                                setState({
                                                    ...state,
                                                    payload: newPayload,
                                                });
                                                if (onChange) {
                                                    onChange(newPayload)
                                                }
                                            }}
                                            renderInput={(params) => <TextField {...params} 
                                            error={parseBoolean(field.required) ? state.payload[field.apiName] == null : false}
                                            required={parseBoolean(field.required) ? true : false} variant="standard" fullWidth={true}/>}/>
                                        </LocalizationProvider>
                                    </Grid>);
                                    break;
                                case 'DATETIME':
                                    result.push(<Grid key={l} size={12}>
                                        <LocalizationProvider dateAdapter={AdapterDayjs}>
                                    <DesktopDateTimePicker key={field.apiName} name={field.apiName}
                                        label={field.label}
                                        value={dayjs(state.payload[field.apiName])}
                                        minDate={dayjs('0-01-01')}
                                        onChange={(value, test) => {
                                            const newPayload = {
                                                ...state.payload,
                                                [field.apiName]: value.format('YYYY-MM-DD HH:mm:ss')
                                            };
                                            setState({
                                                ...state,
                                                payload: newPayload,
                                            });
                                            if (onChange) {
                                                onChange(newPayload)
                                            }
                                        }}
                                        renderInput={(params) => <TextField {...params} 
                                        error={parseBoolean(field.required) ? state.payload[field.apiName] == null : false}
                                        required={parseBoolean(field.required) ? true : false} variant="standard" fullWidth={true}/>}/>
                                    </LocalizationProvider>
                                    </Grid>);
                                    break;
                            }
                        }
                    }
                }
            }
        }
        return result;
    };

    const getLoading = () => {
        if (state.loading) {
            return <Stack spacing={2} sx={{mt: 2}}>
                <Skeleton variant="text" />
                <Stack direction={'row'}>
                    <Skeleton variant="circular" width={40} height={40} />
                    <Skeleton variant="text" width="100%" sx={{ml: 1}}/>
                </Stack>
                <Skeleton variant="rounded" height={20} />
                <Skeleton variant="rounded" height={20} />
                <Skeleton variant="rounded" height={20} />
                <Skeleton variant="rounded" height={30} />
                <Skeleton variant="rounded" height={30} />
                <Skeleton variant="rounded" height={40} />
          </Stack>
        }
    };

    const handleObjectChange = (event, newValue) => {
        let object = null;
        if (newValue) {
          object = {
              id: newValue.id,
              apiName: newValue.apiName,
              label: newValue.title || newValue.label,
            };
        }
          setState({
            ...state,
            object,
            objectsLoading: false,
            reloadRecords: true,
            loading: true,
            noRecord: false,
            record: null,
          });
      };

      const getContent = () => {
        if (!state.loading) {
            return <Grid container spacing={2} sx={{mt: 1}}>
                { state.mode == 'view' ? getViewFields() : getEditFields()}
            </Grid>;
        }
      };

      const getNoRecord = () => {
        if (!state.loading && (state.noRecord || state.metadata == null || state.record == null)) {
            return <Box sx={{ display: 'flex', width: '100%', height: 'calc(100% - 80px)', verticalAlign: 'middle', alignItems: 'center', justifyContent: 'center' }}>
                <Typography variant="subtitle1" gutterBottom>
                    No record selected.
                </Typography>
            </Box>;
        }
      };

      const onRecordChange = (apiName, record, metadata) => {
        setState({
            ...state,
            metadata,
            record,
            objectType: getObjectType(metadata, record),
            noRecord: record == null,
            loading: false,
        });
      };

      const onUpdate = () => {
        if (state.mode == 'view') {
            setState({
                ...state,
                mode: 'edit',
                payload: JSON.parse(JSON.stringify(state.record)),
            });
        } else {
            if (valid()) {
                setState({
                    ...state,
                    loading: true,
                    updatePayload: true,
                });
            } else {
                app.toast('warning', 'all required fields are not specified.');
            }
        }
      };

      const onView = () => {
        setState({
            ...state,
            mode: 'view',
        });
      };

      const getUpdateViewButton = () => {
        return [<ToggleButtonGroup
                key={'amperWidgetRecordUpdateButtonGroup'}
                exclusive
                aria-label="text alignment"
                sx={{ml: 1, mt: '5px'}}
                >
                    <ToggleButton
                        sx={{mt: '3px'}}
                        value="left"
                        disabled={state.record == null}
                        onClick={onUpdate}
                        key={'amperWidgetRecordUpdateButton'}
                        size="medium"
                        color="primary"
                        aria-label="Update">
                            <Tooltip title="Update">
                                <AdjustIcon sx={{ fontSize: 25}} color="primary"/>
                            </Tooltip>
                    </ToggleButton>,
                    <ToggleButton
                        sx={{mt: '3px'}}
                        value="right"
                        onClick={onView}
                        disabled={state.record == null || state.mode == 'view'}
                        key={'amperWidgetRecordViewButton'}
                        size="medium"
                        color="primary"
                        aria-label="View">
                            <Tooltip title="View">
                                <RemoveRedEyeIcon sx={{ fontSize: 25}} color={ state.record == null || state.mode == 'view' ? 'inactive' : 'primary'}/>
                            </Tooltip>
                    </ToggleButton>
                </ToggleButtonGroup>];
    };

    return (
        <Paper variant="outlined" sx={{p: 1}} style={{height: 'calc(100% - 17px)', overflowY: 'auto'}}>
            <Box sx={{
                display: 'flex',
                flexDirection: 'row',
                bgcolor: 'background.paper',
                }}>
                    <Box sx={{ flexGrow: 1 }}>
                        <Stack direction={'row'} sx={{mt: 1}}>
                            <Autocomplete
                                key="selectObject"
                                sx={{ width: 175}}
                                onChange={handleObjectChange}
                                options={state.objects}
                                value={state.object}
                                isOptionEqualToValue={(option, value) => option.id === value.id}
                                getOptionLabel={item => item.title || item.label}
                                renderInput={(params) => <TextField variant="standard" {...params} name="object" label="Select object" />}
                                />
                            <ReferenceField variant="standard" sx={{ml:2, width: 220}} key="selectRecord"
                            label={state.object != null ? 'Select ' + state.object.label + ' record' : "Select record"} 
                            multiple={false} objectId={state.object != null ? state.object.id: null}
                            metadata={state.metadata == null || state.object == null ? true : (state.metadata.Object.id != state.object.id)}
                            onChange={onRecordChange} reload={state.reloadRecords} record={state.record}/>
                        </Stack>
                    </Box>
                    <Box sx={{ display: 'flex', flexGrow: 0, textAlign: 'center', verticalAlign: 'middle', alignItems: 'center', justifyContent: 'center' }}>
                        {getUpdateViewButton()}
                    </Box>
            </Box>
            {getLoading()}
            {getContent()}
            {getNoRecord()}
        </Paper>
    );
};

export default forwardRef(RecordDetails)