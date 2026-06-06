import AmperConstatns from '../../../../util/AmperConstants';

export default function NumberRenderer(props) {
    const { value } = props;
    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    return cachedValue || value;
}