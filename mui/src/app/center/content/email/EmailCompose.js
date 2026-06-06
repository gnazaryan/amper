import React, { forwardRef, useImperativeHandle } from 'react';
import { AppContext } from '../../../../App';
import Paper from '@mui/material/Paper';
import MUIRichTextEditor from "mui-rte";
import { convertFromHTML, convertToRaw } from 'draft-js';
import { ContentState, EditorState } from 'draft-js';
import { stateToHTML } from 'draft-js-export-html';
import {stateFromHTML} from 'draft-js-import-html';
import TextField from '@mui/material/TextField';
import Stack from '@mui/material/Stack';
import Link from '@mui/material/Link';
import Box from '@mui/material/Box';
import { debounce } from '@mui/material/utils';
import { post } from '../../../data/Submit';
import HostManager from '../../../../HostManager';
import { getAddress } from './EmailHelper';
import Convenience from '../../../help/Convenience';
import EmailEditor from './EmailEditor'
import LinearProgress from '@mui/material/LinearProgress';
import FilesRoot from '../../../components/drive/FilesRoot';

import "./styles.css";

function EmailCompose({email, draft, sentCallback, sending}, ref) {

    const app = React.useContext(AppContext);

    const [id, setId] = React.useState(draft != null ? draft.id : null);
    let year = -1;
    let month = -1;
    if (id != null) {
        const idParts = id.split('_')
        const secs = parseInt(idParts[0]);
        const date = new Date(secs * 1000);
        year = date.getFullYear();
        month = date.getMonth() + 1;
    }
    const to = draft != null && draft.envelope != null ? getAddress(draft.envelope.to) : '';
    const cc = draft != null && draft.envelope != null ? getAddress(draft.envelope.cc) : '';
    const bcc = draft != null && draft.envelope != null ? getAddress(draft.envelope.bcc) : '';
    const subject = draft != null && draft.envelope != null ? draft.envelope.subject : '';
    const content = draft != null ? draft.bodyHtml : '';
    const [state, setState] = React.useState({
        cc: Convenience.hasValue(cc),
        bcc: Convenience.hasValue(bcc),
        email: {
            from: email,
            id: id,
            to: to,
            cc: cc,
            bcc: bcc,
            subject: subject,
            content: content,
        },
        showAttachments: true,
      });

      React.useEffect(() => {
        if (sending) {
            post(`${HostManager.myHost()}email/send`,
            {
                ...state.email,
                id: id,
            },
            sendHandler, sendHandler);
        }
      }, [sending]);

      const sendHandler = (result) => {
        if (result.success) {
            app.toast('info', 'Email sent')
            if (sentCallback) {
                sentCallback();
            }
        } else {
            app.toast('warning', result.error)
        }
      };

      useImperativeHandle(ref, () => ({
        send() {
            /**/
        }
      }));


      const saveDraftDebounceTrigger = (newState, id) => {
        save(newState, id);
      };
      const saveDraftHandler = React.useCallback(debounce(saveDraftDebounceTrigger, 2000), []);

      React.useEffect(() => {
        saveDraftHandler(state, id);
      }, [state.email.to, state.email.cc, state.email.bcc, state.email.subject, state.email.content]);
      

      const draftSaveHandler = (result) => {
        if (result.success) {
            app.toast('info', 'Draft saved')
            setId(result.id);
        } else {
            app.toast('warning', result.error)
        }
      };

    const save = (newState, newId) => {
        post(`${HostManager.myHost()}email/saveDraft`,
        {
            ...newState.email,
            id: newId,
        },
        draftSaveHandler, draftSaveHandler);
    };

    const onCcChange = (event) => {
        const {name, value} = event.target;
        const newState = {
            ...state,
            email: {
                ...state.email,
                cc: value,
            },
        };
        setState(newState);
    };

    const getCC = () => {
        if (state.cc) {
            return <TextField id="outlined-basic" label="Cc" variant="filled" fullWidth onChange={onCcChange} value={state.email.cc}/>;
        }
    };
    
    const onBccChange = (event) => {
        const {name, value} = event.target;
        const newState = {
            ...state,
            email: {
                ...state.email,
                bcc: value,
            },
        };
        setState(newState);
    };

    const getBCC = () => {
        if (state.bcc) {
            return <TextField id="outlined-basic" label="Bcc" variant="filled" fullWidth onChange={onBccChange} value={state.email.bcc}/>;
        }
    };

    const activateCC = () => {
        setState({
            ...state,
            cc: true,
        });
    };

    const activateBCC = () => {
        setState({
            ...state,
            bcc: true,
        });
    };

    const onToChange = (event) => {
        const {name, value} = event.target;
        const newState = {
            ...state,
            email: {
                ...state.email,
                to: value,
            },
        };
        setState(newState);
    };

    const onSubjectChange = (event) => {
        const {name, value} = event.target;
        const newState = {
            ...state,
            email: {
                ...state.email,
                subject: value,
            },
        };
        setState(newState);
    };

    const onContentChange = (contentHtml) => {
        const newState = {
            ...state,
            email: {
                ...state.email,
                content: contentHtml,
            },
        }
        setState(newState);
    };

    const getCCBCCActivation = () => {
        const result = [];
        if(!state.cc) {
            result.push(<Link onClick={activateCC} color="inherit" sx={{cursor: 'pointer'}}>cc</Link>);
        }
        if(!state.bcc) {
            result.push(<Link onClick={activateBCC} color="inherit" sx={{cursor: 'pointer'}}>bcc</Link>);
        }
        if (result.length > 0) {
            return <Stack spacing={1} sx={{pr: 1}}>{result}</Stack>
        }
    };

    
    const blocksFromHTML = convertFromHTML(content);
    const contentState = ContentState.createFromBlockArray(
        blocksFromHTML.contentBlocks,
        blocksFromHTML.entityMap,
        );
    const rawContentState = convertToRaw(contentState);

    return <div style={{width: 'calc(100% - 20px)', height: 'calc(100% - 20px)',}}><Paper sx={{width: 'calc(100% - 20px)', height: 'calc(100% - 250px)', mt: '5px', ml: '5px', mr: '5px'}} elevation={3}>
        {sending ? <LinearProgress /> : ''}
        <Stack
            direction="row"
            spacing={1}>
            <TextField id="outlined-basic" label="To" variant="filled" fullWidth onChange={onToChange} value={state.email.to}/>
            {getCC()}
            {getBCC()}
            {getCCBCCActivation()}
        </Stack>
        <TextField id="outlined-basic" label="Subject" variant="standard" fullWidth onChange={onSubjectChange} value={state.email.subject}/>
        <EmailEditor content={content} onContentChange={onContentChange}></EmailEditor>
    </Paper>
    <Paper sx={{width: 'calc(100% - 20px)', height: '250px', m: '5px'}} elevation={3}>
        {state.showAttachments ? <FilesRoot root={'/__system__/Email/' + email + '/Drafts/success/' + year + '/' + month + '/' + id} viewLevel={0}></FilesRoot> : ''}
    </Paper>
    </div>
}

export default forwardRef(EmailCompose)