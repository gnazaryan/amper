import React, { useState } from 'react';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import Chip from '@mui/material/Chip';

function ObjectTypeValue(props) {
    const { item, applyValue, apply } = props;

      const handleObjectTypeChange = (event, newValue) => {
        const originalValue = [];
        if (newValue && newValue.length > 0) {
          for (let l = 0; l < newValue.length; l++) {
            originalValue.push({
              id: newValue[l].id,
              apiName: newValue[l].apiName,
              label : newValue[l].label,
            });
          }
        }
        applyValue({ ...item, value: newValue.map(objectType => objectType.apiName), originalValue: originalValue});
      };
    
      const onBlurHandler = () => {
        apply();
      };
      
      return (
        <Autocomplete
                variant="standard"
                value={item.originalValue}
                onChange={handleObjectTypeChange}
                onBlur={onBlurHandler}
                multiple
                renderTags={(value, getTagProps) =>
                    value.map((option, index) => (
                      <Chip key={option.id} variant="outlined" sx={{maxHeight: '25px', height: '25px'}} label={option.label} />
                    ))
                  }
                options={props.objectTypes}
                renderInput={(params) => <TextField variant="standard" {...params} label="Value" name="objectType"/>}
            />
      );
}

export const objectTypeOnlyOperators = (applyCallback, objectTypes) => {
    return [{
        label: 'Is any of',
        value: 'hasAnyOf',
        getApplyFilterFn: () => {},
        InputComponent: ObjectTypeValue,
        InputComponentProps: { type: 'objectType', objectTypes, apply: applyCallback},
    },{
        label: 'is none of',
        value: 'hasNoneOf',
        getApplyFilterFn: () => {},
        InputComponent: ObjectTypeValue,
        InputComponentProps: { type: 'objectType', objectTypes, apply: applyCallback},
    }];
};