<script>

  import config from '../../config.json';
  import Maze from './lib/Maze.svelte';
  import Pellets from './lib/Pellets.svelte';

  var socket = new WebSocket(`ws://${config.ServerIP}:${config.WebSocketPort}`);
  let message = 'Offline';

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

  var grid = [];
  for (let i=0; i<31; i++) {
    grid[i] = [];
    for (let j=0; j<28; j++) {
      grid[i][j] = true;
    }
  }

</script>

<!--<h1> {message} </h1>-->

<Maze></Maze>
<Pellets bind:grid={grid}/>