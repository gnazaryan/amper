import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import DialogActions from '@mui/material/DialogActions';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import UserSelect from './UserSelect';
import { sessionManager } from '../../../SessionManager';

export default function UserDialog({open, close, chat, onSelectionChange}) {

    const closeDialog = () => {
        close();
    };

    const startChat = () => {
        chat();
    };

    return <Dialog onClose={closeDialog} open={open}>
    <DialogTitle>Select people to chat</DialogTitle>
    <DialogContent>
        <DialogContentText>
            Select amper user(s) then click chat
        </DialogContentText>
        <Box sx={{ height: '100%', width: '100%' }}>
            <UserSelect onSelectionChange={onSelectionChange}></UserSelect>            
        </Box>
    </DialogContent>
    <DialogActions>
        <Button onClick={closeDialog}>Cancel</Button>
        <Button onClick={startChat} disabled={false}>Chat</Button>
    </DialogActions>
</Dialog>;
}