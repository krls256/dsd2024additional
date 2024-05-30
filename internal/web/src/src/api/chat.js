

const WSStart = (accessToken, addMessage, setSender) => {
    const socket = new WebSocket(`ws://localhost:6900/ws?Bearer=${accessToken}`);

    socket.addEventListener("message", (event) => {
        const data = JSON.parse(event.data)

        addMessage(data.data)
    });


    setSender({send: (nickname, message) => {
            socket.send(JSON.stringify({
                "action": "push",
                "payload": {
                    "nickname": nickname,
                    "message": message
                }
            }))
        }})
}

export {WSStart}