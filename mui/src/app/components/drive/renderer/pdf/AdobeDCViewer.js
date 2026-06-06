import React, { useEffect, useState } from 'react';
import ViewSDKClient from "./ViewSDKClient";
import { sessionManager } from '../../../../../SessionManager';

const AdobeDCViewer = ({ url, metadata, directory}) => {

    useEffect(() => {
        load();
    });

  const load = () => {
    const viewSDKClient = new ViewSDKClient();
    viewSDKClient.ready().then(() => {
      viewSDKClient.previewFile(
        "pdf-div",
        {
          defaultViewMode: "FIT_WIDTH",
          showAnnotationTools: true,
          showLeftHandPanel: true,
          showPageControls: true,
          showDownloadPDF: true,
          showPrintPDF: true,
        },
        url,
        {metadata, directory, user: sessionManager.getUser(),}
      );
    });
  };

  return (<div>
    <div style = {{height:"90vh"}} id="pdf-div">
    </div>
  </div>);
};
export default AdobeDCViewer;