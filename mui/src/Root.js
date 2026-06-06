import PropTypes from 'prop-types';
import { BrowserRouter } from 'react-router-dom';
import { StaticRouter } from 'react-router-dom/server';
import { ThemeProvider, StyledEngineProvider, createTheme, adaptV4Theme } from '@mui/material/styles';
import App from './App';

function Root() {

    function Router(props) {
        const { children } = props;
        if (typeof window === 'undefined') {
          return <StaticRouter location="/">{children}</StaticRouter>;
        }
      
        return <BrowserRouter>{children}</BrowserRouter>;
      }
      Router.propTypes = {
        children: PropTypes.node,
      };

      const theme = createTheme(adaptV4Theme({
        palette: {
          primary: {
            main: '#2196f3',
            label: '#8d8e8f',
            text: '#000000',
            warn: '',
            borderRadius: 1,
            contrastText: '#ffffff',
            dirty: '#b0b0b0',
            selectedBackground: '#d9ecfc',
          },
          secondary: {
            main: '#ffffff',
            contrastText: '#2196f3',
            gridText: '#424242',
            menuText: '#424242',
          },
          inactive: {
            main: '#c8c9cc',
            contrastText: '#808080',
          }
        },
        typography: {
          fontWeightLight: 'bold'
        },
        components: {
          MuiTooltip: {
              styleOverrides: {
                  tooltip: {
                      fontSize: '1em'
                  }
              }
          }
        }
      }));
      return (
        <StyledEngineProvider injectFirst>
          <ThemeProvider theme={theme}>
              <Router>
                  <App></App>
              </Router>
          </ThemeProvider>
        </StyledEngineProvider>
      );
}

export default Root;