import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import EditIcon from '@mui/icons-material/Edit';
import PhoneIcon from '@mui/icons-material/Phone';
import PhoneIphoneIcon from '@mui/icons-material/PhoneIphone';
import FacebookIcon from '@mui/icons-material/Facebook';
import Link from '@mui/material/Link';
import XIcon from '@mui/icons-material/X';
import BusinessIcon from '@mui/icons-material/Business';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import { sessionManager } from '../../../SessionManager';
import EmailIcon from '@mui/icons-material/Email';

export default function Overview() {
    const user = sessionManager.getUser();
    const onInfoEditClicked = () => {

    };

    const getInfoSectionViewMode = () => {
        return <table>
            <tr>
                <td width="15px"><BusinessIcon color="primary"/></td>
                <td><Link href="https://www.facebook.com/grigornazaryan/">Yerevan, Armenia</Link></td>
            </tr>
            <tr>
                <td width="15px"><EmailIcon color="primary"/></td>
                <td><Link href={"mailto: " + user.email}>{user.email}</Link></td>
            </tr>
            <tr>
                <td width="15px"><PhoneIcon color="primary"/></td>
                <td>+374 95 332343</td>
            </tr>
            <tr>
                <td width="15px"><PhoneIphoneIcon color="primary"/></td>
                <td>+374 95 332343</td>
            </tr>
            <tr>
                <td width="15px"><LinkedInIcon color="primary"/></td>
                <td><Link href="https://www.linkedin.com/in/grigornazaryan/">linkedin</Link></td>
            </tr>
            <tr>
                <td width="15px"><FacebookIcon color="primary"/></td>
                <td><Link href="https://www.facebook.com/grigornazaryan/">facebook</Link></td>
            </tr>
            <tr>
                <td width="15px"><XIcon color="primary"/></td>
                <td><Link href="https://x.com/GNazaryan">twitter</Link></td>
            </tr>
        </table>;
    };

    const getInfoSection = () => {
        return <Paper elevation={3} sx={{p: 1, m:3, height: '300px', width: '100%'}}>
            <Box sx={{display: 'flex', flexDirection: 'row'}}>
                <Box sx={{display: 'flex', flexGrow: 1}}>
                    <Typography variant="h5" gutterBottom>
                        Contacts
                    </Typography>
                </Box>
                <Box sx={{display: 'flex'}} >
                    <IconButton aria-label="edit" onClick={() => {onInfoEditClicked()}}>
                        <EditIcon color="primary"/>
                    </IconButton>
                </Box>
            </Box>
            <Divider sx={{mb: 2}}/>
            {getInfoSectionViewMode()}
        </Paper>;
    };

    return <Box sx={{display: 'flex', flexDirection: 'row', width: '100%'}}>
        <Box sx={{display: 'flex', width: '40%'}}>
            {getInfoSection()}
        </Box>
        <Box sx={{display: 'flex', width: '60%'}} >
        </Box>
    </Box>;
}