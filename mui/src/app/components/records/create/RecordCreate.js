import React, { useState, useImperativeHandle, forwardRef } from 'react';
import Grid from '@mui/material/Grid2';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';
import { DesktopDatePicker } from '@mui/x-date-pickers/DesktopDatePicker';
import dayjs from 'dayjs';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { ReferenceField } from '../fields/ReferenceField';
import { ObjectTypeField } from '../fields/ObjectTypeField';
import { DesktopDateTimePicker } from '@mui/x-date-pickers/DesktopDateTimePicker';
import LinearProgress from '@mui/material/LinearProgress';

function RecordCreate({ metadata, onChange}, ref) {

    const getFieldsMap = () => {
        const result = {};
        for (let i = 0; i < metadata.Fields.length; i++) {
            result[metadata.Fields[i].id] = metadata.Fields[i];
        }
        return result;
    };

    const [state, setState] = useState({
        objectType: metadata.ObjectTypes && metadata.ObjectTypes.length > 0 ? metadata.ObjectTypes[0] : null,
        fieldsMap: getFieldsMap(),
        payload: {
            'objectType_sys': metadata.ObjectTypes && metadata.ObjectTypes.length > 0 ? metadata.ObjectTypes[0].apiName : null,
        },
        loading: false,
    });

    useImperativeHandle(ref, () => ({
        onSubmit() {
            setState({
                ...state,
                loading: true,
            })
        },
        onSubmitComplete() {
            setState({
                ...state,
                loading: false,
            })
        },
    }));
    const isPayloadValid = (payload) => {
        let result = true;
        if (state.objectType && state.objectType.objectTypeFields, metadata.Fields && metadata.Fields.length > 0) {
            for (let i = 0; i < state.objectType.objectTypeFields.length; i++) {
                const objectTypeField = state.objectType.objectTypeFields[i];
                if (state.fieldsMap[objectTypeField.fieldId] != null) {
                    const field = state.fieldsMap[objectTypeField.fieldId];
                    if (!NON_MODIFIABLE_FIELDS.includes(field.apiName) && field.required == 1 && (payload == null ||payload[field.apiName] == null)) {
                        result = false;
                    }
                }
            }
        }
        return result;
    };

    const onReferenceChange = (name, value) => {
        var newPayload = {
            ...state.payload,
            [name]: value.identifier_sys,
        };
        setState({
            ...state,
            payload: newPayload,
        });
        if (onChange) {
            onChange(newPayload, isPayloadValid(newPayload))
        }
    };

    const onObjectTypeChange = (event, value) => {
        var newPayload = {
            ...state.payload,
            'objectType_sys': value.apiName,
        };
        setState({
            ...state,
            payload: newPayload,
            objectType: value,
        });
        if (onChange) {
            onChange(newPayload, isPayloadValid())
        }
    };

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
            onChange(newPayload, isPayloadValid(newPayload))
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
            onChange(newPayload, isPayloadValid(newPayload))
        }
    };
    
    const NON_MODIFIABLE_FIELDS = ['id', 'identifier_sys']
    const getFields = () => {
        const result = [];
        if (state.objectType && state.objectType.objectTypeFields, metadata.Fields && metadata.Fields.length > 0) {
            for (let i = 0; i < state.objectType.objectTypeFields.length; i++) {
                const objectTypeField = state.objectType.objectTypeFields[i];
                if (state.fieldsMap[objectTypeField.fieldId] != null) {
                    const field = state.fieldsMap[objectTypeField.fieldId];
                    if (field && !NON_MODIFIABLE_FIELDS.includes(field.apiName)) {
                        let fieldElemetn = null;
                        let switchValue = field.apiName == 'objectType_sys' ? field.apiName : field.type; 
                        switch(switchValue) {
                            case 'objectType_sys': 
                            fieldElemetn = <ObjectTypeField key={field.apiName}
                                                            name={field.apiName} 
                                                            objectTypes={metadata.ObjectTypes}
                                                            required={field.required == 1 ? true : false}
                                                            label={field.label}
                                                            onChange={onObjectTypeChange}
                                                            value={state.objectType} 
                                                            error={field.required == 1 ? state.payload[field.apiName] == null : false}/>
                            break;
                            case 'TEXT':
                                fieldElemetn = <TextField key={field.apiName}
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
                                                        fullWidth={true}/>;
                            break;
                            case 'BOOLEAN':
                                fieldElemetn = <FormControlLabel key={field.apiName}
                                                        name={field.apiName}
                                                        value="start"
                                                        sx={{mt:2,}}
                                                        control={<Switch color="primary" onChange={onChangeSwitch} checked={state.payload[field.apiName] || false}/>}
                                                        label={field.label} 
                                                        required={field.required == 1 ? true : false} 
                                                        labelPlacement="start"/>
                            break;
                            case 'REFERENCE':
                                fieldElemetn = <ReferenceField key={field.apiName}
                                                        name={field.apiName}
                                                        objectId={field.objectReference}
                                                        multiple={false}
                                                        required={field.required == 1 ? true : false} 
                                                        label={field.label}
                                                        onChange={onReferenceChange} 
                                                        variant="standard"
                                                        fullWidth={true}/>;
                            break;
                            case 'NUMBER':
                                fieldElemetn = <TextField key={field.apiName}
                                                        name={field.apiName}
                                                        onChange={onInputChange}
                                                        inputProps={{ max: field.textLength }}
                                                        value={state.payload[field.apiName] || ''}
                                                        type={'number'}
                                                        required={field.required == 1 ? true : false}
                                                        error={field.required == 1 ? state.payload[field.apiName] == null : false}
                                                        label={field.label} 
                                                        variant="standard" 
                                                        fullWidth={true}/>;
                            break;
                            case 'DATE':
                                fieldElemetn = <LocalizationProvider dateAdapter={AdapterDayjs}>
                                    <DesktopDatePicker key={field.apiName} name={field.apiName}
                                        label={field.label}
                                        value={dayjs(state.payload[field.apiName] || '')}
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
                                                onChange(newPayload, isPayloadValid(newPayload))
                                            }
                                        }}
                                        renderInput={(params) => <TextField {...params} 
                                        error={field.required == 1 ? state.payload[field.apiName] == null : false}
                                        required={field.required == 1 ? true : false} variant="standard" fullWidth={true}/>}/>
                                    </LocalizationProvider>
                            break;
                            case 'DATETIME':
                                fieldElemetn = <LocalizationProvider dateAdapter={AdapterDayjs}>
                                    <DesktopDateTimePicker key={field.apiName} name={field.apiName}
                                        label={field.label}
                                        value={dayjs(state.payload[field.apiName] || '')}
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
                                                onChange(newPayload, isPayloadValid(newPayload))
                                            }
                                        }}
                                        renderInput={(params) => <TextField {...params} 
                                        error={field.required == 1 ? state.payload[field.apiName] == null : false}
                                        required={field.required == 1 ? true : false} variant="standard" fullWidth={true}/>}/>
                                    </LocalizationProvider>
                            break;
                        }
                        if (fieldElemetn) {
                            result.push(<Grid key={i} size={12}>
                                {fieldElemetn}
                            </Grid>);    
                        }
                    }
                }
            }
        }
        return result;
    };

    return <Box sx={{ flexGrow: 1, mt: 1 }}>
        {state.loading ? <LinearProgress sx={{mb: 2}} /> : ''}
        <Grid container spacing={2}>
            {getFields()}
        </Grid>
    </Box>;
};

export default forwardRef(RecordCreate)