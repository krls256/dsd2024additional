import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';

import Link from '@mui/material/Link';
import Paper from '@mui/material/Paper';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import { ThemeProvider } from '@mui/material/styles';
import {useEffect, useState} from "react";
import {GetProfile, Login, UpsertProfile} from "../../api/api";
import {Alert} from "@mui/material";

const ChatWindow = ({theme, accessToken}) => {
    const [profile, setProfile] = useState(null)
    const [success, setSuccess] = useState(null)
    const [failure, setFailure] = useState(null)


    const reset = () => {
        setFailure(null)
        setSuccess(null)
    }

    useEffect(() => {
        GetProfile(accessToken, ({nickname, about_me}) => {
            setSuccess("Successfully updated")
            setProfile({
                nickname: nickname,
                about_me: about_me
            })
        }, setFailure)
    }, []);


    const handleSubmit = (event) => {
        event.preventDefault();
        reset()
        const data = new FormData(event.currentTarget);

        UpsertProfile(accessToken, data.get("nickname"), data.get("about_me"),({nickname, about_me}) => {
            setSuccess("Successfully updated")
            setProfile({
                nickname: nickname,
                about_me: about_me
            })
        }, setFailure)


        // Login(data.get('email'), data.get('password'), goChat, setFailure)
    };

    return (
        <ThemeProvider theme={theme}>
            {
                failure ? <Alert severity="error">{failure}</Alert> : null
            }
            {
                success ? <Alert severity="success">{success}</Alert> : null
            }

            <Grid container component="main" sx={{ height: '100vh' }}>
                <CssBaseline />
                <Grid
                    item
                    xs={false}
                    sm={4}
                    md={7}
                >
                    <h1>Chat</h1>
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
                        <Box component="form" noValidate onSubmit={handleSubmit} sx={{ mt: 1 }}>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                id="nickname"
                                label="Nickname"
                                name="nickname"
                                autoFocus
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                name="about_me"
                                label="About Me"
                                id="about_me"
                            />
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                sx={{ mt: 3, mb: 2 }}
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