import React from 'react';
import FilesRoot from "../../../components/drive/FilesRoot";

export default function Files({expanded}) {
  return <FilesRoot expanded={expanded} viewLevel={10}/>;
}
