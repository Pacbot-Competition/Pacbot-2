<script>

  import config from '../../config.json';
  import Maze from './lib/Maze.svelte';
  import Pellets from './lib/Pellets.svelte';

  var socket = new WebSocket(`ws://${config.ServerIP}:${config.WebSocketPort}`);
  socket.binaryType = 'arraybuffer';
  let message = 'Offline';

  var grid = [];

  for (let row = 0; row < 31; row++) {
    grid[row] = [];
    for (let col = 0; col < 28; col++) {
      grid[row][col] = 0;
    }
  }

  socket.addEventListener('open', (_) => {
    message = 'Online!';
    console.log('WebSocket connection established');
    socket.send('Hello, server!');
  });

  socket.addEventListener('message', (event) => {
    if (event.data instanceof ArrayBuffer) {
      // binary frame
      let view = new DataView(event.data);
      
      if (view) {
        for (let row = 0; row < 31; row++) {
          let binRow = view.getUint32(4*row, false);
          for (let col = 0; col < 28; col++) {
            let superPellet = ((row === 3) || (row === 23)) && ((col === 1) || (col === 26));
            grid[row][col] = ((binRow >> col) & 1) ? (superPellet ? 2 : 1) : 0;
          }
        }
      }

      grid = grid;
    }
  });

  socket.addEventListener('close', (_) => {
    console.log('WebSocket connection closed');
  });

</script>

<!--<h1> {message} </h1>-->

<Maze></Maze>

<Pellets bind:grid/>