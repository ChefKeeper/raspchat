<template lang="pug">
.h-split-pane
  .rows
    .row-top(:style="{height: topHeight, display: topDisplay}")
      .row-container
        slot(name="top")
    .row-bottom(:style="{height: bottomHeight, display: bottomDisplay}")
      .row-container
        slot(name="bottom")
</template>

<script>
export default {
  name: 'HorizontalSplitPane',
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
    topHeight() {
      return Math.floor(this.ratio * 100) + '%'
    },

    bottomHeight() {
      return Math.floor((1 - this.ratio) * 100) + '%'
    },

    topDisplay() {
      return this.ratio < this.hideLimit ? 'none' : 'block';
    },

    bottomDisplay() {
      return 1 - this.ratio < this.hideLimit ? 'none' : 'block';
    }
  },

  data() {
    return {}
  }
}
</script>

<style lang="scss">
.h-split-pane {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;

  .rows {
    position: relative;
    display: block;
    clear: both;
    width: 100%;
    height: 100%;

    .row-top, .row-bottom {
      position: absolute;
      left: 0;
      right: 0;

      .row-container {
        position: relative;
        display: block;
        width: 100%;
        height: 100%;
      }
    }

    .row-top {
      top: 0;
    }

    .row-bottom {
      bottom: 0;
    }
  }
}
</style>
