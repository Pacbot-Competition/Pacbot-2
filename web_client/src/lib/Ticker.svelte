<style>

  .ticker-div {
    position: absolute;
    opacity: 0.8;
    width:  var(--grid-size);
    height: var(--grid-size);
    top:    calc(10.5 * var(--grid-size));
    left:   calc(1.5 * var(--grid-size));
  }

  circle {
    fill: transparent;
    stroke: yellow;
    stroke-width: 0.5;
  }

  path {
    fill: yellow;
  }

  svg.pie {
    width: 230px;
    height: 230px;
  }

</style>

<script>
  export let gridSize;
  export let currTicks;
  let updateTicks = 12;
  $: degrees = 360 * (currTicks % updateTicks) / updateTicks;
  $: cosine = Math.cos(Math.PI / 180 * degrees);
  $: sine = Math.sin(Math.PI / 180 * degrees);
  $: longArcFlag = (degrees > 180) ? 1 : 0;
</script>

<div class='ticker-div' style:--grid-size='{gridSize}px'>
  <div class='ticker'>
    <svg class="pie">
      <circle cx="{gridSize}" cy="{gridSize}" r="{gridSize}"></circle>
      <path d="M{gridSize} {gridSize} 
              L{gridSize} 0
              A{gridSize} {gridSize} 0 {longArcFlag} 1 {gridSize + gridSize * sine} {gridSize - gridSize * cosine} 
              z"></path>
    </svg>
  </div>
</div>