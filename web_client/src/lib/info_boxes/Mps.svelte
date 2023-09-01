<!-- MPS = Messages Per Second -->

<style>

  /* Control of header margins */
  h2 {
    margin: 0;
  }

  /* Information box for MPS Counter */
  .mps-counter {

    /* Positioning */
    position: absolute;
    text-align: center;

    /* Formatting */
    background-color: rgba(0,0,0,0.3);

    /* Grid-size related CSS */
    width:       calc(5   * var(--grid-size));
    height:      calc(3   * var(--grid-size));
    line-height: calc(3   * var(--grid-size));
    left:        calc(23  * var(--grid-size));
    top:         calc(10  * var(--grid-size));
    font-size:   calc(0.9 * var(--grid-size));
  }

</style>

<script>
  export let gridSize;
  export let mpsAvg;
</script>

<h2
  class='mps-counter'
  style:--grid-size='{gridSize}px'
>
  MPS: {mpsAvg}
</h2>