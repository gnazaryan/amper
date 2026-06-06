import Box from '@mui/material/Box';
import UserProfile from "../../../components/profile/Profile";

export default function Profile({expanded}) {
  return (
    <Box sx={{ height: '100%', width: 'calc(100% - 25px)'}}>
         <UserProfile expanded={expanded}></UserProfile>
    </Box>
    );
}
