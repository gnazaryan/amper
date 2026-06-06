import * as React from 'react';
import TextField from '@mui/material/TextField';

function NumberInputValue(props) {
    const { item, applyValue, apply } = props;

    const handleFilterChange = (event) => {
      if (event.target.value && event.target.value.length > 0){
        applyValue({ ...item, value: parseInt(event.target.value) });
        apply();
      }
    };
    
    const onChange = (event) => {
      if (event.target.value && event.target.value.length > 0){
        applyValue({ ...item, value: parseInt(event.target.value) });
      }
    };
      return (
        <TextField
          onChange={onChange}
          label="Value"
          type="number"
          value={item.value}
          inputProps={{
              min: 0,
              onBlur: (event) => {
                handleFilterChange(event)
              }
          }}
          variant="standard"
        />
      );
}

export const numberOnlyOperators = (applyCallback) => {
    return [{
      label: 'Greater then',
      value: 'greaterThen',
      getApplyFilterFn: ()=>{},
      InputComponent: NumberInputValue,
      InputComponentProps: { type: 'number', apply: applyCallback },
    }, {
      label: 'Less then',
      value: 'lessThen',
      getApplyFilterFn: ()=>{},
      InputComponent: NumberInputValue,
      InputComponentProps: { type: 'number', apply: applyCallback },
    }, {
      label: 'Equals',
      value: 'equalsNumber',
      getApplyFilterFn: () => {},
      InputComponent: NumberInputValue,
      InputComponentProps: { type: 'number', apply: applyCallback },
    }, {
      label: 'Not equals',
      value: 'notEqualsNumber',
      getApplyFilterFn: () => {},
      InputComponent: NumberInputValue,
      InputComponentProps: { type: 'number', apply: applyCallback },
    },{
      label: 'Is not empty',
      value: 'isNotEmptyNumber',
      getApplyFilterFn: applyCallback,
    },{
        label: 'Is empty',
        value: 'isEmptyNumber',
        getApplyFilterFn: applyCallback,
    }];
};