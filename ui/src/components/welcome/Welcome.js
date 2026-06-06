import React from 'react';
import './Welcome.css';

function Welcome(props) {
  return (
    	<div className="welcomeMainContainer">
			<img src="/images/logo1.png" width={"60px"} height={"41px"}/>
			<div className="welcomeTextContainer">Welcome to Amper {props.name || ''}</div>
		</div>
  );
}

export default Welcome;
