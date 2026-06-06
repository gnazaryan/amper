import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';

export const ObjectTypeField = ({ name, value, label, required, error, onChange, objectTypes }) => {

    return (
        <Autocomplete
                variant="standard"
                value={value}
                onChange={onChange}
                multiple={false}
                options={objectTypes}
                fullWidth={true}
                isOptionEqualToValue={(option, value) => option.id === value.id}
                renderInput={(params) => <TextField variant="standard" {...params} label={label} name={name} error={error} required={required}/>}
            />
      );
}