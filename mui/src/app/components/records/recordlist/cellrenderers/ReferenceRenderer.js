import Link from '@mui/material/Link';
import AmperConstatns from '../../../../util/AmperConstants';

export default function ReferenceRenderer(props) {
    const { value } = props;
    const nameValue = props.row[props.field + '_name_sys']

    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    const cachedValueComplete = props.colDef.getCacheValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    const dirty = cachedValue != null;

    return (
        <Link
            href="#"
            sx={{color: dirty ? 'primary.contrastText' : 'primary.main', textDecoration: 'underline', textDecorationColor: dirty ? 'primary.contrastText' : 'primary.main'}}
            underline="hover">
            {cachedValue != null ? cachedValueComplete[AmperConstatns.SYSTEM_FIELDS.NAME] : (nameValue || value)}
        </Link>
    );
}