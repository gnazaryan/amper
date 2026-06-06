import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import DashboardIcon from '@mui/icons-material/Dashboard';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Box from '@mui/material/Box';
import SpeedIcon from '@mui/icons-material/Speed';
import { Link as RouterLink, useLocation } from 'react-router-dom'
import { breadcrumbs } from '../Breadcrambs';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';

export default function Dashboard({expanded}) {
    const rout = '/dashboard';
    const location = useLocation();
    const open = location.pathname.startsWith(rout);
    const getExpanded = () => {
        return (
            <Box>
                <ListItemButton component={RouterLink} to={rout} state={{expanded}} >
                    <ListItemIcon>
                        <DashboardIcon sx={{fontSize: '35px'}} color={open ? 'primary' : 'inherit'}/>
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -2, color: open ? 'secondary.contrastText' : 'secondary.menuText' }} primary="Dashboard" />
                    {open ? <ExpandLess sx={{ color: 'secondary.contrastText' }} /> : <ExpandMore sx={{ color: 'secondary.menuText' }} />}
                </ListItemButton>
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        <ListItemButton component={RouterLink} to={breadcrumbs.dashboard.overview.path} state={{expanded}} selected={location.pathname === breadcrumbs.dashboard.overview.path} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <SpeedIcon  color={location.pathname.startsWith(breadcrumbs.dashboard.overview.path) ? 'primary' : 'inherit'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(breadcrumbs.dashboard.overview.path) ? 'secondary.contrastText' : 'secondary.menuText' }} primary={breadcrumbs.dashboard.overview.label} />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={breadcrumbs.dashboard.add.path} state={{expanded}} selected={location.pathname === breadcrumbs.dashboard.add.path} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <AddCircleOutlineIcon  color={location.pathname.startsWith(breadcrumbs.dashboard.add.path) ? 'primary' : 'inherit'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(breadcrumbs.dashboard.add.path) ? 'secondary.contrastText' : 'secondary.menuText' }} primary={breadcrumbs.dashboard.add.label} />
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
                    <DashboardIcon sx={{fontSize: '35px', ml: '-2px'}} color={open ? 'primary' : 'inherit'}/>
                </ListItemIcon>
            </ListItemButton>
        );
    }

    return (
        expanded ? getExpanded() : getCollapsed()
    );
}
