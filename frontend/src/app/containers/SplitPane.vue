<template lang="pug">
.split-pane
  .columns
    .column-left
      slot(name="left")
    .column-right
      slot(name="right")
</template>

<script>
export default {
  name: 'SplitPane',
  data() {
    return {}
  }
}
</script>
