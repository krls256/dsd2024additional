import './App.css';
import {SignInSide} from "./components/SignInSide";
import {SignUpSide} from "./components/SignUpSide";
import {ChatWindow} from "./components/Chat";
import {useState} from "react";
import {createTheme} from '@mui/material/styles';

const SignIn = 1
const SignUp = 2
const Chat = 3

const defaultTheme = createTheme();


function App() {
    const [state, setState] = useState(SignIn)
    const [accessToken, setAccessToken] = useState(null)

    const GoSignIn = () => {
        setState(SignIn)
    }

    const GoSignUp = () => {
        setState(SignUp)
    }

    const GoChat = () => {
        setState(Chat)
    }

    return (
        <div className="App">
            {state === SignIn ? <SignInSide goSignUp={GoSignUp} goChat={GoChat} theme={defaultTheme}
                                            setAccessToken={setAccessToken}/> : null}
            {state === SignUp ? <SignUpSide goSignIn={GoSignIn} goChat={GoChat} theme={defaultTheme}/> : null}
            {state === Chat ? <ChatWindow theme={defaultTheme} accessToken={accessToken}/> : null}
        </div>
    );
}

export default App;
