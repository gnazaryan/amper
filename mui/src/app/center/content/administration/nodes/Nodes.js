import Box from '@mui/material/Box';
import NodeList from '../../../../components/adminstration/nodes/NodeList';

export default function Nodes() {
  return (
        <Box sx={{ height: '100%', width: 'calc(100% - 25px)'}}>
            <NodeList></NodeList>
        </Box>
    );
}
