import React, { useState, forwardRef, useRef } from 'react';
import Grid from '@mui/material/Grid2';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import { DecimalPrecision } from '../../amper/Instruments'; 
import Divider from '@mui/material/Divider';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import { formatVersion } from '../../amper/Instruments';
import Modal from '@mui/material/Modal';
import CircularProgress from '@mui/material/CircularProgress';

function FileDetail({metadata, version, onVersionChange}, ref) {

    const [state, setState] = useState(() => ({
        version: {
            label: formatVersion(metadata.version),
            value: metadata.version,
        },
        loading: false,
    }));

    const getFileMetadatas = () => {
        const result = [];
        const lastModified = new Date(metadata.lastModified);
        result.push(<Grid key={'version'} size={1} md={1}>
            <Typography variant="subtitle2" gutterBottom>
               Version
            </Typography>
        </Grid>);
        result.push(<Grid key={'version_value'} size={1} md={2}>
            {getVersionSelector()}
        </Grid>);
        const displayValues = [
            /*{
                key: 'version',
                label: 'Version',
                value: (metadata.version.major + '.' + metadata.version.minor + '.' + metadata.version.patch)
            },*/{
                key: 'size',
                label: 'Size',
                value: (DecimalPrecision.ceil(((metadata.size / 1024) / 1024), 2) + ' mb')
            },{
                key: 'type',
                label: 'Type',
                value: metadata.type
            },{
                key: 'lastModified',
                label: 'Last modified',
                value: (lastModified.toLocaleDateString() + ' ' + lastModified.toLocaleTimeString())
            },{
                key: 'rendition',
                label: 'Has rendition',
                value: (metadata.rendition ? ('Yes "' + metadata.renditionType + '"') : 'No')
            },/*{
                key: 'thumbnail',
                label: 'Has thumbnail',
                value: (metadata.thumbnail ? 'Yes "image/png"' : 'No')
            }*/
        ];
        for (let i = 0; i < displayValues.length; i++) {
            if (displayValues[i].display !== false) {
                result.push(<Grid key={displayValues[i].key} size={1} md={1}>
                    <Typography variant="subtitle2" gutterBottom>
                       {displayValues[i].label}
                    </Typography>
                </Grid>);
                result.push(<Grid key={displayValues[i].key + '_value'} size={1} md={2}>
                    <Typography variant="subtitle2" gutterBottom>
                        {displayValues[i].value}
                    </Typography>
                </Grid>);
            }
        }
        return result;
    };

    const getFileExifMetadatas = () => {
        const result = [];
        if (metadata.exifMetadata && metadata.exifMetadata.length > 0) {
            for (let i = 0; i < metadata.exifMetadata.length; i++) {
                const exifMetadata = metadata.exifMetadata[i];
                result.push(<Grid key={exifMetadata.id + '_' + i} xs={1} md={1}>
                    <Typography variant="subtitle2" gutterBottom sx={{overflowWrap: 'break-word'}}>
                       {exifMetadata.name}
                    </Typography>
                </Grid>);
                result.push(<Grid key={exifMetadata.id + '_' + i + '_value'} xs={1} md={2}>
                    <Typography variant="subtitle2" gutterBottom sx={{overflowWrap: 'break-word'}}>
                        {exifMetadata.value}
                    </Typography>
                </Grid>);
            }
        } else {
            result.push(<Grid key={'noExifData'} xs={1} md={3}>
                <Typography variant="caption" gutterBottom>
                   No exif data
                </Typography>
            </Grid>);
        }
        return result;
    };

    const handleVersionChange = (event, value, arg) => {
        debugger
    };

    const getVersionSelector = () => {
        const versionOptions = [];
        if (metadata.availableVersions) {
            for (let i = 0; i < metadata.availableVersions.length; i++) {
                const formattedVersion = formatVersion(metadata.availableVersions[i]);
                versionOptions.push({
                    label: formattedVersion,
                    value: metadata.availableVersions[i],
                });
            }
        }
        return (
            <Autocomplete
              disablePortal
              value={state.version}
              options={versionOptions}
              isOptionEqualToValue={(option, value) => option.label === value.label}
              onChange={(event, selectedOption) => {
                if (selectedOption && (selectedOption.value.major != state.version.value.major || 
                    selectedOption.value.minor != state.version.value.minor || 
                    selectedOption.value.patch != state.version.value.patch)) {
                    setState({
                        ...state,
                        version: selectedOption,
                        loading: true,
                    });
                    if (onVersionChange) {
                        onVersionChange(selectedOption.value);
                    }
                }
              }}      
              sx={{ width: 150 }}
              size="small"
              renderInput={(params) => <TextField {...params} size="small"/>}
            />
          );
    };

    const getProgress = () => {
        return <Box sx={{ display: 'flex', width: '100%', height: (window.innerHeight - 350) + 'px', verticalAlign: 'middle', alignItems: 'center', justifyContent: 'center' }}>
            <CircularProgress />
        </Box>;
    };

    const getContent = () => {
        return [
            <Box sx={{overflowX: 'hidden', overflowY: 'auto', height: '150px',}}>
                <Grid key="fileMetadataList" container columns={2} spacing={0}>
                    {getFileMetadatas()}
                </Grid>
            </Box>,
            <Divider sx={{mb: 1}}/>,
            <Typography variant="button" gutterBottom>
                Exif metadata {'{}'}
            </Typography>,
            <Box key="fileExifMetadataList" sx={{overflowX: 'hidden', overflowY: 'auto', height: (window.innerHeight - 350) + 'px',}}>
                <Grid container columns={3} spacing={0}>
                    {getFileExifMetadatas()}
                </Grid>
            </Box>];
    };

    return <Box key="fileDetail" sx={{ flexGrow: 1, m: 1 }}>
        <Typography key="fileDetailName" variant="h6" gutterBottom>
            {metadata.name}
            </Typography>
        {state.loading ? getProgress() : getContent()}
    </Box>;
};

export default forwardRef(FileDetail)