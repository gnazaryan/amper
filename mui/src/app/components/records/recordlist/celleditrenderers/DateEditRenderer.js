import React, {  } from 'react';
import AmperConstatns from '../../../../util/AmperConstants';
import Box from '@mui/material/Box';
import { DesktopDatePicker } from '@mui/x-date-pickers/DesktopDatePicker';
import dayjs from 'dayjs';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { TextField } from '@mui/material';

export default function DateEditRenderer(props) {
    const { value } = props;
    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    const finalValue = cachedValue || value;
    return <LocalizationProvider dateAdapter={AdapterDayjs}>
    <Box sx={{ml: '4px', mr: '4px'}}><DesktopDatePicker name={props['field']}
        value={dayjs(finalValue)}
        minDate={dayjs('0-01-01')}
        onChange={(value, test) => {
            const valueFormatted = value.format('YYYY-MM-DD');
            props.colDef.cacheValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field'], valueFormatted)
        }}
        renderInput={(params) => <TextField {...params} 
        variant="standard" fullWidth={true}/>}/></Box>
    </LocalizationProvider>;
}