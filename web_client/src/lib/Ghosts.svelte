<style>
  
  /* Ghost "sprite", designed in CSS */
  .ghost {

    /* Positioning */
    position: absolute;

    /* Formatting */
    clip-path: polygon(-100% -100%, 0% 100%, 18% 90%, 33% 100%, 50% 90%, 67% 100%, 82% 90%, 100% 100%, 200% -100%);
    border-top-left-radius: 50%;
    border-top-right-radius: 50%;
    display: flex;

    /* Grid-size related CSS */
    width: var(--grid-size);
    height: var(--grid-size);
  }

  /* Eyes */
  .left-eye, .right-eye {

    /* Positioning */
    position: absolute;
    top: 20%;
    left: 15%;
    justify-content: center;
    align-items: center;

    /* Formatting */
    width: 28%;
    height: 40%;
    border-radius: 50%;
    background-color: white;
    display: flex;
  }

  /* Right eye is horizontally translated */
  .right-eye {
    transform: translateX(calc(0.4*var(--grid-size)));
  }

  /* Ghost iris (blue part of the eye) */
  .iris {
    height: 60%;
    border-radius: 100%;
    background-color: blue;
    flex: 0.5;
  }

  /* Eye directions */
  .up    { transform: translate(0%, -30%); }
  .down  { transform: translate(0%,  30%); }
  .left  { transform: translate(-50%, 0%); }
  .right { transform: translate(50%,  0%); }

  /* Ghost colors */
  .red    { background-color: red; }
  .pink   { background-color: hotpink; }
  .blue   { background-color: cyan; }
  .orange { background-color: orange; }

  /* Frightened modifiers */
  .fr { background-color: blue; }
  .fr .left-eye, .fr .right-eye { background-color: transparent; }
  .fr .iris { background-color: white; height: 40%; flex: 0.5; }
  .fr .up, .fr .down, .fr .left, .fr .right { transform: none; }

  /* Recovering modifiers */
  .rc { background-color: white; }
  .rc .iris { background-color: red; }

  /* Ghost mouth - only shown when frightened */
  .mouth { visibility: none; }

  .fr .mouth {

    /* Positioning */
    position: absolute;
    top: 60%;

    /* Formatting */
    background-color: white;
    width: 100%;
    height: 40%;
    clip-path: polygon(14% 50%, 26% 30%, 38% 50%, 50% 30%, 62% 50%, 74% 30%, 86% 50%, 90% 40%, 74% 20%, 62% 40%, 50% 20%, 38% 40%, 26% 20%, 10% 40%);
  }

  .rc .mouth { background-color: red; }

</style>

<script>
  export let gridSize;
  export let redRow;
  export let redCol;
  export let pinkRow;
  export let pinkCol;
  export let blueRow;
  export let blueCol;
  export let orangeRow;
  export let orangeCol;

  $: redPosX = gridSize * redCol;
  $: redPosY = gridSize * redRow;
  let redLookDir = " left";

  $: pinkPosX = gridSize * pinkCol;
  $: pinkPosY = gridSize * pinkRow;
  let pinkLookDir = " right";

  $: bluePosX = gridSize * blueCol;
  $: bluePosY = gridSize * blueRow;
  let blueLookDir = " up";

  $: orangePosX = gridSize * orangeCol;
  $: orangePosY = gridSize * orangeRow;
  let orangeLookDir = " down";

  let frightened = false;
  let recovering = false;
  $: frightenedModifer = frightened ? (recovering ? 'fr rc' : 'fr') : '';

</script>

<div class="ghost red {frightenedModifer}"
     style:--grid-size="{gridSize}px"
     style:left="{redPosX}px"
     style:top="{redPosY}px">
  
  <div class="left-eye">
    <div class="iris {redLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {redLookDir}"/>
  </div>

  <div class="mouth"/>

</div>

<div class="ghost pink {frightenedModifer}" 
     style:--grid-size="{gridSize}px"
     style:left="{pinkPosX}px"
     style:top="{pinkPosY}px">
  
  <div class="left-eye">
    <div class="iris {pinkLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {pinkLookDir}"/>
  </div>

  <div class="mouth"/>

</div>

<div class="ghost blue {frightenedModifer}" 
     style:--grid-size="{gridSize}px"
     style:left="{bluePosX}px"
     style:top="{bluePosY}px">
  
  <div class="left-eye">
    <div class="iris {blueLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {blueLookDir}"/>
  </div>

  <div class="mouth"/>

</div>

<div class="ghost orange {frightenedModifer}" 
     style:--grid-size="{gridSize}px"
     style:left="{orangePosX}px"
     style:top="{orangePosY}px">
  
  <div class="left-eye">
    <div class="iris {orangeLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {orangeLookDir}"/>
  </div>

  <div class="mouth"/>

</div>