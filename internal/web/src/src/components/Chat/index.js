import * as React from 'react';
import {Fragment, useEffect, useState} from 'react';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Paper from '@mui/material/Paper';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import {ThemeProvider} from '@mui/material/styles';
import {GetProfile, UpsertProfile} from "../../api/api";
import {Alert} from "@mui/material";
import {WSStart} from "../../api/chat";

const ChatWindow = ({theme, accessToken}) => {
    const [success, setSuccess] = useState(null)
    const [failure, setFailure] = useState(null)

    const [nickname, setNickname] = useState("")
    const [aboutMe, setAboutMe] = useState("")
    const [message, setMessage] = useState("")

    const [messages, setMessages] = useState([])
    const [sendMessage, setSendMessage] = useState({send: () => {
    }})


    const reset = () => {
        setFailure(null)
        setSuccess(null)
    }

    const addMessage = (message) => {
        messages.push(message)
        setMessages([...messages])
    }

    useEffect(() => {
        GetProfile(accessToken, ({nickname, about_me}) => {
            setSuccess("Successfully updated")

            setNickname(nickname)
            setAboutMe(about_me)
        }, setFailure)

        WSStart(accessToken, addMessage, setSendMessage)
    }, []);


    const handleSubmit = (event) => {
        event.preventDefault();
        reset()
        const data = new FormData(event.currentTarget);

        UpsertProfile(accessToken, data.get("nickname"), data.get("about_me"), ({nickname, about_me}) => {
            setSuccess("Successfully updated")
            setNickname(nickname)
            setAboutMe(about_me)
        }, setFailure)
    };


    return (
        <ThemeProvider theme={theme}>
            {
                failure ? <Alert severity="error">{failure}</Alert> : null
            }
            {
                success ? <Alert severity="success">{success}</Alert> : null
            }

            <Grid container component="main" sx={{height: '100vh'}}>
                <CssBaseline/>
                <Grid
                    item
                    xs={false}
                    sm={4}
                    md={7}
                >
                    <div style={{
                        height: "100vh",
                        display: "flex",
                        flexDirection: "column-reverse",
                        justifyContent: "start"
                    }}>
                    <Box component="form" noValidate onSubmit={(e) => {
                        e.preventDefault();
                        sendMessage.send(nickname, message)
                    }} sx={{mt: 1}}>
                        <TextField style={{minWidth: "100%"}}
                            id="outlined-basic" label="Message" value={message} onChange={(e) => setMessage(e.target.value)} variant="outlined"/>
                    </Box>
                        <div style={{
                            overflow: "scroll",
                            maxHeight: "calc(100% - 300px)"
                        }}>
                            {
                                messages.map(({message, nickname}) => {
                                    if (message === "") {
                                        return null
                                    }

                                    return (
                                        <div style={{textAlign: "start", padding: "0.5rem 1rem"}}>
                                            {nickname}:{message}
                                        </div>
                                    )
                                })
                            }
                        </div>
                    </div>

                </Grid>
                <Grid item xs={12} sm={8} md={5} component={Paper} elevation={6} square>
                    <Box
                        sx={{
                            my: 8,
                            mx: 4,
                            display: 'flex',
                            flexDirection: 'column',
                            alignItems: 'center',
                        }}
                    >
                        <Typography component="h1" variant="h5">
                            Profile
                        </Typography>
                        <Box component="form" noValidate onSubmit={handleSubmit} sx={{mt: 1}}>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                id="nickname"
                                label="Nickname"
                                name="nickname"
                                value={nickname}
                                onChange={(e) => {
                                    setNickname(e.target.value);
                                }}
                                autoFocus
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                name="about_me"
                                label="About Me"
                                value={aboutMe}
                                onChange={(e) => {
                                    setAboutMe(e.target.value);
                                }}
                                id="about_me"
                            />
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                sx={{mt: 3, mb: 2}}
                            >
                                Update
                            </Button>
                        </Box>
                    </Box>
                </Grid>
            </Grid>
        </ThemeProvider>
    );
}

export {
    ChatWindow
}