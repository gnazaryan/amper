import Box from '@mui/material/Box';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import Dashboard from './menu/Dashboard';
import Configuration from './menu/Configuration';
import Administration from './menu/Administration';
import Drive from './menu/Drive';
import Chat from './menu/Chat'
import Profile from './menu/Profile';
import AmperConstatns from '../util/AmperConstants';
import Communication from './menu/Email';

export default function LeftMenuPanel({expanded}) {
  return (
    <Box height={'calc(100% - 2px)'}
      sx={{
        display: 'flex',
        bgcolor: 'background.paper',
        mr: '1px',
        ml: '2px',
        borderRadius: (theme) => (theme.palette.primary.borderRadius)}}>
          <Collapse orientation="horizontal" in={expanded} collapsedSize={60}>
            <Box sx={{ borderRadius: (theme) => (theme.palette.primary.borderRadius), width: (expanded ? AmperConstatns.LEFT_MENU_WIDTH : 60), height: 'calc(100% - 8px)' }}>
              <List
                sx={{borderRadius: (theme) => (theme.palette.primary.borderRadius), width: (expanded ? AmperConstatns.LEFT_MENU_WIDTH : 60), height: 'calc(100% - 8px)' }}
                component="nav"
                aria-labelledby="nested-list-subheader">
                <Profile expanded={expanded}/>
                <Dashboard expanded={expanded}/>
                <Chat expanded={expanded}/>
                <Communication expanded={expanded}/>
                <Drive expanded={expanded}/>
                <Configuration expanded={expanded}/>
                <Administration expanded={expanded}/>
              </List>
            </Box>
          </Collapse>
      </Box>
    );
}
