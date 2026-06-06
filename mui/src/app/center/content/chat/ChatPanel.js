import Box from '@mui/material/Box';
import Chat from '../../../components/chat/Chat';

export default function ChatPanel() {
  return (
    <Box sx={{ height: '100%', width: 'calc(100% - 25px)'}}>
         <Chat></Chat>
    </Box>
    );
}
