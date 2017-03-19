<template lang="pug">
#main
  v-split-pane(:ratio="0.2")
    p(slot="left") Left Panel
    h-split-pane(slot="right", :ratio="0.9")
      chat-log-pane(slot="top") Hello
      textarea(slot="bottom", style="width: 100%; height: 100%")
</template>

<script>
import VerticalSplitPane from '@/components/VerticalSplitPane'
import BottomUpScrollPane from '@/components/BottomUpScrollPane'
import HorizontalSplitPane from '@/components/HorizontalSplitPane'

export default {
  name: 'ChatView',
  components: {
    'v-split-pane': VerticalSplitPane,
    'h-split-pane': HorizontalSplitPane,
    'chat-log-pane': BottomUpScrollPane
  }
}
</script>

<style lang="scss">
 #main {
    margin: 0;
    height: 100%;
    font-family: 'Open Sans', sans-serif;
  }
</style>
