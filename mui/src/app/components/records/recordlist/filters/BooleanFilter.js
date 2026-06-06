import React from 'react';
import Box from '@mui/material/Box';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

function BooleanInputValue(props) {
    const { item, applyValue, apply, focusElementRef } = props;

    const handleFilterChange = (event) => {
        applyValue({ ...item, value: event.target.checked});
        apply();
      };
      return (
        <Box
          sx={{
            display: 'inline-flex',
            flexDirection: 'row',
            alignItems: 'center',
            pl: '10px',
            mt: '15px'
          }}
        >
            <Stack direction="row" spacing={1} alignItems="center">
            <Typography>Passive</Typography>
                <Switch
                    placeholder="Filter value"
                    onChange={handleFilterChange}
                    defaultChecked={item.value}
                />
                <Typography>Active</Typography>
            </Stack>
        </Box>
      );
}

export const booleanOnlyOperators = (applyCallback) => {
  return [{
      label: 'Is',
      value: 'is',
      getApplyFilterFn: () => {},
      InputComponent: BooleanInputValue,
      InputComponentProps: { type: 'checkbox', apply: applyCallback },
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