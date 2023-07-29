<style>

  .ticker-box {

    /* Positioning */
    position: absolute;
    text-align: center;
    z-index: 1;

    /* Formatting */
    background-color: rgba(0,0,0,0.3);

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
    top:    calc(10.5 * var(--grid-size) - var(--pad));
    left:   calc(1.5  * var(--grid-size) - var(--pad));
  }

  circle {
    fill: transparent;
    stroke: yellow;
    stroke-width: 1;
  }

  path {
    fill: yellow;
  }

</style>

<script>
  export let gridSize;
  export let currTicks;
  let updateTicks = 12;
  const pad = 0.5;
  $: degrees = 360 * (currTicks % updateTicks) / updateTicks;
  $: cosine = Math.cos(Math.PI / 180 * degrees);
  $: sine = Math.sin(Math.PI / 180 * degrees);
  $: longArcFlag = (degrees > 180) ? 1 : 0;
</script>

<div class='ticker-box' style:--grid-size='{gridSize}px'/>
<svg class='ticker' style:--grid-size='{gridSize}px' style:--pad='{pad}px'>
  <circle cx="{gridSize+pad}" cy="{gridSize+pad}" r="{gridSize}" />
  <path d="M {gridSize+pad} {gridSize+pad} 
           L {gridSize+pad} {pad}
           A {gridSize} {gridSize} 0 {longArcFlag} 1 {gridSize + gridSize * sine + pad} {gridSize - gridSize * cosine + pad} 
           z" />
</svg>