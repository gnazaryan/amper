import Chip from '@mui/material/Chip';
import CheckBoxIcon from '@mui/icons-material/CheckBox';
import CheckBoxOutlineBlankIcon from '@mui/icons-material/CheckBoxOutlineBlank';
import {parseBoolean} from '../../../../util/BooleanUtil';
import AmperConstatns from '../../../../util/AmperConstants';

export default function BooleanRenderer(props) {
    const { hasFocus, value } = props;
    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    const valueBoolean = parseBoolean(cachedValue != null ? cachedValue : value);
    const dirty = cachedValue != null;
    return <Chip
                color={dirty ? 'secondary' : (valueBoolean ? 'primary' : 'inactive')} variant="outlined"
                icon={valueBoolean ? <CheckBoxIcon color={dirty ? 'secondary' : 'primary'}/> : <CheckBoxOutlineBlankIcon color={dirty ? 'secondary' : 'inactive'}/>}
                label={valueBoolean ? 'Active' : 'Inactive'}/>;
}