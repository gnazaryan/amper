import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import Collapse from '@mui/material/Collapse';
import List from '@mui/material/List';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import Box from '@mui/material/Box';
import TuneIcon from '@mui/icons-material/Tune';
import { Link as RouterLink, useLocation } from 'react-router-dom'
import { breadcrumbs } from '../Breadcrambs';
import { sessionManager } from '../../../SessionManager';
import { gettStoreValue } from '../../amper/Instruments';
import Avatar from '@mui/material/Avatar';
import Convenience from '../../help/Convenience';
import AccessibilityIcon from '@mui/icons-material/Accessibility';
import LocalActivityIcon from '@mui/icons-material/LocalActivity';

export default function Profile({expanded}) {
    const settingsRout = breadcrumbs.profile.path + '/' + breadcrumbs.profile.settings.key;
    const aboutRout = breadcrumbs.profile.path + '/' + breadcrumbs.profile.about.key;
    const overviewRout = breadcrumbs.profile.path + '/' + breadcrumbs.profile.overview.key;

    const location = useLocation();
    const open = location.pathname.startsWith(breadcrumbs.profile.path);

    const getImageSource = () => {
        const photo = gettStoreValue('user_photo');
        if (Convenience.hasValue(photo)) {
          return 'data:image/png;base64,' + photo;
        }
        return '/static/images/avatar/2.jpg';
      };

    const user = sessionManager.getUser();
    const getExpanded = () => {
        return (
            <Box>
                <ListItemButton component={RouterLink} to={breadcrumbs.profile.path} state={{expanded}}>
                    <ListItemIcon>
                        <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main' }} alt={user.firstName + ' ' + user.lastName} src={getImageSource()} />
                    </ListItemIcon>
                    <ListItemText sx={{ ml: -1, color: open ? 'secondary.contrastText' : 'secondary.menuText' }} primary={sessionManager.getUser().firstName + ' ' + sessionManager.getUser().lastName} />
                    { open ? <ExpandLess sx={{ color: 'secondary.contrastText' }} /> : <ExpandMore sx={{ color: 'secondary.menuText' }} />}
                </ListItemButton>
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding>
                        <ListItemButton component={RouterLink} to={overviewRout} state={{expanded}} selected={location.pathname === overviewRout} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <LocalActivityIcon color={location.pathname.startsWith(overviewRout) ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(overviewRout) ? 'secondary.contrastText' : 'secondary.menuText'}} primary="Overview" />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={aboutRout} state={{expanded}} selected={location.pathname === aboutRout} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <AccessibilityIcon color={location.pathname.startsWith(aboutRout) ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(aboutRout) ? 'secondary.contrastText' : 'secondary.menuText'}} primary="About" />
                        </ListItemButton>
                        <ListItemButton component={RouterLink} to={settingsRout} state={{expanded}} selected={location.pathname === settingsRout} sx={{ pl: 4 }}>
                            <ListItemIcon>
                                <TuneIcon color={location.pathname.startsWith(settingsRout) ? 'primary' : 'secondary.menuText'}/>
                            </ListItemIcon>
                            <ListItemText sx={{ ml: -2, color: location.pathname.startsWith(settingsRout) ? 'secondary.contrastText' : 'secondary.menuText'}} primary="Settings" />
                        </ListItemButton>
                    </List>
                </Collapse>
            </Box>
        );
    };

    const getCollapsed = () => {
        return (
        <ListItemButton component={RouterLink} to={breadcrumbs.profile.path} state={{expanded}}>
            <ListItemIcon>
                <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', ml: '-3px' }} alt={user.firstName + ' ' + user.lastName} src={getImageSource()} />
            </ListItemIcon>
        </ListItemButton>);
    };

  return expanded ? getExpanded() : getCollapsed();
}
