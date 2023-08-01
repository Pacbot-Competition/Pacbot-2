<style>

  .ghost-svg {
    position: absolute;
    width:  var(--grid-size);
    height: var(--grid-size);
  }

  /* Blinky */
  path.red { fill: red; }

  /* Pinky */
  path.pink { fill: hotpink; }

  /* Inky */
  path.cyan { fill: turquoise; }

  /* Clyde */
  path.orange { fill: orange; }

  /* Eye */
  ellipse.white {
    fill: white;
  }

  /* Iris */
  ellipse.blue {
    fill: blue;
  }

</style>

<script>
  export let gridSize;
  export let color;
  export let rowState;
  export let colState;

  /* TODO: Implement SVG versions of frightened sprites */
  let frightenedCounter = 0;
  let frightenedModifer = 0;

  /* 
    The last 5 bits of each state byte are the position, 
    while the first 2 bits of are the signed direction
  */

  // Using the & operator to pick out the 5 lowest bits
  $: posX = colState & 0b11111
  $: posY = rowState & 0b11111
  
  // The below code is a sign-extention trick, taking advantage of 32-bit
  // integer representations in JavaScript
  $: dirX = ((colState >> 6) << 30) >> 30
  $: dirY = ((rowState >> 6) << 30) >> 30

</script>

<!-- SVG Sprite of Ghost -->
<svg class='ghost-svg' 
     style:--grid-size='{~~gridSize+1}px'
     style:top='{posY*gridSize}px' 
     style:left='{posX*gridSize}px'>

  <!-- Body of ghost -->
  <path d='M {0} {gridSize/2}
           A {gridSize/2} {gridSize/2} 0 0 1 {gridSize} {gridSize/2} 
           L {gridSize} {gridSize}
           L {0.82 * gridSize} {0.9 * gridSize}
           L {0.67 * gridSize} {      gridSize}
           L {0.50 * gridSize} {0.9 * gridSize}
           L {0.33 * gridSize} {      gridSize}
           L {0.18 * gridSize} {0.9 * gridSize}
           L {0    * gridSize} {      gridSize}
           z' class={color}/>

  <!-- Left eye -->
  <ellipse cx='{0.30 * gridSize}' 
           cy='{0.40 * gridSize}' 
           rx='{0.14 * gridSize}' 
           ry='{0.20 * gridSize}' 
           class='white'/>

  <!-- Right eye -->
  <ellipse cx='{0.70 * gridSize}' 
           cy='{0.40 * gridSize}' 
           rx='{0.14 * gridSize}' 
           ry='{0.20 * gridSize}' 
           class='white'/>
  
  <!-- Left iris -->
  <ellipse cx='{(0.30 + 0.06*dirX) * gridSize}' 
           cy='{(0.40 + 0.09*dirY) * gridSize}' 
           rx='{0.07 * gridSize}' 
           ry='{0.10 * gridSize}' 
           class='blue'/>

  <!-- Right iris -->
  <ellipse cx='{(0.70 + 0.06*dirX) * gridSize}' 
           cy='{(0.40 + 0.09*dirY) * gridSize}' 
           rx='{0.07 * gridSize}' 
           ry='{0.10 * gridSize}' 
           class='blue'/>
  
</svg>