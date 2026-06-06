import React, { useState, useImperativeHandle, forwardRef, useEffect, useMemo } from 'react';
import Slide from '@mui/material/Slide';
import { Snackbar } from '@mui/material';
import Paper from '@mui/material/Paper';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import IconButton from '@mui/material/IconButton';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import Typography from '@mui/material/Typography';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ListItem from '@mui/material/ListItem';
import LinearProgress from '@mui/material/LinearProgress';
import makeStyles from '@mui/styles/makeStyles';
import {post, postBlob} from '../../data/Submit';
import HostManager from '../../../HostManager';

const LINES_TO_SHOW = 1;
const useStyles = makeStyles({
    multiLineEllipsis: {
        overflow: "hidden",
        textOverflow: "ellipsis",
        display: "-webkit-box",
        "-webkit-line-clamp": LINES_TO_SHOW,
        "-webkit-box-orient": "vertical"
    }
});

function FilesBar({}, ref) {

    const initialState = () => {
        return {
            open: false,
            expanded: false,
            processing: false,
        };
    };
    const [state, setState] = useState(initialState);
    const cache = useMemo(() => { return {
        progressFiles: [],
        addProgressFile: (progressFile) => {
            cache.progressFiles.push(progressFile);
        },
    };}, [true])
    const [expanded, setExpanded] = useState(true);
    const memo = useMemo(() => { return {expanded: state.expanded};})

    useImperativeHandle(ref, () => ({
        upload(directory, files, callback, upversion, metadata) {
            if (files != null) {
                for (let i = 0; i < files.length; i++) {
                    const progressFile = {
                        file: files[i],
                        progress: 0,
                        complete: false,
                        index: 0,
                        directory: directory,
                        callback: callback,
                        upversion: (upversion == true),
                        originalMetadata: metadata,
                    };
                    cache.addProgressFile(progressFile);
                }
                setState({
                    ...state,
                    open: true,
                    expanded: true,
                    processing: true,
                });
            }
        }
    }));

    const CHUNK_SIZE = 1024*1000;
    const uploadFiles = () => {
        let index = 0;
        let progressFile = null;
        const progressFiles = cache.progressFiles;
        for (index = 0; index < progressFiles.length; index++) {
            if (progressFiles[index].complete === false && progressFiles[index].failed != true) {
                progressFile = progressFiles.splice(index, 1)[0];
                break;
            }
        }
        if (progressFile != null && progressFile.file != null) {
            const file = progressFile.file;
            const fileReader = new FileReader();
            fileReader.onload = () => {
                const chunk = fileReader.result;

                const api = progressFile.upversion ? 'files-v1/upversion' : 'files-v1/upload';

                const formData = new FormData()
                formData.append('chunk', new Blob([chunk]));
                formData.append('name', file.name);
                formData.append('size', file.size);
                formData.append('type', file.type);
                formData.append('directory', progressFile.directory);

                if (progressFile.upversion) {
                    formData.append('id', progressFile.originalMetadata.id);
                    formData.append('major', progressFile.originalMetadata.version.major);
                    formData.append('minor', progressFile.originalMetadata.version.minor);
                    formData.append('patch', progressFile.originalMetadata.version.patch);
                    if (progressFile.metadata != null) {
                        formData.append('newId', progressFile.metadata.id);
                    }
                } else {
                    if (progressFile.metadata != null) {
                        formData.append('id', progressFile.metadata.id);
                    }
                }
                postBlob(`${HostManager.amperHost()}${api}`, formData, (result) => {
                    if (result.success) {
                        progressFile.metadata = result.metadata;
                        progressFile.index = progressFile.index + chunk.byteLength;

                        progressFile.progress = (progressFile.index / file.size) * 100;
                        progressFile.complete = progressFile.index >= file.size;
                        cache.progressFiles.splice(index, 0, progressFile);
                        
                        if (progressFile.callback && progressFile.index >= file.size) {
                            progressFile.callback(progressFile.directory, result.metadata);
                        }
                        setState({
                            ...state,
                            processing: true,
                            expanded: memo.expanded,
                        });
                    } else {
                        progressFile.failed = true;
                        progressFile.error = result.error
                        cache.progressFiles.splice(index, 0, progressFile);
                        setState({
                            ...state,
                            processing: true,
                            expanded: memo.expanded,
                        });
                    }
                }, (result) => {
                    progressFile.failed = true;
                    progressFile.error = result.error
                    cache.progressFiles.splice(index, 0, progressFile);
                    setState({
                        ...state,
                        processing: true,
                        expanded: memo.expanded,
                    });
                });
            };
            const seek = () => {
                let slice = file.slice(progressFile.index, progressFile.index + CHUNK_SIZE);
                fileReader.readAsArrayBuffer(slice);    
            };
            seek();
        } else {
            setState({
                ...state,
                processing: false,
                open: false,
            });
        }
    };

    useEffect(() => {
        if (state.processing) {
            setTimeout(uploadFiles, 1000);
            setState({
                ...state,
                processing: false,
            });
        }
    }, [state.processing]);

    const handleClose = () => {
        
    };

    const classes = useStyles();
    const getFile = (name, progress, index) => {
        return <Box sx={{ display: 'flex', alignItems: 'center' }} key={index}>
            <Box sx={{ minWidth: 138, maxWidth: 138, width: 138, pr: '2px'}}>
                <Typography variant="body2" color="text.secondary" className={classes.multiLineEllipsis}>
                    {name}
                </Typography>
            </Box>
            <Box sx={{ minWidth: 100, maxWidth: 100, mr: 1 }}>
                <LinearProgress variant="determinate" value={Math.round(progress)}/>
            </Box>
            <Box sx={{ minWidth: 35 }}>
                <Typography variant="body2" color="text.secondary">{`${Math.round(progress)}%`}
                </Typography>
            </Box>
      </Box>
    };

    const getFiles = () => {
        const result = [];
            for (let i = 0; i < cache.progressFiles.length; i++) {
                const progressFile = cache.progressFiles[i];
                result.push(getFile(progressFile.file.name, progressFile.progress, i));
            }
        return result;
    };

    const onExpand = () => {
        setExpanded(expanded => !expanded)
    };
    
  return (
    <Snackbar
        onClose={handleClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        open={state.open}
        message={state.snackBarMessage}>
            <Accordion expanded={expanded} onChange={onExpand}>
                <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    style={{height: 40, minHeight: 40}}
                    aria-controls="progressbarContent"
                    id="progressbar">
                        <Typography sx={{ml: -1}}>File progress</Typography> 
                </AccordionSummary>
                <AccordionDetails sx={{}}>
                    <Box sx={{width: 300, maxHeight: 100, overflowY: 'auto', ml: 0, mr: -1, p: 0}}>
                        {getFiles()}
                    </Box>
                </AccordionDetails>
            </Accordion>
    </Snackbar>
    );
};

export default forwardRef(FilesBar)