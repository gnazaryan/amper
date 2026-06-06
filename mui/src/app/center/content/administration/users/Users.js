import Box from '@mui/material/Box';
import UserList from '../../../../components/adminstration/users/UserList';

export default function Users() {
  return (
        <Box sx={{ height: '100%', width: 'calc(100% - 25px)'}}>
            <UserList></UserList>
        </Box>
    );
}
