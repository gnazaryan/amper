import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import CloudQueueIcon from '@mui/icons-material/CloudQueue';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Box from '@mui/material/Box';
import { Link as RouterLink, useLocation } from 'react-router-dom'
import { breadcrumbs } from '../Breadcrambs';
import FolderIcon from '@mui/icons-material/Folder';
import ShareIcon from '@mui/icons-material/Share';

export default function Drive({expanded}) {
    const rout = '/drive';
    const location = useLocation();
    const open = location.pathname.startsWith(rout);
    const getExpanded = () => {
        return (
            <Box>
                <ListItemButton component={RouterLink} to={rout} state={{expanded}} >
                    <ListItemIcon>
                        <CloudQueueIcon sx={{fontSize: '35px'}} color={open ? 'primary' : 'inherit'}/>
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -2, color: open ? 'secondary.contrastText' : 'secondary.menuText' }} primary="Drive" />
                    {open ? <ExpandLess sx={{ color: 'secondary.contrastText' }} /> : <ExpandMore sx={{ color: 'secondary.menuText' }} />}
                </ListItemButton>
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        <ListItemButton component={RouterLink} to={breadcrumbs.drive.files.path} state={{expanded}} selected={location.pathname === breadcrumbs.drive.files.path} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <FolderIcon  color={location.pathname.startsWith(breadcrumbs.drive.files.path) ? 'primary' : 'inherit'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(breadcrumbs.drive.files.path) ? 'secondary.contrastText' : 'secondary.menuText' }} primary={breadcrumbs.drive.files.label} />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={breadcrumbs.drive.shared.path} state={{expanded}} selected={location.pathname === breadcrumbs.drive.shared.path} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <ShareIcon  color={location.pathname.startsWith(breadcrumbs.drive.shared.path) ? 'primary' : 'inherit'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(breadcrumbs.drive.shared.path) ? 'secondary.contrastText' : 'secondary.menuText' }} primary={breadcrumbs.drive.shared.label} />
                        </ListItemButton>
                    </List>
                </Collapse>
            </Box>
        );
    }

    const getCollapsed = () => {
        return (
            <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
                <ListItemIcon>
                    <CloudQueueIcon sx={{fontSize: '35px', ml: '-2px'}} color={open ? 'primary' : 'inherit'}/>
                </ListItemIcon>
            </ListItemButton>
        );
    }

    return (
        expanded ? getExpanded() : getCollapsed()
    );
}
