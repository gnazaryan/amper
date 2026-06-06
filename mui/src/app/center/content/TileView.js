import React, { useState, useEffect, useRef } from 'react';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid2';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import { CardActionArea } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom'
import Stack from '@mui/material/Stack';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import SearchIcon from '@mui/icons-material/Search';
import Convenience from "../../../app/help/Convenience";
import {registerResize} from "../../amper/Instruments";
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';


export default function TileView({expanded, items, order, removeHandler}) {

    const getPartsInTile = (width, expanded) => {
        const itemsInPage = (width - (expanded ? 200 : 60)) / 300;
        return 12 / itemsInPage;
    };
    const [state, setState] = useState({
        dashboardViewHeight: window.innerHeight - 250,
        search: undefined,
        parts: getPartsInTile(window.innerWidth, expanded),
        menuOpen: false
    });
    const [anchorEl, setAnchorEl] = React.useState(null);
    const menuOpen = Boolean(anchorEl);

    const handleResize = (height, width, expanded) => {
        setState({
          ...state,
          dashboardViewHeight: height - 250,
          parts: getPartsInTile(width, expanded),
        })
      };
    
      useEffect(() => {
        registerResize(handleResize, expanded)
      }, []);

      const handleSearchChange = (event) => {
        const {
          target: { value, name },
        } = event;
        setState({
          ...state,
          search: value,
        });
      };

      const removeOnClick = (event) => {
        const {
            currentTarget: {dataset: { name, value}},
          } = event;
          removeHandler(name, value);
      };

      const handleMenuClose = (event) => {
        event.stopPropagation();
        setState({
            ...state,
            menuOpen: false,
        });
      };

      const handleMenuClick = (event) => {
        setAnchorEl(event.currentTarget);
        setState({
            ...state,
            menuOpen: true,
        });
      };

      const getRemoveButton = (key, label, primary) => {
        if (!primary) {
            return (
                <IconButton
                    onClick={handleMenuClick}
                    style={{position: 'absolute', right: 0, top: 0, zIndex: 1001}}
                    size="large">
                    <MoreHorizIcon color="primary" fontSize="small" />
                    <Menu
                        anchorEl={anchorEl}
                        open={state.menuOpen}
                        onClose={handleMenuClose}
                        MenuListProps={{
                        'aria-labelledby': 'basic-button',
                        }}
                    >
                        <MenuItem data-name={key} data-value={label} onClick={removeOnClick}>
                            <ListItemIcon>
                                <DeleteIcon color="primary" fontSize="small" />
                            </ListItemIcon>
                            <ListItemText>Remove</ListItemText>
                        </MenuItem>
                    </Menu>
            </IconButton>
            );
        }
      };

      const getCard = (key, icon, label, description, path, primary) => {
            return (
            <Card variant="outlined" sx={{ width: 250, height: 250, overflow: 'visible'}}>
                <CardActionArea sx={{ width: 250, height: 250, zIndex: 1000}} component="div">
                    {getRemoveButton(key, label, primary)}
                    <CardActionArea component={RouterLink} to={path} state={{expanded}} sx={{ width: 250, height: 250, zIndex: 1000}}>
                        <Box display="flex"
                            justifyContent="center"
                            alignItems="center">
                            {icon}
                        </Box>
                        <CardContent>
                        <Typography gutterBottom variant="h5" component="div" sx={{mt: -2, textAlign: 'center', color: 'secondary.menuText' }}>
                            {label}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{mt: 0}}>
                            {description}
                        </Typography>
                        </CardContent>
                    </CardActionArea>
                </CardActionArea>
            </Card>);
      };

      const cardsGrid = [];
      if (order) {
        for (let i = 0; i < order.length; i++) {
            const value = items[order[i]];
            if ((!Convenience.hasValue(state.search) || value.label.toLowerCase().includes(state.search.toLowerCase()) || value.description.toLowerCase().includes(state.search.toLowerCase()))) {
                cardsGrid.push(
                    <Grid key={value.path} item size={state.parts}>
                        {getCard(value.key,
                            React.cloneElement(value.icon, {sx: {fontSize: '100px'}, color: 'primary'}),
                            value.label, value. description, value.path, value.primary)}
                    </Grid>
                );
            }
        }
      }
      for (const [key, value] of Object.entries(items)) {
            if (value.key && (!order || !order.includes(value.key)) && (!Convenience.hasValue(state.search) || value.label.toLowerCase().includes(state.search.toLowerCase()) || value.description.toLowerCase().includes(state.search.toLowerCase()))) {
                cardsGrid.push(
                    <Grid key={value.path} item size={state.parts}>
                        {getCard(value.key,
                            React.cloneElement(value.icon, {sx: {fontSize: '100px'}, color: 'primary'}),
                            value.label, value.description, value.path, value.primary)}
                    </Grid>
                );
            }
      }
  return (
    <Box sx={{ width: 'calc(100% - 23px)'}}>
        <Stack direction="row" sx={{ mb: 2, mr: 1}}>
            <TextField
                placeholder="Search..."
                onChange={handleSearchChange}
                InputProps={{
                startAdornment: (
                    <InputAdornment position="start">
                        <SearchIcon />
                    </InputAdornment>
                ),
                }}
                sx={{minWidth: '100%'}}
                variant="standard"
            />
            </Stack>
        <Box sx={{ height: state.dashboardViewHeight, overflowX: "hidden", overflowY: "auto",}}>
          <Grid container spacing={2}>
              {cardsGrid}
          </Grid>
        </Box>
    </Box>
    );
}
