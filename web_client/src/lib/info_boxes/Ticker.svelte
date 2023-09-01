<style>

  /* Information box for ticker */
  .ticker-box {

    /* Positioning */
    position: absolute;
    text-align: center;
    z-index: 1;

    /* Formatting */
    background-color: rgba(0,0,0,0.3);
    border: none;
    cursor: pointer;

    /* Grid-size related CSS */
    width:       calc(5   * var(--grid-size));
    height:      calc(3   * var(--grid-size));
    line-height: calc(3   * var(--grid-size));
    left:        calc(0   * var(--grid-size));
    top:         calc(10  * var(--grid-size));
    font-size:   calc(0.9 * var(--grid-size));
  }
  
  /* Info box outline (clicked), should be none */
  .ticker-box:focus-visible {
    outline: none;
  }
  
  /* Ticker shape */
  .ticker {
    
    /* Positioning */
    position: absolute;
    z-index: 2;

    /* Formatting */
    opacity: 0.8;

    /* Grid-size related CSS */
    width:  calc(2   * var(--grid-size) + 2 * var(--pad));
    height: calc(2   * var(--grid-size) + 2 * var(--pad));
    top:    calc(0.5 * var(--grid-size) - var(--pad));
    left:   calc(1.5 * var(--grid-size) - var(--pad));
  }

  /* Ticker circular outline */
  circle {
    fill: transparent;
    stroke: var(--color);
    stroke-width: var(--pad);
  }

  /* Filled circular ticker path */
  path {
    fill: var(--color);
  }

</style>

<script>

  // Grid size attributes
  export let gridSize;
  $: pad = gridSize/20;

  // Simulated key event (space, toggle pause) when clicking the ticker box
  export let togglePause;

  // Decide the color of the ticker based on the game-mode
  export let Modes;
  export let gameMode;
  let modeColor = 'yellow'
  $: {
    if (gameMode === Modes.Paused) {
      modeColor = 'gray';
    } else if (gameMode === Modes.Scatter) {
      modeColor = 'lime';
    } else if (gameMode === Modes.Chase) {
      modeColor = 'yellow';
    } else if (gameMode === Modes.Offline) {
      modeColor = '#333';
    } else { // The ticker should not be red ever - if it is, there's a bug
      modeColor = 'red';
    }
  }

  // Math to calculate relevant degree measures, lengths, and flags 
  export let modTicks;
  export let updatePeriod;
  const degToRad = Math.PI / 180
  $: degrees     = 360 * modTicks / updatePeriod;
  $: cosine      = Math.cos(degToRad * degrees);
  $: sine        = Math.sin(degToRad * degrees);
  $: longArcFlag = (degrees > 180) ? 1 : 0; // Reflexive angle condition

</script>

<!-- Interactive button to toggle pausing -->
<button
  class='ticker-box'
  style:--grid-size='{gridSize}px'
  on:click={() => togglePause()}
>

  <!-- Ticker shape -->
  <svg
    class='ticker'
    style:--grid-size='{gridSize}px'
    style:--pad='{pad}px'
    style:--color='{modeColor}'
  >
    
    <!-- Circular outline of ticker -->
    <circle
      cx='{gridSize+pad}'
      cy='{gridSize+pad}'
      r='{gridSize}'
    />

    <!-- Filled ticker region -->
    <path 
      d=' M {gridSize+pad} {gridSize+pad} 
          L {gridSize+pad} {pad}
          A {gridSize} {gridSize} 0 
            {longArcFlag} 1 
            {gridSize + gridSize * sine + pad} 
            {gridSize - gridSize * cosine + pad} 
          z' 
    />
  </svg>
</button>