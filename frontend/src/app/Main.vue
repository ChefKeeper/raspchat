<template lang="pug">
#main
  split-pane
    p(slot="left") Left Panel
    p(slot="right") Right Panel
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
