import React from 'react';
import { Button } from '@material-ui/core';
import logo from './logo.svg';
import './App.css';
import { useAppStore } from './AppContext';

const App: React.FC = () => {
  const { state, dispatch } = useAppStore();
  return (
    <div className="App">
      <header className="App-header">
        <p>Authenticated? {state.user.authenticated ? "true" : "false"}</p>
        <Button variant="contained" color="primary" onClick={() => dispatch({
          type: "authenticated",
          user: {
            username: state.user.authenticated ? "" : "heylookafakeusername",
            authenticated: !state.user.authenticated
          }
        })}>
          {state.user.authenticated ? `Logout: ${state.user.username}` : "Login"}
        </Button>
      </header>
    </div>
  );
}

export default App;
