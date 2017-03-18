<template lang="pug">
.split-pane
  .columns
    .column-left(v-bind:style="{width: leftWidth, display: leftDisplay}")
      .column-container
        slot(name="left")
    .column-right(v-bind:style="{width: rightWidth, display: rightDisplay}")
      .column-container
        slot(name="right")
</template>


<style lang="scss">

@import '../../styles/split-pane';

#main {
    width: 100%;
    height: 100%;
    @include panel-skin(#304ffe, #ffffff);
}
</style>

<script>
export default {
  name: 'SplitPane',
  props: {
    ratio: {
      type: Number,
      default: 0.5,
      validator: (v) => v >= 0 && v <= 1
    },

    hideLimit: {
      type: Number,
      default: 0.1,
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
