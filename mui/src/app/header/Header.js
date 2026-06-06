import React, { forwardRef, useImperativeHandle, useState } from 'react';
import AmperBar from './AmperBar';
import MenueBar from './MenueBar';
import Box from '@mui/material/Box';
import Snackbar from '@mui/material/Snackbar';
import MuiAlert from '@mui/material/Alert';
import Slide from '@mui/material/Slide';

function Header({hooks, props}, ref) {
  const {expand, collapse, logOut} = hooks;
  const {expanded} = props;

  const [state, setState] = useState({
    snackBarOpen: false,
  });

  useImperativeHandle(ref, () => ({
    toast(severity, message) {
      handleSnackBarOpen(severity, message);
    },
}));

  const handleSnackbarClose = (event, reason) => {
    //ignore clickaway calls, sincethis causes issues with double alerts
    if (reason !== 'clickaway') {
      setState({
        ...state,
        snackBarOpen: false,
        snackBarMessage: '',
      });
    }
  };

  const TransitionUp = (props) => {
    return <Slide {...props} direction="up" />;
  };

  const handleSnackBarOpen = (severity, message) => {
    setState({
      ...state,
      snackBarOpen: true,
      snackBarMessage: message,
      snackBarSeverity: severity,
    });
  };
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'row',
        bgcolor: 'background.paper',
        }}>
        <AmperBar expand={expand} collapse={collapse} expanded={expanded}/>
        <MenueBar expand={expand} expanded={expanded} logOut={logOut}/>
        <Snackbar
              onClose={handleSnackbarClose}
              autoHideDuration={6000}
              anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
              open={state.snackBarOpen}
              TransitionComponent={TransitionUp}
              message={state.snackBarMessage}
            >
              <MuiAlert onClose={handleSnackbarClose} severity={state.snackBarSeverity} variant="filled" sx={{ maxWidth: 500 }}>
                {state.snackBarMessage}
              </MuiAlert>
          </Snackbar>
    </Box>
  );
}

export default forwardRef(Header)