import { postBlob } from "../../../../data/Submit";
import HostManager from "../../../../../HostManager";
import { requireExternal } from "../../../../amper/Instruments";
import Convenience from "../../../../help/Convenience";
import { sessionManager } from "../../../../../SessionManager";

class ViewSDKClient {
    constructor() {
      this.readyPromise = new Promise((resolve) => {
        if (window.AdobeDC) {
          resolve();
        } else {
          requireExternal("https://documentservices.adobe.com/view-sdk/viewer.js");
          document.addEventListener("adobe_dc_view_sdk.ready", () => {
            resolve();
          });
        }
      });
    }
    ready() {
      return this.readyPromise;
    }
    previewFile(divId, viewerConfig, url, {metadata, directory, user}) {
      this.metadata = metadata;
      this.user = user;
      this.directory = directory;
      const settings = sessionManager.getSettings();
      const config = {
        clientId: settings.adobeLicenseKey || '', ///enter lient id here
        sendAutoPDFAnalytics: false,
        loggingUri: false,
      };
      if (divId) {
        config.divId = divId;
      }
      this.config = config;
      this.adobeDCView = new window.AdobeDC.View(config);
      const previewFilePromise = this.adobeDCView.previewFile(
        {
          content: {
            location: {
              url: url,
            },
          },
          metaData: {
            fileName: this.metadata.name,
            id: "6d07d124-ac85-43b3-a867-36930f502ac6",
          },
        },
        viewerConfig
      );
      this.registerCallbacks()
      return previewFilePromise;
    }
    registerCallbacks() {
        this.registerUserProfile(this.user);
        this.registerSaveApiHandler(this.user, this.metadata, this.directory, this.config);
    }
    registerUserProfile(user) {
        this.adobeDCView.registerCallback(
            window.AdobeDC.View.Enum.CallbackType.GET_USER_PROFILE_API,
            function() {
               return new Promise((resolve, reject) => {
                  resolve({
                     code: window.AdobeDC.View.Enum.ApiResponseCode.SUCCESS,
                     data: {
                        userProfile: {
                            name: user.firstName + ' ' + user.lastName,
                            firstName: user.firstName,
                            lastName: user.lastName,
                            email: user.email
                        }
                     }
                  });
               });
            },
         {});
    }
    registerSaveApiHandler(user, metadata, directory) {
      const saveApiHandler = (adobeMetaData, content, options) => {
        return new Promise((resolve) => {
          const formData = new FormData()
          formData.append('chunk', new Blob([content.buffer]));
          formData.append('id', metadata.id);
          formData.append('name', metadata.name);
          formData.append('major', metadata.version.major);
          formData.append('minor', metadata.version.minor);
          formData.append('patch', metadata.version.patch);
          formData.append('directory', directory);
          formData.append('size', content.buffer.byteLength);
          postBlob(`${HostManager.amperHost()}files-v1/upversion`, formData, (result) => {
                let response = null;
                if (result.success) {
                    response = {
                        code: window.AdobeDC.View.Enum.ApiResponseCode.SUCCESS,
                        data: {
                          metaData: Object.assign(adobeMetaData, {
                            updatedAt: new Date().getTime(),
                          }),
                        },
                      };
                } else {
                    response = {
                        code: window.AdobeDC.View.Enum.ApiResponseCode.FAIL,
                    };
                }
                resolve(response);
            }, (result) => {
                resolve({
                    code: window.AdobeDC.View.Enum.ApiResponseCode.FAIL,
                });
            });
        });
      };
      this.adobeDCView.registerCallback(
        window.AdobeDC.View.Enum.CallbackType.SAVE_API,
        saveApiHandler,
        {}
      );
    }
    registerEventsHandler() {
      this.adobeDCView.registerCallback(
        window.AdobeDC.View.Enum.CallbackType.EVENT_LISTENER,
        (event) => {
          console.log(event);
        },
        {
          enablePDFAnalytics: true,
        }
      );
    }
  }
  export default ViewSDKClient;