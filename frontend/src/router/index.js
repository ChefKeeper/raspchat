import Vue from 'vue'
import Router from 'vue-router'
import ChatPane from '@/components/ChatPane'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Home',
      component: ChatPane
    }
  ]
})
