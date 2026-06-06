import React, { useRef, useEffect } from 'react';
import useState from 'react-usestateref'
import Box from '@mui/material/Box';
import CloseIcon from '@mui/icons-material/Close';
import IconButton from '@mui/material/IconButton';
import Message from './Message';
import { formatDate } from '../../amper/Instruments';
import { Tooltip } from '@mui/material';
import Avatar from '@mui/material/Avatar';
import Typography from '@mui/material/Typography';
import Convenience from '../../help/Convenience';
import ChatEditor from './ChatEditor';
import Attachment from './Attachment';
import HostManager from '../../../HostManager';
import { download } from '../../data/Submit';
import Chip from '@mui/material/Chip';
import EmotionPicker from './emotion/EmotionPicker';
import MoodIcon from '@mui/icons-material/Mood';
import { post } from '../../data/Submit';
import { sessionManager } from '../../../SessionManager';
import './ChatReplies.css';
import { AppContext } from '../../../App';
import { v4 as uuidv4 } from 'uuid';
import { base64DecToArr, UTF8ArrToStr } from '../../amper/Instruments';

export default function ChatReplies({open, close, message, messageUser, messageType, from, to, onMessageReaction, participants, onChatReply, chatItem}) {
    const messageEndRef = useRef(null);
    const app = React.useContext(AppContext);
    app.registerServerUpdate('chatReply', (updates) => {
        applyUpdates(updates);
    });
    const [state, setState, stateRef] = useState({
        reactionsEl: null,
        newMessages: [],
        chatContent: [],
        participants: {},
        attachments: [],
        chatContentLoading: true,
        reservedMessageId: null,
    });

    useEffect(() => {
        const timer = setTimeout(() => {
            messageEndRef.current?.scrollIntoView();
        }, 500);
        return () => clearTimeout(timer);
    }, [state.scrollToMessageId]);

    const getChatHost = () => {
        if (message.replyBatchId) {
            const replyToBatchIdDecoded = UTF8ArrToStr(base64DecToArr(message.replyBatchId));
            const replyToBatchIdDecodedSplited = replyToBatchIdDecoded.split('_');
            if (replyToBatchIdDecodedSplited.length == 8) {
                const targetInstanceId = parseInt(replyToBatchIdDecodedSplited[0]);
                return HostManager.amperHostById(targetInstanceId);
            }
        } else {
            if (messageType === 'channel') {
                let instanceId = chatItem.amperId;
                return HostManager.amperHostById(instanceId);
            } else {
                return HostManager.myHost();
            }
        }
    };

    useEffect(() => {
        if (state.chatContentLoading && message.replyBatchId != null) {
            const replyToBatchIdDecoded = UTF8ArrToStr(base64DecToArr(message.replyBatchId));
            const replyToBatchIdDecodedSplited = replyToBatchIdDecoded.split('_');
            if (replyToBatchIdDecodedSplited.length == 8) {
                const targetInstanceId = parseInt(replyToBatchIdDecodedSplited[0]);
                post(`${HostManager.amperHostById(targetInstanceId)}chat/fetchReplies`, {
                    batchId: message.replyBatchId,
                    messageType,
                }, (result) => {
                    if (result.data.length > 0) {
                        setState({
                            ...state,
                            chatContent: [
                                ...result.data,
                                ...state.chatContent,
                            ],
                            participants: {
                                ...result.participants,
                                ...state.participants,
                            },
                            scrollToMessageId: result.data[result.data.length - 1].id,
                            chatContentLoading: false,
                        });
                    }
                }, (result) => {
                    app.toast('info', `not able to retrieve the chat replies content`);
                });
            }
        }
    }, [state.chatContentLoading]);

    const reactionButtonRef = useRef(null);

    const messageTime = formatDate(message.dateTime);
    const getImageSource = (user) => {
        if (Convenience.hasValue(user.photo)) {
            return 'data:image/png;base64,' + user.photo;
        }
        return '/static/images/avatar/2.jpg';
    };

    const applyUpdates = (updates) => {
        let activeChatContent = [
            ...state.chatContent,
        ];
        const newActiveMessages = [];
        const participants = state.participants;
        let updated = false;
        for (let i = 0; i < updates.length; i++) {
            const update = updates[i];
            if (isCurrentActiveChatUpdate(update)) {
                if (update.updateType === 'newMessage') {
                    if (state.chatContent.length == 0 || 
                        (update.message.dateTime >= state.chatContent[state.chatContent.length - 1].dateTime 
                            && update.message.id !== state.chatContent[state.chatContent.length - 1].id)) {
                        newActiveMessages.push(update.message);
                        if (update.message.fromUser != null) {
                            participants[update.message.fromUser.id] = update.message.fromUser;
                        }
                        updated = true;
                    }
                } else if (update.updateType === 'edit') {
                    for (let l = 0; l < activeChatContent.length; l++) {
                        if (activeChatContent[l].id == update.message.id) {
                            activeChatContent[l].text = update.value;
                            updated = true;
                        }
                    }
                } else if (update.updateType === 'reaction') {
                    for (let l = 0; l < activeChatContent.length; l++) {
                        if (activeChatContent[l].id == update.message.id) {
                            if (update.opperationType === 'add') {
                                if (activeChatContent[l].reactions == null) {
                                    activeChatContent[l].reactions = {};
                                }
                                if (activeChatContent[l].reactions[update.value] == null) {
                                    activeChatContent[l].reactions[update.value] = [parseInt(update.from)];
                                } else {
                                    activeChatContent[l].reactions[update.value].push(parseInt(update.from));
                                }
                            } else if (update.opperationType === 'remove') {
                                const from = parseInt(update.from);
                                activeChatContent[l].reactions[update.value] = activeChatContent[l].reactions[update.value].filter(userId => userId != from);
                            }
                            if (update.message.fromUser != null) {
                                participants[update.message.fromUser.id] = update.message.fromUser;
                            }
                            updated = true;
                        }
                    }
                } else if (update.updateType === 'remove') {
                    activeChatContent = activeChatContent.filter(message => {
                        return message.id !== update.message.id;
                    });
                    updated = true;
                }
            }
        }
        if (updated) {
            const newState = {
                ...state,
                chatContent: [
                    ...activeChatContent,
                    ...newActiveMessages,
                ],
                participants,
            };
            if (newActiveMessages.length > 0) {
                newState.scrollToMessageId = newActiveMessages[newActiveMessages.length - 1].id
            }
            setState(newState);
        }
    };

    const isCurrentActiveChatUpdate = (update) => {
        if (update.message && update.message.batchId === message.replyBatchId) {
            return true;
        }
    };

    const send = (text) => {
        const newMessage = {
            from: from,
            to: to + "",
            id: state.attachments.length > 0 ? state.reservedMessageId : uuidv4(),
            text,
            attachments: state.attachments.map((attachment) => {
                return {
                    id: attachment.id,
                    name: attachment.name,
                    directory: attachment.directory,
                };
            }),
        };
        setState({
            ...stateRef.current,
            scrollToMessageId: newMessage.id,
            newMessages: [
                ...stateRef.current.newMessages,
                newMessage,
            ],
        });
        post(`${getChatHost()}chat/sendReply`, {
            batchId: message.replyBatchId,
            from: (from + ''),
            to: (to + ''),
            repliesToMessageId: message.id,
            repliesToMessageBatchId: message.batchId,
            repliesToMessageBatchId1: message.batchId1,
            repliesToMessageType: messageType,
            text: newMessage.text,
            id: newMessage.id,
            attachments: newMessage.attachments,
            participants: Object.keys(participants),
        }, (result) => {
            const pendingMessages = stateRef.current.newMessages;
            const excludedPendingMessages = pendingMessages.filter((message) => {
                return message.id !== result.message.id;
            });
            let participants = state.participants;
            if (result.message.fromUser) {
                participants = {
                    ...participants,
                    [result.message.fromUser.id]: result.message.fromUser,
                };
            }
            setState({
                ...stateRef.current,
                chatContent: [
                    ...stateRef.current.chatContent,
                    result.message,
                ],
                participants: participants,
                newMessages: excludedPendingMessages,
                attachments: [],
            });
            onChatReply(message.id, from, result.message.batchId);
        }, (result) => {
            app.toast('info', `not able to send the message`);
        });
    };

    const onAttachmentDownload = (id, fileName) => {
        download(`${HostManager.myHost()}chat/download?id=` + encodeURIComponent(id) + '&fileName=' + encodeURIComponent(fileName));
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
            let instanceId = chatItem.amperId;
            return HostManager.amperHostById(instanceId);
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
            participants: {
                ...state.participants,
                [sessionManager.getUser().id]: sessionManager.getUser(),
            },
            reactionsEl: null,
        });
        post(`${getHost()}chat/update`, {
            batchId: message.batchId,
            batchId1: message.batchId1,
            messageId: message.id,
            messageType,
            opperationType,
            value: emoticon,
            updateType: 'reaction',
            from: (from + ''),
            to: (to + ''),
        }, (result) => {
            onMessageReaction(message.id, reactions);
        }, (result) => {
            app.toast('info', `not able to sync the reaction with server`);
        });
    };
    
    const onMouseOverMessage = (event) => {
        reactionButtonRef.current.style.visibility = "visible";
    };

    const onMouseOutMessage = (event) => {
        reactionButtonRef.current.style.visibility = "hidden";
    };

    const onMessageRemply = () => {
        //do nothing, since there is no reply to the replied message
    };

    const onMessageUpdate = (updatedMessage, value) => {
        if (state.chatContent != null) {
            const result = [];
            const messages = state.chatContent;
            for (let i = 0; i < messages.length; i++) {
                const message = messages[i];
                if (message.id === updatedMessage.id) {
                    message.text = value;
                }
                result.push(message);
            }
            setState({
                ...state,
                chatContent: result,
            });
        }
    };

    const onMessageRemove = (removedMessage) => {
        if (state.chatContent != null) {
            const result = [];
            const messages = state.chatContent;
            for (let i = 0; i < messages.length; i++) {
                const message = messages[i];
                if (message.id !== removedMessage.id) {
                    result.push(message);
                }
            }
            setState({
                ...state,
                chatContent: result,
            });
        }
    };

    const onReplyMessageReaction = (messageId, reactions) => {
        const messages = state.chatContent;
        for (let i = 0; i < messages.length; i++) {
            if (messages[i].id === messageId) {
                messages[i].reactions = reactions;
            }
        }
        setState({
            ...state,
            participants: {
                ...state.participants,
                [sessionManager.getUser().id]: sessionManager.getUser(),
            },
            chatContent: messages,
        });
    };

    const onUploadComplete = (reservedMessageId, metadata, directory) => {
        setState({
            ...stateRef.current,
            reservedMessageId,
            attachments: [
                ...stateRef.current.attachments,
                {
                    ...metadata,
                    directory,
                }
            ]
        });
    };

    const onAttachmentClicked = (event) => {
        let files = event.currentTarget.files;
        const reservedMessageId = state.reservedMessageId != null ? state.reservedMessageId :  uuidv4();
        let directory = '__system__/Chat/' + from + '/' + reservedMessageId;

        app.upload(directory, files, (dir, metadata) => {
            onUploadComplete(reservedMessageId, metadata, directory);
        }, false, null);
    };

    const onAttachmentRemove = (metadataId, directory) => {
        post(`${HostManager.amperHost()}files-v1/removeFiles`, {
            root: directory,
            ids: [metadataId],
        }, (result) => {
        }, (result) => {
        });
        setState({
            ...state,
            attachments: state.attachments.filter(attachment => attachment.id != metadataId),
        });
    };

    const getChatContent = () => {
        const result = [];
        const allParticipants = {
            ...participants,
            ...state.participants,
        }
        if (state.chatContent != null) {
            const sortedMessages = state.chatContent.sort((a,b) => a.dateTime - b.dateTime);
            for (let i = 0; i < sortedMessages.length; i++) {
                const message = sortedMessages[i];
                const user = allParticipants[message.from];
                result.push(<Message key={message.id} 
                    host={getChatHost()}
                    from={from} 
                    to={to} 
                    participants={allParticipants}
                    onReply={onMessageRemply} 
                    batchId={message.batchId} 
                    batchId1={message.batchId1} 
                    left={(from + '') == message.from} 
                    message={message}
                    user={user} 
                    messageType={messageType}
                    onMessageUpdate={onMessageUpdate}
                    onMessageRemove={onMessageRemove}
                    onMessageReaction={onReplyMessageReaction}
                    replyMessage={true}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>);
                }
            }
        }

        if (state.newMessages != null) {
            for (let i = 0; i < state.newMessages.length; i++) {
                const message = state.newMessages[i];
                result.push(<Message key={message.id} 
                    host={getChatHost()} 
                    message={message} 
                    participants={allParticipants} 
                    user={sessionManager.getUser()} 
                    left={true} messageTime={null} 
                    replyMessage={false}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>)
                }
            }
        }
        return result;
    };

    const getToolBoxButtons = (viceOrder) => {
        const result = [];
    
        result.push(
            <IconButton ref={reactionButtonRef} sx={{visibility: 'hidden'}} color="primary" aria-label="Reply" onClick={onReactionClicked}>
                <MoodIcon />
            </IconButton>);

        if (message.reactions != null) {
            for (const [emoticon, usersSet] of Object.entries(message.reactions)) {
                if (usersSet.length > 0) {
                    result.push(<Chip style={{cursor: 'pointer'}} label={usersSet.length > 1 ? (emoticon + ' ' + usersSet.length) : emoticon} variant="outlined" sx={{ml:'1px', mr:'1px', mt: '3px'}}/>);
                }
            }   
        }
        return result;
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
            return <Box sx={{display: 'flex', ml: 1, overflowY: 'hidden', overflowX: 'auto'}} height="70px">
                {attachmentElements}
            </Box>;
        }
    };

    const getHeader = () => {
        return (<Box sx={{mb:0}}>
            <Box sx={{display: 'flex'}}>
                <Box sx={{flexGrow: 0}}>
                    <Box sx={{display: 'flex'}}>
                        <Tooltip title={messageUser.firstName + ' ' + messageUser.lastName}>
                            <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 1 }} src={getImageSource(messageUser)} />
                        </Tooltip>
                    </Box>
                </Box>
                <Box sx={{flexGrow: 1, textAlign: 'left', mr: 3, mt: 1}}>
                    <Box className={'messageUserName'}>
                        {messageUser.firstName + ' ' + messageUser.lastName} <Typography variant="caption">
                            {messageTime}
                        </Typography>
                    </Box>
                    <Box dangerouslySetInnerHTML={{ __html: message.text }}>
                    </Box>
                </Box>
            </Box>
            {getAttachmentsLeftSided()}
            <Box className="messageToolBox" sx={{display: 'flex', ml: 5, overflowY: 'auto'}}>
                <Box sx={{flexGrow: 1}}>
                    {getToolBoxButtons()}
                </Box>
                <Box sx={{flexGrow: 0, mr: 10}}>                
                </Box>
            </Box>
        </Box>);
    };

    const getAttachments = () => {
        if (state.attachments.length > 0) {
            const attachmentElements = []
            for (let i = 0; i < state.attachments.length; i++) {
                attachmentElements.push(
                    <Attachment metadata={state.attachments[i]} onRemove={onAttachmentRemove}>
                    </Attachment>
                );
            }
            return <Box sx={{display: 'flex', ml: 1, overflowY: 'hidden', overflowX: 'auto'}} minHeight="50px">
                {attachmentElements}
            </Box>;
        }
    };

    return  (<Box sx={{width: open ? '50%' : '0%', height: 'calc(100vh - 126px)', overflowY: 'auto', ml: 1, display: 'flex', flexDirection: 'column'}} onMouseOver={onMouseOverMessage} onMouseOut={onMouseOutMessage}>
            <Box sx={{height: '40px', flexShrink: 0}}>
                <IconButton aria-label="close" onClick={close}>
                    <CloseIcon />
                </IconButton>
            </Box>
            {Boolean(state.reactionsEl) && <EmotionPicker onSelect={onEmoticonReacted} el={state.reactionsEl} onClose={onReactionClose}></EmotionPicker>}
            <Box sx={{maxHeight: '400px', overflowY: 'auto', flexShrink: 0}}>
                {getHeader()}
            </Box>
            <Box sx={{display: 'flex', flexDirection: 'column', mt: '2px', pt: '2px', height: '100%'}} className="chatRepliesMEssageContent">
                {getChatContent()}
            </Box>
            {getAttachments()}
            <Box sx={{ flexShrink: 0}}>
                <ChatEditor content={''} send={send} onAttachmentClicked={onAttachmentClicked} sendEnabled={state.attachments.length > 0}></ChatEditor>
            </Box>
        </Box>);
}