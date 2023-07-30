<style>

  :root {
    font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;

    color-scheme: light dark;
    color: rgba(255, 255, 255, 0.87);
    background-color: #242424;

    font-synthesis: none;
    text-rendering: optimizeLegibility;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    -webkit-text-size-adjust: 100%;
  }

  .maze-space {

    /* Positioning */
    position: absolute;
    top: 5vh;
    left: 5vw;
  }

</style>

<script>

  // Imports
  import config from '../../config.json';
  import Maze from './lib/Maze.svelte';
  import Pellets from './lib/Pellets.svelte';
  import Pacman from './lib/Pacman.svelte';
  import Ghosts from './lib/Ghosts.svelte';
  import MpsCounter from './lib/MpsCounter.svelte';
  import Ticker from './lib/Ticker.svelte';

  // Creating a websocket client
  var socket = new WebSocket(`ws://${config.ServerIP}:${config.WebSocketPort}`);
  socket.binaryType = 'arraybuffer';
  let socketOpen = false;

  /* 
    This generates an empty array of pellet states 
    (0 = none, 1 = pellet, 2 = super)
  */
  let pelletGrid = [];
  for (let row = 0; row < 31; row++) {
    pelletGrid[row] = [];
    for (let col = 0; col < 28; col++) {
      pelletGrid[row][col] = 0;
    }
  }

  /*
    We use a circular queue (with a fixed max capacity to keep track of the times
    of the most recent messages. For every message we receive, we should add this
    time to the queue and remove all times longer than a millisecond ago. The 
    length of this array will be the MPS (messages per second), which should be
    synced with the frame rate of the game engine if there is no lag.
  */
  const MPS_BUFFER_SIZE = 60; // For higher fps, replace this
  let mpsBuffer = new Array(MPS_BUFFER_SIZE);
  let mpsIdxLeft = 0;
  let mpsIdxRight = 1;
  let mpsAvg = 0;
  mpsBuffer[0] = Date.now();

  // Keep track of the number of ticks elapsed (from the game engine)
  let currTicks = 0;

  // Handling a new connection
  socket.addEventListener('open', (_) => {
    console.log('WebSocket connection established');
    socket.send('Hello, server!');
    socketOpen = true;
  });

  // When the ticker is clicked, send a message
  let paused = false;
  $: {if (socketOpen) {
    console.log(paused); socket.send(paused);
  }}

  // Message events
  socket.addEventListener('message', (event) => {
    if (event.data instanceof ArrayBuffer) {

      // Increment the internal count of game engine ticks
      ++currTicks;

      // Log the time
      const ts = Date.now();
      mpsBuffer[mpsIdxRight] = ts;
      mpsAvg++;
      mpsIdxRight = (mpsIdxRight + 1) % MPS_BUFFER_SIZE;
      // Leeway of about 2%
      while (ts - mpsBuffer[mpsIdxLeft] > 1020 && mpsIdxLeft != mpsIdxRight) {
        mpsIdxLeft = (mpsIdxLeft + 1) % MPS_BUFFER_SIZE;
        mpsAvg--;
      }

      // Binary frame for manually parsing the input
      let view = new DataView(event.data);

      if (view) {

        // Keep track of the byte index we are reading
        let byteIdx = 0;

        // Get the game mode
        console.log(view.getUint8(byteIdx, false));
        //byteIdx++;

        // Parse pellet data
        for (let row = 0; row < 31; row++) {
          const binRow = view.getUint32(byteIdx, false);
          byteIdx += 4;
          for (let col = 0; col < 28; col++) {
            let superPellet = ((row === 3) || (row === 23)) && ((col === 1) || (col === 26));
            pelletGrid[row][col] = ((binRow >> col) & 1) ? (superPellet ? 2 : 1) : 0;
          }
        }

        // Trigger an update for the pellets
        pelletGrid = pelletGrid;
      }
    }
  });

  socket.addEventListener('close', (_) => {
    socketOpen = false;
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
  <Maze {gridSize} />
  <Pellets {pelletGrid} {gridSize} />
  <Pacman {gridSize} {pacmanRow} {pacmanCol} />
  <Ghosts {gridSize} {redRow} {redCol} {pinkRow} {pinkCol} {blueRow} {blueCol} {orangeRow} {orangeCol} />
  <MpsCounter {gridSize} {mpsAvg} />
  <Ticker {gridSize} {currTicks} bind:paused/>
</div>