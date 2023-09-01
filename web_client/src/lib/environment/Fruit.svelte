<style>

  /* Cherry SVG Sprite */
  .cherry-svg {
    position: absolute;
    width:  var(--grid-size);
    height: var(--grid-size);
    transform: scale(1.5);
  }

  /* Cherry */
  .cherry {
    fill: red;
    stroke: #242424;
    stroke-width: calc(var(--grid-size)/50);
  }

  /* Stem */
  .stem {
    fill: transparent;
    stroke: #dd9751;
    stroke-width: calc(var(--grid-size)/25);
  }

  /* Indentation */
  .indentation {
    fill: transparent;
    stroke: darkred;
    stroke-width: calc(var(--grid-size)/50);
  }

  /* Reflection */
  .reflection {
    fill: #ffbbbb;
    stroke: none;
  }

</style>

<script>

  // Grid size, same as for other components
  export let gridSize;

  // Fruit state
  let fruitRowState = 17 | 32;
  let fruitColState = 13;

  // Using the & operator to pick out the 5 lowest bits
  $: posX = fruitColState & 0b11111
  $: posY = fruitRowState & 0b11111

  // Hide the fruit if bit 5 (32) of either coordinate is set
  $: showFruit = ((fruitRowState | fruitColState) & 0b100000) ? false : true;

</script>

{#if showFruit}
  <svg class='cherry-svg'
      style:--grid-size='{~~gridSize+1}px'
      style:left='{gridSize * posX}px'
      style:top='{gridSize * posY}px'
  >

    <circle
      cx='{5*gridSize/16}'
      cy='{5*gridSize/8}'
      r='{gridSize/5}'
      class='cherry'
    />

    <ellipse
      cx='{6*gridSize/16}'
      cy='{17*gridSize/32}'
      rx='{gridSize/16}'
      ry='{gridSize/32}'
      class='indentation'
      style:transform-origin='{6*gridSize/16}px {17*gridSize/32}px'
      style:transform='rotate(35deg)'
    />

    <ellipse
      cx='{3*gridSize/16}'
      cy='{21*gridSize/32}'
      rx='{gridSize/32}'
      ry='{2*gridSize/32}'
      class='reflection'
      style:transform-origin='{3*gridSize/16}px {21*gridSize/32}px'
      style:transform='rotate(-20deg)'
    />

    <circle
      cx='{5*gridSize/8}'
      cy='{6*gridSize/8}'
      r='{gridSize/5}'
      class='cherry'
    />

    <ellipse
      cx='{10*gridSize/16}'
      cy='{10*gridSize/16}'
      rx='{gridSize/16}'
      ry='{gridSize/32}'
      class='indentation'
      style:transform-origin='{10*gridSize/16}px {10*gridSize/16}px'
      style:transform='rotate(5deg)'
    />

    <ellipse
      cx='{4*gridSize/8}'
      cy='{25*gridSize/32}'
      rx='{gridSize/32}'
      ry='{2*gridSize/32}'
      class='reflection'
      style:transform-origin='{4*gridSize/8}px {25*gridSize/32}px'
      style:transform='rotate(-15deg)'
    />

    <circle
      cx='{27*gridSize/32}'
      cy='{5*gridSize/32}'
      r='{gridSize/32}'
      style:fill='yellow'
    />

    <path
      d=' M {6*gridSize/16} {17*gridSize/32}
          A {gridSize} {gridSize} 0
            0 1
            {13 * gridSize / 16} 
            {3  * gridSize / 16}
          V {2  * gridSize / 16}
          H {14 * gridSize / 16}
          V {3  * gridSize / 16}
          H {13 * gridSize / 16}
          A {gridSize} {gridSize} 0
            0 0
            {10 * gridSize / 16} 
            {10 * gridSize / 16}
          '
      class='stem'
    />

  </svg>
{/if}