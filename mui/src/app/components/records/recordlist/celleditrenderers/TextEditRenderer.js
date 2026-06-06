import React, {  } from 'react';
import AmperConstatns from '../../../../util/AmperConstants';
import { TextField } from '@mui/material';

export default function TextEditRenderer(props) {
    const { hasFocus, value } = props;
    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    const onChange = (event) => {
        const value = event.target.value;
        props.colDef.cacheValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field'], value)
    };

    return <TextField autoFocus={true} variant="outlined" fullWidth={true} onChange={onChange} value={cachedValue || value}/>;
}