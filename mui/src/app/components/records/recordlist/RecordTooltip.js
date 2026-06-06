import Tooltip from '@mui/material/Tooltip';

export default function RecordTooltip(props) {
    return <Tooltip {...props}>{props.children}</Tooltip>
}