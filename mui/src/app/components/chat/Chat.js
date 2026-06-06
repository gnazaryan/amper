import React, { useRef, useEffect } from 'react';
import useState from 'react-usestateref'
import Box from '@mui/material/Box';
import EditIcon from '@mui/icons-material/Edit';
import IconButton from '@mui/material/IconButton';
import Link from '@mui/material/Link';
import UserDialog from '../user/UserDialog';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';
import Convenience from '../../help/Convenience';
import "./Chat.css"
import ChatEditor from './ChatEditor';
import Typography from '@mui/material/Typography';
import Message from './Message';
import { sessionManager } from '../../../SessionManager';
import { v4 as uuidv4 } from 'uuid';
import { post } from '../../data/Submit';
import HostManager from '../../../HostManager';
import { AppContext } from '../../../App';
import Button from '@mui/material/Button';
import ChatReplies from './ChatReplies';
import Attachment from './Attachment';
import GroupsIcon from '@mui/icons-material/Groups';
import { truncate } from '../../amper/Instruments';
import TagIcon from '@mui/icons-material/Tag';
import OfflineBoltIcon from '@mui/icons-material/OfflineBolt';
import ElectricBoltIcon from '@mui/icons-material/ElectricBolt';

export default function Chat() {
    const app = React.useContext(AppContext);
    const messageContainerRef = useRef(null);
    const messageEndRef = useRef(null);

    app.registerServerUpdate('chat', (updates) => {
        applyUpdates(updates);
    });
    /*let timer = -1;
    const onChatContentScroll = (event) => {
        if (messageContainerRef.current) {
            if (messageContainerRef.current.scrollTop === 0) {
                if (timer != -1) {
                    clearTimeout(timer)
                }
                timer = setTimeout(() => {

                }, 500);
            }
        }
    };*/

    const [state, setState, stateRef] = useState({
        channelGroups: [],
        threads: [],
        directs: [],
        activeChat: null,
        userDialogOpen: false,
        newMessages: {
            directs: {},
            threads: {},
            channels:{},
        },
        chatStateLoaded: false,
        chatContentLoading: false,
        chatContentPage: 0,
        chatContent: [],
        participants: {},
        reservedMessageId: null,
        attachments: [],
    });

    const [chatRepliesState, setChatRepliesState] = useState({
        open: false,
        message: null,
    });

    useEffect(() => {
        const timer = setTimeout(() => {
            messageEndRef.current?.scrollIntoView();
            /*if (messageContainerRef.current != null) {
                let height = 0;
                for (const child of messageContainerRef.current.children) {
                    height+=child.clientHeight;
                    if (child.id=='scrollToElement') {
                        break;
                    }
                }
                messageContainerRef.current?.scrollTo({
                    top: height,
                    left: 0,
                });
            }*/
        }, 500);
        return () => clearTimeout(timer);
    }, [state.scrollToMessageId]);

    useEffect(() => {
        if (!state.chatStateLoaded) {
            post(`${HostManager.myHost()}chat/state`, {
            }, (result) => {
                const directs = mergeDirectsWithExistingDirects(state.directs, result.directs);
                setState({
                    ...state,
                    directs,
                    threads: result.threads,
                    channelGroups: result.channelGroups,
                });
            }, (result) => {debugger
                app.toast('info', `not able to retrieve the chat state`);
            });
        }
    }, [state.chatStateLoaded]);
    
    useEffect(() => {
        if (state.chatContentLoading) {
            const type = state.activeChat.type;
            let historyItem = null;
            if (state.activeChat.chatItem 
                && state.activeChat.chatItem.chatHistory
                && state.activeChat.chatItem.chatHistory.historyItems.length >= state.chatContentPage) {
                    historyItem = state.activeChat.chatItem.chatHistory.historyItems[state.activeChat.chatItem.chatHistory.historyItems.length - state.chatContentPage - 1];
            }
            if (historyItem != null) {
                post(`${HostManager.myHost()}chat/fetch`, {
                    from: state.activeChat.chatItem.chatHistory.from + "",
                    to: state.activeChat.chatItem.chatHistory.to + "",
                    type: type,
                    id: historyItem.id,
                    IncludeLatest: state.chatContentPage == 0,
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
                    } else {
                        loadMoreChat();
                    }
                }, (result) => {
                    app.toast('info', `not able to retrieve the chat content`);
                });
            }
        }
    }, [state.chatContentLoading, state.chatContentPage]);

    const applyUpdates = (updates) => {
        let activeChatContent = [
            ...state.chatContent,
        ];
        const newActiveMessages = [];
        let updated = false;
        let directs = state.directs;
        let threads = state.threads;
        let channelGroups = state.channelGroups;
        const participants = state.participants;
        for (let i = 0; i < updates.length; i++) {
            const update = updates[i];
            if (isCurrentActiveChatUpdate(update)) {
                if (update.updateType === 'newMessage') {
                    if (state.chatContent.length == 0 || 
                        (update.message.dateTime >= state.chatContent[state.chatContent.length - 1].dateTime 
                            && update.message.id !== state.chatContent[state.chatContent.length - 1].id)) {
                        let host = '';
                        let toParty = '';
                        if (update.messageType === 'direct') {
                            host = HostManager.myHost();
                            toParty = update.from;
                            //in case a new message is received, check to see if the batch id is in the client chat
                            //if not add it to the chat history's history items
                            //this behaviour is only for direct chat case, because for the thread and channel
                            //the solution is provided as part of update type 'newBatch'
                            for (let i = 0; i < directs.length; i++) {
                                if (parseInt(directs[i].chatHistory.to) === parseInt(update.from)) {
                                    let historyItemFound = false;
                                    for (let l = 0; l < directs[i].chatHistory.historyItems.length; l++) {
                                        if (directs[i].chatHistory.historyItems[l].id === update.message.batchId) {
                                            historyItemFound = true;
                                        }
                                    }
                                    if (!historyItemFound) {
                                        directs[i].chatHistory.historyItems.push({
                                            id: update.message.batchId,
                                            full: false,
                                        });
                                    }
                                    break
                                }
                            }
                        } else if (update.messageType === 'thread') {
                            const to = update.message.to;
                            let instanceId = null;
                            const toSplit = to.split('_');
                            if (toSplit.length == 2) {
                                instanceId = parseInt(toSplit[0]);
                                host = HostManager.amperHostById(instanceId)
                            }
                            toParty = update.message.to;
                        }
                        //skip marking unread for channel chats
                        if (update.messageType !== 'channel') {
                            markChatUnread(host, toParty, update.messageType);
                        }
                        if (update.message.fromUser != null) {
                            participants[update.message.fromUser.id] = update.message.fromUser;
                        }
                        newActiveMessages.push(update.message);
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
                            updated = true;
                            if (update.message.fromUser != null) {
                                participants[update.message.fromUser.id] = update.message.fromUser;
                            }
                        }
                    }
                } else if (update.updateType === 'remove') {
                    activeChatContent = activeChatContent.filter(message => {
                        return message.id !== update.message.id;
                    });
                    updated = true;
                } else if (update.updateType === 'reply') {
                    for (let l = 0; l < activeChatContent.length; l++) {
                        if (activeChatContent[l].id == update.message.id) {
                            const replyUserId = parseInt(update.value)
                            if (!activeChatContent[l].replies) {
                                activeChatContent[l].replies = {};
                            }
                            if (activeChatContent[l].replies[replyUserId] > 0) {
                                activeChatContent[l].replies[replyUserId] = activeChatContent[l].replies[replyUserId] + 1;
                            } else {
                                activeChatContent[l].replies[replyUserId] = 1;
                            }
                            updated = true;
                        }
                    }
                } else if (update.updateType === 'replyInitialisation') {
                    for (let l = 0; l < activeChatContent.length; l++) {
                        if (activeChatContent[l].id == update.message.id) {
                            activeChatContent[l].replyBatchId = update.value;
                            updated = true;
                        }
                    }
                }
            } else {
                if (update.updateType === 'newMessage') {
                    if (update.messageType === 'direct') {
                        let found = false;
                        for (let i = 0; i < directs.length; i++) {
                            if (directs[i].chatHistory.to === update.from) {
                                if (directs[i].chatHistory.unreadMessages == null) {
                                    directs[i].chatHistory.unreadMessages = 1;
                                } else {
                                    directs[i].chatHistory.unreadMessages++;
                                }
                                directs[i].chatHistory.lastUpdateTime = Date.now();
                                updated = true;
                                found = true;
                                //in case a new message is received, check to see if the batch id is in the client chat
                                //if not add it to the chat history's history items
                                //this behaviour is only for direct chat case, because for the thread and channel
                                //the solution is provided as part of update type 'newBatch'
                                let historyItemFound = false;
                                for (let l = 0; l < directs[i].chatHistory.historyItems.length; l++) {
                                    if (directs[i].chatHistory.historyItems[l].id === update.message.batchId) {
                                        historyItemFound = true;
                                    }
                                }
                                if (!historyItemFound) {
                                    directs[i].chatHistory.historyItems.push({
                                        id: update.message.batchId,
                                        full: false,
                                    });
                                }
                            }
                        }
                        if (!found) {
                            directs = [
                                {
                                    user: update.users[0],
                                    chatHistory: update.chatHistory,
                                },
                                ...directs
                            ];
                            updated = true;
                        }
                    } else if (update.messageType === 'thread') {
                        let found = false;
                        for (let i = 0; i < threads.length; i++) {
                            if (threads[i].chatHistory.to === update.to) {
                                if (threads[i].chatHistory.unreadMessages == null) {
                                    threads[i].chatHistory.unreadMessages = 1;
                                } else {
                                    threads[i].chatHistory.unreadMessages++;
                                }
                                threads[i].chatHistory.lastUpdateTime = Date.now();
                                updated = true;
                                found = true;
                            }
                        }
                        if (!found) {
                            threads = [
                                {
                                    users: update.users,
                                    chatHistory: update.chatHistory,
                                },
                                ...threads
                            ];
                            updated = true;
                        }
                    }
                }
            }
            //process updates for both cases active and not active
            if (update.updateType === 'newBatch') {
                if (update.messageType === 'channel') {
                    channelGroupsLoop:
                    for (let i = 0; i < channelGroups.length; i++) {
                        for (let l = 0; l < channelGroups[i].channels.length; l++) {
                            if (channelGroups[i].channels[l].channelId == parseInt(update.to)) {
                                channelGroups[i].channels[l].chatHistory.historyItems.push({
                                    id: update.value,
                                    full: false,
                                });
                                updated = true;
                                break channelGroupsLoop;
                            }
                        }
                    }
                } else if (update.messageType === 'thread') {
                    for (let i = 0; i < threads.length; i++) {
                        if (update.to === threads[i].chatHistory.to) {
                            threads[i].chatHistory.historyItems.push({
                                id: update.value,
                                full: false,
                            });
                            updated = true;
                            break
                        }
                    }
                }
            }
        }
        if (updated) {
            const newState = {
                ...state,
                directs,
                threads,
                channelGroups,
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
        if (state.activeChat != null && state.activeChat.type === update.messageType && update.messageType === 'direct') {
            if (state.activeChat.chatItem.chatHistory.to === update.from) {
                return true;
            }
        } else if (state.activeChat != null && state.activeChat.type === update.messageType && update.messageType === 'thread') {
            if (state.activeChat.chatItem.chatHistory.to === update.to) {
                return true;
            }
        } else if (state.activeChat != null && state.activeChat.type === update.messageType && update.messageType === 'channel') {
            if (state.activeChat.chatItem.chatHistory.to === update.to) {
                return true;
            }
        }
        return false;
    };

    const loadMoreChat = () => {
        if (state.activeChat.chatItem.chatHistory.historyItems.length > (state.chatContentPage + 1)) {
            setState({
                ...state,
                chatContentLoading: true,
                chatContentPage: state.chatContentPage + 1,
            });
        } else {
            setState({
                ...state,
                chatContentLoading: false,
            });
        }
    };

    const mergeDirectsWithExistingDirects = (existingDirects, remoteDirects) => {
        const exclusionDirects = [];
        for (let i = 0; i < existingDirects.length; i++) {
            const existingDirect = existingDirects[i];
            let exists = false;
            for (let l = 0; l < remoteDirects.length; l++) {
                if (existingDirect.chatHistory.from == remoteDirects[l].chatHistory.from
                    && existingDirect.chatHistory.to == remoteDirects[l].chatHistory.to
                ) {
                    exists = true;
                    break;
                }
            }
            if (!exists) {
                exclusionDirects.push(existingDirect);
            }
        }
        return [
            ...exclusionDirects,
            ...remoteDirects,
        ];
    };
    
    const openUserDialog = () => {
        setState({
            ...state,
            userDialogOpen: true,
        })
    };

    const closeUserDialog = () => {
        setState({
            ...state,
            userDialogOpen: false,
        });
    };

    const onSelectionChange = (users) => {
        setState({
            ...state,
            selectedUsers: users,
        })
    };

    const startChat = () => {
        const users = state.selectedUsers;
        if (users.length === 1) {
            const newState =  {
                ...state,
                activeChat: {
                    type: 'direct',
                    chatItem: {
                        chatHistory: {
                            from: sessionManager.getUser().id + '',
                            to: users[0].id + '',
                            historyItems: [],
                        },
                        user: users[0]
                    },
                },
                userDialogOpen: false,
            };
            let found = false;
            for (let i = 0; i < state.directs.length; i++) {
                if (newState.directs[i].user.id === users[0].id) {
                    newState.directs[i].chatHistory.lastUpdateTime = Date.now();
                    found = true;
                }
            }
            if (!found) {
                newState.directs = [
                    {
                        chatHistory: {
                            from: sessionManager.getUser().id + '',
                            to: users[0].id + '',
                            historyItems: [],
                            lastUpdateTime: Date.now(),
                        },
                        user: users[0],
                    },
                    ...newState.directs
                ];
            }
            setState(newState);
        } else if (users.length > 1) {
            let threadId = null;
            let historyItems = [];
            const thread = getThreadIfExists(users)
            if (thread != null) {
                threadId = thread.chatHistory.to;
                historyItems = thread.chatHistory.historyItems;
            } else {
                const user = sessionManager.getUser()
                threadId = user.amperId + '_' + uuidv4();
            }
            const newState =  {
                ...state,
                activeChat: {
                    type: 'thread',
                    chatItem: {
                        chatHistory: {
                            from: sessionManager.getUser().id + '',
                            to: threadId,
                            historyItems: historyItems,
                        },
                        users: [
                            ...users,
                            sessionManager.getUser(),
                        ],
                    },
                },
                userDialogOpen: false,
            };

            if (thread == null) {
                newState.threads = [
                    {
                        chatHistory: {
                            from: sessionManager.getUser().id + '',
                            to: threadId,
                            historyItems: historyItems,
                            lastUpdateTime: Date.now(),
                        },
                        users: users,
                    },
                    ...newState.threads,
                ];
            } else {
                const threads = state.threads;
                for (let i = 0; i < threads; i++) {
                    if (threads[i].chatHistory.to === thread.chatHistory.to) {
                        threads[i].chatHistory.lastUpdateTime = Date.now();
                    }
                }
                newState.threads = threads;
            }
            setState(newState);
        }
    };

    const getThreadIfExists = (users) => {
        for (let i = 0; i < state.threads.length; i++) {
            const thread = state.threads[i];
            const threadUserIDs = thread.users.map(user => user.id)
            const userIds = users.map(user => user.id)
            var isEqual = (JSON.stringify(threadUserIDs.sort()) === JSON.stringify(userIds.sort()));
            if (isEqual) {
                return thread;
            }
        }
        return null;
    };

    const getImageSource = (user) => {
        if (Convenience.hasValue(user.photo)) {
            return 'data:image/png;base64,' + user.photo;
        }
        return '/static/images/avatar/2.jpg';
    };

    const selectChat = (chatItem) => {
        const newState = {
            ...state,
            activeChat: chatItem,
            chatContentLoading: true,
            chatContentPage: 0,
            chatContent: [],
            participants: {},
            reservedMessageId: null,
            attachments: [],
        };
        if (chatItem.type === 'direct') {
            const host = HostManager.myHost();
            for (let i = 0; i < newState.directs.length; i++) {
                if (newState.directs[i].chatHistory.to === chatItem.chatItem.chatHistory.to) {
                    newState.directs[i].chatHistory.unreadMessages = 0;
                }
            }
            setState(newState);
            markChatUnread(host, chatItem.chatItem.chatHistory.to, chatItem.type)
        } else if (chatItem.type === 'thread') {
            const to = chatItem.chatItem.chatHistory.to;
            let instanceId = null;
            const toSplit = to.split('_');
            let host = '';
            if (toSplit.length == 2) {
                instanceId = parseInt(toSplit[0]);
                host = HostManager.amperHostById(instanceId)
            }
            for (let i = 0; i < newState.threads.length; i++) {
                if (newState.threads[i].chatHistory.to === chatItem.chatItem.chatHistory.to) {
                    newState.threads[i].chatHistory.unreadMessages = 0;
                }
            }
            setState(newState);
            markChatUnread(host, chatItem.chatItem.chatHistory.to, chatItem.type)
        } else if (chatItem.type === 'channel') {
            setState(newState);
        }
    };

    const markChatUnread = async (host, to, type) => {
        post(`${host}chat/markUnread`, {
            to,
            type,
        }, (result) => {
        }, (result) => {
            app.toast('info', `not able to mark the message as unread`);
        });
    };

    const getDirectsAndThreads = () => {
        const result = [];
        const directsAndThreads = [];
        for (let i = 0; i < state.directs.length; i++) {
            directsAndThreads.push({
                type: 'direct',
                ...state.directs[i],
            });
        }
        for (let i = 0; i < state.threads.length; i++) {
            directsAndThreads.push({
                type: 'thread',
                ...state.threads[i],
            });
        }
        const directsAndThreadsSorted = directsAndThreads.sort((a,b) => b.chatHistory.lastUpdateTime - a.chatHistory.lastUpdateTime);
        for (let i = 0; i < directsAndThreadsSorted.length; i++) {
            const directOrThread = directsAndThreadsSorted[i];
            let id = "";
            let label = "";
            let avatar = null;
            let active = false;
            let unreadMessage = directOrThread.chatHistory.unreadMessages;
            if (directOrThread.type === 'direct') {
                const user = directOrThread.user;
                label = user.firstName + ' ' + user.lastName + (unreadMessage > 0 ? (' (' + unreadMessage + ')') : '')
                id = user.id
                avatar = <ListItemAvatar><Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 2 }} alt={user.firstName + ' ' + user.lastName} src={getImageSource(user)} /></ListItemAvatar>;
                active = state.activeChat != null && state.activeChat.type === 'direct' && state.activeChat.chatItem.user.id == user.id
            } else if (directOrThread.type === 'thread') {
                const users = directOrThread.users;
                if (!directOrThread.label) {
                    for (let l = 0; l < users.length; l++) {
                        label = label + (l > 0 ? ', ' : '') + users[l].firstName;
                    }
                    label = truncate(label, 25, false);
                } else {
                    label = directOrThread.label;
                }
                label = label + (unreadMessage > 0 ? (' (' + unreadMessage + ')') : '');
                avatar = <ListItemAvatar><Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 2 }}><GroupsIcon/></Avatar></ListItemAvatar>;
                id = directOrThread.chatHistory.to;
                active = state.activeChat != null && state.activeChat.chatItem.chatHistory.to === directOrThread.chatHistory.to;
            }

            result.push(<ListItem key={id} disablePadding
                onClick={() => {
                        selectChat({
                            type: directOrThread.type,
                            chatItem: directOrThread,
                        })
                    }
                }
             sx={{backgroundColor: active ? 'primary.selectedBackground' : 'white'}}>
                <ListItemButton>
                    {avatar != null && avatar }
                    <ListItemText id={id} primary={<span style={{fontWeight: unreadMessage > 0 ? '900' : '400'}}>{label}</span>}/>
                </ListItemButton>
              </ListItem>);
        }
        return result;
    };

    const send = (text) => {
        const from = sessionManager.getUser().id.toString();
        const message = {
            from: from,
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
        const type = state.activeChat.type;
        if (type === 'direct') {
            const to = state.activeChat.chatItem.user.id.toString();
            message.to = to;
            let messages = [];
            if (state.newMessages.directs[state.activeChat.chatItem.user.id] != null) {
                messages = state.newMessages.directs[state.activeChat.chatItem.user.id];
            }
            const oldState = {
                ...state,
                scrollToMessageId: message.id,
                newMessages: {
                    ...state.newMessages,
                    directs: {
                        ...state.newMessages.directs,
                        [state.activeChat.chatItem.user.id]: [
                            ...messages,
                            message,
                        ]
                    },
                },
                attachments: []
            };
            setState(oldState);
            post(`${HostManager.myHost()}chat/send`, {
                from,
                to,
                type,
                text: message.text,
                id: message.id,
                attachments: message.attachments,
            }, (result) => {
                const pendingMessages = oldState.newMessages.directs[state.activeChat.chatItem.user.id];
                const excludedPendingMessages = pendingMessages.filter((message) => {
                    return message.id !== result.message.id;
                });
                //iterate and try to find if a new batch was initiated by the send function, 
                //if a new batch id was rolled by the backend, add the batch id to the chat history
                const directs = state.directs;
                for (let i = 0; i < directs.length; i++) {
                    if (parseInt(directs[i].chatHistory.to) === parseInt(to)) {
                        let found = false;
                        for (let l = 0; l < directs[i].chatHistory.historyItems.length; l++) {
                            if (directs[i].chatHistory.historyItems[l].id === result.message.batchId) {
                                found = true;
                            }
                        }
                        if (!found) {
                            directs[i].chatHistory.historyItems.push({
                                id: result.message.batchId,
                                full: false,
                            });
                        }
                        break
                    }
                }
                setState({
                    ...oldState,
                    chatContent: [
                        ...state.chatContent,
                        result.message,
                    ],
                    newMessages: {
                        ...state.newMessages,
                        directs: {
                            ...state.newMessages.directs,
                            [state.activeChat.chatItem.user.id]: excludedPendingMessages,
                        },
                    },
                    directs,
                    attachments: [],
                });
            }, (result) => {
                app.toast('info', `not able to send the message`);
            });
        } else if (type === 'thread') {
            const to = state.activeChat.chatItem.chatHistory.to;
            let instanceId = null;
            const toSplit = to.split('_');
            if (toSplit.length == 2) {
                instanceId = parseInt(toSplit[0]);
            } else {
                return;
            }
            message.to = to;
            let messages = [];
            if (state.newMessages.threads[to] != null) {
                messages = state.newMessages.threads[to];
            }
            const oldState = {
                ...state,
                scrollToMessageId: message.id,
                newMessages: {
                    ...state.newMessages,
                    threads: {
                        ...state.newMessages.threads,
                        [to]: [
                            ...messages,
                            message,
                        ]
                    },
                },
                attachments: []
            };
            setState(oldState);
            const participants = state.activeChat.chatItem.users.map(user => user.id);
            post(`${HostManager.amperHostById(instanceId)}chat/send`, {
                from,
                to,
                type,
                text: message.text,
                id: message.id,
                participants,
                attachments: message.attachments,
            }, (result) => {
                const pendingMessages = oldState.newMessages.threads[to];
                const excludedPendingMessages = pendingMessages.filter((message) => {
                    return message.id !== result.message.id;
                });
                setState({
                    ...oldState,
                    chatContent: [
                        ...state.chatContent,
                        result.message,
                    ],
                    newMessages: {
                        ...state.newMessages,
                        threads: {
                            ...state.newMessages.directs,
                            [to]: excludedPendingMessages,
                        },
                    },
                    attachments: [],
                });
            }, (result) => {
                app.toast('info', `not able to send the message`);
            });
        }  else if (type === 'channel') {
            const to = state.activeChat.chatItem.chatHistory.to;
            let instanceId = state.activeChat.chatItem.amperId;
            message.to = to;
            let messages = [];
            if (state.newMessages.channels[to] != null) {
                messages = state.newMessages.channels[to];
            }
            const oldState = {
                ...state,
                scrollToMessageId: message.id,
                newMessages: {
                    ...state.newMessages,
                    channels: {
                        ...state.newMessages.channels,
                        [to]: [
                            ...messages,
                            message,
                        ]
                    },
                },
                attachments: []
            };
            setState(oldState);
            post(`${HostManager.amperHostById(instanceId)}chat/send`, {
                from,
                to,
                type,
                text: message.text,
                id: message.id,
                attachments: message.attachments,
            }, (result) => {
                const pendingMessages = oldState.newMessages.channels[to];
                const excludedPendingMessages = pendingMessages.filter((message) => {
                    return message.id !== result.message.id;
                });
                setState({
                    ...oldState,
                    chatContent: [
                        ...state.chatContent,
                        result.message,
                    ],
                    participants: {
                        ...state.participants,
                        [result.message.fromUser.id]: result.message.fromUser,
                    },
                    newMessages: {
                        ...state.newMessages,
                        channels: {
                            ...state.newMessages.directs,
                            [to]: excludedPendingMessages,
                        },
                    },
                    attachments: [],
                });
            }, (result) => {
                app.toast('info', `not able to send the message`);
            });
        }
    };

    const onMessageRemply = (message, user) => {
        const from = sessionManager.getUser().id;
        let participants = {};
        let to = null;
        if (state.activeChat.type === "direct") {
            to = state.activeChat.chatItem.user.id;
            participants[sessionManager.getUser().id.toString()] = sessionManager.getUser();
            participants[state.activeChat.chatItem.user.id.toString()] = state.activeChat.chatItem.user;
        } else if (state.activeChat.type === "thread") {
            to = state.activeChat.chatItem.chatHistory.to;
            for (let i = 0; i < state.activeChat.chatItem.users.length; i++) {
                const user = state.activeChat.chatItem.users[i];
                participants[user.id.toString()] = user;
            }
        } else if (state.activeChat.type === "channel") {
            to = state.activeChat.chatItem.chatHistory.to;
            participants = state.participants;
        }
        setChatRepliesState({
            ...chatRepliesState,
            message: message,
            messageUser: user,
            messageType: state.activeChat.type,
            from: from,
            to: to,
            participants: participants,
            chatItem: state.activeChat.chatItem,
            open: true,
        });
    };

    const closeMessageReply = () => {
        setChatRepliesState({
            ...chatRepliesState,
            open: false,
        });
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

    const onChatReply = (messageId, from, batchId) => {
        const messages = state.chatContent;
        const fromInt = parseInt(from);
        for (let i = 0; i < messages.length; i++) {
            if (messages[i].id === messageId) {
                if (!messages[i].replies) {
                    messages[i].replies = {};
                }
                if (messages[i].replies[fromInt] > 0) {
                    messages[i].replies[fromInt] = messages[i].replies[fromInt] + 1;
                } else {
                    messages[i].replies[fromInt] = 1;
                }
                messages[i].replyBatchId = batchId;
            }
        }
        setState({
            ...state,
            chatContent: messages,
        });
    };

    const onMessageReaction = (messageId, reactions) => {
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
        let directory = '__system__/Chat';
        if (state.activeChat.type === 'direct') {
            directory = directory + '/' + state.activeChat.chatItem.user.id + '/' + reservedMessageId
        } else if (state.activeChat.type === 'thread') {
            directory = directory + '/' + state.activeChat.chatItem.chatHistory.from + '/' + reservedMessageId
        } else if (state.activeChat.type === 'channel') {debugger
            directory = directory + '/' + state.activeChat.chatItem.chatHistory.from + '/' + reservedMessageId
        }
        app.upload(directory, files, (dir, metadata) => {
            onUploadComplete(reservedMessageId, metadata, directory);
        }, false, null);
    };

    const getChannelChatContent = () => {
        const result = [];
        if (state.chatContent != null) {
            const sortedMessages = state.chatContent.sort((a,b) => a.dateTime - b.dateTime);
            for (let i = 0; i < sortedMessages.length; i++) {
                const message = sortedMessages[i];
                let user = null;
                if (message.fromUser != null) {
                    user = message.fromUser;
                } else {
                    user = state.participants[parseInt(message.from)]
                }
                result.push(<Message key={message.id}
                    participants={state.participants}
                    host={HostManager.amperHostById(state.activeChat.chatItem.amperId)}
                    from={sessionManager.getUser().id} 
                    to={state.activeChat.chatItem.chatHistory.to} 
                    onReply={onMessageRemply} 
                    batchId={message.batchId} 
                    batchId1={message.batchId1} 
                    left={sessionManager.getUser().id == message.from} 
                    message={message}
                    user={user} 
                    messageType={state.activeChat.type}
                    onMessageUpdate={onMessageUpdate}
                    onMessageRemove={onMessageRemove}
                    onMessageReaction={onMessageReaction}
                    replyMessage={false}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>);
                }
            }
        }

        const messages = state.newMessages.channels[state.activeChat.chatItem.chatHistory.to];
        if (messages != null) {
            for (let i = 0; i < messages.length; i++) {
                const message = messages[i];
                result.push(<Message key={message.id}
                    participants={state.participants}
                    host={HostManager.amperHostById(state.activeChat.chatItem.amperId)}
                    message={message} 
                    user={sessionManager.getUser()} 
                    left={sessionManager.getUser().id == message.from} 
                    messageTime={null} 
                    replyMessage={false}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>)
                }
            }
        }
        return result;
    };
    const getThreadChatContent = () => {
        const result = [];
        const participants = {};
        for (let i = 0; i < state.activeChat.chatItem.users.length; i++) {
            const user = state.activeChat.chatItem.users[i];
            participants[user.id.toString()] = user;
        }
        if (state.chatContent != null) {
            const sortedMessages = state.chatContent.sort((a,b) => a.dateTime - b.dateTime);
            for (let i = 0; i < sortedMessages.length; i++) {
                const message = sortedMessages[i];
                let user = null;
                for (let i = 0; i < state.activeChat.chatItem.users.length; i++) {
                    if (state.activeChat.chatItem.users[i].id == message.from) {
                        user = state.activeChat.chatItem.users[i];
                    }
                }
                
                result.push(<Message key={message.id} 
                    from={sessionManager.getUser().id} 
                    to={state.activeChat.chatItem.chatHistory.to} 
                    onReply={onMessageRemply} 
                    batchId={message.batchId} 
                    batchId1={message.batchId1} 
                    left={sessionManager.getUser().id == message.from} 
                    message={message}
                    user={user} 
                    messageType={state.activeChat.type}
                    onMessageUpdate={onMessageUpdate}
                    onMessageRemove={onMessageRemove}
                    onMessageReaction={onMessageReaction}
                    replyMessage={false}
                    participants={participants}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>);
                }
            }
        }

        const messages = state.newMessages.threads[state.activeChat.chatItem.chatHistory.to];
        if (messages != null) {
            for (let i = 0; i < messages.length; i++) {
                const message = messages[i];
                result.push(<Message key={message.id} 
                    message={message} 
                    user={sessionManager.getUser()} 
                    left={sessionManager.getUser().id == message.from} 
                    messageTime={null} 
                    replyMessage={false}
                    participants={participants}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>)
                }
            }
        }

        return result;
    };

    const getDirectChatContent = () => {
        const result = [];
        const participants = {
            [sessionManager.getUser().id.toString()]: sessionManager.getUser(),
            [state.activeChat.chatItem.user.id.toString()]: state.activeChat.chatItem.user,
        };
        if (state.chatContent != null) {
            const sortedMessages = state.chatContent.sort((a,b) => a.dateTime - b.dateTime);
            for (let i = 0; i < sortedMessages.length; i++) {
                const message = sortedMessages[i];
                const user = state.activeChat.chatItem.user.id == message.from ? state.activeChat.chatItem.user : sessionManager.getUser();
                
                result.push(<Message key={message.id} 
                    from={sessionManager.getUser().id} 
                    to={state.activeChat.chatItem.user.id} 
                    onReply={onMessageRemply} 
                    batchId={message.batchId} 
                    batchId1={message.batchId1} 
                    left={sessionManager.getUser().id == message.from} 
                    message={message}
                    user={user} 
                    messageType={state.activeChat.type}
                    onMessageUpdate={onMessageUpdate}
                    onMessageRemove={onMessageRemove}
                    onMessageReaction={onMessageReaction}
                    replyMessage={false}
                    participants={participants}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>);
                }
            }
        }

        const messages = state.newMessages.directs[state.activeChat.chatItem.user.id];
        if (messages != null) {
            for (let i = 0; i < messages.length; i++) {
                const message = messages[i];
                result.push(<Message key={message.id} 
                    message={message} 
                    user={sessionManager.getUser()} 
                    left={sessionManager.getUser().id == message.from} 
                    messageTime={null} 
                    replyMessage={false}
                    participants={participants}/>);
                if (message.id === state.scrollToMessageId) {
                    result.push(<div id="scrollToElement" ref={messageEndRef}></div>)
                }
            }
        }

        return result;
    };

    const getChatContent = () => {
        const result = [];
        if (state.activeChat.chatItem.chatHistory.historyItems.length > (state.chatContentPage + 1)) {
            result.push(<Box sx={{height:'40px', mt: 1, textAlign: 'center'}}>
                <Button variant="outlined" onClick={loadMoreChat}>Load more...</Button>
            </Box>);
        }
        if (state.activeChat.type === "direct") {
            result.push(...getDirectChatContent());
        } else if (state.activeChat.type === "thread") {
            result.push(...getThreadChatContent());
        } else if (state.activeChat.type === "channel") {
            result.push(...getChannelChatContent());
        }
        return result;
    };

    const getChatReplies = () => {
        if (chatRepliesState.open) {
            return <ChatReplies
                        key={chatRepliesState.message.id}
                        open={chatRepliesState.open} 
                        message={chatRepliesState.message} 
                        messageUser={chatRepliesState.messageUser} 
                        close={closeMessageReply} 
                        messageType={chatRepliesState.messageType}
                        from={chatRepliesState.from}
                        to={chatRepliesState.to}
                        onMessageReaction={onMessageReaction}
                        participants={chatRepliesState.participants}
                        onChatReply={onChatReply}
                        chatItem={chatRepliesState.chatItem}>
                </ChatReplies>;
        }
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

    const getRightContent = () => {
        if (state.activeChat != null) {
            return <Box className="chatRightPanel">
                        <Box width={chatRepliesState.open ? '50%' : '100%'} height="100%" className={"chatOuterContent"}>
                            <Box ref={messageContainerRef} className={"chatInnerContent"} sx={{height: state.attachments.length > 0 ? 'calc(100vh - 390px)' : 'calc(100vh - 340px)'}}>
                                {getChatContent()}
                            </Box>
                            {getAttachments()}
                            <Box width="100%">
                                <ChatEditor content={''} send={send} onAttachmentClicked={onAttachmentClicked} sendEnabled={state.attachments.length > 0}></ChatEditor>
                            </Box>
                        </Box>
                        {getChatReplies()}
                    </Box>;
        } else {
            return <Box width="calc(100%)" height="100%" sx={{textAlign: 'center'}}>
                <Typography variant="subtitle1" gutterBottom sx={{verticalAlign: 'middle', lineHeight: '700px'}}>
                    Select a chat or start a new conversation...
                </Typography>
            </Box>;
        }
    };

    const getChannels = (channels) => {
        const result = [];
        for (let i = 0; i < channels.length; i++) {
            const active = state.activeChat != null && state.activeChat.chatItem.chatHistory.to === channels[i].chatHistory.to;
            result.push(<ListItem key={channels[i].channelId} disablePadding
                onClick={() => {
                        selectChat({
                            type: 'channel',
                            chatItem: channels[i],
                        });
                    }
                }
             sx={{backgroundColor: active ? 'primary.selectedBackground' : 'white'}}>
                <ListItemButton>
                    <ListItemAvatar><Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 2 }}><ElectricBoltIcon/></Avatar></ListItemAvatar>
                    <ListItemText id={channels[i].channelId} primary={<span>{channels[i].label}</span>}/>
                </ListItemButton>
              </ListItem>);
        }
        return result;
    };

    const getChannelGroups = () => {
        const result = [];
        for (let i = 0; i < state.channelGroups.length; i++) {
            result.push(<Box width="calc(100% - 22px)" maxHeight="250px" minHeight="50px" className="chatLeftChannelPanel">
                <span className="groupTagLabel">
                    {state.channelGroups[i].label}
                </span>
                <Box width="100%" maxHeight="250px" minHeight="50px" sx={{overflowX: 'hidden', overflowY: 'auto'}}>
                    <List dense sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper', mt: 1}}>
                        {getChannels(state.channelGroups[i].channels)}
                    </List>
                </Box>
            </Box>);
        }
        return result;
    };

    return (
        <Box style={{}} width="100%" height="100%" className="chatMainPanel">
            <UserDialog open={state.userDialogOpen} close={closeUserDialog} chat={startChat} onSelectionChange={onSelectionChange}></UserDialog>
            <Box width="250px" height="100%">
                <Box className="newMessagePanel">
                    <IconButton color="primary" onClick={openUserDialog}>
                        <EditIcon />
                    </IconButton>
                </Box>
                <Box width="100%" height="calc(100% - 40px)" sx={{overflowY: 'auto'}}>
                    {getChannelGroups()}
                    <Box width="calc(100% - 22px)" maxHeight="450px" minHeight="250px" className="chatLeftChannelPanel">
                        <span className="groupTagLabel">
                            Directs & Threads
                        </span>
                        <Box width="100%" maxHeight="450px" minHeight="250px" style={{overflowX: 'hidden', overflowY: 'auto'}}>
                            <List dense sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper', mt: 1}}>
                                {getDirectsAndThreads()}
                            </List>
                        </Box>
                    </Box>
                </Box>
            </Box>
            <Box width="calc(100% - 250px)" height="100%">
                {getRightContent()}
            </Box>
        </Box>
    );
}