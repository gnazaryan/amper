import React, { useState, useEffect } from 'react';
import DataStore from '../../../../app/data/DataStore';
import HostManager from '../../../../HostManager';
import Autocomplete from '@mui/material/Autocomplete';
import { debounceLatest } from '../../../amper/Instruments';
import { TextField } from '@mui/material';
import Convenience from '../../../help/Convenience';
import {Chip} from '@mui/material';

export const ReferenceField = ({ name, required, objectId, label, onChange, multiple, variant, record, sx, metadata, reload }) => {

    const searchUpdateDebounce = (getTerm, state) => {
        setState({
            ...state,
            searchTerm: getTerm(),
            loading: true,
            loadingMask: true,
        })
    };

    const [state, setState] = useState({
        options : [],
        values: multiple ? [] : null,
        start: 0,
        loading: true,
        loadingMask: true,
        searchUpdate: debounceLatest(searchUpdateDebounce, 1000)
    });

    const options = state.options;
    if (record != null && record.identifier_sys) {
        let hasItem = false;
        for (let i = 0; i < options.length; i++) {
            if (options[i].identifier_sys === record.identifier_sys) {
                hasItem = true;
                break;
            }
        }
        if (!hasItem) {
            options.push(record);
        }
    }

    useEffect(() => {
        if (objectId != null && state.loading || (objectId != null && reload != null && reload == true)) {
          getDataStore().load((result)=> {
              setState({
                ...state,
                loading: false,
                loadingMask: false,
                options: result.data || [],
                metadata: result.metadata ? result.metadata : state.metadata,
              })
          });
        }
      }, [state.loading, objectId, reload]);

      const LIMIT = 100;
    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}records/fetch`,
            requestMethod: "POST",
            parameters: {
                objectId: objectId,
                limit: LIMIT,
                start: state.start,
                search: Convenience.hasValue(state.searchTerm) ? JSON.stringify({term: [{apiName: 'name_sys', operator: 'contains', value: state.searchTerm}]}) : "",
                metadata: metadata != null ? metadata : false,
            }
        });
    };

    const onValueChange = (event, newValue) => {
        setState({
            ...state,
            loading: false,
            values: newValue,
        });
        if (onChange) {
            onChange(name, newValue, state.metadata);
        }
    };

    const onBluer = () => {

    };

    const onTermChange = (event) => {
        state.searchUpdate(event.target.value, state);
        setState({
            ...state,
            loadingMask: true,
            options : [],
        })
    };
    
    return (
        <Autocomplete
            sx={{...sx}}
            multiple={multiple != null ? multiple : true}
            onChange={onValueChange}
            onBlur={onBluer}
            options={options}
            getOptionLabel={(option) => (option['name_sys'])}
            isOptionEqualToValue={(option, value) => {
                return option.id === value.id;
            }}
            loading={state.loadingMask}
            fullWidth= {true}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip key={option.id} variant="outlined" sx={{maxHeight: '25px', height: '25px'}} label={option['name_sys']} />
              ))
            }
            value={reload === true ? record : (record || state.values)}
            renderInput={(params) => (
                <TextField
                    {...params}
                    name={name}
                    required={required}
                    error={required ? (multiple ? state.values.length < 1 : !(state.values != null || record != null)) : false}
                    label={label}
                    onChange={onTermChange}
                    variant={variant || 'standard'}
                />
            )}
        />
    );
};