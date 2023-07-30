<style>

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
  .ticker {
    
    /* Positioning */
    position: absolute;
    z-index: 2;

    /* Formatting */
    opacity: 0.8;

    /* Grid-size related CSS */
    width:  calc(2    * var(--grid-size) + 2 * var(--pad));
    height: calc(2    * var(--grid-size) + 2 * var(--pad));
    top:    calc(0.5 * var(--grid-size) - var(--pad));
    left:   calc(1.5  * var(--grid-size) - var(--pad));
  }

  circle {
    fill: transparent;
    stroke: yellow;
    stroke-width: var(--pad);
  }

  path {
    fill: yellow;
  }

</style>

<script>

  // Grid size attributes
  export let gridSize;
  $: pad = gridSize/20;

  // Pausing event when the ticker gets clicked
  export let paused;
  function togglePause() {
    paused = !paused;
  }

  // Math to calculate the degree measures and flags for the ticker object
  export let currTicks;
  let updateTicks = 12;
  $: degrees = 360 * (currTicks % updateTicks) / updateTicks;
  $: cosine = Math.cos(Math.PI / 180 * degrees);
  $: sine = Math.sin(Math.PI / 180 * degrees);
  $: longArcFlag = (degrees > 180) ? 1 : 0;

</script>

<button class='ticker-box' style:--grid-size='{gridSize}px' on:click={() => togglePause()}>
  <svg class='ticker' style:--grid-size='{gridSize}px' style:--pad='{pad}px'>
    <circle cx="{gridSize+pad}" cy="{gridSize+pad}" r="{gridSize}"/>
    <path d="M {gridSize+pad} {gridSize+pad} 
            L {gridSize+pad} {pad}
            A {gridSize} {gridSize} 0 {longArcFlag} 1 {gridSize + gridSize * sine + pad} {gridSize - gridSize * cosine + pad} 
            z" />
  </svg>
</button>