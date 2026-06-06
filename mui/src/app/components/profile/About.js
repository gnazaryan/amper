import React, { useState, useEffect } from 'react';
import Paper from '@mui/material/Paper';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import EditIcon from '@mui/icons-material/Edit';
import ChatEditor from '../chat/ChatEditor';
import { post } from '../../data/Submit';
import HostManager from '../../../HostManager';
import { AppContext } from '../../../App';
import { Gauge, gaugeClasses } from '@mui/x-charts/Gauge';
import Chip from '@mui/material/Chip';
import SkillAutocomplete from './SkillAutocomplete'
import AddIcon from '@mui/icons-material/Add';
import Button from '@mui/material/Button';
import CancelIcon from '@mui/icons-material/Cancel';
import SaveIcon from '@mui/icons-material/Save';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Grid from '@mui/material/Grid2';
import Convenience from '../../help/Convenience';
import Avatar from '@mui/material/Avatar';

export default function About({data, onUpdate, expanded}) {
    const app = React.useContext(AppContext);
    const [state, setState] = useState({
        about_me: data?.about_me,
        responsibilities: data?.responsibilities,
        skills: data?.skills
    });

    useEffect(() => {
        setState({
            about_me: data?.about_me,
            responsibilities: data?.responsibilities,
            skills: data?.skills
        });
      }, [data?.about_me, data?.responsibilities, data?.skills]);

    const onEditClicked = (id) => {
        setState({
            ...state,
            editSection: id,
        });
    };

    const saveSection = (id, value) => {
        post(`${HostManager.myHost()}profile/saveDetail`, {
            name: id,
            value: value,
        }, (result) => {
            if (result.success) {
                setState({
                    ...state,
                    editSection: null,
                    [id]: value,
                });
            } else {
                app.toast('warning', result.error)
            }
        }, (result) => {
            if (result) {
                app.toast('warning', result.error)
            }
        });
    };

    const cancelSaveSection = () => {
        setState({
            ...state,
            editSection: null,
        });
    };

    const getImageSource = (base64Image) => {
        if (Convenience.hasValue(base64Image)) {
          return 'data:image/png;base64,' + base64Image;
        }
        return '/static/images/avatar/2.jpg';
      };

      const CARD_WIDTH = 350;

    const getRelationship = (user) => {
        const cardCount = Math.floor((document.body.clientWidth - (expanded ? 350 : 200)) / CARD_WIDTH);
        const size = Math.round(12 / cardCount);
        console.log(size);
        return (<Grid size={size}>
            <Card sx={{ display: 'flex', m: 1, minWidth: CARD_WIDTH, maxWidth: CARD_WIDTH}}>
                <CardContent sx={{ flex: '1 0 auto' }}>
                  <Typography component="div" variant="h6">
                    {user.firstName + ' ' + user.lastName}
                  </Typography>
                  <Typography
                    variant="subtitle1"
                    component="div"
                    sx={{ color: 'text.secondary' }}>
                    {user.username}
                  </Typography>
                </CardContent>
                <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 1, width: 150, height: 150 }} src={getImageSource(user.photo)} />
            </Card>
            </Grid>
          );
        
    };

    const RelationshipSection = function({id, title, users}) {
        return <Paper elevation={3} sx={{p: 1, m:3}}>
            <Box sx={{display: 'flex', flexDirection: 'row'}}>
                <Box sx={{display: 'flex', flexGrow: 1}}>
                    <Typography variant="h5" gutterBottom>
                        {title}
                    </Typography>
                </Box>
                <Box sx={{display: 'flex'}} >
                </Box>
            </Box>
            <Divider sx={{mb: 2}}/>
            <Box sx={{display: 'flex', flexDirection: 'row', overflowX: 'auto'}}>
                <Grid container spacing={2} sx={{width: '100%'}}>
                    {users.map(user => getRelationship(user))}
                </Grid>
            </Box>
        </Paper>;
    };

    const AboutSection = function({id, title, content}) {
        return <Paper elevation={3} sx={{p: 1, m:3}}>
            <Box sx={{display: 'flex', flexDirection: 'row'}}>
                <Box sx={{display: 'flex', flexGrow: 1}}>
                    <Typography variant="h5" gutterBottom>
                        {title}
                    </Typography>
                </Box>
                <Box sx={{display: 'flex'}} >
                    <IconButton aria-label="edit" onClick={() => {onEditClicked(id)}}>
                        <EditIcon color="primary"/>
                    </IconButton>
                </Box>
            </Box>
            <Divider sx={{mb: 2}}/>
            {state.editSection === id ? 
                <ChatEditor content={content} showSaveCancel={true} save={(value) => {saveSection(id, value)}} cancel={cancelSaveSection}></ChatEditor> : 
                (!content ? <Box sx={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}><Typography variant="subtitle1" gutterBottom>This section is empty</Typography></Box>: <Box dangerouslySetInnerHTML={{ __html: content }}></Box>)}
        </Paper>;
    }

    const onSkillChange = (index, skills) => {
        const skillsSection = JSON.parse(state.skills);
        skillsSection[index].skills = skills;
        setState({
            ...state,
            skills: JSON.stringify(skillsSection),
        });
    };
    
    const sumSkillSectionScore = (index, value) => {
        const skillsSection = JSON.parse(state.skills);
        skillsSection[index].score += value;
        if (skillsSection[index].score > 100) {
            skillsSection[index].score = 100;
        } else if (skillsSection[index].score < 0) {
            skillsSection[index].score = 0;
        }
        setState({
            ...state,
            skills: JSON.stringify(skillsSection),
        });
    };

    const getSkillSectionEdit = (index, score, skills) => {
        return <Box sx={{display: 'flex', flexDirection: 'row'}}>
        <Box sx={{display: 'flex'}} style={{height: '70px', minWidth: '170px'}}>
            <IconButton aria-label="add" size="medium" onClick={() => {sumSkillSectionScore(index, -5)}}>
                <RemoveCircleOutlineIcon color="primary" fontSize="inherit" />
            </IconButton>
            <Gauge
                value={parseInt(score)}
                startAngle={-110}
                endAngle={110}
                sx={{
                    [`& .${gaugeClasses.valueText}`]: {
                    fontSize: 12,
                    transform: 'translate(0px, 0px)',
                    },
                }}
                text={
                    ({ value, valueMax }) => `${value} / ${valueMax}`
                }/>
                <IconButton aria-label="add" size="medium" onClick={() => {sumSkillSectionScore(index, 5)}}>
                    <AddCircleOutlineIcon color="primary" fontSize="inherit" />
                </IconButton>
        </Box>
        <Box style={{display: 'flex', justifyContent: 'center', alignItems: 'center', width: '100%'}}>
            <SkillAutocomplete defoultOptions={skills} onChange={(updatedSkills) => {onSkillChange(index, updatedSkills)}}/>
        </Box>
    </Box>
    };

    const getSkillSection = (percentile, skills) => {
        const skillsChips = [];
        for (let i = 0; i < skills.length; i++) {
            skillsChips.push(<Chip label={skills[i]} variant="outlined" sx={{m: '2px'}}/>)
        }
        return <Box sx={{display: 'flex', flexDirection: 'row'}}>
            <Box sx={{display: 'flex'}} style={{height: '70px', minWidth: '100px'}}>
                <Gauge
                    value={parseInt(percentile)}
                    startAngle={-110}
                    endAngle={110}
                    sx={{
                        [`& .${gaugeClasses.valueText}`]: {
                        fontSize: 12,
                        transform: 'translate(0px, 0px)',
                        },
                    }}
                    text={
                        ({ value, valueMax }) => `${value} / ${valueMax}`
                    }/>
            </Box>
            <Box style={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                <Box style={{display: 'inline-block', justifyContent: 'center', alignItems: 'center'}} sx={{ml: 3}}>
                    {skillsChips}
                </Box>
            </Box>
        </Box>
    };

    const addNewSkillSet = () => {
        const skills = state.skills ? JSON.parse(state.skills) : [];
        skills.push({
            score: 75,
            skills: [],
        })
        setState({
            ...state,
            skills: JSON.stringify(skills),
        });
    };

    const SkillsSection = function({id, title, content}) {
        let empty = false;
        if (!content) {
            content = '[]'
            empty = true;
        }
        //[{score: 75, skills: ["java"]}]
        const skillSectionsResult = [];
        const skillSections = JSON.parse(content);
        for (let i = 0; i < skillSections.length; i++) {
            const skillSection = skillSections[i];
            skillSectionsResult.push(state.editSection === id ? getSkillSectionEdit(i, skillSection.score, skillSection.skills) : getSkillSection(skillSection.score, skillSection.skills));
        }

        return <Paper elevation={3} sx={{p: 1, m:3}}>
            <Box sx={{display: 'flex', flexDirection: 'row'}}>
                <Box sx={{display: 'flex', flexGrow: 1}}>
                    <Typography variant="h5" gutterBottom>
                        {title}
                    </Typography>
                </Box>
                <Box sx={{display: 'flex'}} >
                    <IconButton aria-label="edit" onClick={() => {onEditClicked(id)}}>
                        <EditIcon color="primary"/>
                    </IconButton>
                </Box>
            </Box>
            <Divider sx={{mb: 2}}/>
            {(empty && state.editSection !== id) && <Box sx={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}><Typography variant="subtitle1" gutterBottom>This section is empty</Typography></Box>}
            {skillSectionsResult}
            {state.editSection === id && <Box sx={{justifyContent: 'center', display: 'flex'}}>
                <IconButton aria-label="add" size="large" onClick={addNewSkillSet}>
                    <AddIcon color="primary" fontSize="inherit" />
                </IconButton>
            </Box>}
            {state.editSection === id && <Box sx={{justifyContent: 'right', display: 'flex'}}>
                <Button variant="contained" endIcon={<CancelIcon />} sx={{mr:1}} onClick={cancelSaveSection}>
                    Cancel
                </Button>
                <Button variant="contained" endIcon={<SaveIcon />} onClick={() => {saveSection(id, state.skills)}}>
                    Save
                </Button>
            </Box>}
        </Paper>;
    };

    return <Box>
        {data && data.reporterUsers && data.reporterUsers.length > 0 && <RelationshipSection id="reporters" title={'Reporters'} users={data.reporterUsers}/>}
        {data && data.managerUsers && data.managerUsers.length > 0 && <RelationshipSection id="reporters" title={'Managers'} users={data.managerUsers}/>}
        <AboutSection id="about_me" title={'About me'} content={state?.about_me}/>
        <AboutSection id="responsibilities" title={'Responsabilities'} content={state?.responsibilities}/>
        <SkillsSection id="skills" title={'Skills'} content={state?.skills}/>
    </Box>
}