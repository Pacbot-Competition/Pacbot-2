<style>

  .maze-space {
    position: absolute;
    top: 5vh;
    left: 5vw;
  }

</style>

<script>

  import config from '../../config.json';
  import Maze from './lib/Maze.svelte';
  import Pellets from './lib/Pellets.svelte';
  import Pacman from './lib/Pacman.svelte';
  import Ghosts from './lib/Ghosts.svelte';

  var socket = new WebSocket(`ws://${config.ServerIP}:${config.WebSocketPort}`);
  socket.binaryType = 'arraybuffer';
  let message = 'Offline';

  let grid = [];

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

  let innerWidth = 0;
  let innerHeight = 0;

  let gridSize;
  $: gridSize = 0.9 * ((innerHeight * 28 < innerWidth * 31) ? (innerHeight / 31) : (innerWidth / 28));

  let pacmanRow = 23;
  let pacmanCol = 14;

  let redRow = 11;
  let redCol = 14;

  let pinkRow = 14;
  let pinkCol = 14;

  let blueRow = 14;
  let blueCol = 12;

  let orangeRow = 14;
  let orangeCol = 16;

</script>

<svelte:window bind:innerWidth bind:innerHeight />

<div class='maze-space'>
  <Maze {gridSize}/>
  <Pellets {grid} {gridSize}/>
  <Pacman {gridSize} {pacmanRow} {pacmanCol}/>
  <Ghosts {gridSize} {redRow} {redCol} {pinkRow} {pinkCol} {blueRow} {blueCol} {orangeRow} {orangeCol}/>
</div>