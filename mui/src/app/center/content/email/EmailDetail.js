import React from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Stack from '@mui/material/Stack';
import Paper from '@mui/material/Paper';
import Avatar from '@mui/material/Avatar';
import { getFromName, getFromEmail, getDate } from './EmailHelper';
import Tooltip from '@mui/material/Tooltip';
import Convenience from '../../../help/Convenience';
import { isSeen, FlagSeen } from './EmailHelper';
import { post } from '../../../data/Submit';
import HostManager from '../../../../HostManager';
import { AppContext } from '../../../../App';
import EmailsRoot from './EmailsRoot';
import FilesRoot from '../../../components/drive/FilesRoot';

export default function EmailDetail({email, box}) {

    const app = React.useContext(AppContext);
    
    const dateTime = Date.parse(email.envelope.date);
    let year = '';
    let month = '';
    if (email.envelope && email.envelope.date) {
        const dateTime = new Date(email.envelope.date)
        year = dateTime.getFullYear();
        month = dateTime.getMonth() + 1;
    }
    React.useEffect(() => {
        if (!isSeen(email)) {
            post(`${HostManager.myHost()}email/flag`, {
              emails: [{
                email: email.email,
                id: email.id,
              }],
              flags: [FlagSeen],
              box: box,
            }, (result) => {
                if (!result.success) {
                    app.toast('warning', result.error)
                } else {
                    if (email.flags == null) {
                        email.flags = [];
                    }
                    email.flags.push(FlagSeen)
                }
            }, (result) => {
                if (result) {
                    app.toast('warning', result.error)
                }
            });
        }
      }, [email.id]);

    const getBody = (email) => {
        if (Convenience.hasValue(email.bodyHtml)) {
            return <Box sx={{height: 'calc(100% - 90px)', width: '100%', mt: 2}} dangerouslySetInnerHTML={{__html: email.bodyHtml}}>
            </Box>;
        } else if (Convenience.hasValue(email.body)) {
            return <Box sx={{height: 'calc(100% - 90px)', width: 'calc(100% - 40px)', mt: 2, p: '20px', whiteSpace: 'pre-wrap', overflow: 'hidden'}}>
                {email.body}
            </Box>;;
        }
    };

    const fromName = getFromName(email);
    const fromEmail = getFromEmail(email);
    const date = getDate(email);
    return <Paper sx={{width: 'calc(100% - 20px)', m: '10px'}} elevation={3}>
        <Stack direction="row" spacing={2} sx={{pt: 1, pb: 1, ml: 2, height: '40px'}}>
            <Box sx={{flexGrow: 1, width: 'calc(100% - 150px)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}}>
                <Typography variant="h5" gutterBottom>
                    {email.envelope.subject}
                </Typography>
            </Box>
            <Box sx={{flexGrow: 0, width: '150px'}}></Box>
        </Stack>
        <Stack direction="row">
            <Box sx={{flexGrow: 1}}>
                <Stack direction="row" spacing={2} sx={{ pt: 2, pl: 2, ml: 2, height: '50px', width: '400px', border: '1px solid #ccc!important', borderRadius: '15px'}}>
                    <Box sx={{width: '50px', height: '50px'}}>
                        <Avatar sx={{ bgcolor: 'primary.main', color: 'secondary.main' }} alt={fromName} />
                    </Box>
                    <Box>
                        <Tooltip title={<Typography variant="body1">{fromName}</Typography>}>
                            <Typography variant="body1" gutterBottom sx={{mt: -1, width: '350px', whiteSpace: 'nowrap', 'overflow': 'hidden', textOverflow: 'ellipsis'}}>
                                {fromName}
                            </Typography>
                        </Tooltip>
                        <Tooltip title={<Typography variant="body1">{fromEmail}</Typography>}>
                            <Typography variant="body1" gutterBottom sx={{mt: '-3px', width: '350px', whiteSpace: 'nowrap', 'overflow': 'hidden', textOverflow: 'ellipsis'}}>
                                {fromEmail}
                            </Typography>
                        </Tooltip>
                    </Box>
                </Stack>
            </Box>
            <Box sx={{flexGrow: 0, mr:2, mt:1}}>{date}</Box>
        </Stack>
        {getBody(email)}
        <FilesRoot root={'/__system__/Email/' + email.email + '/' + box + '/success/' + year + '/' + month + '/' + email.id} viewLevel={0}></FilesRoot>
    </Paper>;
}