import React from 'react';
import logo from './logo.svg';
import './App.css';
import { useAppStore } from './AppContext';

const App: React.FC = () => {
  const { state, dispatch } = useAppStore();
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>

        <p>Authenticated? {state.user.authenticated ? "true" : "false"} {state.user.username}</p>
        <button onClick={() => dispatch({
          type: "authenticated",
          user: {
            username: state.user.authenticated ? "" : "heylookafakeusername",
            authenticated: !state.user.authenticated
          }
        })}>
          Fake {state.user.authenticated ? "Logout" : "Login"}
        </button>
      </header>
    </div>
  );
}

export default App;
