import React, { useEffect, useRef, useState } from 'react';
import Box from '@mui/material/Box';

const ImageViewer = ({ url, detail}) => {
    
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
                  <Box component="img" src={url} sx={{height: "auto", maxWidth: '100%', backgroundColor: '000000', maxHeight: (window.innerHeight - 100) +'px', width: 'auto' }} />
                </Box>
              </Box>
              <Box sx={{ flexGrow: 0, width: '4.5in' }}>
                {detail}
              </Box>
        </Box>
};
export default ImageViewer;