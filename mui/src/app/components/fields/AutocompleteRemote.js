import * as React from 'react';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import CircularProgress from '@mui/material/CircularProgress';
import { post } from '../../data/Submit';
import HostManager from '../../../HostManager';

const AutocompleteRemote = ({name, value, required, error, url, parameters, label, keyIdentifier, labelIdentifier, onChange}) => {
    const [open, setOpen] = React.useState(false);
    const [options, setOptions] = React.useState([]);
    const loading = open && options.length === 0;

    React.useEffect(() => {
        if (loading) {
            post(url, {
                ...parameters
              }, (result) => {
                setOptions(result.data || []);
              }, (result) => {
                setOptions([]);
              })
        }
    }, [loading]);

    return <Autocomplete
        fullWidth
        open={open}
        onOpen={() => {
            setOpen(true);
        }}
        onClose={() => {
            setOpen(false);
        }}
        sx={{mt: 0, mb: 0}}
        isOptionEqualToValue={(option, value) => option[keyIdentifier] === value[keyIdentifier]}
        getOptionLabel={(option) => option[labelIdentifier]}
        options={options}
        loading={loading}
        value={value}
        onChange={(event, value) => {
            onChange(name, event, value);
        }}
        renderInput={(params) => (
            <TextField
              name={name}
              {...params}
              label={label || 'Label'}
              sx={{mt: 1, mb: 0}}
              required={required}
              error={error}
              variant="filled"
              InputProps={{
                ...params.InputProps,
                endAdornment: (
                  <React.Fragment>
                    {loading ? <CircularProgress color="inherit" size={20} /> : null}
                    {params.InputProps.endAdornment}
                  </React.Fragment>
                ),
              }}
            />
          )}
        />
};
export default AutocompleteRemote;