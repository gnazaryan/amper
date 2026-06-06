import React from 'react';
import FilesRoot from "../../../components/drive/FilesRoot";
import EmailsRoot from './EmailsRoot';

export default function Email({id, expanded}) {
  return <EmailsRoot id={id} expanded={expanded}/>;
}
