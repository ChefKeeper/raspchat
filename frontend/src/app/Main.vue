<template lang="pug">
#main
  split-pane(ratio="0.2")
    p(slot="left") Left Panel
    p(slot="right") Right Panel
</template>

<script>
import SplitPane from './components/SplitPane.vue'

export default {
  name: 'Main',
  components: {
    'split-pane': SplitPane
  }
}
</script>
