import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Box from '@mui/material/Box';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import { Link as RouterLink, useLocation } from 'react-router-dom'
import PersonIcon from '@mui/icons-material/Person';
import SecurityIcon from '@mui/icons-material/Security';
import { breadcrumbs } from '../Breadcrambs';
import TuneIcon from '@mui/icons-material/Tune';
import StorageIcon from '@mui/icons-material/Storage';
import ChatIcon from '@mui/icons-material/Chat';
import InterpreterModeIcon from '@mui/icons-material/InterpreterMode';

export default function Administration({expanded}) {
    const rout = '/administration';
    const settingsRoute = breadcrumbs.administration.settings.path;
    const usersnRoute = '/administration/users';
    const sequrityProfilesRoute = '/administration/profiles';
    const nodesRoute  = '/administration/nodes';
    const chatRoute = breadcrumbs.administration.chat.path;
    const relationshipRoute = breadcrumbs.administration.relationship.path;

    const location = useLocation();
    const open = location.pathname.startsWith(rout);


    const getCollapsed = () => {
        return (
            <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
                <ListItemIcon>
                    <AdminPanelSettingsIcon sx={{fontSize: '35px', ml: '-2px'}} color={open ? 'primary' : 'inherit'}/>
                </ListItemIcon>
            </ListItemButton>
        );
    };

    const getExpanded = () => {
        return (
            <Box>
                <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
                    <ListItemIcon>
                        <AdminPanelSettingsIcon sx={{fontSize: '35px'}} color={open ? 'primary' : 'inherit'}/>
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -2, color: open ? 'secondary.contrastText' : 'secondary.menuText' }} primary="Administration" />
                    {location.pathname.startsWith(rout) ? <ExpandLess sx={{ color: 'secondary.contrastText' }} /> : <ExpandMore sx={{ color: 'secondary.menuText' }} />}
                </ListItemButton>
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        <ListItemButton component={RouterLink} to={settingsRoute} state={{expanded}} selected={location.pathname === settingsRoute} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <TuneIcon color={location.pathname.startsWith(settingsRoute) ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(settingsRoute) ? 'secondary.contrastText' : 'secondary.menuText'}} primary="Settings" />
                        </ListItemButton>
                    </List>
                    <List component="div" disablePadding>
                        <ListItemButton component={RouterLink} to={usersnRoute} state={{expanded}} selected={location.pathname === usersnRoute} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <PersonIcon color={location.pathname.startsWith(usersnRoute) ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(usersnRoute) ? 'secondary.contrastText' : 'inherit'  }} primary="Users" />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={sequrityProfilesRoute} state={{expanded}} selected={location.pathname === sequrityProfilesRoute} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <SecurityIcon color={location.pathname === sequrityProfilesRoute ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(sequrityProfilesRoute) ? 'secondary.contrastText' : 'inherit' }} primary="Profiles" />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={nodesRoute} state={{expanded}} selected={location.pathname === nodesRoute} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <StorageIcon color={location.pathname === nodesRoute ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(nodesRoute) ? 'secondary.contrastText' : 'inherit' }} primary="Nodes" />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={chatRoute} state={{expanded}} selected={location.pathname === chatRoute} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <ChatIcon color={location.pathname === chatRoute ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(chatRoute) ? 'secondary.contrastText' : 'inherit' }} primary="Chat" />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={relationshipRoute} state={{expanded}} selected={location.pathname === relationshipRoute} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <InterpreterModeIcon color={location.pathname === relationshipRoute ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(relationshipRoute) ? 'secondary.contrastText' : 'inherit' }} primary="Relationship" />
                        </ListItemButton>
                    </List>
                </Collapse>
            </Box>
    );
    };
    
  return expanded ? getExpanded() : getCollapsed();
}
