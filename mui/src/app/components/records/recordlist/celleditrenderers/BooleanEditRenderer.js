import React, { useState, useEffect } from 'react';
import FormControlLabel from '@mui/material/FormControlLabel';
import {parseBoolean} from '../../../../util/BooleanUtil';
import Switch from '@mui/material/Switch';
import AmperConstatns from '../../../../util/AmperConstants';

export default function BooleanEditRenderer(props) {
    const { hasFocus, value } = props;
    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    const valueBoolean = parseBoolean(cachedValue != null ? cachedValue : value);
    const onChange = (event) => {
        const value = event.target.checked === true ? 1 : 0;
        props.colDef.cacheValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field'], value)
    };
    return <FormControlLabel label={valueBoolean ? 'Inactivate' : 'Activate'} sx={{color: 'primary.text'}} control={<Switch defaultChecked={valueBoolean} onChange={onChange}/>}></FormControlLabel>;
}