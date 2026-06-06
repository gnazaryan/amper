import React, { useState, useEffect } from 'react';
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Chip from '@mui/material/Chip';
import DataStore from '../../../../data/DataStore';
import HostManager from '../../../../../HostManager';
import { debounceLatest } from '../../../../amper/Instruments';
import Convenience from '../../../../help/Convenience';

function ReferenceInputValue(props) {
    const { item, applyValue, apply } = props;

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
        values: [],
        startId: 0,
        objectId: props.objectReferenceId,
        loading: true,
        loadingMask: true,
        searchUpdate: debounceLatest(searchUpdateDebounce, 1000)
    });

    useEffect(() => {
        if (state.loading) {
          setState({
            ...state,
            loading: false,
            loadingMask: false,
          })
          getDataStore().load((result)=> {
              setState({
                ...state,
                loading: false,
                loadingMask: false,
                options: result.data || [],
                metadata: state.metadata ? state.metadata : result.metadata,
                searchTerm: null,
              })
          });
        }
      });

      const LIMIT = 100;
    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}records/fetch`,
            requestMethod: "POST",
            parameters: {
                objectId: state.objectId,
                limit: LIMIT,
                startId: state.startId,
                search: Convenience.hasValue(state.searchTerm) ? JSON.stringify({term: [{apiName: 'name_sys', operator: 'contains', value: state.searchTerm}]}) : "",
                metadata: false,
            }
        });
    };

      const handleFilterChange = (event, newValue) => {
        applyValue({ ...item, value: newValue.map(record => record['identifier_sys']), originalValue: newValue});
          setState({
              ...state,
              loadingMask: false,
              loading: false,
              values: newValue,
          });
      };
    
      const onBluer = () => {
        apply();
      };


      const onChange = (event) => {
        state.searchUpdate(event.target.value, state);
        setState({
            ...state,
            loadingMask: true,
            options : [],
        })
      };

      return (
        <Autocomplete
            multiple
            onChange={handleFilterChange}
            onBlur={onBluer}
            options={state.options}
            getOptionLabel={(option) => (option['name_sys'])}
            isOptionEqualToValue={(option, value) => option.id === value.id}
            loading={state.loadingMask}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip key={option.id} variant="outlined" sx={{maxHeight: '25px', height: '25px'}} label={option['name_sys']} />
              ))
            }
            value={item.originalValue}
            renderInput={(params) => (
            <TextField
                {...params}
                label="Values"
                onChange={onChange}
                variant="standard"
            />
            )}
        />
      );
}

export const referenceOperators = (applyCallback, objectReferenceId) => {
    return [{
        label: 'Is any of',
        value: 'hasAnyOf',
        getApplyFilterFn: () => {},
        InputComponent: ReferenceInputValue,
        InputComponentProps: { type: 'referenceautocomplete', objectReferenceId: objectReferenceId, apply: applyCallback},
    },{
        label: 'Is none of',
        value: 'hasNoneOf',
        getApplyFilterFn: () => {},
        InputComponent: ReferenceInputValue,
        InputComponentProps: { type: 'referenceautocomplete', objectReferenceId: objectReferenceId, apply: applyCallback},
    },{
        label: 'Is not empty',
        value: 'isNotEmpty',
        getApplyFilterFn: applyCallback,
    },{
        label: 'Is empty',
        value: 'isEmpty',
        getApplyFilterFn: applyCallback,
    }];
};