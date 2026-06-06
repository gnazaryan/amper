import * as React from 'react';
import Chip from '@mui/material/Chip';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import DataStore from '../../data/DataStore';
import HostManager from '../../../HostManager';
import Box from '@mui/material/Box';
import Convenience from '../../help/Convenience';
import Avatar from '@mui/material/Avatar';
import { debounceLatest } from '../../amper/Instruments';
import { sessionManager } from '../../../SessionManager';

export default function UserSelect({onSelectionChange, includeSelf, singleSelect}) {

    const [value, setValue] = React.useState([]);

    const debounceReload = (arg0, search) => {
        setState({
            ...state,
            loading: true,
            search: search.split(" "),
            input: search,
        });
    };

    const [state, setState] = React.useState(() => {
        return {
            loading: true,
            start: 0,
            limit: 50,
            data: [],
            search: [],
            debounceReload: debounceLatest(debounceReload, 1000),
            input: '',
        };
    });

    React.useEffect(() => {
        if (state.loading) {
            getDataStore().load((result)=> {
                const currentUser = sessionManager.getUser();
                let data = null;
                if (includeSelf === true) {
                  data = result.data;
                } else {
                  data = result.data.filter(user => user.id != currentUser.id);
                }

                setState({
                  ...state,
                  loading: false,
                  data: data || [],
                });
            });
        }
    }, [state.loading]);

    const getDataStore = () => {
        return new DataStore({
            url: `${HostManager.amperHost()}users/fetch`,
            requestMethod: "POST",
            parameters: {
                start: state.start,
                limit: state.limit,
                search: state.search,
            }
        });
    }

    const getImageSource = (user) => {
        if (Convenience.hasValue(user.photo)) {
          return 'data:image/png;base64,' + user.photo;
        }
        return '/static/images/avatar/2.jpg';
      };

    return (<Autocomplete
        multiple
        id="fixed-tags-demo"
        value={value}
        sx={{mt: 1}}
        onChange={(event, newValue) => {
          let finalValue = newValue;
          if (singleSelect && newValue.length > 0) {
            finalValue = [newValue[newValue.length - 1]];
          }
          setValue([
            ...finalValue,
          ]);
          onSelectionChange(finalValue);
        }}
        options={state.data}
        getOptionLabel={(option) => (option.firstName + ' ' + option.lastName)}
        renderTags={(tagValue, getTagProps) =>
          tagValue.map((option, index) => {
            const { key, ...tagProps } = getTagProps({ index });
            return (
              <Chip
                avatar={<Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 2 }} alt={option.firstName + ' ' + option.lastName} src={getImageSource(option)} />}
                key={key}
                label={(option.firstName + ' ' + option.lastName)}
                {...tagProps}
              />
            );
          })
        }
        style={{ width: 500 }}
        renderInput={(params) => {
          return <TextField
            onChange={(event) => {
                const {name, value} = event.target;
                state.debounceReload(null, value);
            }}
            onBlur={() => {
                debounceReload(null, '');
            }}
            {...params} 
            value={state.input}
            label="People" 
            placeholder="Search" />
        }}
        renderOption={(props, option) => {
            const { key, ...optionProps } = props;
            return (
              <Box
                key={key}
                component="li"
                sx={{ '& > img': { mr: 2, flexShrink: 0 }, m: 0 }}
                {...optionProps}
              >
                <Avatar sx={{ bgcolor: 'secondary.main', color: 'primary.main', mr: 2 }} alt={option.firstName + ' ' + option.lastName} src={getImageSource(option)} />

                {option.firstName} {option.lastName}
              </Box>
            );
          }}
      />);
}