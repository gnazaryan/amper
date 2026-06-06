import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import SettingsIcon from '@mui/icons-material/Settings';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Box from '@mui/material/Box';
import DataObjectIcon from '@mui/icons-material/DataObject';
import { Link as RouterLink, useLocation } from 'react-router-dom'

export default function Configuration({expanded}) {
    const rout = '/configuration';
    const objectsRout = '/configuration/objects';
    const location = useLocation();
    const open = location.pathname.startsWith(rout);

    const getExpanded = () => {
        return (
            <Box>
                <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
                    <ListItemIcon>
                        <SettingsIcon sx={{fontSize: '35px'}} color={open ? 'primary' : 'inherit'}/>
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -2, color: open ? 'secondary.contrastText' : 'secondary.menuText' }} primary="Configuration" />
                    { open ? <ExpandLess sx={{ color: 'secondary.contrastText' }} /> : <ExpandMore sx={{ color: 'secondary.menuText' }} />}
                </ListItemButton>
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        <ListItemButton component={RouterLink} to={objectsRout} state={{expanded}} selected={location.pathname === objectsRout} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <DataObjectIcon color={location.pathname.startsWith(objectsRout) ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(objectsRout) ? 'secondary.contrastText' : 'secondary.menuText'}} primary="Objects" />
                        </ListItemButton>
                    </List>
                </Collapse>
            </Box>
        );
    };

    const getCollapsed = () => {
        return (
        <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
            <ListItemIcon>
                <SettingsIcon sx={{fontSize: '35px', ml: '-2px'}} color={open ? 'primary' : 'inherit'}/>
            </ListItemIcon>
        </ListItemButton>);
    };

  return expanded ? getExpanded() : getCollapsed();
}
