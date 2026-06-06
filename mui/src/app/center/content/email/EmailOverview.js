import React from 'react';
import TileView from '../TileView';
import { sessionManager } from '../../../../SessionManager';
import AlternateEmailIcon from '@mui/icons-material/AlternateEmail';

export default function EmailOverview({id, expanded}) {
  const emailTiles = {}
  const user = sessionManager.getUser();
  
  if (user.emails && user.emails.length > 0) {
    for (let i = 0; i < user.emails.length; i++) {
      const email = user.emails[i];
      
      emailTiles[email.email] = {
        key: email.email,
        path: '/email/' + email.email,
        label: email.label,
        icon: <AlternateEmailIcon color={'primary'} fontSize="medium"/>,
        description:  'To explore your emails for ' + email.email + ', click on this tile to navigate into the mailboxes.',
        primary: true
      }
    }
  }
  return <TileView expanded={expanded} items={emailTiles}></TileView>;
}
