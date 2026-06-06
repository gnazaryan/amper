import React, { useState, useRef } from 'react';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import KeyboardDoubleArrowRightIcon from '@mui/icons-material/KeyboardDoubleArrowRight'
import KeyboardDoubleArrowLeftIcon from '@mui/icons-material/KeyboardDoubleArrowLeft';
import Typography from '@mui/material/Typography';
import KeyboardDoubleArrowUpIcon from '@mui/icons-material/KeyboardDoubleArrowUp';
import KeyboardDoubleArrowDownIcon from '@mui/icons-material/KeyboardDoubleArrowDown';
import Paper from '@mui/material/Paper';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import ButtonGroup from '@mui/material/ButtonGroup';
import Button from '@mui/material/Button';
import DeleteIcon from '@mui/icons-material/Delete';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import SettingsIcon from '@mui/icons-material/Settings';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import OpenWithIcon from '@mui/icons-material/OpenWith';
import Slide from '@mui/material/Slide';
import Tooltip from '@mui/material/Tooltip';

export default function Widget(props) {
    const {onRemove, onExpand, onShrink, onHighten, onLower, onConfigure, widget, configuration} = props;
    const [state, setState] = useState({
        menuOpen: false,
        repositionOpen: false
    });
    const containerRef = useRef(null);

    const [anchorEl, setAnchorEl] = React.useState(null);

    const expand = () => {
        onExpand(widget);
    };

    const shrink = () => {
        onShrink(widget);
    };

    const highten = () => {
        onHighten(widget);
    };

    const lower = () => {
        onLower(widget);
    };
    
    const more = (event) => {
        setAnchorEl(event.currentTarget);
        setState({
            ...state,
            menuOpen: true,
        });
    };

    const remove = () => {
        onRemove(widget);
    };

    const moreClose = (event) => {
        event.stopPropagation();
        setState({
            ...state,
            menuOpen: false,
        });
    };

    const configure = (event) => {
        onConfigure();
    };

    const onMouseOver = (event) => {
        setState({
            ...state,
            repositionOpen: true,
        });
    };

    const onMouseOut = () => {
        setState({
            ...state,
            repositionOpen: false,
        });
    };

    const onMouseOverGroup = () => {
        setState({
            ...state,
            repositionOpen: true,
        });
    };

    const onMouseOutGroup = (event) => {
        setState({
            ...state,
            repositionOpen: false,
        });
    };

    const getNotConfigured = () => {
        return <Box sx={{ display: 'flex', width: '100%', height: '100%', verticalAlign: 'middle', alignItems: 'center', justifyContent: 'center' }}>
            <Typography variant="subtitle1" gutterBottom>
                Widget is not configured.
            </Typography>
        </Box>;
    };

    let conditionalProps = state.repositionOpen ? {onMouseOver: onMouseOverGroup} : {};
  return (
        <Paper sx={{height: props.height, position: 'relative', overflow: 'hidden'}} variant="outlined" >
            <Stack direction="row" spacing={1} sx={{ mb: 1 }} bgcolor="secondary.main" style={{borderTopLeftRadius: '4px', borderTopRightRadius: '4px'}}>
                
                    <Typography sx={{ml: 1, mt: 1, mr: 2, flexGrow: 1, textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap', color: 'primary.main'}} variant="subtitle1" gutterBottom>
                        <Tooltip title={widget.description}><span>{widget.label}</span></Tooltip>
                        {/*<Tooltip title={widget.description}><InfoIcon color="primary" style={{width: '15px', height: '15px'}} sx={{ml: 2, mb: -0.4}}/></Tooltip>*/}
                    </Typography>
                
                <ButtonGroup variant="outlined" {...conditionalProps} onMouseOut={onMouseOutGroup}  color="primary" ref={containerRef} sx={{overflow: 'visible', maxHeight: '30px'}}>
                    <Slide direction="left" in={state.repositionOpen} container={containerRef.current}>
                        <Box sx={{mr: -5}}>
                            <Button onClick={shrink} color="primary" title="Shrink" size="small" sx={{mt:'7px', height: '30px'}}>
                                <KeyboardArrowLeftIcon />
                            </Button>
                            <Button onClick={expand} color="primary" title="Expand" size="small" sx={{mt:'7px', height: '30px'}}>
                                <KeyboardArrowRightIcon />
                            </Button>
                            <Button onClick={shrink} color="primary" title="Shrink" size="small" sx={{mt:'7px', height: '30px'}}>
                                <KeyboardDoubleArrowLeftIcon />
                            </Button>
                            <Button onClick={expand} color="primary" title="Expand" size="small" sx={{mt:'7px', height: '30px'}}>
                                <KeyboardDoubleArrowRightIcon />
                            </Button>
                            <Button onClick={highten} color="primary" title="Highten" size="small" sx={{mt:'7px', height: '30px'}}>
                                <KeyboardDoubleArrowDownIcon/>
                            </Button>
                            <Button onClick={lower} color="primary" title="Lower" size="small" sx={{mt:'7px', height: '30px'}}>
                                <KeyboardDoubleArrowUpIcon/>
                            </Button>
                        </Box>
                    </Slide>
                    <Slide direction="left" in={!state.repositionOpen} container={containerRef.current}>
                        <Button onMouseOver={onMouseOver} color="primary" title="Lower" size="small" sx={{mt:'7px', height: '30px'}}>
                            <OpenWithIcon/>
                        </Button>
                    </Slide>
                </ButtonGroup>
                <Button onClick={more} color="primary" title="More" style={{width: '10px'}}>
                    <MoreHorizIcon sx={{mr: 0, height: '30px'}}/>
                    <Menu
                        anchorEl={anchorEl}
                        open={state.menuOpen}
                        onClose={moreClose}
                        MenuListProps={{
                        'aria-labelledby': 'basic-button',
                        }}
                    >
                        <MenuItem onClick={configure}>
                            <ListItemIcon>
                                <SettingsIcon color="primary" fontSize="small" />
                            </ListItemIcon>
                            <ListItemText>Configure</ListItemText>
                        </MenuItem>
                        <MenuItem onClick={remove}>
                            <ListItemIcon>
                                <DeleteIcon color="primary" fontSize="small" />
                            </ListItemIcon>
                            <ListItemText>Remove</ListItemText>
                        </MenuItem>
                    </Menu>
                </Button>
            </Stack>
            <Box sx={{m: 1, height: 'calc(100% - 60px)'}}>
                {configuration.configured === true ? props.children : getNotConfigured()}
                {props.children[0]}
            </Box>
            
            {/*<span style={{borderBottomRightRadius: 0, position: 'absolute', width: '15px', height: '15px', borderBottom:'4px double', cursor: 'nwse-resize', borderRight: '4px double', right: -2, bottom: -2}}></span>*/}
        </Paper>
    );
}
