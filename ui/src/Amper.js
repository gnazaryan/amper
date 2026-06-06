import React from 'react';
import './Amper.css';
import Convenience from "./components/help/Convenience";
import Login from './components/login/Login.js';
import Activation from './activation/Activation';
import UI from './components/ui/ui.js';
import './components/Main.css'
import {sessionManager} from "./SessionManager";

export default class Amper extends React.Component {

  constructor(props) {
    super(props);
    this.state = {sessionId: sessionManager.getSessionId()};
  }

  updateState(user) {
      sessionManager.setUser(user);
      this.setState({
          sessionId: user.sessionId,
      });
  }

    handleLogOut() {
        this.setState({
            sessionId: null,
        });
    }

    render() {
        let face = 0;
        const activationCode = Convenience.getUrlParameterValue('activationCode');
        if (activationCode != null) {
            face = 2;
        } else if (this.state.sessionId) {
            face = 1;
        }
        return (
            <div className="amper">
                {(function(face, me) {
                    switch(face) {
                        case 0:
                        return <Login parent={me}></Login>;
                        case 1:
                        return<UI parent={me}></UI>;
                        case 2:
                            return <Activation activationCode={activationCode} parent={me}></Activation>;
                    }
                })(face, this)}
            </div>
        );
    }
}
