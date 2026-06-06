import React from 'react';
import { render } from 'react-dom';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Grid from '@mui/material/Grid2';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import InboxIcon from '@mui/icons-material/MoveToInbox';
import DraftsIcon from '@mui/icons-material/Drafts';
import SendIcon from '@mui/icons-material/Send';
import DeleteIcon from '@mui/icons-material/Delete';
import { useLocation, useNavigate } from 'react-router-dom'
import { AppContext } from '../../../../App';
import { post } from '../../../data/Submit';
import HostManager from '../../../../HostManager';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Convenience from '../../../help/Convenience';
import Checkbox from '@mui/material/Checkbox';
import Typography from '@mui/material/Typography';
import makeStyles from '@mui/styles/makeStyles';
import useTheme from "@mui/material/styles/useTheme";
import Stack from '@mui/material/Stack';
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import EmailDetail from './EmailDetail';
import { getFrom, getDate, getSubject, isSeen } from './EmailHelper';
import Pagination from '@mui/material/Pagination';
import LinearProgress from '@mui/material/LinearProgress';
import { sessionManager } from '../../../../SessionManager';
import AttachmentIcon from '@mui/icons-material/Attachment';
import Chip from '@mui/material/Chip';
import CreateIcon from '@mui/icons-material/Create';
import Button from '@mui/material/Button';
import EmailCompose from './EmailCompose';

const emailItemStyles = makeStyles(theme => {
  return {
    box: {
        backgroundColor: theme.palette.secondary.main,
        '&:hover': {
          backgroundColor: theme.palette.primary.main,
          color: theme.palette.secondary.main,
        }
    },
    typography: {
      color: 'inherit',
    }
  }
});

export default function EmailsRoot({id, expanded}) {//mvenrokssylhqyhn
  const app = React.useContext(AppContext);
  const theme = useTheme();

  const location = useLocation();
  const {pathname} = location;
  const navigate = useNavigate();
  const pathparts = pathname.split('/');

  const [state, setState] = React.useState({
    checked: true,
    loaded: true,
    mailboxesLoaded: false,
    start: 0,
    limit: 50,
    search: null,
    data: [],
    selected: {},
    view: 'list',
    pointer: {
      pages: {

      }
    },
  });

  const composerRef = React.useRef();

  if (app) {
    app.registerRefresh(id, () => {
        reload();
    });    
}

  const getEmail = () => {
    return pathparts[2];
  };

  const getBox = () => {
    return state.selectedMailbox;
  };

  const checkHandler = (result) => {
    if (!result.success) {
          app.toast('warning', result.error)
    }
    setState({
      ...state,
      checked: true,
    });
  };

  const mailboxesHandler = (result) => {
    if (!result.success) {
      app.toast('warning', result.error)
    }
    setState({
      ...state,
      mailboxesLoaded: true,
      mailboxes: result.data,
      selectedMailbox: result.data[0].label,
      loaded: false,
      checked: false,
    });
  };

  React.useEffect(() => {
    if (!state.mailboxesLoaded) {
        post(`${HostManager.myHost()}email/mailboxes`, {
          email: getEmail(),
        }, mailboxesHandler, mailboxesHandler);
    }
  }, [state.mailboxesLoaded]);

  React.useEffect(() => {
    if (!state.checked) {
        post(`${HostManager.myHost()}email/check`, {
          email: getEmail(),
        }, checkHandler, checkHandler);
    }
  }, [state.checked]);

  React.useEffect(() => {
    if (!state.loaded) {
      load();
    }
  }, [state.loaded]);

  const reload = () => {
    setState({
      ...state,
      checked: false,
    })
  };

  const load = () => {
    post(`${HostManager.myHost()}email/fetch`, {
      email: getEmail(),
      box: getBox(),
      start: state.start,
      limit: state.limit,
      search: state.search,
      pointer: state.pointer,
    }, (result) => {
        if (result.success) {
            let pointer = null;
            if (state.totalCount == result.count) {
              const lastId = (result.data != null && result.data.length > 0) ? result.data[result.data.length - 1].id : null;
              pointer = {
                pages: {
                  ...state.pointer.pages,
                [state.limit/(state.limit - state.start)]: lastId,
                }
              };
            } else {
              pointer = {
                pages: {
                },
              };
            }
            setState({
                ...state,
                loaded: true,
                data: result.data || [],
                totalCount: result.count,
                pointer,
            });
        } else {
            app.toast('warning', result.error);
        }
    }, (result) => {
        if (result) {
            app.toast('warning', result.error);
        }
    });
  };

  const getBasePath = () => {
    const pathparts = pathname.split('/');
    return '/email/' + pathparts[2];
  };

  const handleMailboxClick = (mailbox) => {
    setState({
        ...state,
        selectedMailbox: mailbox,
        loaded: false,
        view: 'list',
      });
  };

  const emailItemClasses = emailItemStyles();

  const getEmailLeftMenu = () => {
    const user = sessionManager.getUser();
    const mailboxesList = [];
    if (state.mailboxesLoaded) {
      const mailboxes = state.mailboxes;
      for (let m = 0; m < mailboxes.length; m++) {
        const label = mailboxes[m].label;
        mailboxesList.push(
          <ListItemButton onClick={() => {handleMailboxClick(label)}}>
            <ListItemIcon>
              <InboxIcon />
            </ListItemIcon>
            <ListItemText primary={label} sx={{ml: -2}}/>
          </ListItemButton>
        );
      }
    }
    return <List
        sx={{ width: '100%', maxHeight: 'calc(100% - 20px)', overflowX: 'auto', overflowX: 'hidden', bgcolor: 'background.paper', ml: -1}}
      /*subheader={
          <ListSubheader component="div">
            Nested List Items
        </ListSubheader>
        }*/
      >
        {mailboxesList}
      </List>;
  };

  const selectAll = (event) => {
    const selected = {};
    if (event.target.checked) {
      for (let i = 0; i < state.data.length; i++) {
        selected[state.data[i].id] = state.data[i];
      }
    }
    setState({
      ...state,
      selected
    });
  };

  const isSelected = (email) => {
    return state.selected[email.id] != null;
  };

  const onEmailClick = (event, email) => {
    if (state.selectedMailbox === 'Drafts') {
      setState({
        ...state,
        view: 'compose',
        draft: email,
      });
    } else {
      setState({
        ...state,
        view: 'detail',
        email: email,
      });
    }
  };

  const onBackClick = () => {
    setState({
      ...state,
      view: 'list',
      email: null,
      loaded: false,
    });
  };

  const onEmailSelect = (event, email) => {
    if (state.selected[email.id] == null) {
      setState({
        ...state,
        selected: {
          ...state.selected,
          [email.id]: email,
        }
      });
    } else {
      setState({
        ...state,
        selected: {
          ...state.selected,
          [email.id]: null,
        }
      });
    }
  };

  const hasSelection = () => {
    for (let value in state.selected) {
      if (state.selected[value] != null) {
        return true;
      }
    }
    return false;
  };

  const moveToTrash = () => {
    const selection = [];
    if (state.view == 'list') {
      for (let value in state.selected) {
        if (state.selected[value] != null) {
          selection.push(state.selected[value]);
        }
      } 
    } else if (state.view == 'detail') {
      selection.push(state.email);
    }
    if (selection.length > 0) {
      post(`${HostManager.myHost()}email/move`, {
        emails: selection,
        from: getBox(),
        to: 'trash',
      }, (result) => {
          if (result.success) {
              setState({
                  ...state,
                  selected: {},
                  loaded: false,
                  view: 'list',
                  email: null,
              });
          } else {
              app.toast('warning', result.error)
          }
      }, (result) => {
          if (result) {
              app.toast('warning', result.error)
          }
      });
    }
  };

  const onPageChange = (event, page) => {
    setState({
      ...state,
      start: ((page - 1) * 50),
      limit: ((page - 1) * 50 + 50),
      loaded: false,
      data: [],
    })
  };

  const getMailboxEmpty = () => {
    return <Box sx={{textAlign: 'center', padding: '200px 0'}}>
      <Typography variant="h5" gutterBottom>
        {decodeURI(getBox())} empty
      </Typography>
    </Box>;
  };

  let dragElement = null;
  const onEmailDragStart = (event, email) => {
    if (dragElement) {
      dragElement.remove();
    }
    if (hasSelection()) {

    } else {
      dragElement = getDragEmailsElement([email]);
      
      document.body.appendChild(dragElement);
    }
    event.dataTransfer.setDragImage(dragElement, 0, 0)
  };

  const send = () => {
    setState({
      ...state,
      sending: true,
    })
  };

  const sentCallback = () => {
    setState({
      ...state,
      view: 'list',
      sending: false,
    })
  };

  const getDragEmailsElement = (emails) => {
    const dragElement = document.createElement("div");
    dragElement.style.position = 'relative';
    dragElement.style.borderRadius = '25px';
    for (let i = 0; i < 3; i++) {
      if (i < emails.length) {
        const from = getFrom(emails[i]);
        const subject = getSubject(emails[i]);
        const dragItem = document.createElement("div");
        dragItem.style.position = 'absolute';
        dragItem.style.verticalAlign = 'middle';
        dragItem.style.textAlign = 'center';
        dragItem.style.border = 'solid 1px'
        dragItem.style.width = '400px';
        dragItem.style.borderRadius = '25px';
        dragItem.style.height = '50px';
        dragItem.style.backgroundColor = 'white';
        dragItem.style.color = theme.palette.primary.main;
        dragItem.innerText = subject;
        dragElement.appendChild(dragItem);
      }
    }
    return dragElement;
  };

  const getEmailList = () => {
    return <TableContainer component={Paper} sx={{ml: 1, mt: 1, mb: 1,}}>
      <Table width="100%" style={{'table-layout': 'fixed', cursor: 'pointer'}}>
        <TableBody>
            {state.data.map((email) => {
              const selected = isSelected(email);
              const from = getFrom(email);
              const subject = getSubject(email);
              const date = getDate(email);
              const seen = isSeen(email);
              const hasAttachment = email.attachments != null ? email.attachments.length > 0 : false;
              let attachmentLable = '';
              let attachementIconTabelCell = null
              
              if (hasAttachment) {
                for (let i = 0; i < email.attachments.length; i++) {
                  attachmentLable+=((i > 0 ? ', ' : '') + email.attachments[i].Name)
                }
                attachementIconTabelCell = <Tooltip title={<Typography variant="body1" sx={{fontWeight: 'inherit'}}>{attachmentLable}</Typography>}><Chip icon={<AttachmentIcon/>} label={attachmentLable} /></Tooltip>;
              }
              return <TableRow
                key={email.id}
                onDragStart={(event) => onEmailDragStart(event, email)}
                draggable="true"
                sx={{ '&:last-child td, &:last-child th': { border: 0 }, fontWeight: seen ? 'normal' : 'bold' }} className={emailItemClasses.box} style={selected ? {backgroundColor: theme.palette.primary.main, color: theme.palette.secondary.main} : {}}>
                <TableCell padding="checkbox" className={emailItemClasses.typography} sx={{fontWeight: 'inherit'}}>
                  <Checkbox
                    color={selected ? "secondary" : "secondary"}
                    className={emailItemClasses.typography}
                    onClick={(event)=>onEmailSelect(event, email)}
                    checked={selected}
                  />
                </TableCell>
                <TableCell onClick={(event)=>onEmailClick(event, email)} align="left" width="130px" sx={{fontWeight: 'inherit', 'white-space': 'nowrap', 'overflow': 'hidden', 'text-overflow': 'ellipsis'}} className={emailItemClasses.typography}>
                  <Tooltip title={<Typography variant="body1" sx={{fontWeight: 'inherit'}} className={emailItemClasses.typography}>{from}</Typography>}><Typography variant="body1" className={emailItemClasses.typography} sx={{fontWeight: 'inherit'}}>{from}</Typography></Tooltip>
                </TableCell>
                <TableCell onClick={(event)=>onEmailClick(event, email)} align="left" sx={{fontWeight: 'inherit', 'white-space': 'nowrap', 'overflow': 'hidden', 'text-overflow': 'ellipsis'}} className={emailItemClasses.typography}><Tooltip title={<Typography variant="body1" sx={{fontWeight: 'inherit'}}>{subject}</Typography>}><Typography variant="body1" sx={{fontWeight: 'inherit'}}>{subject}</Typography></Tooltip>{attachementIconTabelCell}</TableCell>
                <TableCell onClick={(event)=>onEmailClick(event, email)} align="left" width="150px" sx={{fontWeight: 'inherit', 'white-space': 'nowrap', 'overflow': 'hidden', 'text-overflow': 'ellipsis'}} className={emailItemClasses.typography}><Tooltip title={<Typography variant="body1" sx={{fontWeight: 'inherit'}}>{date}</Typography>}><Typography variant="body1" sx={{fontWeight: 'inherit'}}>{date}</Typography></Tooltip></TableCell>
              </TableRow>
            })}
          </TableBody>
      </Table>
    </TableContainer>;
  };


  const getView = () => {
    if (state.view == 'detail') {
      return <EmailDetail email={state.email} box={getBox()}></EmailDetail>
    } else if (state.view == 'compose') {
      return <EmailCompose ref={composerRef} email={getEmail()} draft={state.draft} sentCallback={sentCallback} sending={state.sending}></EmailCompose>
    } else {
      if (state.data.length > 0) {
        return getEmailList();
      } else if (state.loaded != false) {
        return getMailboxEmpty();
      }
    }
  };

  const getPagingBar = () => {
    return <Pagination onChange={onPageChange} count={(Math.ceil(state.totalCount/50))} variant="outlined" color="primary" shape="rounded" siblingCount={0}/>;
  };

  const composeEmail = () => {
    setState({
      ...state,
      view: 'compose',
      draft: null
    });
  };

  const getEmailTopMenu = () => {
    const leftResult = [];
    const rightResult = [];
    if (state.view == 'detail' || state.view == 'compose') {
      leftResult.push(<IconButton sx={{}} onClick={onBackClick} aria-label="Back">
        <ArrowBackIcon/>
      </IconButton>);
    } else {
      leftResult.push(<Checkbox sx={{'& .MuiSvgIcon-root': { fontSize: 28 }}} onClick={selectAll}/>);
      rightResult.push(getPagingBar());
    }

    if (state.view == 'compose') {
      rightResult.push(<Button variant="contained" sx={{mr: 1}} onClick={send} disabled={state.sending} endIcon={<SendIcon />}>
        Send
      </Button>);
    }
    return [<Box sx={{flexGrow: 1}}>
      {leftResult}
      </Box>, <Box sx={{flexGrow: 0}}>
        {rightResult}
      </Box>];
  };

  const getLeftTopMenu = () => {
    if (state.view != 'compose') {
      return <Button variant="contained" onClick={composeEmail} startIcon={<CreateIcon />}>
        Compose
      </Button>;
    }
  };

  return (
    <Box sx={{ width: 'calc(100% - 25px)', height: '100%' }}>
      <Grid container columns={48} spacing={0} sx={{height: '100%'}}>
        <Grid size={8} sx={{flexGrow: 1, height: '40px'}}>
          <Box sx={{ height: '50px'}}>
            {getLeftTopMenu()}
          </Box>
        </Grid>
        <Grid size={40} sx={{flexGrow: 1, height: '40px'}}>
          <Stack direction="row" spacing={2} sx={{pl: '10px'}}>
            {getEmailTopMenu()}
          </Stack>
        </Grid>
        <Grid size={8} sx={{flexGrow: 1, height: 'calc(100% - 40px)'}}>
          <Box sx={{ height: '100%'}}>
            {state.mailboxesLoaded == false ? <LinearProgress sx={{mr: 1}}/> : ''}
            {getEmailLeftMenu()}
          </Box>
        </Grid>
        <Grid size={40} sx={{flexGrow: 1, height: 'calc(100% - 40px)'}}>
            <Box sx={{ height: '100%', overflowX: 'hidden', overflowY: 'auto'}}>
              {state.loaded == false || state.checked == false ? <LinearProgress /> : ''}
              {getView()}
            </Box>
        </Grid>
      </Grid>
    </Box>
    );
}
