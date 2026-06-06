import DashboardIcon from '@mui/icons-material/Dashboard';
import SpeedIcon from '@mui/icons-material/Speed';
import SettingsIcon from '@mui/icons-material/Settings';
import DataObjectIcon from '@mui/icons-material/DataObject';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import PersonIcon from '@mui/icons-material/Person';
import SecurityIcon from '@mui/icons-material/Security';
import EditIcon from '@mui/icons-material/Edit';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import AdjustIcon from '@mui/icons-material/Adjust';
import ViewQuilt from '@mui/icons-material/ViewQuilt';
import CloudQueueIcon from '@mui/icons-material/CloudQueue';
import FolderIcon from '@mui/icons-material/Folder';
import ShareIcon from '@mui/icons-material/Share';
import TuneIcon from '@mui/icons-material/Tune';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import LocalActivityIcon from '@mui/icons-material/LocalActivity';
import EmailIcon from '@mui/icons-material/Email';
import StorageIcon from '@mui/icons-material/Storage';
import ChatIcon from '@mui/icons-material/Chat';
import AccessibilityIcon from '@mui/icons-material/Accessibility';
import InterpreterModeIcon from '@mui/icons-material/InterpreterMode';

export const breadcrumbDashboard = {
  key: 'dashboard',
  path: '/dashboard',
  label: 'Dashboard',
  icon: <DashboardIcon color={'primary'} fontSize="medium"/>,
  description: 'This page can display a set of user defined dashboards, a dashboard contains widgets that let you customize views for rendering amper data',
  primary: true,
  overview: {
    key: 'overview',
    path: '/dashboard/overview',
    label: 'Overview',
    icon: <SpeedIcon color={'primary'} fontSize="medium"/>,
    description: 'This is a defoult overview page displaying statistical and usage metric data for the current amper instance and user',
    primary: true,
  },
  add: {
    key: 'add',
    path: '/dashboard/add',
    label: 'Add',
    icon: <AddCircleOutlineIcon color={'primary'} fontSize="medium"/>,
    description: 'Use this action to define a new dashboard and configure the containing widgets to design views with rich structured data',
    primary: true,
  }
};

export const breadcrumbs = {
    resetDashboard: () => {
      for (const [key, value] of Object.entries(breadcrumbDashboard)) {
        if (value.dashboard) {
          delete breadcrumbDashboard[key];
        }
      }
    },
    addDashboard: (dashboard) => {
      breadcrumbs.dashboard[dashboard.id] = {
        key: dashboard.id,
        path: '/dashboard/item/' + dashboard.id,
        label: dashboard.label,
        icon: <ViewQuilt color={'primary'} fontSize="medium"/>,
        description: dashboard.description,
        dashboard: true,
      }
    },
    addCrumb: (path, crumbIcon) => {
      const pathParts = path.split('/');
      if (pathParts.length > 0) {
        let parent = breadcrumbs;
        for (let i = 0; i < pathParts.length; i++) {
          const currentLevel = parent[pathParts[i]];
          if (currentLevel != null && pathParts.length - 1 > i) {
            if (currentLevel[encodeURI(pathParts[i + 1])] == null) {
              currentLevel[encodeURI(pathParts[i + 1])] = {
                key: encodeURI(pathParts[i + 1]),
                path: currentLevel.path + '/' + encodeURI(pathParts[i + 1]),
                label: pathParts[i + 1],
                icon: crumbIcon,
              };
            }
            parent = currentLevel;
          }
        }
      }
    },
    alternativePaths: {},
    dashboard: breadcrumbDashboard,
    profile: {
      key: 'profile',
      path: '/profile',
      label: 'Profile',
      icon: <AccountCircleIcon/>,
      overview: {
        key: 'overview',
        path: '/profile/overview',
        label: 'Overview',
        icon: <LocalActivityIcon/>,
      },
      about: {
        key: 'about',
        path: '/profile/about',
        label: 'About',
        icon: <AccessibilityIcon/>,
      },
      settings: {
        key: 'settings',
        path: '/profile/settings',
        label: 'Settings',
        icon: <TuneIcon/>,
      },
    },
    chat: {
      key: 'chat',
      path: '/chat',
      label: 'Chat',
      icon: <ChatIcon/>,
      description: 'The page is designed for internal amper wide user chating, this would require to have the emails and credentials configured first.',
      primary: true,
    },
    email: {
      key: 'email',
      path: '/email',
      label: 'Email',
      icon: <EmailIcon/>,
      description: 'This page lets you read, write and manage emails received or sent by you, this would require to have the emails and credentials configured first.',
      primary: true,
    },
    drive: {
      key: 'drive',
      path: '/drive',
      label: 'Drive',
      icon: <CloudQueueIcon/>,
      primary: true,
      files: {
        key: 'files',
        path: '/drive/files',
        label: 'Files',
        icon: <FolderIcon/>,
        description: 'This page lets you add, manage and view all files, all content is your private ownership and can be accessed with your full permissions',
        primary: true,
      },
      shared: {
        key: 'shared',
        path: '/drive/shared',
        label: 'Shared Files',
        icon: <ShareIcon/>,
        description: 'This section lets you view and manage all files shared with you, all content is private and can be accessed only by owner provided permissions',
        primary: true,
      },
    },
    configuration: {
      key: 'configuration',
      path: '/configuration',
      label: 'Configuration',
      icon: <SettingsIcon/>,
      primary: true,
      objects: {
        key: 'objects',
        path: '/configuration/objects',
        label: 'Objects',
        icon: <DataObjectIcon/>,
        description: 'This section lets you add, manage and view all abstract data types defined for the current instance, which would require permissions to perform the actions',
        primary: true,
        edit: {
          key: 'edit',
          path: '/configuration/objects/edit',
          label: 'Manage',
          icon: <EditIcon/>,
          primary: true,
        }
      },
    },
    administration: {
      key: 'administration',
      path: '/administration',
      label: 'Administration',
      icon: <AdminPanelSettingsIcon/>,
      primary: true,
      settings: {
        key: 'settings',
        path: '/administration/settings',
        label: 'Settings',
        icon: <TuneIcon/>,
        description: 'Use settings panel to adjust global configuratuion for the amper application such as defining the data location on disc, setting license keys and etc.',
        primary: true,
      },
      relationship: {
        key: 'relationship',
        path: '/administration/relationship',
        label: 'Relationship',
        icon: <InterpreterModeIcon/>,
        description: 'This setting interface lets you define relationship between employees, use it to set the managers of the employees',
        primary: true,
      },
      users: {
        key: 'users',
        path: '/administration/users',
        label: 'Users',
        icon: <PersonIcon/>,
        description: 'The page is responsable for managing the users of the application, you should have permission to add, remove or modify an existing user',
        primary: true,
        new: {
          key: 'new',
          path: '/administration/users/new',
          label: 'New user',
          icon: <PersonAddIcon/>,
          description: 'The page is responsable for add a new user to the amper application, this action requires permission to manage amper users',
          primary: true,
        },
      },
      profiles: {
        key: 'profiles',
        path: '/administration/profiles',
        label: 'Profiles',
        icon: <SecurityIcon/>,
        description: 'This view is designed to adjust the permissions for the application users, each user is assigned a profile which describes the user permissions',
        primary: true,
      },
      nodes: {
        key: 'nodes',
        path: '/administration/nodes',
        label: 'Nodes',
        icon: <StorageIcon/>,
        description: 'This sections is for managing the distributed nodes of Amper instances, use it to add new nodes for supporting amper data scale',
        primarty: true,
        new: {
          key: 'new',
          path: '/administration/nodes/new',
          label: 'New instance',
          icon: <AddCircleOutlineIcon/>,
          description: 'The page is responsable for adding a new node instance to the amper application, this action requires permission to manage amper nodes',
          primary: true,
        },
        update: {
          key: 'update',
          path: '/administration/nodes/update',
          label: 'Update instance',
          icon: <AdjustIcon/>,
          description: 'The page is responsable for updating an existing node instance, this action requires permission to manage amper nodes',
          primary: true,
        },
      },
      chat: {
        key: 'chat',
        path: '/administration/chat',
        label: 'Chat',
        icon: <ChatIcon/>,
        description: 'The panel is designed to let configure chat channel groups, channels and assign users to channels',
        primarty: true,
        newChannel: {
          key: 'newChannel',
          path: '/administration/chat/newChannel',
          label: 'New channel',
          icon: <AddCircleOutlineIcon/>,
          description: 'The dialog is responsible for adding new channel, this action requires permission to manage amper nodes',
          primary: true,
        },
        newGroup: {
          key: 'newGroup',
          path: '/administration/chat/newGroup',
          label: 'New group',
          icon: <AddCircleOutlineIcon/>,
          description: 'The dialog is responsible for adding new channel groups, this action requires permission to manage amper nodes',
          primary: true,
        },
      }
    },
  };