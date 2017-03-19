<template lang="pug">
div.scroll-container(:class="{streched: fill}")
  div.content-wrapper
    div.content-container
      slot
</template>

<script>
export default {
  props: {
    fill: {
      type: Boolean,
      default: true
    }
  }
}
</script>

<style lang="scss" scoped>
.scroll-container {
  position: relative;
  display: block;
  clear: both;
  overflow: auto;

  &.streched {
    width: 100%;
    height: 100%;
  }

  .content-wrapper {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
  }

  .content-container {
    position: relative;
    display: block;
  }
}
</style>
