import AmperIcon from '../icons/Amper';
import Box from '@mui/material/Box';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import Typography from '@mui/material/Typography';
import Collapse from '@mui/material/Collapse';
import AmperConstatns from '../util/AmperConstants';

export default function AmperBar({expand, collapse, expanded}) {

  const getMenuIconButton = () => {
    const result = [];
    if (expanded) {
      result.push(
        <Typography key={'amperMenuText'} variant="h6" component="div" sx={{ ml: 1, mb: -1, fontSize: '25px', flexGrow: 1 }}>
          {expanded ? 'amper' : ''}
        </Typography>);
      result.push(<IconButton
            key={'amperMenuButton'}
            size="medium"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: -3.5}}
            onClick={collapse}
          >
            <MenuIcon color='secondary' fontSize="large"/>
        </IconButton>)
    }
    return result;
  };

  return (
    <Box >
      <Collapse orientation="horizontal" in={expanded} collapsedSize={64}>
      <Box sx={{ borderRadius: (theme) => (theme.palette.primary.borderRadius), width: expanded ? AmperConstatns.LEFT_MENU_WIDTH + 4 : 64}}> 
        <AppBar position="static">
          <Toolbar style={{minHeight:60}}>
              <IconButton
                size="medium"
                edge="start"
                color="inherit"
                aria-label="menu"
                sx={{ }}>
                  <AmperIcon color='secondary' sx={{m: -1}} style={{fontSize: '45px'}}/>
            </IconButton>
            {getMenuIconButton()}
          </Toolbar>
        </AppBar>
      </Box>
      </Collapse>
    </Box>
  );
}
