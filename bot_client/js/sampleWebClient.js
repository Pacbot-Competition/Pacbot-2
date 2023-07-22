config = require('../../config.json');
const WebSocket = require('ws');

var socket = new WebSocket(`ws://${config.ServerIP}:${config.WebSocketPort}`);

socket.addEventListener('open', (_) => {
  message = 'Online!';
  console.log('WebSocket connection established');
  socket.send('Hello, server!');
});

socket.addEventListener('message', (event) => {
  console.log(`Received: ${event.data}`);
  message = event.data;
});

socket.addEventListener('close', (_) => {
  console.log('WebSocket connection closed');
});