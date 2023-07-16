<script>

  import config from '../../config.json';

  var socket = new WebSocket(`ws://${config.ServerIP}:${config.WebSocketPort}`);
  let message = 'Offline';

  socket.addEventListener('open', (_) => {
    message = 'Online!';
    console.log('WebSocket connection established');
    socket.send('Hello, server!');
  });

  socket.addEventListener('message', (event) => {
    // console.log(`Received: ${event.data}`);
    message = event.data;
  });

  socket.addEventListener('close', (_) => {
    console.log('WebSocket connection closed');
  });

</script>

<h1> {message} </h1>