import React, { useState } from 'react';
import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Chip from '@mui/material/Chip';

function AnyAutocompleteInputValue(props) {
    const { item, applyValue, apply } = props;
    const [state, setState] = useState({
        options : [],
    });
      const handleFilterChange = (event, newValue) => {
        const result = [];
        for (let i = 0; i < newValue.length; i++) {
            if (newValue[i].term) {
                result.push(newValue[i].term);
            }
        }
        applyValue({ ...item, value: result, originalValue: newValue});
      };
    
      const handleBlure = () => {
        apply();
      };

      const onChange = (event) => {
        setState({
            options: [{
                term: event.target.value,
                filter: () => {}
            }]
        });
      };

      return (
        <Autocomplete
            multiple
            onChange={handleFilterChange}
            onBlur={handleBlure}
            options={state.options}
            value={item.originalValue}
            getOptionLabel={(option) => option.term}
            renderTags={(value, getTagProps) =>
                value.map((option, index) => (
                  <Chip variant="outlined" sx={{maxHeight: '25px', height: '25px'}} label={option.term} />
                ))
              }
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

function AnyInputValue(props) {
    const { item, applyValue, apply } = props;

      const onBlure = (event) => {
        apply();
      };
    
      const onChange = (event) => {
        if (event.target.value && event.target.value.length > 0){
          applyValue({ ...item, value: event.target.value });
        }
      };
      
      return (
        <TextField
            label="Value"
            variant="standard"
            onChange={onChange}
            value={item.value}
            inputProps={{
                onBlur: (event) => {
                    onBlure(event)
                }
            }}
        />
      );
}

export const anyOperators = (applyCallback) => {
    return [{
        label: 'Contains',
        value: 'contains',
        getApplyFilterFn: ()=>{},
        InputComponent: AnyInputValue,
        InputComponentProps: { type: 'any', apply: applyCallback },
    },{
        label: 'Is any of',
        value: 'hasAnyOf',
        getApplyFilterFn: ()=>{},
        InputComponent: AnyAutocompleteInputValue,
        InputComponentProps: { type: 'anyautocomplete', apply: applyCallback },
    },{
      label: 'is none of',
      value: 'hasNoneOf',
      getApplyFilterFn: () => {},
      InputComponent: AnyAutocompleteInputValue,
      InputComponentProps: { type: 'anyautocomplete', apply: applyCallback},
    },{
        label: 'Equals',
        value: 'equals',
        getApplyFilterFn: ()=>{},
        InputComponent: AnyInputValue,
        InputComponentProps: { type: 'any', apply: applyCallback },
    },{
        label: 'Starts with',
        value: 'startsWith',
        getApplyFilterFn: ()=>{},
        InputComponent: AnyInputValue,
        InputComponentProps: { type: 'any', apply: applyCallback },
    },{
        label: 'Ends with',
        value: 'endsWith',
        getApplyFilterFn: ()=>{},
        InputComponent: AnyInputValue,
        InputComponentProps: { type: 'any', apply: applyCallback },
    },{
        label: 'Is empty',
        value: 'isEmpty',
        getApplyFilterFn: applyCallback,
    },{
        label: 'Is not empty',
        value: 'isNotEmpty',
        getApplyFilterFn: applyCallback,
    }];
};