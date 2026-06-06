import React, { useState, useRef, useEffect } from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import EmailIcon from '@mui/icons-material/Email';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Box from '@mui/material/Box';
import { post } from '../../data/Submit'
import { Link as RouterLink, useLocation } from 'react-router-dom'
import { sessionManager } from '../../../SessionManager';
import { truncate } from '../../amper/Instruments';
import Tooltip from '@mui/material/Tooltip';
import AlternateEmailIcon from '@mui/icons-material/AlternateEmail';

export default function Communication({expanded}) {
    const rout = '/email';
    const emailRout = '/communication/email';
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

    const getEmailRout = (item) => {
        return rout + '/' + item.email;
    };

    const getEmails = () => {
        const result = [];
        const user = sessionManager.getUser();
        if (user.emails != null) {
            for (let i = 0; i < user.emails.length; i++) {
                const emailRout = getEmailRout(user.emails[i]);
                result.push(<Tooltip arrow title={user.emails[i].label + ' (' + user.emails[i].email + ')'}>
                <ListItemButton component={RouterLink} to={emailRout} state={{expanded}} selected={location.pathname === emailRout} sx={{ pl: 4 }}>
                    <ListItemIcon>
                        <AlternateEmailIcon color={location.pathname.startsWith(emailRout) ? 'primary' : 'secondary.menuText'}/>
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(emailRout) ? 'secondary.contrastText' : 'secondary.menuText'}} primary={truncate(user.emails[i].label + ' (' + user.emails[i].email + ')', 20, false)} />
                </ListItemButton></Tooltip>);
            }
        }
        return result;
    };

    const getExpanded = () => {
        return (
            <Box>
                <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
                    <ListItemIcon>
                        <EmailIcon sx={{fontSize: '35px'}} color={open ? 'primary' : 'inherit'}/>
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -2, color: open ? 'secondary.contrastText' : 'secondary.menuText' }} primary="Email" />
                    { open ? <ExpandLess sx={{ color: 'secondary.contrastText' }} /> : <ExpandMore sx={{ color: 'secondary.menuText' }} />}
                </ListItemButton>
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        {getEmails()}
                    </List>
                </Collapse>
            </Box>
        );
    };

    const getCollapsed = () => {
        return (
        <ListItemButton component={RouterLink} to={rout} state={{expanded}}>
            <ListItemIcon>
                <EmailIcon sx={{fontSize: '35px', ml: '-2px'}} color={open ? 'primary' : 'inherit'}/>
            </ListItemIcon>
        </ListItemButton>);
    };

  return expanded ? getExpanded() : getCollapsed();
}
