import React, { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Grid from '@mui/material/Grid2';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import {post} from '../../../../data/Submit'
import HostManager from "../../../../../HostManager";
import { useLocation } from 'react-router-dom'
import Convenience from "../../../../help/Convenience";
import LinearProgress from '@mui/material/LinearProgress';

export default function Details() {

    const [state, setState] = useState({
        loading: true,
        detailForm: {
            title: '',
            titlePlural: '',
            apiName: '',
        }
    });
    const {search} = useLocation();
    const objectId = Convenience.getUrlParameterValueFromQuery(search, 'objectId')

    useEffect(() => {
        if (state.loading) {
            post(`${HostManager.amperHost()}entities/getEntity`, {
                entityId: parseInt(objectId),
            }, (result) => {
                setState({
                    ...state,
                    loading: false,
                    detailForm: {
                        title: result.entity.title,
                        titlePlural: result.entity.titlePlural,
                        apiName: result.entity.apiName,
                    },
                })
                }, (result) => {
                setState({
                    ...state,
                    loading: false,
                });
            });
        }
    });

    const getProgressBar = () => {
        if (state.loading) {
            return <LinearProgress sx={{mb: 1, mr: 6}}/>;
        }
    };

    const save = () => {
        setState({
            ...state,
            loading: true,
        })
        post(`${HostManager.amperHost()}entities/edit`, {
            Id: parseInt(objectId),
            ...state.detailForm
        }, (result) => {
            setState({
                ...state,
                loading: false,
            })
            }, (result) => {
            setState({
                ...state,
                loading: false,
            });
        });
    };

    const handleObjectLabelChange = (event) => {
      const {
        target: { value, name },
      } = event;
      const detailForm = state.detailForm;
      detailForm[name] = value;
      setState({
        ...state,
        detailForm : detailForm
      });
    };

    return (
        <Box sx={{ m: 0, ml: -3, mr: -3, height: 'calc(100% - 10)', width: '100%', flexGrow: 1}}>
            {getProgressBar()}
            <Grid container spacing={3}>
                <Grid item size={3}>
                    <TextField fullWidth label="Label" name="title" value={state.detailForm.title} onChange={handleObjectLabelChange} variant="filled" color="primary" autoFocus size="large"/>
                </Grid>
                <Grid item size={9}>
                </Grid>
                <Grid item size={3}>
                    <TextField fullWidth label="Label plural" name="titlePlural" value={state.detailForm.titlePlural} onChange={handleObjectLabelChange} variant="filled" color="primary" autoFocus size="large"/>
                </Grid>
                <Grid item size={9}>
                </Grid>
                <Grid item size={3}>
                    <TextField InputProps={{
                            readOnly: true,
                         }} value={state.detailForm.apiName} fullWidth label="Api name" variant="filled" color="primary" autoFocus size="large"/>
                </Grid>
                <Grid item size={9}>
                </Grid>
                <Grid item size={3} sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                }}>
                    <Button variant="contained" size='large' onClick={save} disabled={state.loading} sx={{alignSelf: 'end'}}>Save</Button>
                </Grid>
            </Grid>
        </Box>
        
    );
}