import AmperConstatns from '../../../../util/AmperConstants';

export default function DateRenderer(props) {
    const { value } = props;
    const cachedValue = props.colDef.getPayloadValue(props.row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], props['field']);
    return cachedValue || value;
}