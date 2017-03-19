<template lang="pug">
.v-split-pane
  .columns
    .column-left(:style="{width: leftWidth, display: leftDisplay}")
      .column-container
        slot(name="left")
    .column-right(:style="{width: rightWidth, display: rightDisplay}")
      .column-container
        slot(name="right")
</template>

<script>
export default {
  name: 'VerticalSplitPane',
  props: {
    ratio: {
      type: Number,
      default: 0.5,
      validator: (v) => v >= 0 && v <= 1
    },

    hideLimit: {
      type: Number,
      default: 0.01,
      validator: (v) => v >= 0 && v <= 1
    }
  },

  computed: {
    leftWidth() {
      return Math.floor(this.ratio * 100) + '%'
    },

    rightWidth() {
      return Math.floor((1 - this.ratio) * 100) + '%'
    },

    leftDisplay() {
      return this.ratio < this.hideLimit ? 'none' : 'block';
    },

    rightDisplay() {
      return 1 - this.ratio < this.hideLimit ? 'none' : 'block';
    }
  },

  data() {
    return {}
  }
}
</script>

<style lang="scss">
@import '../assets/styles/core.mixins';

.v-split-pane {
  @extend .reset;
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  text-align: left;
  vertical-align: top;

  .columns {
    @extend .reset;
    position: relative;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;

    .column-left {
      @extend .reset;
      position: absolute;
      left: 0;
      top: 0;
      bottom: 0;
      width: 50%;
    }

    .column-right {
      @extend .reset;
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
      width: 50%;
    }

    .column-right .column-container, .column-left .column-container {
      @extend .reset;
      position: relative;
      height: 100%;
      width: 100%;
    }
  }
}
</style>
