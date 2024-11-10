<style>

  /* Information box for ticker */
  .reset-box {

    /* Positioning */
    position: absolute;
    text-align: center;
    z-index: 1;

    /* Formatting */
    background-color: rgba(82, 165, 37, 0.3);
    border: none;
    cursor: pointer;

    /* Grid-size related CSS */
    width:       calc(5   * var(--grid-size));
    height:      calc(3   * var(--grid-size));
    line-height: calc(3   * var(--grid-size));
    left:        calc(0   * var(--grid-size));
    top:         calc(5  * var(--grid-size));
    font-size:   calc(0.9 * var(--grid-size));
  }

  /* Info box outline (clicked), should be none */
  .reset-box:focus-visible {
    outline: 5px solid rgb(255, 7, 7); /* or any color and width you prefer */
} 

.reset-box:active {
    outline: 2px solid #FF0000; /* Red outline when clicked */
  }

</style>

<script>

  // Grid size attributes
  export let gridSize;
  $: pad = gridSize/20;

  // Simulated key event (space, toggle pause) when clicking the ticker box
  export let toggleReset;

</script>

<!-- Interactive button to toggle pausing -->
<button
  class='reset-box'
  style:--grid-size='{gridSize}px'
  on:click={() => toggleReset()}
>

</button>