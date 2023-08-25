<style>

  /* Control of header margins */
  h2.score-box {
    margin: 0;
  }

  /* Information box for levels and score */
  .score-box {

    /* Positioning */
    position: absolute;
    text-align: center;
    z-index: 1;

    /* Formatting */
    background-color: rgba(0,0,0,0.3);
    border: none;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    /* Grid-size related CSS */
    width:       calc(5   * var(--grid-size));
    height:      calc(3   * var(--grid-size));
    line-height: calc(1.1 * var(--grid-size));
    left:        calc(0   * var(--grid-size));
    top:         calc(16  * var(--grid-size));
    font-size:   calc(0.8 * var(--grid-size));
  }

</style>

<script>
  export let gridSize;
  export let currScore;
  export let currLevel;
</script>

<h2
  class='score-box'
  style:--grid-size="{gridSize}px"
>
  <div>Level {currLevel}</div>
  <div>{currScore}</div>
</h2>