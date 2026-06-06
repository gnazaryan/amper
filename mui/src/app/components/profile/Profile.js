import React, { useState, useRef, useEffect } from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import { useNavigate, useLocation } from 'react-router-dom'
import { postBlob } from '../../data/Submit';
import HostManager from '../../../HostManager';
import Typography from '@mui/material/Typography';
import EditIcon from '@mui/icons-material/Edit';
import IconButton from '@mui/material/IconButton';
import { sessionManager } from '../../../SessionManager';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import LocalActivityIcon from '@mui/icons-material/LocalActivity';
import SettingsIcon from '@mui/icons-material/Settings';
import AccessibilityIcon from '@mui/icons-material/Accessibility';
import Settings from './Settings';
import Overview from './Overview';
import About from './About'
import { breadcrumbs } from '../../center/Breadcrambs';
import OpenWithIcon from '@mui/icons-material/OpenWith';
import DoneIcon from '@mui/icons-material/Done';
import ZoomInIcon from '@mui/icons-material/ZoomIn';
import ZoomOutIcon from '@mui/icons-material/ZoomOut';
import { post } from '../../data/Submit'
import { AppContext } from '../../../App';
import { setStoreValue } from '../../amper/Instruments';

export default function Profile({expanded}) {

    const app = React.useContext(AppContext);
    const navigate = useNavigate();
    const location = useLocation();
    const [state, setState] = useState(() => {
        return {
            coverUrl: `${HostManager.myHost()}profile/viewCover`,
            photoUrl: `${HostManager.myHost()}profile/viewPhoto`,
            tab: 0,
            mode: 0,
            zoom: 400,
            loaded: false,
            data: {},
        };
    });
      
    useEffect(() => {
        if (!state.loaded) {
            post(`${HostManager.myHost()}profile/state`, {
            }, (result) => {
                if (result.success) {
                    setState({
                        ...state,
                        loaded: true,
                        data: result.data,
                    });
                } else {
                    app.toast('warning', result.error)
                }
            }, (result) => {
                if (result) {
                    app.toast('warning', result.error)
                }
            });
        }
    }, [state.loaded]);
    

    const profileBoxRef = useRef(null);
    const logRef = useRef(null);

    const activeTab = (() => {
        if (location.pathname.startsWith(breadcrumbs.profile.overview.path)) {
            return 0
        } else if (location.pathname.startsWith(breadcrumbs.profile.about.path)) {
            return 1
        } else if (location.pathname.startsWith(breadcrumbs.profile.settings.path)) {
            return 2
        }
        return 0;
    })();
    const switchTab = (event, newTab) => {
        switch(newTab) {
            case 0:
                navigate(breadcrumbs.profile.overview.key);
                break;
            case 1:
                navigate(breadcrumbs.profile.about.key);
                break;
            case 2:
                navigate(breadcrumbs.profile.settings.key);
                break;
        }
    };

    const uploadCover = (event) => {
        let files = event.currentTarget.files;
        if (files.length > 0) {
            var reader = new FileReader();
            setState({
                ...state,
                coverUrl: '',
            })
            reader.onloadend = function (evt) {
                if (evt.target.readyState == FileReader.DONE) {
                    var arrayBuffer = evt.target.result,
                    content = new Uint8Array(arrayBuffer);
                    postBlob(`${HostManager.myHost()}profile/updateCover`, content, (result) => {
                        if (result.success) {
                            setState({
                                ...state,
                                coverUrl: `${HostManager.myHost()}profile/viewCover`
                            })
                        } else {
                        }
                    }, (result) => {

                    });
                }
            };
            reader.readAsArrayBuffer(files[0]);
        }
    };

    const uploadPhoto = (event) => {
        let files = event.currentTarget.files;
        if (files.length > 0) {
            var reader = new FileReader();
            setState({
                ...state,
                photoUrl: '',
            });
            reader.onloadend = function (evt) {
                if (evt.target.readyState == FileReader.DONE) {
                    var arrayBuffer = evt.target.result,
                    content = new Uint8Array(arrayBuffer);
                    postBlob(`${HostManager.myHost()}profile/updatePhoto`, content, (result) => {
                        if (result.success) {
                            setState({
                                ...state,
                                photoUrl: `${HostManager.myHost()}profile/viewPhoto`
                            });
                        } else {
                        }
                    }, (result) => {

                    });
                }
            };
            reader.readAsArrayBuffer(files[0]);
        }
    };

    const repositionPhoto = () => {
        setState({
            ...state,
            mode: 1,
        });
    };

    const doneRepositionPhoto = () => {
        post(`${HostManager.myHost()}profile/adjustPhoto`, {
            Width: parseInt(profileBoxRef.current.clientWidth),
            Height: parseInt(profileBoxRef.current.clientHeight),
            PositionX: parseInt(profileBoxRef.current.style.left),
            PositionY: parseInt(profileBoxRef.current.style.top),
        }, (result) => {
            if (result.success) {
                setStoreValue("user_photo", result.value);
                window.location.reload();
            } else {
                app.toast('warning', result.error)
            }
        }, (result) => {
            if (result) {
                app.toast('warning', result.error)
            }
        });
    };

    const zoomInRepositionPhoto = () => {
        setState({
            ...state,
            zoom: state.zoom + 10,
        });
    };

    const zoomOutRepositionPhoto = () => {
        setState({
            ...state,
            zoom: state.zoom - 10,
        });
    };

    const onMouseMove = (event) => {
        profileBoxRef.current.style.top = parseInt(profileBoxRef.current.style.top) + event.movementY + 'px';
        profileBoxRef.current.style.left = parseInt(profileBoxRef.current.style.left) + event.movementX  + 'px';
        //logRef.current.innerHTML = parseInt(profileBoxRef.current.style.top) + ' : ' + parseInt(profileBoxRef.current.style.left);
    };

    const onDragStart = (event) => {
        return false;
    };

    const onDragEnd = () => {
        
    };

    const onMouseDown = (event) => {
        profileBoxRef.current.addEventListener('mousemove', onMouseMove);
        profileBoxRef.current.ondragstart = () => {
            return false;
        };
        profileBoxRef.current.onmouseup = () => {
            profileBoxRef.current.removeEventListener('mousemove', onMouseMove);
            profileBoxRef.current.onmouseup = null;
        };
    };

    const onMouseUp = () => {
    };

    const getProfilePictureZoom = () => {
        return state.data && state.data.configuration && state.data.configuration.profile && state.data.configuration.profile.picture ? state.data.configuration.profile.picture.Width : state.zoom;
    };
    
    const getProfilePictureTop = () => {
        return state.data && state.data.configuration && state.data.configuration.profile && state.data.configuration.profile.picture ? state.data.configuration.profile.picture.PositionY : -50;
    };

    const getProfilePictureLeft = () => {
        return state.data && state.data.configuration && state.data.configuration.profile && state.data.configuration.profile.picture ? state.data.configuration.profile.picture.PositionX : 0;
    };

    const onAboutUpdate = (detail) => {
        setState({
            ...state,
            data: {
                ...state.data,
                detail,
            }
        });
    };

    const getActiveTabPanel = () => {
        switch(activeTab) {
            case 0:
                return <Overview></Overview>
            case 1:
                return <About data={state.data.detail} onUpdate={onAboutUpdate} expanded={expanded}></About>
            case 2:
                return <Settings data={state.data.configuration}></Settings>
        }
    };

    const user = sessionManager.getUser();

    const getProfileButtons = () => {
        const result = [];

        if (state.mode == 0) {
            result .push(<IconButton color="primary" component="label" style={{position: 'absolute', bottom: 15, left: 15, zIndex: 1000000}}>
                <OpenWithIcon onClick={repositionPhoto}/>
            </IconButton>);
            result.push(<IconButton color="primary" component="label" style={{position: 'absolute', bottom: 15, right: 15, zIndex: 1000000}}>
                <EditIcon />
                <input hidden accept="image/png, image/jpeg" multiple type="file" onChange={uploadPhoto}/>
            </IconButton>);
        } else if (state.mode == 1) {
            result .push(<IconButton color="primary" component="label" style={{position: 'absolute', bottom: 15, left: 15, zIndex: 1000000}}>
                <DoneIcon onClick={doneRepositionPhoto}/>
            </IconButton>);
            result .push(<IconButton color="primary" component="label" style={{position: 'absolute', top: 15, left: 15, zIndex: 1000000}}>
                <ZoomInIcon onClick={zoomInRepositionPhoto}/>
            </IconButton>);
            result .push(<IconButton color="primary" component="label" style={{position: 'absolute', top: 15, right: 15, zIndex: 1000000}}>
                <ZoomOutIcon onClick={zoomOutRepositionPhoto}/>
            </IconButton>);
            //result.push(<span style={{position: 'absolute', bottom: 'calc(50% - 20px)', left: 'calc(50% - 20px)', fontSize: '40px', color: 'white'}}><DragIndicatorIcon style={{fontSize: '40px', color: 'white'}} color="white"/></span>);
        }
        
        return result;
    };

  return (
        <Box style={{overflowX: 'hidden', overflowY: 'auto'}} width="100%" height={'100%'}>
            <Box style={{height: '300px', position: 'relative'}}>
                <Box sx={{
                    display: 'flex',
                    flexDirection: 'row',
                    borderRadius: 15,
                    textAlign: 'center',
                    backgroundColor: '#d2e7f7',
                    justifyContent: 'center',
                    height: '100%'
                    }}>
                    <Box component="img" src={state.coverUrl}  sx={{height: '100%', width: '100%', objectFit: 'cover', objectPosition: '0px, 0px', borderRadius: 15}}>
                    </Box>  
                </Box>
                <Button variant="outlined" component="label" style={{position: 'absolute', top: 5, right: 5}}>
                    Update cover
                    <input hidden accept="image/png, image/jpeg" multiple type="file" onChange={uploadCover}/>
                </Button>
            </Box>
            <Box style={{position: 'relative', height: '150px'}}>
                <Box style={{position: 'absolute', borderRadius: 100, top: -70, left: 5, backgroundColor: '#ffffff', overflow: 'hidden', width: '200px', height: '200px', zIndex: 100}}>
                    <Box style={{cursor: state.mode == 1 ? 'move' : 'default', position: 'absolute', top: getProfilePictureTop() + 'px', left: getProfilePictureLeft() + 'px', width: getProfilePictureZoom() + 'px'}}
                    ref={profileBoxRef}
                    onMouseDown={state.mode == 1 ? onMouseDown : () => {}}
                    onDragStart={onDragStart}
                    onDragEnd={onDragEnd}
                    draggable="true" component="img" src={state.photoUrl}
                    sx={{ }}></Box>
                    {getProfileButtons()}
                </Box>
                <Box sx={{position: 'absolute', top: 0, left: 220}}>
                    <Typography variant="h6">
                        {user.firstName + ' ' + user.lastName}
                    </Typography>
                    <Typography variant="subtitle1">
                        {user.username + ' '}
                    </Typography>
                    <Typography variant="subtitle1">
                        {user.email + ' '}
                    </Typography>
                    {/*<div ref={logRef}>0</div>*/}
                </Box>
            </Box>
            <Box  style={{position: 'relative', width: '100%', height: 'auto'}}>
                <Tabs style={{width: '100%'}} value={activeTab} onChange={switchTab} centered>
                    <Tab key="overview" icon={<LocalActivityIcon />} iconPosition="start" label="Overview" />
                    <Tab key="about" icon={<AccessibilityIcon />} iconPosition="start" label="About" />
                    <Tab key="settings" icon={<SettingsIcon />} iconPosition="start" label="Settings"/>
                </Tabs>
                {getActiveTabPanel()}
            </Box>
        </Box>
    );
}