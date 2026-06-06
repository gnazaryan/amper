import React, { useState, useRef, useEffect } from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ChatIcon from '@mui/icons-material/Chat';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Box from '@mui/material/Box';
import { Link as RouterLink, useLocation } from 'react-router-dom'

export default function Chat({expanded}) {
    const rout = '/chat';
    const location = useLocation();
    const open = location.pathname.startsWith(rout);

    const [state, setState] = useState(() => {
        return {
            loaded: false,
            data: {},
        };
    });

    useEffect(() => {
        if (!state.loaded) {
        }
    }, [state.loaded]);

    const getExpanded = () => {
        return (
            <Box>
                <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
                    <ListItemIcon>
                        <ChatIcon sx={{fontSize: '35px'}} color={open ? 'primary' : 'inherit'}/>
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -2, color: open ? 'secondary.contrastText' : 'secondary.menuText' }} primary="Chat" />
                    { open ? <ExpandLess sx={{ color: 'secondary.contrastText' }} /> : <ExpandMore sx={{ color: 'secondary.menuText' }} />}
                </ListItemButton>
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                    </List>
                </Collapse>
            </Box>
        );
    };

    const getCollapsed = () => {
        return (
        <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
            <ListItemIcon>
                <ChatIcon sx={{fontSize: '35px', ml: '-2px'}} color={open ? 'primary' : 'inherit'}/>
            </ListItemIcon>
        </ListItemButton>);
    };

  return expanded ? getExpanded() : getCollapsed();
}
