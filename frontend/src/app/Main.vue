<template lang="pug">
split-pane.root
  p(slot="left") Hello
  p(slot="right") World
</template>

<script>
import SplitPane from './containers/SplitPane.vue'

export default {
  name: 'Main',
  components: {
    'split-pane': SplitPane
  }
}
</script>
