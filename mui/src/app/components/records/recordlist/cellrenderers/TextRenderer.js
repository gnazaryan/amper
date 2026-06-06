import AmperConstatns from '../../../../util/AmperConstants';

export default function TextRenderer (props) { 
  const { row, field } = props;
  const cachedValue = props.colDef.getPayloadValue(row[AmperConstatns.SYSTEM_FIELDS.IDENTIFIER], field);
  return cachedValue || props.value;
};