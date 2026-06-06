import * as React from 'react';
import Box from '@mui/material/Box';
import Tab from '@mui/material/Tab';
import TabContext from '@mui/lab/TabContext';
import TabList from '@mui/lab/TabList';
import TabPanel from '@mui/lab/TabPanel';
import LabelIcon from '@mui/icons-material/Label';
import TypeSpecimenIcon from '@mui/icons-material/TypeSpecimen';
import ListAltIcon from '@mui/icons-material/ListAlt';
import Details from './Details';
import ObjectTypes from './ObjectTypes';
import Fields from './Fields';

export default function Edit({toast}) {

    const [value, setValue] = React.useState('0');
    const handleChange = (event, newValue) => {
        setValue(newValue);
      };
    return (
        <Box sx={{height: '100%', width: 'calc(100% - 35px)'}}>
            <TabContext value={value}>
                    <Box sx={{ borderBottom: 1, mt: -3, borderColor: 'divider' }}>
                        <TabList centered onChange={handleChange} aria-label="lab API tabs example">
                            <Tab icon={<LabelIcon color="inherit" />} iconPosition="start" label="Details" value="0"/>
                            <Tab icon={<TypeSpecimenIcon color="inherit" />} iconPosition="start" label="Object Types" value="1"/>
                            <Tab icon={<ListAltIcon color="inherit" />} iconPosition="start" label="Fields" value="2"/>
                        </TabList>
                    </Box>
                    <TabPanel value="0" sx={{ height: 'calc(100% - 72px)', width: '100%'}}><Details/></TabPanel>
                    <TabPanel value="1" sx={{ height: 'calc(100% - 72px)', width: '100%'}}><ObjectTypes/></TabPanel>
                    <TabPanel value="2" sx={{ height: 'calc(100% - 72px)', width: '100%'}}><Fields toast={toast}/></TabPanel>
            </TabContext>
        </Box>
    );
}
