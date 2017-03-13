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

<script>
export default {
  name: 'SplitPane',
  props: {
    split: {
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
      return Math.floor(this.split * 100) + '%'
    },

    rightWidth() {
      return Math.floor((1 - this.split) * 100) + '%'
    },

    leftDisplay() {
      return this.split < this.hideLimit ? 'none' : 'block';
    },

    rightDisplay() {
      return 1 - this.split < this.hideLimit ? 'none' : 'block';
    }
  },

  data() {
    return {}
  }
}
</script>
