<template lang="pug">
#main
  split-pane(:ratio="0.2")
    p(slot="left") Left Panel
    p(slot="right") Right Panel
</template>

<style lang="scss">
 #main {
    margin: 0;
    height: 100%;
    font-family: 'Open Sans', sans-serif;
  }
</style>

<script>
import SplitPane from '@/components/SplitPane.vue'

export default {
  name: 'chat-pane',
  components: {
    'split-pane': SplitPane
  }
}
</script>
