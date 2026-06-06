import React, {  } from 'react';
import AmperConstatns from '../../../../util/AmperConstants';
import Box from '@mui/material/Box';
import Tooltip from '@mui/material/Tooltip';
import { DesktopDateTimePicker } from '@mui/x-date-pickers/DesktopDateTimePicker';
import dayjs from 'dayjs';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { TextField } from '@mui/material';

export default function DateTimeEditRenderer(props) {
    const { value } = props;
    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    const finalValue = cachedValue || value;
    return <LocalizationProvider dateAdapter={AdapterDayjs}>
    <Box sx={{ml: '4px', mr: '4px'}}><DesktopDateTimePicker name={props['field']}
        value={dayjs(finalValue)}
        minDate={dayjs('0-01-01')}
        onChange={(value) => {
            const valueFormatted = value.format('YYYY-MM-DD HH:mm:ss');
            props.colDef.cacheValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field'], valueFormatted)
        }}
        renderInput={(params) => <TextField {...params} 
        variant="standard" fullWidth={true}/>}/></Box>
    </LocalizationProvider>;
}