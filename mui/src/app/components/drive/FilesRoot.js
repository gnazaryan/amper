import React, { useState, useEffect, useRef, useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom'
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import LinearProgress from '@mui/material/LinearProgress';
import Grid from '@mui/material/Grid2';
import FileUploadIcon from '@mui/icons-material/FileUpload';
import { AppContext } from '../../../App';
import HostManager from '../../../HostManager';
import DataStore from '../../data/DataStore';
import Paper from '@mui/material/Paper';
import { Typography } from '@mui/material';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import { styled } from "@mui/material/styles";
import makeStyles from '@mui/styles/makeStyles';
import Tooltip from '@mui/material/Tooltip';
import CreateNewFolderIcon from '@mui/icons-material/CreateNewFolder';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import {post, download} from '../../data/Submit';
import FolderIcon from '@mui/icons-material/Folder';
import { breadcrumbs } from '../../center/Breadcrambs';
import FormHelperText from '@mui/material/FormHelperText';
import FormControl from '@mui/material/FormControl';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import DeleteIcon from '@mui/icons-material/Delete';
import DownloadIcon from '@mui/icons-material/Download';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Convenience from '../../help/Convenience';
import Radio from '@mui/material/Radio';
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import Popover from '@mui/material/Popover';
import CircularProgress from '@mui/material/CircularProgress';
import TreeView from '@mui/lab/TreeView';
import TreeItem from '@mui/lab/TreeItem';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import ContentCutIcon from '@mui/icons-material/ContentCut';
import ContentPasteIcon from '@mui/icons-material/ContentPaste';
import { debounceLatest } from '../../amper/Instruments';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import FileView from './FileView';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import PreviewIcon from '@mui/icons-material/Preview';


const LINES_TO_SHOW = 1;
const useStyles = makeStyles({
    multiLineEllipsis: {
        overflow: "hidden",
        textOverflow: "ellipsis",
        display: "-webkit-box",
        "-webkit-line-clamp": LINES_TO_SHOW,
        "-webkit-box-orient": "vertical"
    },
    preventSelection: {
        "-webkit-user-select": "none", /* Safari */
        "-ms-user-select": "none", /* IE 10 and IE 11 */
        "user-select": "none" /* Standard syntax */
    },
    onDropTarget: {
        borderStyle: 'solid;',
        borderColor: '#2196f3;',
        borderWidth: '3px;'
        //background: '#2196f3;', /* Old browsers */
        //background: '-moz-radial-gradient(center, ellipse cover,  #2196f3 0%, #ffffff 100%);', /* FF3.6+ */
        //background: '-webkit-gradient(radial, center center, 0px, center center, 100%, color-stop(0%,#2196f3), color-stop(100%,#ffffff));', /* Chrome,Safari4+ */
        //background: '-webkit-radial-gradient(center, ellipse cover,  #2196f3 0%,#ffffff 100%);', /* Chrome10+,Safari5.1+ */
        //background: '-o-radial-gradient(center, ellipse cover,  #2196f3 0%,#ffffff 100%);', /* Opera 12+ */
        //background: '-ms-radial-gradient(center, ellipse cover,  #2196f3 0%,#ffffff 100%);', /* IE10+ */
        //background: 'radial-gradient(ellipse at center,  #2196f3 0%,#ffffff 100%);', /* W3C */
        //filter: "progid:DXImageTransform.Microsoft.gradient( startColorstr='#2196f3', endColorstr='#ffffff',GradientType=1 );" /* IE6-9 fallback on horizontal gradient */
    }
});

const thumbnailStyles = makeStyles(theme => ({
    box: {
        backgroundColor: '#ffffff',
        '&:hover': {
            cursor: 'pointer',
            "& $buttonFlip": {
                color: "#ffffff",
                visibility: 'visible',
            },
            "& $title": {
                backgroundColor: '#ffffff',
            },
        }
    },
    buttonFlip: () => ({
        color: 'primary.main',
    }),
    title: () => ({
        backgroundColor: 'secondary.main',
    })
  }));
const fileStyles = makeStyles(theme => ({
    box: {
        backgroundColor: '#ffffff',
        '&:hover': {
            cursor: 'pointer',
            backgroundColor: '#2196f3',
            "& $icon": {
                color: "#ffffff"
            },
            "& $typography": {
                color: "#ffffff"
            },
            "& $typographyFlip": {
                color: "#2196f3"
            },
            "& $buttonFlip": {
                color: "#ffffff",
                visibility: 'visible',
            },
            "& $title": {
                backgroundColor: '#ffffff',
            },
        }
    },
    icon: () => ({
        color: 'secondary.main',
    }),
    typography: () => ({
        color: 'secondary.main',
    }),
    typographyFlip: () => ({
        color: 'primary.main',
    }),
    buttonFlip: () => ({
        color: 'primary.main',
    }),
    title: () => ({
        backgroundColor: 'secondary.main',
    })
  }));

  /**
   * 
   * @param {
   *    root - describes the directory of the file/folder viewer, if not supplied then will host the root folder
   *    viewLevel - values {0, 1,  and etc.} the level of visible features on the file/folder viewer
   * } 
   * @returns 
   */
export default function FilesRoot({id, expanded, root, viewLevel}) {
    id = id || 'AmperDrive';
    const app = React.useContext(AppContext);
    const location = useLocation();
    const {pathname} = location;
    const navigate = useNavigate();
    let directory = '/';
    if (root == null && pathname.indexOf(breadcrumbs.drive.files.path) > -1) {
        directory = pathname.replace(breadcrumbs.drive.files.path, '');
        directory = decodeURI(directory);
        if (directory.length < 1) {
            directory = '/';
        }
    } else if (root != null) {
        directory = root;
    }
    const debounceReload = (arg0, originalState) => {
        setState({
            ...stateRef.current,
            loading: true,
        });
    };

    const initialState = () => {
        return {
            loading: true,
            progress: false,
            data: [],
            newFolderOpen: false,
            discoverPopoverOpen: false,
            discoverPopoverAnchorEl: null,
            menuOpen: false,
            openWithOpen: false,
            openWithAnchorEl: null,
            selected: {},
            fileMenuY: 0,
            fileMenuX: 0,
            discoverPopoverLoading: false,
            discoverData: [],
            directory: directory,
            copy: {},
            cut: {},
            debounceReload: debounceLatest(debounceReload, 10000),
        };
    };
    const [state, setState] = useState(initialState);
    const [sx, setSx] = useState(6);
    const stateRef = useRef();
    stateRef.current= state;


    const myRef = useRef();
    const dropRef = useRef();
    const fileViewRef = useRef();

    useEffect(() => {
        const width = myRef.current.clientWidth;
        setSx(Math.ceil(width / 250));
    });

    useEffect(() => {
        if (state.loading || (directory != state.directory)) {
            getDataStore().load((result) => {
                setState({
                    ...state,
                    data: result.data || [],
                    loading: false,
                    progress: false,
                    directory: directory,
                });
            });
        }
    }, [state.loading, directory]);

    useEffect(() => {
        if (state.discoverPopoverLoading) {
            getDiscoverDataStore().load((result) => {
                setState({
                    ...state,
                    discoverData: result.data || [],
                    discoverPopoverLoading: false,
                });
            });
        }
    }, [state.discoverPopoverLoading]);
    
    if (app) {
        app.registerRefresh(id, () => {
            reload();
        });    
    }

    const getProgressBar = () => {
        return <LinearProgress sx={{mb: 1, mr: 6, visibility: state.loading || state.progress ? 'visible' : 'hidden'}}/>;
    };

    const reload = () => {
        setState({
            ...state,
            loading: true
        });
    };

    const addFile = (uploadDirectory, metadata) => {
        if (metadata && directory === uploadDirectory) {
            setState({
                ...stateRef.current,
                data: [...stateRef.current.data, metadata],
            });    
        }
    };

    const upload = (event) => {
        let files = event.currentTarget.files;
        app.upload(directory, files, ()=>{reload()}, false, null);
    };

    const upversion = (event) => {
        let files = event.currentTarget.files;
        app.upload(directory, files, ()=>{reload()}, true, state.contextFile);
    };

    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.myHost()}files-v1/fetch`,
            requestMethod: "POST",
            parameters: {
                root: directory,
            }
        });
    };

    const getDiscoverDataStore = () => {
        return new DataStore({
            url: `${HostManager.myHost()}files-v1/discover`,
            requestMethod: "POST",
        });
    };
    const classes = useStyles();
    
    const newFolder = () => {
        setState({
            ...state,
            newFolderOpen: true
        });
    };

    const handleNewFolderClose = () => {
        setState({
            ...state,
            newFolderOpen: false
        });
    };

    const inputFieldChange = (event) => {
        const {
            target: { value, name },
          } = event;
          setState({
            ...state,
            [name]: value
          });
    };

    const handleNewFolderSubmit = () => {
        if (state.folderName != null && state.folderName.length > 0) {
            setState({
                ...state,
                progress: true,
            });
            post(`${HostManager.myHost()}files-v1/newDir`, {
                name: state.folderName.trim(),
                root: directory,
            }, (result) => {
                if (result.success) {
                    setState({
                        ...state,
                        loading: true,
                        progress: false,
                        newFolderOpen: false,
                        discoverData: [],
                    });
                } else {
                    app.toast('warning', result.error)
                }
            }, (result) => {
                app.toast('warning', result.error)
            });
        }
    };
    
    const handleMoveFolders = (moveTo, files) => {
        const moveIds = [];
        const selectedEntries = Object.entries(files);
        for (let i = 0; i < selectedEntries.length; i++) {
            const entry = selectedEntries[i];
            if (entry[1] === true) {
                moveIds.push(entry[0]);
            }
        }
        if (moveIds.length > 0) {
            setState({
                ...state,
                progress: true,
            });
            post(`${HostManager.myHost()}files-v1/moveFiles`, {
                directory: moveTo,
                root: directory,
                ids: moveIds,
            }, (result) => {
                if (result.success) {
                    setState({
                        ...state,
                        loading: true,
                        progress: false,
                        selected: {},
                        discoverData: [],
                    });
                } else {
                    app.toast('warning', result.error)
                }
            }, (result) => {
                app.toast('warning', result.error)
            });
        }
    };

    const removeFileFolders = (fileId) => {
        const removeIds = [];
        const selectedEntries = Object.entries(state.selected);
        for (let i = 0; i < selectedEntries.length; i++) {
            const entry = selectedEntries[i];
            if (entry[1] === true) {
                removeIds.push(entry[0]);
            }
        }
        if (Convenience.hasValue(fileId) && !removeIds.includes(fileId)) {
            removeIds.push(fileId);
        }
        if (removeIds.length > 0) {
            setState({
                ...state,
                progress: true,
            });
            post(`${HostManager.myHost()}files-v1/removeFiles`, {
                root: directory,
                ids: removeIds,
            }, (result) => {
                if (result.success) {
                    setState({
                        ...state,
                        loading: true,
                        progress: false,
                        selected: {},
                        discoverData: [],
                        menuOpen: false,
                    });
                } else {
                    app.toast('warning', result.error)
                }
            }, (result) => {
                app.toast('warning', result.error)
            });
        }
    };

    const showFolderTreePopover = (event) => {
        setState({
            ...state,
            discoverPopoverOpen: true,
            discoverPopoverAnchorEl: event.currentTarget,
            discoverPopoverLoading: state.discoverData.length == 0,
        });
    };

    const closeFolderTreePopover = () => {
        setState({
            ...state,
            discoverPopoverOpen: false,
            discoverPopoverAnchorEl: null,
        });
    };

    const onFolderClick = (navDirectory) => {
        let basePath = pathname.trim().replace(/\/+$/, '');
        breadcrumbs.addCrumb(basePath + '/' + navDirectory, <FolderIcon/>)
        navigate(basePath + '/' + navDirectory, {state: {loading: true}});
        setState({
            ...state,
            selected: {},
            progress: true,
        });
    };

    const discoveryFolderClicked = (navDirectory, path) => {
        //remove the ending path seperator '/' to make sure equals catches paths
        const basePath = pathname.trim().replace(/\/+$/, '');
        const toPath = '/drive/files' + path.replace(/\/+$/, '');
        if (toPath !== decodeURI(basePath)) {
            navigate(toPath)
            setState({
                ...state,
                selected: {},
                progress: true,
            });    
        }
    };

    const validDirectoryName = () => {
        return state.folderName != null && state.folderName != '' && /^[^\s^\x00/^[^\s^\x00-\x1f\\?*:"";<>|\/.][^\x00-\x1f\\?*:"";<>|\/]*[^\s^\x00-\x1f\\?*:"";<>|\/.]+$/g.test(state.folderName.trim());
    };
    
    const fileMore = (event, id, name, isFile, processing, file) => {
        if (viewLevel > 0 || viewLevel == null || isFile) {
            event.stopPropagation();
            window.event.returnValue = false;
            event.preventDefault();
            setState({
                ...state,
                anchorEl: event.currentTarget,
                contextId: id,
                contextIsFile: (isFile == true),
                contextFile: file,
                contextProcessing: processing,
                contextName: name,
                menuOpen: true,
                fileMenuX: event.clientX,
                fileMenuY: event.clientY,
            });
        }
    };

    const hideDefaultContext = (event) => {
        window.event.returnValue = false;
        event.preventDefault();
    }

    const fileMoreClose = (event) => {
        event.stopPropagation();
        setState({
            ...state,
            menuOpen: false,
            contextId: null,
            contextName: null,
        });
    };

    function FolderNameHelperText({valid}) {
        const helperText = useMemo(() => {
            if (!valid) {
                return "Folder name is required and can't start or end with '.', ' ', '__progress__', '__file__' or contain any of the following characters '*, \\, :, \", /, >, <, ?, |'";
            }
        }, [valid]);
        
        return <FormHelperText>{helperText}</FormHelperText>;
    };

    let dragTargetClone = null;
    const onDragStart = (event, file) => {
        if (dragTargetClone != null) {
            dragTargetClone.remove();
        }
        const selected = state.selected;
        const identifier = file.isDir === true ? file.name : file.id;
        const selectedCopy = JSON.parse(JSON.stringify(selected));
        selectedCopy[identifier] = true;
        event.dataTransfer.setData("file", JSON.stringify(selectedCopy));
        const selectedEntries = Object.entries(selectedCopy);
        if (selectedEntries.length < 2) {
            if (event.currentTarget.style.backgroundImage !== '') {
                dragTargetClone = event.currentTarget.cloneNode(false);
                dragTargetClone.style.width = '100px';
                dragTargetClone.style.height = '100px';
            } else {
                dragTargetClone = event.currentTarget.children[0].cloneNode(true);
                dragTargetClone.children[0].remove();
                dragTargetClone.style.width = '100px';
                dragTargetClone.style.height = '100px';
            }
            document.body.appendChild(dragTargetClone);
        } else {
            dragTargetClone = document.createElement("span", {id: 'cloneTempItem'});
            dragTargetClone.style.position = 'relative';
            let position = 0;
            for (let i = 0 ; i < selectedEntries.length; i++) {
                if (i == 3) {
                    const moreItems = selectedEntries.length - i;
                    if (moreItems > 0) {
                        const moreItemsElement = document.createElement("span", {id: 'moreItems'});
                        moreItemsElement.style.width = '95px';
                        moreItemsElement.style.position = 'absolute';
                        moreItemsElement.style.top = position + 50 + 'px';
                        moreItemsElement.style.left = position + 'px';
                        moreItemsElement.style.backgroundColor = '#ffffff';
                        moreItemsElement.style.borderRadius = '5px';
                        moreItemsElement.style.padding = '3px';
                        moreItemsElement.innerHTML = '+' + moreItems + ' ' +(moreItems > 1 ? ' more items' : 'more item')
                        dragTargetClone.appendChild(moreItemsElement);
                    }
                    break;
                }
                const identifier = selectedEntries[i];
                const targetElement = document.getElementById(identifier[0])
                let targetElementClone = null;
                if (targetElement.style.backgroundImage !== '') {
                    targetElementClone = targetElement.cloneNode(false);
                } else {
                    targetElementClone = targetElement.children[0].cloneNode(true);
                    targetElementClone.children[0].remove();
                    targetElementClone.children[0].remove();
                    targetElementClone.style.backgroundColor = '#ffffff'
                    targetElementClone.style.borderRadius = '5px'
                }
                targetElementClone.style.width = '100px';
                targetElementClone.style.height = '100px';
                targetElementClone.style.position = 'absolute';
                targetElementClone.style.top = position + 'px';
                targetElementClone.style.left = position + 'px';
                dragTargetClone.appendChild(targetElementClone);
                position += 10;
            }
            document.body.appendChild(dragTargetClone);
        }
        event.dataTransfer.setDragImage(dragTargetClone, 0, 0)
    };

    let uploadDragOverCounter = 0;
    const onUploadDragEnter = (event) => {
        event.stopPropagation();
        event.preventDefault();
        uploadDragOverCounter++;
        let items = event.dataTransfer.items;
        if (items != null && items.length > 0 && items[0].kind === 'file') {
            const dropIndicator = dropRef.current;
            dropIndicator.style.display = "inline";
            
        }
    };
    const onUploadDragOver = (event) => {
        event.stopPropagation();
        event.preventDefault();
    };
    const onUploadDragLeave = (event) => {
        event.stopPropagation();
        event.preventDefault();
        uploadDragOverCounter--
        let items = event.dataTransfer.items;
        if (uploadDragOverCounter == 0 && items != null && items.length > 0 && items[0].kind === 'file') {
            const dropIndicator = dropRef.current;
            dropIndicator.style.display = "none";
            
        }
    };
    const onUploadDragEnd = (event) => {
        event.stopPropagation();
        event.preventDefault();
        uploadDragOverCounter = 0;
        const dropIndicator = dropRef.current;
        dropIndicator.style.display = "none";
    };

    const onUploadDrop = (event) => {
        event.stopPropagation();
        event.preventDefault();
        let files = event.dataTransfer.files;
        if (files != null && files.length > 0) {
            app.upload(directory, files, ()=>{reload()}, false, null);
        }
        uploadDragOverCounter = 0;
        const dropIndicator = dropRef.current;
        dropIndicator.style.display = "none";
    };

    let lastTarget = null;
    const onDrop = (event) => {
        event.stopPropagation();
        event.preventDefault();
        resetFolderStyles(event);
        let files = event.dataTransfer.files;
        if (files != null && files.length > 0) {
            let uploadDirectory = directory;
            if (!uploadDirectory.endsWith('/')) {
                uploadDirectory += '/';
            }
            app.upload(uploadDirectory + event.currentTarget.id, files, ()=>{reload()}, false, null);
        } else {
            const fileJson = event.dataTransfer.getData("file");
            if (fileJson != null && fileJson.length > 0) {
                const files = JSON.parse(fileJson);
                handleMoveFolders(event.currentTarget.id, files)
            }    
        }
        const dropIndicator = dropRef.current;
        dropIndicator.style.display = "none";
    };

    const onDragOver = (event) => {
        event.preventDefault();
        if (event.currentTarget.clientWidth == 190) {
            resetLastTarget();
            event.currentTarget.style.width = (event.currentTarget.clientWidth - 6) + 'px';
            event.currentTarget.style.height = (event.currentTarget.clientHeight - 6) + 'px';            
            //event.currentTarget.zIndex = 2;
            event.currentTarget.classList.add(classes.onDropTarget);
            lastTarget = event.currentTarget;
        }
    };

    let dragOverCounter = 0;
    const onDragEnter = (event) => {
        dragOverCounter++;
    };

    const onDragLeave = (event) => {
        dragOverCounter--;
        event.preventDefault();
        if (dragOverCounter === 0) {
            resetFolderStyles(event);
        }
    };

    const onDragEnd = (event) => {
        dragOverCounter = 0;
        resetLastTarget();
        if (dragTargetClone != null) {
            dragTargetClone.remove();
        }
    };

    const resetFolderStyles = (event) => {
        event.currentTarget.style.width = '190px';
        event.currentTarget.style.height = '190px';
        //event.currentTarget.zIndex = 1;
        event.currentTarget.classList.remove(classes.onDropTarget);
    };

    const resetLastTarget = () => {
        if (lastTarget != null) {
            lastTarget.style.width = '190px';
            lastTarget.style.height = '190px';
            lastTarget.classList.remove(classes.onDropTarget);
        }
    };

    const getNewFolderDialog = () => {
        const valid = validDirectoryName();
        return (
            <Dialog open={state.newFolderOpen} onClose={handleNewFolderClose}>
            <DialogTitle>New Folder</DialogTitle>
            <DialogContent>
              <DialogContentText>
                To add a new folder, provide the name below to be displayed on your brand new directory
              </DialogContentText>
              <FormControl variant="standard" fullWidth>
                <TextField
                    onChange={inputFieldChange}
                    autoFocus
                    margin="dense"
                    name="folderName"
                    label="Folder name"
                    type="text"
                    fullWidth
                    error={!valid}
                    variant="standard"></TextField>
                <FolderNameHelperText valid={valid}/>
              </FormControl>
              
            </DialogContent>
            <DialogActions>
              <Button onClick={handleNewFolderClose}>Cancel</Button>
              <Button disabled={!valid} onClick={handleNewFolderSubmit}>Ok</Button>
            </DialogActions>
          </Dialog>
        );
    };

    const removeFile = () => {
        if (Convenience.hasValue(state.contextId)) {
            removeFileFolders(state.contextId);
        }
    };
    
    const downloadFile = () => {
        if (Convenience.hasValue(state.contextId)) {
            setState({
                ...state,
                menuOpen: false,
            });
            download(`${HostManager.myHost()}files-v1/download?id=` + encodeURIComponent(state.contextId) + '&major=' + state.contextFile.version.major+ '&minor=' + state.contextFile.version.minor + '&patch=' + state.contextFile.version.patch + '&root=' + encodeURIComponent(directory));
        }
    };

    const downloadFileRendition= () => {
        if (Convenience.hasValue(state.contextId)) {
            setState({
                ...state,
                menuOpen: false,
            });
            download(`${HostManager.myHost()}files-v1/download?id=` + encodeURIComponent(state.contextId) + '&major=' + state.contextFile.version.major+ '&minor=' + state.contextFile.version.minor + '&patch=' + state.contextFile.version.patch + '&root=' + encodeURIComponent(directory) + '&rendition=true');
        }
    };

    const pasteItemsNumber = Object.keys(state.copy).length + Object.keys(state.cut).length;
    const pasteFile = () => {
        if (pasteItemsNumber > 0) {
            setState({
                ...state,
                progress: true,
            });
            post(`${HostManager.myHost()}files-v1/pasteFiles`, {
                root: directory,
                copy: JSON.stringify(state.copy),
                cut: JSON.stringify(state.cut),
            }, (result) => {
                if (result.success) {
                    setState({
                        ...state,
                        loading: true,
                        progress: false,
                        selected: {},
                        discoverData: [],
                        cut: {},
                        copy: {},
                        menuOpen: false,
                    });
                } else {
                    app.toast('warning', result.error)
                }
            }, (result) => {
                app.toast('warning', result.error)
            });
        }
    };
    
    const copyFile = () => {
        const copy = {
        };
        if (Convenience.hasValue(state.contextId) && state.contextId !== './') {
            copy[state.contextId] = directory;
            const selectedEntries = Object.entries(state.selected);
            for (let i = 0; i < selectedEntries.length; i++) {
                const entry = selectedEntries[i];
                if (entry[1] === true) {
                    copy[entry[0]] = directory;
                }
            }
        } else if (state.contextId === './') {
            for (let i = 0; i < state.data.length; i++) {
                const item = state.data[i];
                const itemIdentifier = item.isDir === true ? item.name : item.id;
                copy[itemIdentifier] = directory;
            }
        }
        setState({
            ...state,
            copy,
            cut: {},
            menuOpen: false,
        });
    };

    const cutFile = () => {
        const cut = {
        };
        if (Convenience.hasValue(state.contextId) && state.contextId !== './') {
            cut[state.contextId] = directory;
            const selectedEntries = Object.entries(state.selected);
            for (let i = 0; i < selectedEntries.length; i++) {
                const entry = selectedEntries[i];
                if (entry[1] === true) {
                    cut[entry[0]] = directory;
                }
            }
        } else if (state.contextId === './') {
            for (let i = 0; i < state.data.length; i++) {
                const item = state.data[i];
                const itemIdentifier = item.isDir === true ? item.name : item.id;
                cut[itemIdentifier] = directory;
            }
        }
        setState({
            ...state,
            cut,
            copy: {},
            menuOpen: false,
        });
    };

    const handleOpenWithClose = () => {
        setState({
            ...state,
            openWithAnchorEl: null,
            openWithOpen: false,
        });
    };

    const handleOpenWithOpen = (event) => {
        setState({
            ...state,
            openWithAnchorEl: event.currentTarget,
            openWithOpen: true,
        });
    };

    const previousRendition = (metadata) => {
        for (let l = 0; l < state.data.length; l++) {
            if (metadata.id === state.data[l].id) {
                for (let i = l - 1; i >= 0; i--) {
                    if (state.data[i].rendition === true || state.data[i].viewable === true) {
                        fileViewRef.current.view(directory, state.data[i])
                        return;    
                    }
                }
            }
        }
    };
    
    const nextRendition = (metadata) => {
        for (let l = 0; l < state.data.length; l++) {
            if (metadata.id === state.data[l].id) {
                for (let i = l + 1; i < state.data.length; i++) {
                    if (state.data[i].rendition === true || state.data[i].viewable === true) {
                        fileViewRef.current.view(directory, state.data[i])
                        return;    
                    }
                }
            }
        }
    };

    const openWithMenu = () => {
        
        return (<Menu
        anchorEl={state.openWithAnchorEl}
        open={state.openWithOpen}
        onClose={handleOpenWithClose}
        sx={{ml: '1px'}}
        anchorOrigin={{
            vertical: 'top',
            horizontal: 'right',
        }}
        transformOrigin={{
            vertical: 'top',
            horizontal: 'left',
        }}>
            <MenuItem onClick={() => {fileViewRef.current.view(directory, state.contextFile, 'ADOBE_DC')}} disabled={state.contextFile == null || (state.contextFile.type !== "application/pdf" && state.contextFile.renditionType !== "application/pdf")}>
            <ListItemIcon>
                <img src="/images/adobe_console_logo.svg" width="20px" style={{marginLeft: '2px'}}/>
            </ListItemIcon>
            <ListItemText>Adobe DC</ListItemText>
            </MenuItem>
    </Menu>)};
    
    const contextMenueItems = [];

    if (state.contextIsFile) {
        if (viewLevel > 0) {
            contextMenueItems.push(<MenuItem key="open" onClick={() => {fileViewRef.current.view(directory, state.contextFile)}}>
                <ListItemIcon>
                    <PreviewIcon color="primary" fontSize="medium" />
                </ListItemIcon>
                <ListItemText>Open</ListItemText>
            </MenuItem>);
        contextMenueItems.push(<MenuItem key="openWith" onClick={handleOpenWithOpen} disabled={(!state.contextIsFile || state.contextProcessing)}>
                <ListItemIcon>
                    <OpenInNewIcon color="primary" fontSize="medium" />
                </ListItemIcon>
                <ListItemText>Open with</ListItemText>
                <ListItemIcon>
                    <ArrowForwardIosIcon color="primary" fontSize="medium" sx={{ml: 3}}/>
                </ListItemIcon>
            </MenuItem>);
        contextMenueItems.push(<MenuItem key="upversion" disabled={!state.contextIsFile} component="label">
                <input hidden accept="*" type="file" onChange={upversion}/>
                <ListItemIcon>
                    <FileUploadIcon color="primary" fontSize="medium" />
                </ListItemIcon>
                <ListItemText>Upversion</ListItemText>
            </MenuItem>);
        }
        
        contextMenueItems.push(<MenuItem key="download" onClick={downloadFile} disabled={!state.contextIsFile}>
            <ListItemIcon>
                <DownloadIcon color="primary" fontSize="medium" />
            </ListItemIcon>
            <ListItemText>Download {(state.contextFile && state.contextFile.type && state.contextFile.type.indexOf('/') > 0) ? state.contextFile.type.substring(state.contextFile.type.indexOf('/') + 1, state.contextFile.type.length) : (state.contextFile && state.contextFile.type) ? state.contextFile.type : ''}</ListItemText>
        </MenuItem>);
    }
    if (state.contextFile && state.contextFile.rendition) {
        contextMenueItems.push(<MenuItem key="downloadRendition" onClick={downloadFileRendition} disabled={!state.contextIsFile}>
        <ListItemIcon>
            <DownloadIcon color="primary" fontSize="medium" />
        </ListItemIcon>
        <ListItemText>Download {(state.contextFile && state.contextFile.renditionType && state.contextFile.renditionType.indexOf('/') > 0) ? state.contextFile.renditionType.substring(state.contextFile.renditionType.indexOf('/') + 1, state.contextFile.type.length) : (state.contextFile && state.contextFile.renditionType) ? state.contextFile.renditionType : ''}</ListItemText>
        </MenuItem>);
    }
    if (viewLevel > 0) {
        /* For now disable the copy feature, since it is not useful and very hard to implement
        contextMenueItems.push(<MenuItem key="copy" onClick={copyFile} disabled={state.contextProcessing}>
            <ListItemIcon>
                <FileCopyIcon color="primary" fontSize="medium" />
            </ListItemIcon>
            <ListItemText>Copy</ListItemText>
        </MenuItem>);*/
        contextMenueItems.push(<MenuItem key="cut" onClick={cutFile} disabled={state.contextProcessing}>
            <ListItemIcon>
                <ContentCutIcon color="primary" fontSize="medium" />
            </ListItemIcon>
            <ListItemText>Cut</ListItemText>
        </MenuItem>);

        if (state.contextId === './') {
            contextMenueItems.push(<MenuItem key="paste" onClick={pasteFile} disabled={pasteItemsNumber < 1}>
                <ListItemIcon>
                    <ContentPasteIcon color="primary" fontSize="medium" />
                </ListItemIcon>
                <ListItemText>Paste {pasteItemsNumber > 0 ? ('(' + pasteItemsNumber + (pasteItemsNumber == 1 ? ' Item' : ' Items') + ')') : ''}</ListItemText>
            </MenuItem>);
            contextMenueItems.push(<MenuItem key="newFolder" onClick={newFolder}>
                <ListItemIcon>
                    <ContentPasteIcon color="primary" fontSize="medium" />
                </ListItemIcon>
                <ListItemText>New folder</ListItemText>
            </MenuItem>);
        }
        contextMenueItems.push(<MenuItem key="remove" onClick={removeFile} disabled={state.contextId === './'}>
            <ListItemIcon>
                <DeleteIcon color="primary" fontSize="medium" />
            </ListItemIcon>
            <ListItemText>Remove</ListItemText>
        </MenuItem>);
    }
    const fileContextMenu = (<Menu
        anchorReference="anchorPosition"
        anchorPosition={{ top: state.fileMenuY, left: state.fileMenuX }}
        open={state.menuOpen}
        onClose={fileMoreClose}
        onContextMenu={hideDefaultContext}
        MenuListProps={{
        'aria-labelledby': 'basic-button',
        }}>
        {contextMenueItems}
    </Menu>);

    const select = (event, file) => {
        event.stopPropagation();
        const selected = state.selected;
        const identifier = file.isDir === true ? file.name : file.id;
        if (selected[identifier] === true) {
            if (event.shiftKey) {
                let selectToIndex = null;
                for (let i = 0; i < state.data.length; i++) {
                    const item = state.data[i];
                    const itemIdentifier = item.isDir === true ? item.name : item.id;
                    if (itemIdentifier === identifier) {
                        selectToIndex = i;
                        break;
                    }
                }
                if (selectToIndex != null) {
                    for (let l = 0; l <= selectToIndex; l++) {
                        const item = state.data[l];
                        const itemIdentifier = item.isDir === true ? item.name : item.id;
                        delete selected[itemIdentifier];
                    }
                }
            } else {
                delete selected[identifier];
            }
        } else {
            if (event.shiftKey) {
                let lastSelectedIndex = null;
                let selectToIndex = null;
                for (let i = 0; i < state.data.length; i++) {
                    const item = state.data[i];
                    const itemIdentifier = item.isDir === true ? item.name : item.id;
                    if (selected[itemIdentifier] == true) {
                        lastSelectedIndex = i;
                    }
                    if (itemIdentifier === identifier) {
                        selectToIndex = i;
                        break;
                    }
                }
                if (lastSelectedIndex != null &&  selectToIndex != null) {
                    for (let l = lastSelectedIndex; l <= selectToIndex; l++) {
                        const item = state.data[l];
                        const itemIdentifier = item.isDir === true ? item.name : item.id;
                        selected[itemIdentifier] = true;
                    }
                } else {
                    selected[identifier] = true;
                }
            } else {
                selected[identifier] = true;
            }
        }
        
        setState({
            ...state,
            selected,
        });
    };
    const fileClasses = fileStyles();
    const thumbnailClasses = thumbnailStyles();

    const getFileThumbnail = (file, type) => {
        return <Grid size={1} key={file.id} style={{opacity: (state.cut[file.id] != null) ? 0.5 : 1}}>
            <Paper onDoubleClick={() => {fileViewRef.current.view(directory, file)}} onDragStart = {(event) => {onDragStart(event, file)}} onDragEnd={onDragEnd} draggable="true" id={file.id} sx={{width: 190, height: 190, overflowWrap: 'break-word', borderRadius: 2, borderColor: 'primary.main'}} onContextMenu={(event) => {fileMore(event, file.id, file.name, true, false, file)}}
             className={thumbnailClasses.box} elevation={5} style={{backgroundImage: 'url(data:image/png;base64,' + file.thumbnailImage + ')', backgroundRepeat: 'no-repeat', backgroundPosition: 'top'}}>
                <Box sx={{width: 190, height: 190, position: 'relative'}}>
                    <Radio style={{position: 'absolute', left: 0, top: 0}} size="medium" color="primary" onClick={(event)=> select(event, file)} checked={state.selected[file.id] === true}/>
                    <Button onClick={(event) => {fileMore(event, file.id, file.name, true, false, file)}} color="secondary" title="More" style={{width: '10px', borderRadius: 15, position: 'absolute', right: -10, top: 0}}>
                        <MoreHorizIcon sx={{mr: 0, height: '30px',}} className={thumbnailClasses.buttonFlip}/>
                    </Button>
                    <Box variant="outlined" sx={{height: 20, width: 190, position: 'absolute', bottom: 0}} style={{borderBottomRightRadius: 5, borderBottomLeftRadius: 5}} className={thumbnailClasses.title}>
                        <Tooltip title={file.name}>
                            <Typography variant="body2" color="text.secondary" sx={{ml: 1, mr: 1}} className={classes.multiLineEllipsis}>{file.name}</Typography>
                        </Tooltip>
                    </Box>
                </Box>
            </Paper>
        </Grid>;
    };

    const getFileProgress = () => {
        return (
          <Box sx={{ position: 'relative', display: 'inline-flex', width: '100%', alignItems: 'center', justifyContent: 'center' }}>
            <CircularProgress color="primary" disableShrink/>
            <Box
              sx={{
                top: 0,
                left: 0,
                bottom: 0,
                right: 0,
                position: 'absolute',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography variant="body2" component="div" color="text.secondary" className={[classes.multiLineEllipsis, fileClasses.typography].join(' ')}>
                Rendering...
              </Typography>
            </Box>
          </Box>
        );
      };

    const getFile = (file, type) => {
        return <Grid size={1} key={file.id} style={{opacity: (state.cut[file.id] != null) ? 0.5 : 1}}>
            <Paper onDragStart = {(event) => {onDragStart(event, file)}} onDragEnd={onDragEnd} draggable="true" id={file.id} sx={{width: 190, height: 190, overflowWrap: 'break-word', borderRadius: 2, borderColor: 'primary.main'}} onContextMenu={(event) => {fileMore(event, file.id, file.name, true, file.processing, file)}} className={fileClasses.box} elevation={5}>
                <Box variant="outlined" sx={{display: 'flex',  alignItems: 'center', justifyContent: 'center', position: "relative",}}>
                    <Radio style={{position: 'absolute', left: 0, top: 0}} size="medium" color="primary" onClick={(event)=> select(event, file)} checked={state.selected[file.id] === true} className={fileClasses.buttonFlip}/>
                    <Button onClick={(event) => {fileMore(event, file.id, file.name, true, file.processing, file)}} color="secondary" title="More" style={{width: '10px', borderRadius: 15, position: 'absolute', right: -10, top: 0}}>
                        <MoreHorizIcon sx={{mr: 0, height: '30px',}} className={fileClasses.buttonFlip}/>
                    </Button>
                    <InsertDriveFileIcon color="primary" sx={{width: 100, height: 100, mt: 2}} className={fileClasses.icon} />
                    <Typography variant="caption" color="secondary.main" className={classes.multiLineEllipsis + ' ' + fileClasses.typographyFlip} sx={{ml: 1, mr: 1, position: 'absolute', top: 60}}>{type.toUpperCase()}</Typography>
                </Box>
                {file.processing == true ? getFileProgress() : ''}
                <Box variant="outlined" sx={{display: 'flex',  alignItems: 'center', justifyContent: 'center', height: (file.processing == true ? 'calc(100% - 150px)' : 'calc(100% - 80px)')}}>
                    <Tooltip title={file.name}>
                        <Typography variant="body2" color="text.secondary" className={[classes.multiLineEllipsis, fileClasses.typography].join(' ')} sx={{ml: 1, mr: 1}}>{file.name}</Typography>
                    </Tooltip>
                </Box>
            </Paper>
        </Grid>;
    };

    const 
    
    getFilesAndDirectories = () => {
        let hasProcessingFile = false;
        const result = [];
        for (let i = 0; i < stateRef.current.data.length; i++) {
            const file = stateRef.current.data[i];
            if (file.isDir === true) {
                result.push(
                    <Grid key={file.name} size={1} onClick={(()=>{onFolderClick(file.name)})} style={{opacity: (state.cut[file.name] != null) ? 0.5 : 1}}>
                        <Paper id={file.name} onDragStart = {(event) => {onDragStart(event, file)}} onDragEnd={onDragEnd} draggable="true" droppable="true" onDrop={onDrop} onDragOver={onDragOver} onDragEnter={onDragEnter} onDragLeave={onDragLeave} onContextMenu={(event) => {fileMore(event, file.name, file.name, false, false, file)}} sx={{width: 190, height: 190, overflowWrap: 'break-word', borderRadius: 2, borderColor: 'primary.main'}} className={fileClasses.box} style={{zIndex: 1, position: 'relative'}} elevation={5}>
                            <Box variant="outlined" sx={{display: 'flex',  alignItems: 'center', justifyContent: 'center', position: "relative",}}>
                                <Radio style={{position: 'absolute', left: 0, top: 0}} size="medium" color="primary" onClick={(event)=> select(event, file)} checked={state.selected[file.name] === true} className={fileClasses.buttonFlip}/>
                                <Button onClick={(event) => {fileMore(event, file.name, file.name, false, false, file)}} color="secondary" title="More" style={{width: '10px', borderRadius: 15, position: 'absolute', right: -10, top: 0}}>
                                    <MoreHorizIcon sx={{mr: 0, height: '30px',}} className={fileClasses.buttonFlip}/>
                                </Button>
                                <FolderIcon color="primary" sx={{width: 100, height: 100, mt: 2}} className={fileClasses.icon}/>
                            </Box>
                            <Box variant="outlined" sx={{display: 'flex',  alignItems: 'center', justifyContent: 'center', height: 'calc(100% - 80px)'}}>
                                <Tooltip title={file.name}>
                                    <Typography variant="body2" color="text.secondary" className={[classes.multiLineEllipsis, fileClasses.typography].join(' ')} sx={{ml: 1, mr: 1}}>{file.name}</Typography>
                                </Tooltip>
                            </Box>
                        </Paper>
                    </Grid>
                );
            }
        }
        for (let i = 0; i < stateRef.current.data.length; i++) {
            const file = stateRef.current.data[i];
            if (file.processing === true) {
                hasProcessingFile = true;
            }
            if (file.isDir != true) {
                const index = file.name.lastIndexOf('.');
                let type = null;
                if (index > 0) {
                    type = file.name.substring(file.name.lastIndexOf('.') + 1, file.name.length);
                } else if (file.type != null && file.type.length > 0) {
                    type = file.type;
                } else {
                    type = 'UNKNOWN';
                }
                
                result.push(
                    file.thumbnail === true ? getFileThumbnail(file, type) : getFile(file, type)
                );
            }
        }
        if (hasProcessingFile) {
            stateRef.current.debounceReload(null, stateRef.current);
        }
        return result;
    };

    const renderDiscover = (folder) => {
        if (!folder.path) {
            return;
        }
        const label = <Typography sx={{p: '4px'}}>{folder.name}</Typography>;
        return <TreeItem key={folder.name} nodeId={folder.path} 
                label={label} icon={<FolderIcon color='primary' style={{fontSize: 26}}/>}
                onClick={() => {discoveryFolderClicked(folder.name, folder.path)}}>
          {Array.isArray(folder.folders)
            ? folder.folders.map((folder) => renderDiscover(folder))
            : null}
        </TreeItem>
      };

    const getAllFolderNames = (folder, result) => {
        result.push(folder.path);
        if (Array.isArray(folder.folders)) {
            folder.folders.map((folder) => getAllFolderNames(folder, result))
        }
        return result;
    };

    const getDiscoverView = () => {
        const allFolders = getAllFolderNames(state.discoverData, []);
        return (<TreeView
                    onNodeToggle={(event)=>{event.stopPropagation(); return false;}}
                    aria-label="Discover"
                    defaultCollapseIcon={null}
                    defaultExpanded={allFolders}
                    expanded={allFolders}
                    defaultExpandIcon={null}
                    sx={{ maxHeight: 600, flexGrow: 1, maxWidth: 400, overflowY: 'auto', overflowX: 'hidden' }}>
            {renderDiscover(state.discoverData)}
          </TreeView>);
    };

    const getToolBar = () => {
        if (viewLevel < 1) {
            return;
        }
        return <Box sx={{
            display: 'flex',
            flexDirection: 'row',
            bgcolor: 'background.paper',
            }}>
                <Box sx={{ flexGrow: 1 }}>
                    <Stack direction="row" spacing={1} sx={{ mb: 0 }}>
                        <Button component="label" size="medium" startIcon={<FileUploadIcon/>}>
                            Upload
                            <input hidden accept="*" multiple type="file" onChange={upload}/>
                        </Button>
                        <Button component="label" size="medium" startIcon={<CreateNewFolderIcon/>} onClick={newFolder}>
                            New folder
                        </Button>
                        <Button component="label" size="medium" disabled={Object.keys(state.selected).length < 1} startIcon={<DeleteIcon/>} onClick={()=>removeFileFolders()}>
                            Remove
                        </Button>
                    </Stack>
                </Box>
                <Box sx={{ flexGrow: 0 }}>
                    <Stack direction="row" spacing={1} sx={{ mb: 0 }}>
                        <Button component="label" size="medium" endIcon={<AccountTreeIcon/>} onClick={showFolderTreePopover} onMouseOver={showFolderTreePopover}>
                            Navigate
                        </Button>
                    </Stack>
                </Box>
            </Box>;
    };

    return (<Box ref={myRef} sx={{ height: '100%', width: 'calc(100% - 25px)'}}>
            {getToolBar()}
            <Popover
                    id={'discoverPopover'}
                    open={state.discoverPopoverOpen}
                    anchorEl={state.discoverPopoverAnchorEl}
                    onClose={closeFolderTreePopover}
                    anchorReference="anchorPosition"
                    anchorPosition={{ top: 150, left: window.innerWidth }}
                >
                    <Box sx={{
                        minHeight: 100,
                        minWidth: 100,
                        maxHeight: 600,
                        maxWidth: 400,
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                    }}>
                        {state.discoverPopoverLoading ? <CircularProgress/> : getDiscoverView()}
                    </Box>
                    
            </Popover>
        {getProgressBar()}
        {getNewFolderDialog()}
        {fileContextMenu}
        {openWithMenu()}
        <Box sx={{ overflowX: 'hidden', overflowY: 'auto', pl: '5px', pt: 1}} width="100%" height={'calc(100% - 60px)'} onContextMenu={(event) => {fileMore(event, './', './', false, false, null)}} droppable="true" onDrop={onUploadDrop} onDragLeave={onUploadDragLeave} onDragEnter={onUploadDragEnter} onDragOver={onUploadDragOver} onDragEnd={onUploadDragEnd}>
            <Grid key={state.directory} container columns={sx} spacing={3} width="100%">
                {getFilesAndDirectories()}
            </Grid>
            <div ref={dropRef} style={{width: '200px', height: '200px', display: 'none', position: 'absolute', left: 'calc(50% + 100px)', top: 'calc(50% + 100px)', borderRadius: '25px', transform: 'translate(-50%, -50%)',}}>
               <CloudUploadIcon color="primary" sx={{width: 200, height: 200}}/>
            </div>
        </Box>
        <FileView ref={fileViewRef} next={nextRendition} previous={previousRendition}/>
    </Box>);
}
