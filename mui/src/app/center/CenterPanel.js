import * as React from 'react';
import LeftMenuPanel from './LeftMenuPanel';
import RightContentPanel from './RightContentPanel';
import Box from '@mui/material/Box';
import { AppContext } from '../../App';

export default function CenterPanel({props}) {
  const {expanded, toast} = props;
  return (
    <Box height={'100%'} width={'100%'} color="inherit"
      sx={{
        bgcolor: (theme) => (theme.palette.primary.main),
        display: 'flex',
        flexDirection: 'row',
      }}>
        <LeftMenuPanel expanded={expanded}/>
        <RightContentPanel expanded={expanded} toast={toast}/>
    </Box>
  );
}
