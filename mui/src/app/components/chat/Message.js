import React, { useState, useRef } from 'react';
import Box from '@mui/material/Box';
import Avatar from '@mui/material/Avatar';
import CircularProgress from '@mui/material/CircularProgress';
import Convenience from '../../help/Convenience';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import Tooltip from '@mui/material/Tooltip';
import './Message.css'
import Typography from '@mui/material/Typography';
import { formatDate } from '../../amper/Instruments';
import ReplyIcon from '@mui/icons-material/Reply';
import IconButton from '@mui/material/IconButton';
import MoodIcon from '@mui/icons-material/Mood';
import EmotionPicker from './emotion/EmotionPicker';
import { sessionManager } from '../../../SessionManager';
import Chip from '@mui/material/Chip';
import { AppContext } from '../../../App';
import { post } from '../../data/Submit';
import HostManager from '../../../HostManager';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ChatEditor from './ChatEditor';
import Attachment from './Attachment';
import { download } from '../../data/Submit';

export default function Message({host, user, from, to, participants, batchId, batchId1, message, left, onReply, messageBatchId, messageType, onMessageUpdate, onMessageRemove, onMessageReaction, replyMessage}) {
    const app = React.useContext(AppContext);
    const reactionButtonRef = useRef(null);
    const replyButtonRef = useRef(null);
    const moreButtonRef = useRef(null);

    const [state, setState] = useState({
        reactionsEl: null,
        moreMenueEl: null,
        mode: 'view',
    });

    const messageTime = formatDate(message.dateTime);

    const getImageSource = (user) => {
        if (Convenience.hasValue(user.photo)) {
          return 'data:image/png;base64,' + user.photo;
        }
        return '/static/images/avatar/2.jpg';
      };

    const getMessageTime = () => {
        if (messageTime != null) {
            return <Tooltip title={'Sent at ' + messageTime}><AccessTimeIcon size={20} sx={{mr: 1, mt: 1, color: 'primary.main'}}/></Tooltip>;
        } else {
            return <CircularProgress size={20} sx={{mr: 1, mt: 1}}/>;
        }
    };

    const onMouseOverMessage = (event) => {
        reactionButtonRef.current.style.visibility = "visible";
        if (replyButtonRef.current != null) {
            replyButtonRef.current.style.visibility = "visible";
        }
        if (moreButtonRef.current != null) {
            moreButtonRef.current.style.visibility = "visible";
        }
    };

    const onMouseOutMessage = (event) => {
        reactionButtonRef.current.style.visibility = "hidden";
        if (replyButtonRef.current != null) {
            replyButtonRef.current.style.visibility = "hidden";
        }
        if (moreButtonRef.current != null) {
            moreButtonRef.current.style.visibility = "hidden";        
        }
    };

    const onReplyClicked = () => {
        if (onReply) {
            onReply(message, user, messageBatchId);
        }
    };

    const onReactionClicked = (event) => {
        setState({
            ...state,
            reactionsEl: reactionButtonRef.current,
        });
    };

    const onReactionClose = (event) => {
        setState({
            ...state,
            reactionsEl: null,
        });
    };

    const getHost = () => {
        if (messageType === 'direct') {
            return HostManager.myHost();
        } else if (messageType === 'thread') {
            const toSplit = to.split('_');
            if (toSplit.length == 2) {
                const instanceId = parseInt(toSplit[0]);
                return HostManager.amperHostById(instanceId);
            }
        } else if (messageType === 'channel') {
            //since the host property is provided in case of channel
            //it is considered this case will not be reached
        }
    };
    
    const onEmoticonReacted = (emoticon) => {
        if (message.reactions == null) {
            message.reactions = {};
        }
        const reactions = message.reactions;
        if (reactions[emoticon] == null) {
            reactions[emoticon] = [];
        }
        let opperationType = 'add';
        if (reactions[emoticon].includes(sessionManager.getUser().id)) {
            reactions[emoticon] = reactions[emoticon].filter(userId => userId != sessionManager.getUser().id);
            opperationType = 'remove';
        } else {
            reactions[emoticon].push(sessionManager.getUser().id);
        }
        setState({
            ...state,
            reactionsEl: null,
        });
        post(`${host != null ? host : getHost()}chat/` + (replyMessage ? 'updateReply' : 'update'), {
            batchId,
            batchId1,
            messageId: message.id,
            messageType,
            opperationType,
            value: emoticon,
            updateType: 'reaction',
            from: (from + ''),
            to: (to + ''),
            participants: Object.keys(participants),
        }, (result) => {
            onMessageReaction(message.id, reactions);
        }, (result) => {
            app.toast('info', `not able to sync the reaction with server`);
        });
    };

    const onMoreClicked = (event) => {
        setState({
            ...state,
            moreMenueEl: event.currentTarget,
        });
    };
    const onMoreClose = () => {
        setState({
            ...state,
            moreMenueEl: null,
        });
    };

    const onEditClicked = () => {
        setState({
            ...state,
            moreMenueEl: null,
            mode: 'edit',
        });
    };

    const saveMessage = (value) => {
        setState({
            ...state,
            moreMenueEl: null,
            mode: 'view',
        });
        onMessageUpdate(message, value);
        post(`${host != null ? host : getHost()}chat/` + (replyMessage ? 'updateReply' : 'update'), {
            batchId,
            batchId1,
            messageId: message.id,
            messageType,
            opperationType: 'update',
            value: value,
            updateType: 'edit',
            from: (from + ''),
            to: (to + ''),
            participants: Object.keys(participants),
        }, (result) => {
            
        }, (result) => {
            app.toast('info', `not able to sync the reaction with server`);
        });
    };

    const cancelMessageEdit = () => {
        setState({
            ...state,
            mode: 'view',
        });
    };

    const onMessageRemoveClicked = () => {
        onMessageRemove(message);
        post(`${host != null ? host : getHost()}chat/` + (replyMessage ? 'updateReply' : 'update'), {
            batchId,
            batchId1,
            messageId: message.id,
            messageType,
            opperationType: 'remove',
            value: (from + ''),
            updateType: 'remove',
            from: (from + ''),
            to: (to + ''),
            participants: Object.keys(participants),
        }, (result) => {
            
        }, (result) => {
            app.toast('info', `not able to sync the reaction with server`);
        });
    };

    const onAttachmentDownload = (id, fileName) => {
        download(`${HostManager.myHost()}chat/download?id=` + encodeURIComponent(id) + '&fileName=' + encodeURIComponent(fileName));
    };

    const getAttachmentsLeftSided = () => {
        if (message.attachments && message.attachments.length > 0) {
            const attachmentElements = []
            for (let i = 0; i < message.attachments.length; i++) {
                attachmentElements.push(
                    <Attachment metadata={message.attachments[i]} onDownload={onAttachmentDownload}>
                    </Attachment>
                );
            }
            return <Box sx={{display: 'flex', ml: 1, mt: 1, overflowY: 'hidden', overflowX: 'auto'}} height="60px">
                {attachmentElements}
            </Box>;
        }
    };

    const getAttachmentsRightSided = () => {
        if (message.attachments && message.attachments.length > 0) {
            const attachmentElements = []
            for (let i = 0; i < message.attachments.length; i++) {
                attachmentElements.push(
                    <Attachment metadata={message.attachments[i]} onDownload={onAttachmentDownload}>
                    </Attachment>
                );
            }
            return <Box sx={{display: 'flex', ml: 1, mt: 1, overflowY: 'hidden', overflowX: 'auto'}} height="60px">
                <Box sx={{flexGrow: 1}}>
                </Box>
                <Box sx={{flexGrow: 0, display: 'flex'}}>                
                    {attachmentElements}
                </Box>
            </Box>;
        }
    };

    const getReactedUsers = (usersSet) => {
        return usersSet.map(userId => (<div>{participants[userId]?.firstName + ' ' + participants[userId]?.lastName}</div>));
    };

    const getToolBoxButtons = (viceOrder) => {
        const result = [];
        if (viceOrder) {
            if (message.reactions != null) {
                for (const [emoticon, usersSet] of Object.entries(message.reactions)) {
                    if (usersSet.length > 0) {
                        result.push(<Tooltip title={getReactedUsers(usersSet)}>
                                <Chip onClick={() => {onEmoticonReacted(emoticon)}} style={{cursor: 'pointer'}} label={usersSet.length > 1 ? (emoticon + ' ' + usersSet.length) : emoticon} variant="outlined" sx={{ml:'1px', mr:'1px', mt: '3px'}}/>
                            </Tooltip>);
                    }
                }
            }
            result.push(
                <IconButton ref={reactionButtonRef} sx={{visibility: 'hidden'}} color="primary" aria-label="Reply" onClick={onReactionClicked}>
                    <MoodIcon />
                </IconButton>);
                if (!replyMessage) {
                    if (message.replies != null && Object.values(message.replies).length > 0) {
                        const sum = Object.values(message.replies).reduce((partialSum, a) => partialSum + a, 0);
                        result.push(<Chip onClick={onReplyClicked} style={{cursor: 'pointer'}} label={sum + (sum > 1 ? ' Replies' : ' Reply')} variant="outlined" sx={{ml:'1px', mr:'1px', mt: '3px'}}/>);
                    } else {
                        result.push(
                            <IconButton ref={replyButtonRef} sx={{visibility: 'hidden'}} color="primary" aria-label="Reply" onClick={onReplyClicked}>
                                <ReplyIcon />
                            </IconButton>);
                    }
                }
        } else {
            if (!replyMessage) {
                if (message.replies != null && Object.values(message.replies).length > 0) {
                    const sum = Object.values(message.replies).reduce((partialSum, a) => partialSum + a, 0);
                    result.push(<Chip onClick={onReplyClicked} style={{cursor: 'pointer'}} label={sum + (sum > 1 ? ' Replies' : ' Reply')} variant="outlined" sx={{ml:'1px', mr:'1px', mt: '3px'}}/>);
                } else {
                    result.push(
                        <IconButton ref={replyButtonRef} sx={{visibility: 'hidden'}} color="primary" aria-label="Reply" onClick={onReplyClicked}>
                            <ReplyIcon />
                        </IconButton>);
                }
            }
            
            result.push(
                <IconButton ref={reactionButtonRef} sx={{visibility: 'hidden'}} color="primary" aria-label="Reply" onClick={onReactionClicked}>
                    <MoodIcon />
                </IconButton>);
    
            if (message.reactions != null) {
                for (const [emoticon, usersSet] of Object.entries(message.reactions)) {
                    if (usersSet.length > 0) {
                        result.push(<Tooltip title={getReactedUsers(usersSet)}>
                                <Chip onClick={() => {onEmoticonReacted(emoticon)}} style={{cursor: 'pointer'}} label={usersSet.length > 1 ? (emoticon + ' ' + usersSet.length) : emoticon} variant="outlined" sx={{ml:'1px', mr:'1px', mt: '3px'}}/>
                            </Tooltip>);
                    }
                }   
            }
        }
        return result;
    };

    const getLeftSided = () => {
        return (<Box sx={{mb:0, position: 'relative'}} onMouseOver={onMouseOverMessage} onMouseOut={onMouseOutMessage}>
            <IconButton ref={moreButtonRef} sx={{visibility: 'hidden', position: 'absolute', top: 0, right: 0}} color="primary" aria-label="More" onClick={onMoreClicked}>
                <MoreHorizIcon />
            </IconButton>
            <Menu
                id={message.id + '_more'}
                anchorEl={state.moreMenueEl}
                open={Boolean(state.moreMenueEl)}
                onClose={onMoreClose}
                MenuListProps={{
                    'aria-labelledby': 'basic-button',
                }}>
                <MenuItem onClick={onEditClicked}>Edit</MenuItem>
                <MenuItem onClick={onMessageRemoveClicked}>Remove</MenuItem>
            </Menu>
            {Boolean(state.reactionsEl) && <EmotionPicker onSelect={onEmoticonReacted} el={state.reactionsEl} onClose={onReactionClose}></EmotionPicker>}
            <Box sx={{display: 'flex'}}>
                <Box sx={{flexGrow: 0}}>
                    <Box sx={{display: 'flex'}}>
                        {getMessageTime()}
                        <Tooltip title={user.firstName + ' ' + user.lastName}><Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 1 }} src={getImageSource(user)} /></Tooltip>
                    </Box>
                </Box>
                <Box sx={{flexGrow: 1, textAlign: 'left', mr: 3, mt: 1}}>
                    <Box className={'messageUserName'}>
                        {user.firstName + ' ' + user.lastName} <Typography variant="caption">
                            {messageTime}
                        </Typography>
                    </Box>
                    <Box dangerouslySetInnerHTML={{ __html: message.text }}>
                    </Box>
                </Box>
            </Box>
            {getAttachmentsLeftSided()}
            <Box className="messageToolBox" sx={{display: 'flex', ml: 10, overflowY: 'auto'}}>
                <Box sx={{flexGrow: 1}}>
                    {getToolBoxButtons(false)}
                </Box>
                <Box sx={{flexGrow: 0, mr: 10}}>                
                </Box>
            </Box>
        </Box>);
    };

    const getRightSided = () => {
        return (<Box sx={{mb:0}} onMouseOver={onMouseOverMessage} onMouseOut={onMouseOutMessage}>
            {Boolean(state.reactionsEl) && <EmotionPicker onSelect={onEmoticonReacted} el={state.reactionsEl} onClose={onReactionClose}></EmotionPicker>}
            <Box sx={{display: 'flex'}}>
                <Box sx={{flexGrow: 1, textAlign: 'right', mr: 3, mt: 1}}>
                    <Box className={'messageUserName'}>
                        <Typography variant="caption">
                            {messageTime}
                        </Typography> {user.firstName + ' ' + user.lastName}
                    </Box>
                    <Box dangerouslySetInnerHTML={{ __html: message.text }} sx={{ftextAlign: 'left'}}>
                    </Box>
                </Box>
                <Box sx={{flexGrow: 0}}>
                    <Box sx={{display: 'flex'}}>
                        <Tooltip title={user.firstName + ' ' + user.lastName}><Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 1 }} src={getImageSource(user)} /></Tooltip>
                        {getMessageTime()}
                    </Box>
                </Box>
            </Box>
            {getAttachmentsRightSided()}
            <Box className="messageToolBox" sx={{display: 'flex'}}>
                <Box sx={{flexGrow: 1}}>
                </Box>
                <Box sx={{flexGrow: 0, mr: 10, overflowY: 'auto'}}>
                    {getToolBoxButtons(true)}
                </Box>
            </Box>
        </Box>
        );
    };

    const getEditor = () => {
        return <Box sx={{ml: 1, mr: 1}}><ChatEditor content={message.text} showSaveCancel={true} save={saveMessage} cancel={cancelMessageEdit}></ChatEditor></Box>;
    };

    const getView = () => {
        if (state.mode =='view') {
            return (left ? getLeftSided() : getRightSided());
        } else {
            return getEditor();
        }
    };
    return getView();
}