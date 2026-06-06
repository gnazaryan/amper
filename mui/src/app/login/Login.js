import * as React from 'react';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid2';
import Box from '@mui/material/Box';
import AmperIcon from '../icons/Amper';
import AmperTransparentIcon from '../icons/AmperTransparent';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { createTheme, ThemeProvider, StyledEngineProvider, adaptV4Theme } from '@mui/material/styles';
import Convenience from '../help/Convenience.js';
import HostManager from "../../HostManager";
import { useState } from 'react';

export default function Login({hooks}) {
  const {success} = hooks;
  const [state, setState] = useState({
    errorMessage: null,
  });

    const Copyright = (props) => {
        return (
          <Typography variant="body2" color="text.secondary" align="center" {...props}>
            {'Copyright © '}
            <Link color="inherit" href="https://amper.cloud/" underline="hover">
              Amper
            </Link>{' '}
            {new Date().getFullYear()}
            {'.'}
          </Typography>
        );
    }

    const theme = createTheme(adaptV4Theme({
        palette: {
          primary: {
            main: '#73A8EB',
            secondary: '#FFFFFF',
            borderRadius: 3,
            contrastText: '#FFFFFF',
          },
        }
      }));

    const handleSubmit = (event) => {
        event.preventDefault();
        const data = new FormData(event.currentTarget);

        const inputValues = {
          username: data.get('email'),
          password: data.get('password'),
        };

        if (!Convenience.containsNullOrEmpty(inputValues, ['username', 'password'])) {
          setState({
              errorMessage: 'Both username and password are required fields',
          });
          return;
        }
        fetch(`${HostManager.amperHost()}users/login`, {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({
              ...inputValues
          })
        })
        .then(res => res.json())
        .then((result) => {
            if (result) {
                if (result.success) {
                  success(result.user, result.settings);
                } else {
                    setState({
                        errorMessage: result.error,
                    });
                }
            } else {
                setState({
                    errorMessage: 'Something went wrong, please contact your service provider for more details',
                });
            }
        })
      };

    return (
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={theme}>
          <Container component="main" maxWidth="xs">
            <CssBaseline />
            <Box
              sx={{
                marginTop: 14,
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
              }}
            >
              <AmperIcon color='primary' sx={{ width: 70, height: 70, m: 1, color: 'primary.main' }}/>
                <Typography component="h1" variant="h5">
                  Welcome to Amper
                </Typography>
              <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
                <TextField
                  variant="standard"
                  margin="normal"
                  required
                  fullWidth
                  id="email"
                  label="Username"
                  name="email"
                  autoComplete="email"
                  autoFocus />
                <TextField
                  variant="standard"
                  margin="normal"
                  required
                  fullWidth
                  name="password"
                  label="Password"
                  type="password"
                  id="password"
                  autoComplete="current-password"
                  error={state.errorMessage != null}
                  helperText={state.errorMessage} />
                <FormControlLabel control={<Checkbox icon={<AmperTransparentIcon color='primary' sx={{ color: 'primary.main' }}/>} value="remember" color="primary" />}
                  label="Remember me"
                />
                <Button type="submit"
                  fullWidth
                  variant="contained"
                  sx={{ mt: 3, mb: 2 }}
                >
                  Sign In
                </Button>
                <Grid container>
                  <Grid item xs>
                    <Link href="#" variant="body2" underline="hover">
                      Forgot password?
                    </Link>
                  </Grid>
                </Grid>
              </Box>
            </Box>
            <Copyright sx={{ mt: 8, mb: 4 }} />
          </Container>
        </ThemeProvider>
      </StyledEngineProvider>
    );
}
