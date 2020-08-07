import React from 'react';
import { Topbar } from './Topbar';
import { Main } from './Main';
import { TOKEN_KEY } from '../constants';
import '../styles/App.css';

class App extends React.Component {
  state = {
      //check if has token (value) based on key => true
      isLoggedIn: Boolean(localStorage.getItem(TOKEN_KEY)),
  }

  handleLogin = (token) => {
      //key, value
      localStorage.setItem(TOKEN_KEY, token);

      this.setState({
          isLoggedIn: true,
      });
  }

  handleLogout = () => {
      //remove by key
      localStorage.removeItem(TOKEN_KEY);

      this.setState({
          isLoggedIn: false,
      });
  }

  render() {
    return (
        <div className="App">
            <Topbar isLoggedIn={this.state.isLoggedIn} handleLogout={this.handleLogout} />
            <Main isLoggedIn={this.state.isLoggedIn} handleLogin={this.handleLogin} />
        </div>
    )
  }
}

// function App() {
//   return (
//       <div className="App">
//         <Topbar />
//         <Main />
//       </div>
//   );
// }

export default App;
