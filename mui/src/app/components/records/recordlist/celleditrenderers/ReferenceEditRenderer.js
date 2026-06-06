import React, { useState, useEffect } from 'react';
import Link from '@mui/material/Link';
import { ReferenceField } from '../../fields/ReferenceField';
import Convenience from '../../../../help/Convenience';
import AmperConstatns from '../../../../util/AmperConstants';

export default function ReferenceEditRenderer(props) {
    const { hasFocus, value,  } = props;
    let objectId = props.colDef.referenceObjectId;
    const onReferenceChange = (event, value) => {
        props.colDef.cacheValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field'], value[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], value)
    };

    const cacheValue = props.colDef.getCacheValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    let originalVlaue = null
    if (value != null) {
        originalVlaue = {
            [props.field]: value,
            'name_sys': props.row[props.field + '_name_sys']
        };
    }
    return <ReferenceField name={props.field}
                objectId={objectId}
                multiple={false}
                onChange={onReferenceChange}
                record={cacheValue || originalVlaue}
                variant="outlined"
                fullWidth={true}/>;
}