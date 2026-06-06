import * as React from 'react';
import { useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom'
import Box from '@mui/material/Box';
import Header from './app/header/Header';
import CenterPanel from './app/center/CenterPanel';
import Login from './app/login/Login';
import { useLocation } from 'react-router-dom';
import { sessionManager } from "./SessionManager";
import Convenience from "./app/help/Convenience";
import FilesBar from './app/components/drive/FilesBar';
import Activation from './app/activation/Activation';
import { setStoreValue } from './app/amper/Instruments';
import { post } from './app/data/Submit';
import HostManager from './HostManager';

export const AppContext = React.createContext();

function App() {

  const navigate = useNavigate();
  const expanded = sessionManager.isExpanded();
  const initialState = {
      expanded: expanded == true,
      sessionId: sessionManager.getSessionId(),
      snackBarOpen: false,
      snackBarMessage: '',
      snackBarSeverity: 'success',
  };
  const [state, setState] = useState(initialState);
  
  let lastUpdateTime = Date.now();
  useEffect(() => {
    if (state.sessionId) {
      setTimeout(fetchUpdates, 2000);
    }
  }, [state.sessionId]);

  const headerRef = React.useRef();
  const filesRef = React.useRef();
  const expand = () => {
    sessionManager.setSetting('expanded', true);
    setState({
      ...state,
      expanded: true
    })
  };

  const collapse = () => {
    sessionManager.setSetting('expanded', false);
    setState({
      ...state,
      expanded: false
    })
  };

  const success = (user, settings) => {
    setStoreValue("user_photo", user.photo);
    delete user.photo;
    sessionManager.setUser(user);
    sessionManager.setSettings(settings);
    setState({
        ...state,
        sessionId: user.sessionId,
    });
  };

  const activationSuccess = () => {
    navigate('/');
    logOut();
  };

  const logOut = () => {
    sessionManager.invalidateSession();
    setState({
        ...state,
        sessionId: null,
    });
  };

  const handleSnackBarOpen = (severity, message) => {
    setState({
      ...state,
      snackBarOpen: true,
      snackBarMessage: message,
      snackBarSeverity: severity,
    });
  };

  const refreshCallbacks = {};
  const registerRefresh = (id, callback) => {
    refreshCallbacks[id] = callback;
  };

  const refresh = () => {
    for (const [key, value] of Object.entries(refreshCallbacks)) {
      if (value != null) {
        setTimeout(value, 500);
      }
    }
  };

  const toast = (severity, message) => {
    headerRef.current.toast(severity, message);
  };

  const upload = (directory, files, callback, upversion, metadata) => {
    filesRef.current.upload(directory, files, callback, upversion, metadata);
  };

  const serverUpdateCallbacks = useMemo(() => {
    return {};
  }, [true]);
  const registerServerUpdate = (id, callback) => {
    serverUpdateCallbacks[id] = callback;
  };

  document.body.addEventListener("mousemove", (e) => {
      /*if (Date.now() - lastUpdateTime > 10000) {
        fetchUpdates();
      }*/
  });

  const fetchUpdates = () => {
    lastUpdateTime = Date.now();
    try {
      post(`${HostManager.myHost()}updates/fetch`, {
        }, (result) => {
          try {
            if (result.data != null) {
              for (const key in result.data) {
                if (serverUpdateCallbacks[key] != null) {
                  serverUpdateCallbacks[key](result.data[key]);
                }
              }
            }
          } catch ({ name, message }) {
            alert('apply update failed: ' + message)
          } finally {
            setTimeout(fetchUpdates, 2000);
          }
      }, (result) => {
          setTimeout(fetchUpdates, 2000);
          //app.toast('info', `not able to synch to server for updates`);
      });
    } catch ({ name, message }) {
      alert('fetch update failed: ' + message)
    }
  }

  

  const appContext = {
    refresh: refresh,
    registerRefresh: registerRefresh,
    toast: toast,
    upload: upload,
    logOut: logOut,
    registerServerUpdate: registerServerUpdate,
  };

  const getUI = () => {
    return <AppContext.Provider value={appContext}><Box height={'100vh'}
      color="inherit"
        sx={{
          bgcolor: (theme) => (theme.palette.primary.main),
          display: 'flex',
          flexDirection: 'column',
          bgcolor: 'background.paper',
        }}>
          <Header
            props={{
              expanded: state.expanded,
            }}
            ref={headerRef}
            hooks={{
              expand,
              collapse,
              logOut
            }}/>
            <FilesBar ref={filesRef}/>
          <CenterPanel props={{
            expanded: state.expanded,
            toast: handleSnackBarOpen,
          }}/>
      </Box></AppContext.Provider>;
  };
  
  let face = 1;
  const activationCode = Convenience.getUrlParameterValue('activationCode');
  if (activationCode != null) {
      face = 2;
  } else if (!state.sessionId) {
      face = 0;
  }
  return ((function(face) {
        switch(face) {
            case 0:
              return <Login hooks={{success}}/>;
            case 1:
              return getUI();
            case 2:
              return <Activation hooks={{success: activationSuccess}} activationCode={activationCode}></Activation>;
        }
    })(face));
};

export default App;
