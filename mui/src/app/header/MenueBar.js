import * as React from 'react';
import Box from '@mui/material/Box';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import MenuIcon from '@mui/icons-material/Menu';
import RefreshIcon from '@mui/icons-material/Refresh';
import Avatar from '@mui/material/Avatar';
import Tooltip from '@mui/material/Tooltip';
import MenuItem from '@mui/material/MenuItem';
import Menu from '@mui/material/Menu';
import ListItemText from '@mui/material/ListItemText';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import LogoutIcon from '@mui/icons-material/Logout';
import ListItemIcon from '@mui/material/ListItemIcon';
import { sessionManager } from "../../SessionManager";
import { AppContext } from '../../App';
import { gettStoreValue } from '../amper/Instruments';
import Convenience from '../help/Convenience';

export default function MenueBar({ expand, expanded, logOut }) {
  const app = React.useContext(AppContext);

  const getMenuIconButton = () => {
    const result = [];
    if (!expanded) {
      result.push(<IconButton
        key={'amperMenuButton'}
        size="medium"
        edge="start"
        color="inherit"
        aria-label="menu"
        sx={{ mr: 2 }}
        onClick={expand}
      >
        <MenuIcon color='secondary' fontSize="large" />
      </IconButton>)
    }
    result.push(
      <Typography key={'amperMenuText'} variant="h6" component="div" sx={{ flexGrow: 1 }}>

      </Typography>);
    return result;
  };

  const [anchorElUser, setAnchorElUser] = React.useState(null);

  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };

  const handleOpenUserMenu = (event) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleLogOut = () => {
    logOut();
  }

  const user = sessionManager.getUser();

  const getImageSource = () => {
    const photo = gettStoreValue('user_photo');
    if (Convenience.hasValue(photo)) {
      return 'data:image/png;base64,' + photo;
    }
    return '/static/images/avatar/2.jpg';
  };
  
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar style={{ minHeight: 60 }}>
          {getMenuIconButton()}
          <IconButton
            key={'amperMenuButton1'}
            size="medium"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}
            onClick={() => {app.refresh();}}
          >
            <RefreshIcon color='secondary' fontSize="large" />
          </IconButton>
          <Box sx={{ flexGrow: 0 }}>
            <Tooltip title="Profile">
              <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }} size="large">
                <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main' }} alt={user.firstName + ' ' + user.lastName} src={getImageSource()} />
              </IconButton>
            </Tooltip>
            <Menu
              sx={{ mt: '45px' }}
              id="menu-appbar"
              anchorEl={anchorElUser}
              anchorOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              keepMounted
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              open={Boolean(anchorElUser)}
              onClose={handleCloseUserMenu}
            >
              <MenuItem key={'profile'} onClick={handleCloseUserMenu}>
                {/*<ListItemIcon>
                  <AccountCircleIcon color="primary" fontSize="small" />
            </ListItemIcon>*/}
                <ListItemText>{user.firstName + ' ' + user.lastName}</ListItemText>
              </MenuItem>
              <MenuItem key={'logout'} onClick={handleLogOut}>
                <ListItemIcon>
                  <LogoutIcon color="primary" fontSize="small" />
                </ListItemIcon>
                <ListItemText>{'Logout'}</ListItemText>
              </MenuItem>
            </Menu>
          </Box>
        </Toolbar>
      </AppBar>
    </Box>
  );
}
