<style>
  
  .ghost {
    position: absolute;
    width: 50px;
    height: 50px;
    clip-path: polygon(-100% -100%, 0% 100%, 18% 90%, 33% 100%, 50% 90%, 67% 100%, 82% 90%, 100% 100%, 200% -100%);
    border-top-left-radius: 50%;
    border-top-right-radius: 50%;
    display: flex;
  }

  .left-eye, .right-eye {
    position: absolute;
    min-width: 1px;
    min-height: 1px;
    width: 28%;
    height: 40%;
    top: 20%;
    left: 15%;
    border-radius: 50%;
    background-color: white;
    display: flex;
    flex: 1;
    justify-content: center;
    align-items: center;
  }

  .right-eye {
    left: 55%;
  }

  .iris {
    height: 60%;
    border-radius: 100%;
    background-color: blue;
    flex: 0.5;
  }

  .looking-up {
    transform: translate(0%, -30%);
  }

  .looking-down {
    transform: translate(0%, 30%);
  }

  .looking-left {
    transform: translate(-50%, 0%);
  }

  .looking-right {
    transform: translate(50%, 0%);
  }

  .red {
    background-color: red;
  }

  .pink {
    background-color: hotpink;
  }

  .blue {
    background-color: cyan;
  }

  .orange {
    background-color: orange;
  }

  .frightened .left-eye, .frightened .right-eye {
    background-color: transparent;
  }

  .frightened .iris {
    background-color: white;
    height: 40%;
    flex: 0.5;
  }

  .recovering .iris {
    background-color: red;
  }

  .frightened .looking-up, .frightened .looking-down, .frightened .looking-left, .frightened .looking-right {
    transform: none;
  }

  .frightened {
    background-color: blue;
  }

  .recovering {
    background-color: white;
  }

  .mouth {
    visibility: none;
  }

  .frightened .mouth {
    position: absolute;
    background-color: white;
    top: 60%;
    width: 100%;
    height: 40%;
    clip-path: polygon(14% 50%, 26% 30%, 38% 50%, 50% 30%, 62% 50%, 74% 30%, 86% 50%, 90% 40%, 74% 20%, 62% 40%, 50% 20%, 38% 40%, 26% 20%, 10% 40%);
  }

  .recovering .mouth {
    background-color: red;
  }

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
  let redLookDir = " looking-left";

  $: pinkPosX = gridSize * pinkCol;
  $: pinkPosY = gridSize * pinkRow;
  let pinkLookDir = " looking-right";

  $: bluePosX = gridSize * blueCol;
  $: bluePosY = gridSize * blueRow;
  let blueLookDir = " looking-up";

  $: orangePosX = gridSize * orangeCol;
  $: orangePosY = gridSize * orangeRow;
  let orangeLookDir = " looking-down";

  let frightened = false;
  let recovering = false;
  $: frightenedModifer = frightened ? (recovering ? 'frightened recovering' : 'frightened') : '';

</script>

<div class="ghost red {frightenedModifer}" 
     style:width="{gridSize}px" 
     style:height="{gridSize}px" 
     style:left="{redPosX}px"
     style:top="{redPosY}px">
  
  <div class="left-eye">
    <div class="iris {redLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {redLookDir}"/>
  </div>

  <div class="mouth {frightenedModifer}"/>

</div>

<div class="ghost pink {frightenedModifer}" 
     style:width="{gridSize}px" 
     style:height="{gridSize}px" 
     style:left="{pinkPosX}px"
     style:top="{pinkPosY}px">
  
  <div class="left-eye">
    <div class="iris {pinkLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {pinkLookDir}"/>
  </div>

  <div class="mouth {frightenedModifer}"/>

</div>

<div class="ghost blue {frightenedModifer}" 
     style:width="{gridSize}px" 
     style:height="{gridSize}px" 
     style:left="{bluePosX}px"
     style:top="{bluePosY}px">
  
  <div class="left-eye">
    <div class="iris {blueLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {blueLookDir}"/>
  </div>

  <div class="mouth {frightenedModifer}"/>

</div>

<div class="ghost orange {frightenedModifer}" 
     style:width="{gridSize}px" 
     style:height="{gridSize}px" 
     style:left="{orangePosX}px"
     style:top="{orangePosY}px">
  
  <div class="left-eye">
    <div class="iris {orangeLookDir}"/>
  </div>

  <div class="right-eye">
    <div class="iris {orangeLookDir}"/>
  </div>

  <div class="mouth {frightenedModifer}"/>

</div>