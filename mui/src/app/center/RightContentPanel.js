import React, { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import { Routes, Route } from 'react-router-dom';
import Overview from './content/Overview'
import Objects from './content/objects/Objects'
import Edit from './content/objects/edit/Edit'
import Users from './content/administration/users/Users';
import Administration from './content/Administration';
import Configuration from './content/Configuration';
import Settings from './content/administration/settings/Settings';
import Dashboards from './content/dashboard/Dashboards';
import Dashboard from './content/dashboard/Dashboard';
import Profiles from './content/Profiles';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';
import { breadcrumbs } from './Breadcrambs'
import Button from '@mui/material/Button';
import { Link as RouterLink, useLocation } from 'react-router-dom'
import Drive from './content/drive/Drive';
import Files from './content/drive/Files';
import Email from './content/email/Email';
import EmailOverview from './content/email/EmailOverview';
import SharedFiles from './content/drive/SharedFiles';
import { registerResize } from '../amper/Instruments';
import Convenience from '../help/Convenience';
import AmperConstatns from '../util/AmperConstants';
import Profile from './content/profile/Profile';
import Nodes from './content/administration/nodes/Nodes';
import ChatPanel from './content/chat/ChatPanel';
import ChatConfiguration from './content/administration/chat/ChatConfiguration';
import Relationship from './content/administration/relationship/Relationship';

export default function RightContentPanel({expanded, toast}) {
  const {pathname, search} = useLocation();
  const pathParts = pathname.split('/');
  const hierarchy = [];
  const dashboard = <Dashboards expanded={expanded} toast={toast}/>;
  const [state, setState] = useState({
    height: window.innerHeight - AmperConstatns.LEFT_MENU_WIDTH + 5,
    width: expanded ? window.innerWidth - (AmperConstatns.LEFT_MENU_WIDTH + 5) : window.innerWidth - 65,
  });

  const handleResize = (height, width, expanded) => {
    setState({
      ...state,
      height: height - AmperConstatns.LEFT_MENU_WIDTH + 5,
      width: expanded ? width - (AmperConstatns.LEFT_MENU_WIDTH + 5) : width - 65,
    });
  };

  useEffect(() => {
    registerResize(handleResize, expanded);
    return function cleanupListener() {
      window.removeEventListener('resize', handleResize)
    };
  });
  
  useEffect(() => {
    handleResize(window.innerHeight, window.innerWidth, expanded);
  }, [expanded]);

  const getCrumbs = () => {
    const result = [];
    for (let i = 0; i < pathParts.length; i++) {
      const pathPart = pathParts[i];
      if (Convenience.hasValue(pathPart)) {
          let current = null;
          if (hierarchy.length > 0) {
            current = hierarchy[hierarchy.length - 1][pathPart];
          } else if (hierarchy.length == 0) {
            let alternativePath = pathPart;
            if (breadcrumbs.alternativePaths[pathPart] != null) {
              alternativePath = breadcrumbs.alternativePaths[pathPart];
            }
            current = breadcrumbs[alternativePath];
          } else {
            //No matching breadcrumb
            continue;
          }
          //For dashboards and any other view having dynamic url with ids
          //consider the item to be the parent of the dynamic content id
          //skip current path and show the crumb based on id
          if (current == null && pathPart == 'item' &&
           pathParts.length -1 > i && hierarchy.length > 0 && hierarchy[hierarchy.length - 1][pathParts[i + 1]] != null) {
            continue;
          } else if (current == null && pathPart == 'item') {
            break;
          }

          if (current == null && hierarchy.length > 0) {
            current = {
              key: pathPart,
              path: hierarchy[hierarchy.length - 1].path + '/' + pathPart,
              label: decodeURI(pathPart),
              icon: hierarchy[hierarchy.length - 1].icon,
            };
          }
          hierarchy.push(current);
          let path = current.path;
          result.push(hierarchy.length > 1 ? <NavigateNextIcon key={path} sx={{ ml: -1, mr: -1, mt: '15px', height: '20px' }} color={'primary'} /> : '');
          if (i == (pathParts.length - 1) && Convenience.hasValue(search)) {
            path = path + search;
          }
          result.push(<Button key={path} sx={{ mt: 1}} component={RouterLink} variant="outlined" to={path} state={{loading: true}} search={search} size={'small'} startIcon={ current.icon }>
            {current.label}
          </Button>);
      }
    }
    return result;
  };

  return (
          <Box height={'calc(100% - 2px)'} width={'calc(100% - 4px)'}
            sx={{
              display: 'flex',
              flexDirection: 'column',
              bgcolor: 'background.paper',
              mr: '2px',
              ml: '1px',
              borderRadius: (theme) => (theme.palette.primary.borderRadius)}}>
            <Box sx={{display: 'flex'}}>
              <Breadcrumbs width={'100%'} sx={{ ml: 2, mr: 2, mt: 0, mb: 1}} maxItems={11} aria-label="breadcrumb" separator=''>
                {getCrumbs()}
              </Breadcrumbs>
            </Box>
            <Box sx={{display: 'flex', flexGrow: 1, height: state.height, width: state.width , mr: 2, mt: 0, ml: 2, mb: 2 }} >
            <Routes sx={{}}>
                  <Route exact path={breadcrumbs.profile.path + '/*'} element={<Profile expanded={expanded}/>}></Route>
                  
                  <Route exact path={breadcrumbs.dashboard.path} element={dashboard}></Route>
                  <Route exact path={breadcrumbs.dashboard.add.path} element={dashboard}></Route>
                  <Route exact path={breadcrumbs.dashboard.overview.path} element={<Overview expanded={expanded}/>}></Route>
                  <Route exact path={'/dashboard/item/*'} element={<Dashboard expanded={expanded}/>}></Route>

                  <Route exact path={breadcrumbs.chat.path} element={<ChatPanel expanded={expanded}/>}></Route>

                  <Route exact path={breadcrumbs.email.path} element={<EmailOverview expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.email.path+ '/*'} element={<Email expanded={expanded}/>}></Route>

                  <Route exact path={breadcrumbs.drive.path} element={<Drive expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.drive.files.path + '/*'} element={<Files expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.drive.shared.path} element={<SharedFiles expanded={expanded}/>}></Route>


                  <Route exact path={breadcrumbs.configuration.path} element={<Configuration expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.configuration.objects.path} element={<Objects expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.configuration.objects.edit.path} element={<Edit expanded={expanded}/>}></Route>

                  
                  <Route exact path={breadcrumbs.administration.path} element={<Administration expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.administration.settings.path} element={<Settings expanded={expanded}/>}></Route>

                  <Route exact path={breadcrumbs.administration.users.path} element={<Users expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.administration.users.new.path} element={<Users expanded={expanded}/>}></Route>

                  <Route exact path={breadcrumbs.administration.profiles.path} element={<Profiles expanded={expanded}/>}></Route>

                  
                  <Route exact path={breadcrumbs.administration.nodes.path} element={<Nodes expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.administration.nodes.new.path} element={<Nodes expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.administration.nodes.update.path} element={<Nodes expanded={expanded}/>}></Route>

                  <Route exact path={breadcrumbs.administration.chat.path} element={<ChatConfiguration expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.administration.chat.newGroup.path} element={<ChatConfiguration expanded={expanded}/>}></Route>
                  <Route exact path={breadcrumbs.administration.chat.newChannel.path} element={<ChatConfiguration expanded={expanded}/>}></Route>

                  <Route exact path={breadcrumbs.administration.relationship.path} element={<Relationship expanded={expanded}/>}></Route>
            </Routes>
            </Box>
          </Box>
    );
}
