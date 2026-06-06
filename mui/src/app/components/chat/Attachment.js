import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import './Attachment.css'
import { truncate } from '../../amper/Instruments';
import AttachmentIcon from '@mui/icons-material/Attachment';
import Tooltip from '@mui/material/Tooltip';
import ClearIcon from '@mui/icons-material/Clear';
import IconButton from '@mui/material/IconButton';
import DownloadIcon from '@mui/icons-material/Download';

export default function Attachment({metadata, onRemove, onDownload}) {
    const getRemoveDownloadButton = () => {
        if (onRemove != null) {
            return <IconButton aria-label="remove attachment" onClick={() => {
                    onRemove(metadata.id, metadata.directory);
                }}>
                <ClearIcon color="primary"/>
            </IconButton>;
        } else if (onDownload != null) {
            return <IconButton aria-label="download attachment" onClick={() => {
                    onDownload(metadata.id, metadata.name);
                }}>
                <DownloadIcon color="primary"/>
            </IconButton>;
        }
    };

    return <Box sx={{ height: '40px', overflow: 'clip', ml: '3px', display: 'flex', cursor: 'pointer'}} className="attachment">
                <AttachmentIcon sx={{mt: '7px'}}></AttachmentIcon>
                <Tooltip title={metadata.name}>
                    <Typography variant="subtitle1" sx={{mt: '7px', mr: '5px'}} gutterBottom>
                        {truncate(metadata.name, 20, true)}
                    </Typography>
                </Tooltip>
                {getRemoveDownloadButton()}
                
        </Box>;
}