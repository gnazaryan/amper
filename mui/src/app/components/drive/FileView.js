import React, { useState, useImperativeHandle, forwardRef, useRef } from 'react';
import Modal from '@mui/material/Modal';
import Fade from '@mui/material/Fade';
import Backdrop from '@mui/material/Backdrop';
import Box from '@mui/material/Box';
import HostManager from '../../../HostManager';
import PdfRenderer from './renderer/pdf/PdfRenderer';
import FileDetail from './FileDetail';
import { post } from '../../data/Submit';
import ImageViewer from './renderer/image/ImageViewer';
import IconButton from '@mui/material/IconButton';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import ArrowBackIosNewIcon from '@mui/icons-material/ArrowBackIosNew';
import { sessionManager } from '../../../SessionManager';

const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    bgcolor: 'background.paper',
    boxShadow: 24,
    width: 'calc(100% - 100px)',
    height: 'calc(100% - 100px)'
  };

function FileView({next, previous}, ref) {

    const [state, setState] = useState(() => {
        return {
            open: false,
        };
    });

    const detailRef = useRef();

    useImperativeHandle(ref, () => ({
      view(directory, metadata, app) {
        setState({
          ...state,
          open: true,
          directory,
          metadata,
          version: metadata.version,
          app,
        });
      }
    }));

    const handleClose = () => {
        setState({
            ...state,
            open: false,
        });
    };

    const getFileDetail = () => {
      return <FileDetail ref={detailRef} metadata={state.metadata} onVersionChange={onVersionChange}/>;
    };

    const getPdfViewer = (url) => {
      return <PdfRenderer 
        detail={getFileDetail()}
        key={url} 
        url={url} 
        metadata={state.metadata} 
        directory={state.directory} 
        app={state.app}/>
    };

    const onVersionChange = (version) => {
      post(`${HostManager.amperHost()}files/metadata`, {
        directory: state.directory,
        id: state.metadata.id,
        major: version.major,
        minor: version.minor,
        patch: version.patch,
      }, (result) => {
        setState({
          ...state,
          metadata: result.data,
          version,
        });
      }, (result) => {
          
      });
    };

    const getImageViewer = (url) => {
      return <ImageViewer key={url} url={url} detail={getFileDetail()}/>
    };

    const getView = () => {
      if (state.metadata != null && state.metadata.id != null) {
        let type = state.metadata.type;
        if (state.metadata.rendition === true && state.metadata.renditionType !== '?') {
          type = state.metadata.renditionType;
        }
        const sessionId = sessionManager.getSessionId();
        const url = `${HostManager.amperHost()}files-v1/viewFile?id=${encodeURIComponent(state.metadata.id)}&root=${state.directory}&major=${state.version.major}&minor=${state.version.minor}&patch=${state.version.patch}&sessionId=${encodeURIComponent(sessionId)}`;
        switch(type) {
          case 'application/pdf':
            return getPdfViewer(url);
            case 'image/png':
            case 'image/png', 'image/jpeg':
              return getImageViewer(url);            
        }
      }
    };

    return (<Modal
        open={state.open}
        onClose={handleClose}
        closeAfterTransition
        slots={{ backdrop: Backdrop }}
        slotProps={{
          backdrop: {
            timeout: 500,
          },
        }}
      >
        <Fade in={state.open}>
          <Box sx={style}>
            {state.app ? '' : <IconButton color="primary" size='large' style={{position: 'absolute', top: (window.innerHeight - 200) / 2}} onClick={()=>{previous(state.metadata)}}>
              <ArrowBackIosNewIcon sx={{fontSize: '60px'}}/>
            </IconButton>}
            {/*<Typography id="transition-modal-title" variant="h6" component="h2">
              Text in a modal - {state.metadata ? state.metadata.name : ''}
              </Typography>*/}
            {getView()}
            {state.app ? '' : <IconButton color="primary" size='large' style={{position: 'absolute', top: (window.innerHeight - 200) / 2, right: '4.5in'}} onClick={()=>{next(state.metadata)}}>
              <ArrowForwardIosIcon sx={{fontSize: '60px'}}/>
            </IconButton>}
          </Box>
        </Fade>
      </Modal>);
};

export default forwardRef(FileView)