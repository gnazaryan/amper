import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Amper from './Amper';
import * as serviceWorker from './serviceWorker';

window.eventRegistry = {};

(function() {
    const fetch = window.fetch;
    window.fetch = (...args) => (async(args) => {
        var result = await fetch(...args);
        return result;
    })(args);
})();

ReactDOM.render(<Amper />, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
