import React, { useState } from 'react';
import AdobeDCViewer from './AdobeDCViewer';
import PdfJsViewer from './PdfJsViewer';
import Box from '@mui/material/Box';
import FileDetail from '../../FileDetail';

const PdfRenderer = ({ url, metadata, directory, app, onVersionChange, detail}) => {
    
    const [state, setState] = useState(() => {
        return {
         
        };
    });

  const renderApp = () => {
    if ('ADOBE_DC' == app) {
      return <AdobeDCViewer url={url} metadata={metadata} directory={directory}/>;
    } else if ('PDF_JS' == app) {
      return <PdfJsViewer url={url} metadata={metadata} directory={directory}/>
    } else {
      return <Box sx={{
        display: 'flex',
        flexDirection: 'row',
        bgcolor: 'background.paper',
        height: '100%'
        }}>
            <Box sx={{ flexGrow: 1, maxHeight: "100%", maxWidth: '100%', }} >
            <Box sx={{
              display: 'flex',
              flexDirection: 'row',
              textAlign: 'center',
              backgroundColor: '#000000',
              justifyContent: 'center',
              height: '100%'
              }}>
                <embed key={url} src={url + '#view=FitH'} height="100%" width="100%"/>
              </Box>
            </Box>
                <Box sx={{ flexGrow: 0, width: '4.5in' }}>
                  {detail}
                </Box>
        </Box>;
    }
  };

  return (
    renderApp()
  );
};
export default PdfRenderer;